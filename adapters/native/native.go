package native

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// NativeAdapter is a pure in-memory backend. Data lives only for the
// lifetime of the process, which makes it ideal for tests and samples.
type NativeAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps  deps.Deps
	mu    sync.Mutex
	data  map[string][]byte
	locks map[string]time.Time
}

// WriteFactory fills deps.Deps.Write.
func WriteFactory(n *NativeAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		n.data[key] = append([]byte(nil), value...)
		return nil
	}
}

// WriteIfKeyNotExistsFactory fills deps.Deps.WriteIfKeyNotExists.
func WriteIfKeyNotExistsFactory(n *NativeAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		if _, ok := n.data[key]; ok {
			return fmt.Errorf("%w: %s", deps.ErrKeyAlreadyExists, key)
		}
		n.data[key] = append([]byte(nil), value...)
		return nil
	}
}

// WriteIfValueEqualsFactory fills deps.Deps.WriteIfValueEquals.
func WriteIfValueEqualsFactory(n *NativeAdapter) func(key string, value []byte, oldValue []byte) error {
	return func(key string, value []byte, oldValue []byte) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		current, ok := n.data[key]
		if !ok {
			return fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		if !bytes.Equal(current, oldValue) {
			return fmt.Errorf("%w: %s", deps.ErrValueMismatch, key)
		}
		n.data[key] = append([]byte(nil), value...)
		return nil
	}
}

// AppendFactory fills deps.Deps.Append.
func AppendFactory(n *NativeAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		n.data[key] = append(n.data[key], value...)
		return nil
	}
}

// InsertAtFactory fills deps.Deps.InsertAt.
func InsertAtFactory(n *NativeAdapter) func(key string, position int64, value []byte) error {
	return func(key string, position int64, value []byte) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		current := n.data[key]
		if position < 0 || position > int64(len(current)) {
			return fmt.Errorf("keep: position %d out of range for key %s (len %d)", position, key, len(current))
		}
		next := make([]byte, 0, len(current)+len(value))
		next = append(next, current[:position]...)
		next = append(next, value...)
		next = append(next, current[position:]...)
		n.data[key] = next
		return nil
	}
}

// ExistsFactory fills deps.Deps.Exists.
func ExistsFactory(n *NativeAdapter) func(key string) (bool, error) {
	return func(key string) (bool, error) {
		n.mu.Lock()
		defer n.mu.Unlock()
		_, ok := n.data[key]
		return ok, nil
	}
}

// ReadFactory fills deps.Deps.Read.
func ReadFactory(n *NativeAdapter) func(key string) ([]byte, error) {
	return func(key string) ([]byte, error) {
		n.mu.Lock()
		defer n.mu.Unlock()
		value, ok := n.data[key]
		if !ok {
			return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		return append([]byte(nil), value...), nil
	}
}

// ReadAtFactory fills deps.Deps.ReadAt.
func ReadAtFactory(n *NativeAdapter) func(key string, position int64, size int64) ([]byte, error) {
	return func(key string, position int64, size int64) ([]byte, error) {
		n.mu.Lock()
		defer n.mu.Unlock()
		value, ok := n.data[key]
		if !ok {
			return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		if position < 0 || position > int64(len(value)) {
			return nil, fmt.Errorf("keep: position %d out of range for key %s (len %d)", position, key, len(value))
		}
		end := position + size
		if end > int64(len(value)) {
			end = int64(len(value))
		}
		return append([]byte(nil), value[position:end]...), nil
	}
}

// DeleteFactory fills deps.Deps.Delete.
func DeleteFactory(n *NativeAdapter) func(key string) error {
	return func(key string) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		delete(n.data, key)
		return nil
	}
}

// LockFactory fills deps.Deps.Lock.
func LockFactory(n *NativeAdapter) func(key string, ttl int) error {
	return func(key string, ttl int) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		if expiry, ok := n.locks[key]; ok && time.Now().Before(expiry) {
			return fmt.Errorf("%w: %s", deps.ErrKeyLocked, key)
		}
		n.locks[key] = time.Now().Add(time.Duration(ttl) * time.Second)
		return nil
	}
}

// UnLockFactory fills deps.Deps.UnLock.
func UnLockFactory(n *NativeAdapter) func(key string) error {
	return func(key string) error {
		n.mu.Lock()
		defer n.mu.Unlock()
		delete(n.locks, key)
		return nil
	}
}

// New creates a deps.Deps backed by an in-memory adapter. It builds the
// adapter instance and runs every field factory over it, so each closure
// reads the adapter's state at call time. Adding a field to deps.Deps
// means adding its factory call here.
func New() deps.Deps {
	n := &NativeAdapter{
		data:  make(map[string][]byte),
		locks: make(map[string]time.Time),
	}
	n.Deps.Write = WriteFactory(n)
	n.Deps.WriteIfKeyNotExists = WriteIfKeyNotExistsFactory(n)
	n.Deps.WriteIfValueEquals = WriteIfValueEqualsFactory(n)
	n.Deps.Append = AppendFactory(n)
	n.Deps.InsertAt = InsertAtFactory(n)
	n.Deps.Exists = ExistsFactory(n)
	n.Deps.Read = ReadFactory(n)
	n.Deps.ReadAt = ReadAtFactory(n)
	n.Deps.Delete = DeleteFactory(n)
	n.Deps.Lock = LockFactory(n)
	n.Deps.UnLock = UnLockFactory(n)
	return n.Deps
}
