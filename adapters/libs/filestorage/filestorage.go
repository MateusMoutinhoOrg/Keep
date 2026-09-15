package filestorage

import (
	"bytes"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
	storagedeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/storagedeps"
)

// store is the state one bound adapter keeps: the directory every key is
// resolved under, and the mutex the read-modify-write operations take.
type store struct {
	base string
	mu   sync.Mutex
}

// path resolves a key to a file path. Each slash-separated segment is
// escaped on its own, so a key may hold any character without ever
// escaping the base directory or colliding with another key.
func (s *store) path(key string) string {
	parts := []string{s.base}
	for _, segment := range strings.Split(key, "/") {
		if segment == "" {
			continue
		}
		parts = append(parts, url.PathEscape(segment))
	}
	return filepath.Join(parts...)
}

// lockPath is the file holding the lease of a key: the key's own path with
// a suffix no escaped segment can produce.
func (s *store) lockPath(key string) string {
	return s.path(key) + ".keeplock"
}

// Bind fills deps.Deps.Storagedeps with the filesystem implementation: one
// file per key, rooted at the working directory, so what a database writes
// survives the process and can be read with `ls`.
func Bind(deps *deps.Deps) {
	deps.Storagedeps = New(".")
}

// New builds the contract over the given base directory. Bind uses ".", so
// a Props.Path carries the whole of where a database lands; a program that
// wants the tree somewhere else assigns the result of this call itself.
func New(base string) storagedeps.Sandbox {
	s := &store{base: base}
	return storagedeps.Sandbox{
		Write: func(key string, value []byte) error {
			path := s.path(key)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			return os.WriteFile(path, value, 0o644)
		},
		WriteIfKeyNotExists: func(key string, value []byte) (bool, error) {
			path := s.path(key)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return false, err
			}
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
			if errors.Is(err, os.ErrExist) {
				return false, nil
			}
			if err != nil {
				return false, err
			}
			defer file.Close()
			if _, err := file.Write(value); err != nil {
				return false, err
			}
			return true, nil
		},
		WriteIfValueEquals: func(key string, value []byte, oldValue []byte) (bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			current, err := os.ReadFile(s.path(key))
			if errors.Is(err, os.ErrNotExist) {
				return false, nil
			}
			if err != nil {
				return false, err
			}
			if !bytes.Equal(current, oldValue) {
				return false, nil
			}
			if err := os.WriteFile(s.path(key), value, 0o644); err != nil {
				return false, err
			}
			return true, nil
		},
		Append: func(key string, value []byte) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			path := s.path(key)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = file.Write(value)
			return err
		},
		InsertAt: func(key string, position int64, value []byte) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			path := s.path(key)
			current, err := os.ReadFile(path)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			if position < 0 || position > int64(len(current)) {
				return errors.New("keep: position " + strconv.FormatInt(position, 10) +
					" out of range for key " + key)
			}
			next := make([]byte, 0, len(current)+len(value))
			next = append(next, current[:position]...)
			next = append(next, value...)
			next = append(next, current[position:]...)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			return os.WriteFile(path, next, 0o644)
		},
		Exists: func(key string) (bool, error) {
			_, err := os.Stat(s.path(key))
			if errors.Is(err, os.ErrNotExist) {
				return false, nil
			}
			if err != nil {
				return false, err
			}
			return true, nil
		},
		Read: func(key string) ([]byte, bool, error) {
			value, err := os.ReadFile(s.path(key))
			if errors.Is(err, os.ErrNotExist) {
				return nil, false, nil
			}
			if err != nil {
				return nil, false, err
			}
			return value, true, nil
		},
		ReadAt: func(key string, position int64, size int64) ([]byte, bool, error) {
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
			if position < 0 || position > info.Size() {
				return nil, true, errors.New("keep: position " + strconv.FormatInt(position, 10) +
					" out of range for key " + key)
			}
			end := position + size
			if end > info.Size() {
				end = info.Size()
			}
			buffer := make([]byte, end-position)
			if len(buffer) == 0 {
				return buffer, true, nil
			}
			if _, err := file.ReadAt(buffer, position); err != nil {
				return nil, true, err
			}
			return buffer, true, nil
		},
		Delete: func(key string) error {
			err := os.Remove(s.path(key))
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		},
		Lock: func(key string, seconds int) (bool, error) {
			path := s.lockPath(key)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return false, err
			}
			expiry := time.Now().Add(time.Duration(seconds) * time.Second).UnixNano()
			content := []byte(strconv.FormatInt(expiry, 10))
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
			if errors.Is(err, os.ErrExist) {
				raw, readErr := os.ReadFile(path)
				if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
					return false, readErr
				}
				held, _ := strconv.ParseInt(string(raw), 10, 64)
				if time.Now().UnixNano() < held {
					return false, nil
				}
				// The lease expired: take it over.
				if err := os.WriteFile(path, content, 0o644); err != nil {
					return false, err
				}
				return true, nil
			}
			if err != nil {
				return false, err
			}
			defer file.Close()
			if _, err := file.Write(content); err != nil {
				return false, err
			}
			return true, nil
		},
		UnLock: func(key string) error {
			err := os.Remove(s.lockPath(key))
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		},
	}
}
