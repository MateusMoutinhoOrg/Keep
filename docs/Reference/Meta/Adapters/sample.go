//go:build ignore

// This file is an illustrative sample, not part of the build.
package memory

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

// New builds a deps.Deps by implementing every field of the Deps
// contract as a closure. The backing store is the adapter's own,
// opinionated choice.
func New() deps.Deps {
	store := map[string][]byte{}

	return deps.Deps{
		Write: func(key string, value []byte) error {
			store[key] = value
			return nil
		},
		Read: func(key string) ([]byte, error) {
			value, found := store[key]
			if !found {
				return nil, deps.ErrKeyNotFound // expected condition
			}
			return value, nil
		},
		Exists: func(key string) (bool, error) {
			_, found := store[key]
			return found, nil
		},
		Delete: func(key string) error {
			delete(store, key)
			return nil
		},
		// ... every remaining field of the contract
	}
}
