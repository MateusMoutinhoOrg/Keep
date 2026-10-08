package filestorage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
	storagedeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/storagedeps"
)

// Options tunes a store built by NewWithOptions. Its zero value is what New
// builds.
type Options struct {
	// NoSync skips the fsync of every written file and of the directory it
	// is renamed into. A write is still atomic — a reader or a crashed
	// process sees the old value or the new one, never a torn one — but a
	// power loss may undo writes the store had already reported done, and
	// may undo them out of order. Leave it false unless the data can be
	// rebuilt.
	NoSync bool
}

// store is the state one bound adapter keeps: the directory every key is
// resolved under, whether writes are flushed to the disk, and the mutex the
// read-modify-write operations take.
type store struct {
	base string
	sync bool
	mu   sync.Mutex
}

// maxSegment is the longest an encoded segment may grow before it is hashed
// instead: most filesystems refuse a name longer than 255 bytes.
const maxSegment = 200

// writeAttempts bounds how often a write recreates the directory it writes
// into: another writer's prune may remove an empty directory between the
// MkdirAll and the write.
const writeAttempts = 5

// takeoverStale is how old a takeover guard has to be before it is treated
// as left behind by a crashed process.
const takeoverStale = 10 * time.Second

// hexDigits are the digits of the %XX escape. They are upper-case, while
// every byte kept as is is lower-case, so the encoding stays injective on a
// filesystem that ignores case.
const hexDigits = "0123456789ABCDEF"

// encodeSegment turns one key segment into one file name, injectively and
// safely on a filesystem that ignores case:
//
//   - a-z, 0-9, '-' and '_' are kept as they are;
//   - every other byte — an upper-case letter, '.', '/', '%', any non-ASCII
//     byte — becomes %XX, with upper-case hex digits;
//   - the empty segment becomes "%", which no other segment encodes to;
//   - an encoding longer than maxSegment becomes "%%" followed by the
//     SHA-256 of the segment, which no other encoding starts with.
//
// Since no encoded segment holds a '.', none of them is "." or "..", none
// can reach outside the base directory, and none can end in the suffixes the
// store gives its own lock and temporary files.
func encodeSegment(segment string) string {
	if segment == "" {
		return "%"
	}
	var encoded strings.Builder
	for index := 0; index < len(segment); index++ {
		character := segment[index]
		if character >= 'a' && character <= 'z' || character >= '0' && character <= '9' || character == '-' || character == '_' {
			encoded.WriteByte(character)
			continue
		}
		encoded.WriteByte('%')
		encoded.WriteByte(hexDigits[character>>4])
		encoded.WriteByte(hexDigits[character&15])
	}
	if encoded.Len() > maxSegment {
		sum := sha256.Sum256([]byte(segment))
		return "%%" + hex.EncodeToString(sum[:])
	}
	return encoded.String()
}

// path resolves a key to a file path: one directory per segment, each
// encoded on its own by encodeSegment, so a segment may hold any character
// — a slash or a dot included — without ever escaping the base directory or
// colliding with another key. Encoding per segment is what makes ["a/b"] and
// ["a", "b"] two different files, and ["Name"] and ["name"] two different
// files even where the filesystem ignores case.
func (s *store) path(key []string) string {
	parts := make([]string, 0, len(key)+1)
	parts = append(parts, s.base)
	for _, segment := range key {
		parts = append(parts, encodeSegment(segment))
	}
	return filepath.Join(parts...)
}

// lockPath is the file holding the lease of a key: the key's own path with
// a suffix no encoded segment can produce, since none holds a '.'.
func (s *store) lockPath(key []string) string {
	return s.path(key) + ".keeplock"
}

// prune removes the directories a deleted key leaves empty, deepest first,
// so removing a record leaves nothing behind it on disk. It stops at the
// first directory still holding another key — os.Remove refuses a non-empty
// one — and never removes the base directory itself.
func (s *store) prune(key []string) {
	base := s.path(nil)
	for end := len(key) - 1; end > 0; end-- {
		dir := s.path(key[:end])
		if dir == base || os.Remove(dir) != nil {
			return
		}
	}
}

// writeTemp writes value to a new temporary file in dir — creating dir when
// it is missing — flushes it when durable, and returns its path. The name
// starts with a '.', which no encoded segment does.
func (s *store) writeTemp(dir string, value []byte, durable bool) (string, error) {
	var file *os.File
	for attempt := 1; ; attempt++ {
		created, err := os.CreateTemp(dir, ".keep-tmp-*")
		if err == nil {
			file = created
			break
		}
		if !errors.Is(err, os.ErrNotExist) || attempt == writeAttempts {
			return "", err
		}
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return "", err
		}
	}
	name := file.Name()
	if _, err := file.Write(value); err != nil {
		file.Close()
		os.Remove(name)
		return "", err
	}
	if durable {
		if err := file.Sync(); err != nil {
			file.Close()
			os.Remove(name)
			return "", err
		}
	}
	if err := file.Close(); err != nil {
		os.Remove(name)
		return "", err
	}
	return name, nil
}

// syncDir flushes a directory, so a rename or a link into it survives a
// power loss. Some platforms refuse to sync a directory; the file itself is
// already flushed, so that refusal is not an error.
func (s *store) syncDir(dir string) {
	handle, err := os.Open(dir)
	if err != nil {
		return
	}
	_ = handle.Sync()
	handle.Close()
}

// replace stores value under path atomically: written whole to a temporary
// file beside it, flushed, then renamed over it. A reader, and a process
// that crashes part-way, sees the old value or the new one — never an empty
// or a torn file.
func (s *store) replace(key []string, path string, value []byte) error {
	dir := filepath.Dir(path)
	temp, err := s.writeTemp(dir, value, s.sync)
	if err != nil {
		return err
	}
	if err := os.Rename(temp, path); err != nil {
		os.Remove(temp)
		if isDir(path) {
			return prefixError(key)
		}
		return err
	}
	if s.sync {
		s.syncDir(dir)
	}
	return nil
}

// create stores value under path only when nothing is there, atomically and
// with its whole content: written to a temporary file, then hard-linked into
// place, which fails rather than overwrite. created is false when path
// already exists. durable flushes it to the disk first; a lease needs no
// flush, since a power loss ends its holder too.
func (s *store) create(path string, value []byte, durable bool) (bool, error) {
	dir := filepath.Dir(path)
	temp, err := s.writeTemp(dir, value, durable)
	if err != nil {
		return false, err
	}
	defer os.Remove(temp)
	if err := os.Link(temp, path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return false, nil
		}
		return false, err
	}
	if durable {
		s.syncDir(dir)
	}
	return true, nil
}

// read returns the content of path. found is false when nothing is there,
// and when a directory is: a directory is the prefix of other keys, not a
// key of its own.
func read(path string) ([]byte, bool, error) {
	value, err := os.ReadFile(path)
	if err == nil {
		return value, true, nil
	}
	if errors.Is(err, os.ErrNotExist) || isDir(path) {
		return nil, false, nil
	}
	return nil, false, err
}

// isDir reports whether path is a directory.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// prefixError is the failure of a write to a key that other keys extend: on a
// filesystem a name is a file or a directory, never both, so ["a"] cannot
// hold a value while ["a", "b"] does. Keep's own layout never needs both.
func prefixError(key []string) error {
	return errors.New("keep: key " + describe(key) + " is the prefix of other keys and cannot hold a value")
}

// describe renders a key for an error message, in the ["a", "b.txt"] ->
// "a/b.txt" form the whole contract reads keys as.
func describe(key []string) string {
	return strings.Join(key, "/")
}

// Bind fills deps.Deps.StorageDeps with the filesystem implementation: one
// file per key, rooted at the working directory of the process, so what a
// database writes survives the process and can be read with `ls`. Where the
// tree lands therefore depends on where the program is started; a program
// that wants it somewhere fixed assigns New(base) itself.
func Bind(deps *deps.Deps) {
	deps.StorageDeps = New(".")
}

// New builds the contract over the given base directory, flushing every
// write to the disk before it returns. Bind uses ".", so a Props.Path
// carries the whole of where a database lands; a program that wants the tree
// somewhere else assigns the result of this call itself.
func New(base string) storagedeps.Contract {
	return NewWithOptions(base, Options{})
}

// NewWithOptions builds the contract over the given base directory, tuned by
// options.
//
// Files are created 0600 and directories 0700: a database holds whatever a
// program stores, and no other user of the machine should read it.
func NewWithOptions(base string, options Options) storagedeps.Contract {
	s := &store{base: base, sync: !options.NoSync}
	return storagedeps.Contract{
		Write: func(key []string, value []byte) error {
			return s.replace(key, s.path(key), value)
		},
		WriteIfAbsent: func(key []string, value []byte) (bool, error) {
			return s.create(s.path(key), value, s.sync)
		},
		WriteIfValueEquals: func(key []string, value []byte, oldValue []byte) (bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			path := s.path(key)
			current, found, err := read(path)
			if err != nil || !found {
				return false, err
			}
			if !bytes.Equal(current, oldValue) {
				return false, nil
			}
			if err := s.replace(key, path, value); err != nil {
				return false, err
			}
			return true, nil
		},
		Append: func(key []string, value []byte) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			path := s.path(key)
			current, _, err := read(path)
			if err != nil {
				return err
			}
			next := make([]byte, 0, len(current)+len(value))
			next = append(next, current...)
			next = append(next, value...)
			return s.replace(key, path, next)
		},
		InsertAt: func(key []string, position int64, value []byte) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			path := s.path(key)
			current, _, err := read(path)
			if err != nil {
				return err
			}
			if position < 0 || position > int64(len(current)) {
				return errors.New("keep: position " + strconv.FormatInt(position, 10) +
					" out of range for key " + describe(key))
			}
			next := make([]byte, 0, len(current)+len(value))
			next = append(next, current[:position]...)
			next = append(next, value...)
			next = append(next, current[position:]...)
			return s.replace(key, path, next)
		},
		Exists: func(key []string) (bool, error) {
			info, err := os.Stat(s.path(key))
			if errors.Is(err, os.ErrNotExist) {
				return false, nil
			}
			if err != nil {
				return false, err
			}
			return !info.IsDir(), nil
		},
		Read: func(key []string) ([]byte, bool, error) {
			return read(s.path(key))
		},
		ReadAt: func(key []string, position int64, size int64) ([]byte, bool, error) {
			file, err := os.Open(s.path(key))
			if errors.Is(err, os.ErrNotExist) {
				return nil, false, nil
			}
			if err != nil {
				return nil, false, err
			}
			defer file.Close()
			info, err := file.Stat()
			if err != nil {
				return nil, false, err
			}
			if info.IsDir() {
				return nil, false, nil
			}
			if position < 0 || position > info.Size() {
				return nil, true, errors.New("keep: position " + strconv.FormatInt(position, 10) +
					" out of range for key " + describe(key))
			}
			if size < 0 {
				return nil, true, errors.New("keep: negative size " + strconv.FormatInt(size, 10) +
					" for key " + describe(key))
			}
			// Comparing against what remains, rather than adding size to
			// position, is what keeps a huge size from overflowing.
			if remaining := info.Size() - position; size > remaining {
				size = remaining
			}
			buffer := make([]byte, size)
			if len(buffer) == 0 {
				return buffer, true, nil
			}
			if _, err := file.ReadAt(buffer, position); err != nil {
				return nil, true, err
			}
			return buffer, true, nil
		},
		Delete: func(key []string) error {
			path := s.path(key)
			err := os.Remove(path)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				if isDir(path) {
					// A directory is the prefix of other keys: the key
					// itself holds nothing, and nothing is the outcome.
					return nil
				}
				return err
			}
			s.prune(key)
			return nil
		},
		Lock: func(key []string, seconds int) (bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			return s.lock(key, seconds)
		},
		Unlock: func(key []string) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			err := os.Remove(s.lockPath(key))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			s.prune(key)
			return nil
		},
	}
}

// lock takes the lease of key, under the store's mutex. The lease file holds
// its expiry in Unix nanoseconds and is created whole, so a reader never sees
// it empty. An expired lease is taken over by one process at a time: the
// takeover holds a guard file, and replaces the lease only if it still holds
// the expired content it read.
//
// The lease names no holder — the contract's Unlock takes none — so Unlock
// releases it whoever took it.
func (s *store) lock(key []string, seconds int) (bool, error) {
	path := s.lockPath(key)
	content := []byte(strconv.FormatInt(time.Now().Add(time.Duration(seconds)*time.Second).UnixNano(), 10))
	created, err := s.create(path, content, false)
	if err != nil || created {
		return created, err
	}

	held, found, err := read(path)
	if err != nil {
		return false, err
	}
	if !found {
		// Released between the two calls: try once more.
		return s.create(path, content, false)
	}
	if !expired(path, held, seconds) {
		return false, nil
	}

	guard := path + ".takeover"
	file, err := os.OpenFile(guard, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, os.ErrExist) {
		// Another process is taking the lease over. A guard old enough was
		// left by one that crashed doing it, and is cleared for next time.
		if info, statErr := os.Stat(guard); statErr == nil && time.Since(info.ModTime()) > takeoverStale {
			os.Remove(guard)
		}
		return false, nil
	}
	if err != nil {
		return false, err
	}
	file.Close()
	defer os.Remove(guard)

	current, found, err := read(path)
	if err != nil {
		return false, err
	}
	if !found {
		return s.create(path, content, false)
	}
	if !bytes.Equal(current, held) {
		// Someone took it over between our read and our guard.
		return false, nil
	}
	temp, err := s.writeTemp(filepath.Dir(path), content, false)
	if err != nil {
		return false, err
	}
	if err := os.Rename(temp, path); err != nil {
		os.Remove(temp)
		return false, err
	}
	return true, nil
}

// expired reports whether the lease file at path, holding held, has run
// out. Content that is not an expiry — which a lease created whole never
// holds — counts as held for a lease's worth past the file's last change, so
// a corrupt file neither blocks a key forever nor frees it at once.
func expired(path string, held []byte, seconds int) bool {
	expiry, err := strconv.ParseInt(string(held), 10, 64)
	if err == nil {
		return time.Now().UnixNano() >= expiry
	}
	info, statErr := os.Stat(path)
	if statErr != nil {
		return false
	}
	return time.Since(info.ModTime()) > time.Duration(seconds)*time.Second
}
