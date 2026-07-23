//go:build ignore

// This file is an illustrative sample, not part of the build.
package deps

import "errors"

// Sentinel errors every implementation must return (possibly wrapped)
// so the library can tell an expected condition from a real failure.
var (
	ErrKeyNotFound      = errors.New("keep: key not found")
	ErrKeyAlreadyExists = errors.New("keep: key already exists")
)

// Deps holds the injectable functions the library needs from a storage
// backend — one function field per required behavior.
type Deps struct {
	Write               func(key string, value []byte) error
	WriteIfKeyNotExists func(key string, value []byte) error
	Read                func(key string) ([]byte, error)
	Exists              func(key string) (bool, error)
	Delete              func(key string) error
	// ... one field per remaining call of the contract
}
