//go:build ignore

// This file is an illustrative sample, not part of the build.
package memory

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// MemoryAdapter fills deps.Deps against an in-memory map — the
// adapter's own, opinionated choice of backing store.
type MemoryAdapter struct {
	Deps  deps.Deps // the contract this adapter fills
	store map[string][]byte
}

// WriteFactory fills deps.Deps.Write.
func WriteFactory(m *MemoryAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		m.store[key] = value
		return nil
	}
}

// WriteIfKeyNotExistsFactory fills deps.Deps.WriteIfKeyNotExists.
func WriteIfKeyNotExistsFactory(m *MemoryAdapter) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		if _, found := m.store[key]; found {
			return deps.ErrKeyAlreadyExists // expected condition
		}
		m.store[key] = value
		return nil
	}
}

// ReadFactory fills deps.Deps.Read.
func ReadFactory(m *MemoryAdapter) func(key string) ([]byte, error) {
	return func(key string) ([]byte, error) {
		value, found := m.store[key]
		if !found {
			return nil, deps.ErrKeyNotFound // expected condition
		}
		return value, nil
	}
}

// ExistsFactory fills deps.Deps.Exists.
func ExistsFactory(m *MemoryAdapter) func(key string) (bool, error) {
	return func(key string) (bool, error) {
		_, found := m.store[key]
		return found, nil
	}
}

// DeleteFactory fills deps.Deps.Delete.
func DeleteFactory(m *MemoryAdapter) func(key string) error {
	return func(key string) error {
		delete(m.store, key)
		return nil
	}
}

// ... one factory per remaining requirement of the contract

// New builds a deps.Deps backed by an in-memory map. It builds the
// adapter instance and runs every field factory over it, so each
// closure reads the adapter's state at call time.
func New() deps.Deps {
	m := &MemoryAdapter{store: map[string][]byte{}}
	m.Deps.Write = WriteFactory(m)
	m.Deps.WriteIfKeyNotExists = WriteIfKeyNotExistsFactory(m)
	m.Deps.Read = ReadFactory(m)
	m.Deps.Exists = ExistsFactory(m)
	m.Deps.Delete = DeleteFactory(m)
	// ... one assignment per remaining field of deps.Deps
	return m.Deps
}
