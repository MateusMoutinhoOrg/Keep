//go:build ignore

// This file is an illustrative sample, not part of the build.
package memory

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// memory holds the adapter's state. The backing store is the adapter's
// own, opinionated choice.
type memory struct {
	store map[string][]byte
}

// New builds a deps.Deps by returning a struct that implements every
// method of the Deps contract.
func New() deps.Deps {
	return &memory{store: map[string][]byte{}}
}

func (m *memory) Write(key string, value []byte) error {
	m.store[key] = value
	return nil
}

func (m *memory) WriteIfKeyNotExists(key string, value []byte) error {
	if _, found := m.store[key]; found {
		return deps.ErrKeyAlreadyExists // expected condition
	}
	m.store[key] = value
	return nil
}

func (m *memory) Read(key string) ([]byte, error) {
	value, found := m.store[key]
	if !found {
		return nil, deps.ErrKeyNotFound // expected condition
	}
	return value, nil
}

func (m *memory) Exists(key string) (bool, error) {
	_, found := m.store[key]
	return found, nil
}

func (m *memory) Delete(key string) error {
	delete(m.store, key)
	return nil
}

// ... one method per remaining requirement of the contract
