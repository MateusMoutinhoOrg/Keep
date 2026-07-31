package standard

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// StandardAdapter is the filesystem-backed implementation: each key
// becomes a file (path segments split on "/", each segment escaped so
// keys can contain arbitrary characters). Data survives across process
// restarts.
type StandardAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps deps.Deps
	mu   sync.Mutex
	base string
}

func (s *StandardAdapter) path(key string) string {
	parts := []string{s.base}
	for _, seg := range strings.Split(key, "/") {
		if seg == "" {
			continue
		}
		parts = append(parts, url.PathEscape(seg))
	}
	return filepath.Join(parts...)
}

func (s *StandardAdapter) lockPath(key string) string {
	return s.path(key) + ".keeplock"
}

// WriteFactory fills deps.Deps.Write.
func WriteFactory(s *StandardAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		p := s.path(key)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		return os.WriteFile(p, value, 0o644)
	}
}

// WriteIfKeyNotExistsFactory fills deps.Deps.WriteIfKeyNotExists.
func WriteIfKeyNotExistsFactory(s *StandardAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		p := s.path(key)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w: %s", deps.ErrKeyAlreadyExists, key)
		}
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(value)
		return err
	}
}

// WriteIfValueEqualsFactory fills deps.Deps.WriteIfValueEquals.
func WriteIfValueEqualsFactory(s *StandardAdapter) func(key string, value []byte, oldValue []byte) error {
	return func(key string, value []byte, oldValue []byte) error {
		s.mu.Lock()
		defer s.mu.Unlock()
		current, err := os.ReadFile(s.path(key))
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		if err != nil {
			return err
		}
		if !bytes.Equal(current, oldValue) {
			return fmt.Errorf("%w: %s", deps.ErrValueMismatch, key)
		}
		return os.WriteFile(s.path(key), value, 0o644)
	}
}

// AppendFactory fills deps.Deps.Append.
func AppendFactory(s *StandardAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		s.mu.Lock()
		defer s.mu.Unlock()
		p := s.path(key)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(value)
		return err
	}
}

// InsertAtFactory fills deps.Deps.InsertAt.
func InsertAtFactory(s *StandardAdapter) func(key string, position int64, value []byte) error {
	return func(key string, position int64, value []byte) error {
		s.mu.Lock()
		defer s.mu.Unlock()
		p := s.path(key)
		current, err := os.ReadFile(p)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if position < 0 || position > int64(len(current)) {
			return fmt.Errorf("keep: position %d out of range for key %s (len %d)", position, key, len(current))
		}
		next := make([]byte, 0, len(current)+len(value))
		next = append(next, current[:position]...)
		next = append(next, value...)
		next = append(next, current[position:]...)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		return os.WriteFile(p, next, 0o644)
	}
}

// ExistsFactory fills deps.Deps.Exists.
func ExistsFactory(s *StandardAdapter) func(key string) (bool, error) {
	return func(key string) (bool, error) {
		_, err := os.Stat(s.path(key))
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

// ReadFactory fills deps.Deps.Read.
func ReadFactory(s *StandardAdapter) func(key string) ([]byte, error) {
	return func(key string) ([]byte, error) {
		value, err := os.ReadFile(s.path(key))
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		return value, err
	}
}

// ReadAtFactory fills deps.Deps.ReadAt.
func ReadAtFactory(s *StandardAdapter) func(key string, position int64, size int64) ([]byte, error) {
	return func(key string, position int64, size int64) ([]byte, error) {
		f, err := os.Open(s.path(key))
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		if err != nil {
			return nil, err
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if position < 0 || position > info.Size() {
			return nil, fmt.Errorf("keep: position %d out of range for key %s (len %d)", position, key, info.Size())
		}
		end := position + size
		if end > info.Size() {
			end = info.Size()
		}
		buf := make([]byte, end-position)
		if _, err := f.ReadAt(buf, position); err != nil {
			return nil, err
		}
		return buf, nil
	}
}

// DeleteFactory fills deps.Deps.Delete.
func DeleteFactory(s *StandardAdapter) func(key string) error {
	return func(key string) error {
		err := os.Remove(s.path(key))
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
}

// LockFactory fills deps.Deps.Lock.
func LockFactory(s *StandardAdapter) func(key string, ttl int) error {
	return func(key string, ttl int) error {
		p := s.lockPath(key)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		expiry := time.Now().Add(time.Duration(ttl) * time.Second).UnixNano()
		content := []byte(strconv.FormatInt(expiry, 10))
		f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if errors.Is(err, os.ErrExist) {
			raw, readErr := os.ReadFile(p)
			if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
				return readErr
			}
			held, _ := strconv.ParseInt(string(raw), 10, 64)
			if time.Now().UnixNano() < held {
				return fmt.Errorf("%w: %s", deps.ErrKeyLocked, key)
			}
			// The previous lock expired: take it over.
			return os.WriteFile(p, content, 0o644)
		}
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(content)
		return err
	}
}

// UnLockFactory fills deps.Deps.UnLock.
func UnLockFactory(s *StandardAdapter) func(key string) error {
	return func(key string) error {
		err := os.Remove(s.lockPath(key))
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
}

// New creates a deps.Deps that stores files relative to the current
// working directory.
func New() deps.Deps {
	return NewWithBase(".")
}

// NewWithBase creates a deps.Deps that stores all keys under the given
// directory. It builds the adapter instance and runs every field
// factory over it, so each closure reads the adapter's state at call
// time. Adding a field to deps.Deps means adding its factory call here.
func NewWithBase(base string) deps.Deps {
	s := &StandardAdapter{base: base}
	s.Deps.Write = WriteFactory(s)
	s.Deps.WriteIfKeyNotExists = WriteIfKeyNotExistsFactory(s)
	s.Deps.WriteIfValueEquals = WriteIfValueEqualsFactory(s)
	s.Deps.Append = AppendFactory(s)
	s.Deps.InsertAt = InsertAtFactory(s)
	s.Deps.Exists = ExistsFactory(s)
	s.Deps.Read = ReadFactory(s)
	s.Deps.ReadAt = ReadAtFactory(s)
	s.Deps.Delete = DeleteFactory(s)
	s.Deps.Lock = LockFactory(s)
	s.Deps.UnLock = UnLockFactory(s)
	return s.Deps
}
