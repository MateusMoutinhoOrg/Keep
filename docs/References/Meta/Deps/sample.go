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

// Deps is the dependency contract every adapter must satisfy. Each
// method is one injectable behavior the library needs from a storage
// backend, and every one of them addresses a single key.
type Deps interface {
	// Write stores value under key, overwriting any current value.
	Write(key string, value []byte) error
	// WriteIfKeyNotExists stores value only when key is absent,
	// returning ErrKeyAlreadyExists otherwise.
	WriteIfKeyNotExists(key string, value []byte) error
	// Read returns the value of key, or ErrKeyNotFound when absent.
	Read(key string) ([]byte, error)
	// Exists reports whether key currently holds a value.
	Exists(key string) (bool, error)
	// Delete removes key; removing an absent key is not an error.
	Delete(key string) error
	// ... one method per remaining call of the contract
}
