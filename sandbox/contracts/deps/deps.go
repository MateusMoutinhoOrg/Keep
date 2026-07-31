package deps

import "errors"

// Sentinel errors that every Deps implementation must return (possibly
// wrapped) so the library can distinguish expected conditions from real
// failures. Compare with errors.Is.
var (
	ErrKeyNotFound      = errors.New("keep: key not found")
	ErrKeyAlreadyExists = errors.New("keep: key already exists")
	ErrValueMismatch    = errors.New("keep: value mismatch")
	ErrKeyLocked        = errors.New("keep: key is locked")
)

// Deps is the dependency contract every adapter must satisfy. It is a
// struct of function fields, not an interface: an adapter fills every
// field with the behavior it provides, and the library calls those
// fields directly. Each field addresses a single key: the library never
// lists, scans a prefix, or queries a range.
type Deps struct {
	// Write stores value under key, overwriting any current value.
	Write func(key string, value []byte) error
	// WriteIfKeyNotExists stores value only when key is absent,
	// returning ErrKeyAlreadyExists otherwise.
	WriteIfKeyNotExists func(key string, value []byte) error
	// WriteIfValueEquals stores value only when the current value is
	// oldValue, returning ErrValueMismatch otherwise.
	WriteIfValueEquals func(key string, value []byte, oldValue []byte) error
	// Append adds value to the end of the current value of key.
	Append func(key string, value []byte) error
	// InsertAt splices value into the current value of key at position.
	InsertAt func(key string, position int64, value []byte) error
	// Exists reports whether key currently holds a value.
	Exists func(key string) (bool, error)
	// Read returns the value of key, or ErrKeyNotFound when absent.
	Read func(key string) ([]byte, error)
	// ReadAt returns at most size bytes of key starting at position.
	ReadAt func(key string, position int64, size int64) ([]byte, error)
	// Delete removes key; removing an absent key is not an error.
	Delete func(key string) error
	// Lock takes an exclusive lease on key for time seconds, returning
	// ErrKeyLocked when someone else holds it.
	Lock func(key string, time int) error
	// UnLock releases a lease taken by Lock.
	UnLock func(key string) error
}
