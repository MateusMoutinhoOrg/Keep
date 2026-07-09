package native

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

// native is a pure in-memory backend. Data lives only for the lifetime
// of the process, which makes it ideal for tests and samples.
type native struct {
	mu    sync.Mutex
	data  map[string][]byte
	locks map[string]time.Time
}

func New() *native {
	return &native{
		data:  make(map[string][]byte),
		locks: make(map[string]time.Time),
	}
}

func (n *native) Write(key string, value []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.data[key] = append([]byte(nil), value...)
	return nil
}

func (n *native) WriteIfKeyNotExists(key string, value []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.data[key]; ok {
		return fmt.Errorf("%w: %s", deps.ErrKeyAlreadyExists, key)
	}
	n.data[key] = append([]byte(nil), value...)
	return nil
}

func (n *native) WriteIfValueEquals(key string, value []byte, oldValue []byte) error {
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

func (n *native) Append(key string, value []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.data[key] = append(n.data[key], value...)
	return nil
}

func (n *native) InsertAt(key string, position int64, value []byte) error {
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

func (n *native) Exists(key string) (bool, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, ok := n.data[key]
	return ok, nil
}

func (n *native) Read(key string) ([]byte, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	value, ok := n.data[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
	}
	return append([]byte(nil), value...), nil
}

func (n *native) ReadAt(key string, position int64, size int64) ([]byte, error) {
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

func (n *native) Delete(key string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.data, key)
	return nil
}

func (n *native) Lock(key string, ttl int) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if expiry, ok := n.locks[key]; ok && time.Now().Before(expiry) {
		return fmt.Errorf("%w: %s", deps.ErrKeyLocked, key)
	}
	n.locks[key] = time.Now().Add(time.Duration(ttl) * time.Second)
	return nil
}

func (n *native) UnLock(key string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.locks, key)
	return nil
}
