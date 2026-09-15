package memstorage

import (
	"bytes"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
	storagedeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/storagedeps"
)

// store is the state one bound adapter keeps: the map holding every value,
// the map holding every live lease, and the mutex guarding both.
type store struct {
	mu     sync.Mutex
	values map[string][]byte
	leases map[string]time.Time
}

// Bind fills deps.Deps.Storagedeps with the in-memory implementation: a map
// that lives as long as the process, which is what a test or an example
// wants when nothing should be left on disk.
func Bind(deps *deps.Deps) {
	deps.Storagedeps = New()
}

// New builds the contract over a map of its own. Two calls never share
// state, so two databases built from two calls cannot see each other.
func New() storagedeps.Sandbox {
	s := &store{
		values: map[string][]byte{},
		leases: map[string]time.Time{},
	}
	return storagedeps.Sandbox{
		Write: func(key string, value []byte) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			s.values[key] = append([]byte(nil), value...)
			return nil
		},
		WriteIfKeyNotExists: func(key string, value []byte) (bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			if _, ok := s.values[key]; ok {
				return false, nil
			}
			s.values[key] = append([]byte(nil), value...)
			return true, nil
		},
		WriteIfValueEquals: func(key string, value []byte, oldValue []byte) (bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			current, ok := s.values[key]
			if !ok || !bytes.Equal(current, oldValue) {
				return false, nil
			}
			s.values[key] = append([]byte(nil), value...)
			return true, nil
		},
		Append: func(key string, value []byte) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			s.values[key] = append(s.values[key], value...)
			return nil
		},
		InsertAt: func(key string, position int64, value []byte) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			current := s.values[key]
			if position < 0 || position > int64(len(current)) {
				return errors.New("keep: position " + strconv.FormatInt(position, 10) +
					" out of range for key " + key)
			}
			next := make([]byte, 0, len(current)+len(value))
			next = append(next, current[:position]...)
			next = append(next, value...)
			next = append(next, current[position:]...)
			s.values[key] = next
			return nil
		},
		Exists: func(key string) (bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			_, ok := s.values[key]
			return ok, nil
		},
		Read: func(key string) ([]byte, bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			value, ok := s.values[key]
			if !ok {
				return nil, false, nil
			}
			return append([]byte(nil), value...), true, nil
		},
		ReadAt: func(key string, position int64, size int64) ([]byte, bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			value, ok := s.values[key]
			if !ok {
				return nil, false, nil
			}
			if position < 0 || position > int64(len(value)) {
				return nil, true, errors.New("keep: position " + strconv.FormatInt(position, 10) +
					" out of range for key " + key)
			}
			end := position + size
			if end > int64(len(value)) {
				end = int64(len(value))
			}
			return append([]byte(nil), value[position:end]...), true, nil
		},
		Delete: func(key string) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			delete(s.values, key)
			return nil
		},
		Lock: func(key string, seconds int) (bool, error) {
			s.mu.Lock()
			defer s.mu.Unlock()
			if expiry, ok := s.leases[key]; ok && time.Now().Before(expiry) {
				return false, nil
			}
			s.leases[key] = time.Now().Add(time.Duration(seconds) * time.Second)
			return true, nil
		},
		UnLock: func(key string) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			delete(s.leases, key)
			return nil
		},
	}
}
