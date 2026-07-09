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

	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

// standard is the filesystem-backed implementation: each key becomes a
// file (path segments split on "/", each segment escaped so keys can
// contain arbitrary characters). Data survives across process restarts.
type standard struct {
	mu   sync.Mutex
	base string
}

// New stores files relative to the current working directory.
func New() *standard {
	return NewWithBase(".")
}

// NewWithBase stores all keys under the given directory.
func NewWithBase(base string) *standard {
	return &standard{base: base}
}

func (s *standard) path(key string) string {
	parts := []string{s.base}
	for _, seg := range strings.Split(key, "/") {
		if seg == "" {
			continue
		}
		parts = append(parts, url.PathEscape(seg))
	}
	return filepath.Join(parts...)
}

func (s *standard) Write(key string, value []byte) error {
	p := s.path(key)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, value, 0o644)
}

func (s *standard) WriteIfKeyNotExists(key string, value []byte) error {
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

func (s *standard) WriteIfValueEquals(key string, value []byte, oldValue []byte) error {
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

func (s *standard) Append(key string, value []byte) error {
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

func (s *standard) InsertAt(key string, position int64, value []byte) error {
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

func (s *standard) Exists(key string) (bool, error) {
	_, err := os.Stat(s.path(key))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *standard) Read(key string) ([]byte, error) {
	value, err := os.ReadFile(s.path(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
	}
	return value, err
}

func (s *standard) ReadAt(key string, position int64, size int64) ([]byte, error) {
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

func (s *standard) Delete(key string) error {
	err := os.Remove(s.path(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *standard) lockPath(key string) string {
	return s.path(key) + ".keeplock"
}

func (s *standard) Lock(key string, ttl int) error {
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

func (s *standard) UnLock(key string) error {
	err := os.Remove(s.lockPath(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
