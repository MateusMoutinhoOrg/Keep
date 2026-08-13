//go:build ignore

// This file is an illustrative sample, not part of the build.
package deps

import "errors"

// Sentinel errors every adapter must return (possibly wrapped) so the
// library can tell an expected condition from a real failure.
var (
	ErrKeyNotFound      = errors.New("keep: key not found")
	ErrKeyAlreadyExists = errors.New("keep: key already exists")
)

// Deps is the dependency contract every adapter must fill. Each field
// is one injectable behavior the library needs from a storage backend,
// and every one of them addresses a single key. An adapter's New
// constructor assigns a factory-built closure to each field.
type Deps struct {
	// Write stores value under key, overwriting any current value.
	Write func(key string, value []byte) error
	// WriteIfKeyNotExists stores value only when key is absent,
	// returning ErrKeyAlreadyExists otherwise.
	WriteIfKeyNotExists func(key string, value []byte) error
	// Read returns the value of key, or ErrKeyNotFound when absent.
	Read func(key string) ([]byte, error)
	// Exists reports whether key currently holds a value.
	Exists func(key string) (bool, error)
	// Delete removes key; removing an absent key is not an error.
	Delete func(key string) error
	// ... one field per remaining requirement of the contract
}
