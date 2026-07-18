package deps

import "errors"

// Sentinel errors that every Deps implementation must return (possibly
// wrapped) so the database layer can distinguish expected conditions
// from real failures. Compare with errors.Is.
var (
	ErrKeyNotFound      = errors.New("keep: key not found")
	ErrKeyAlreadyExists = errors.New("keep: key already exists")
	ErrValueMismatch    = errors.New("keep: value mismatch")
	ErrKeyLocked        = errors.New("keep: key is locked")
)

// Deps holds the injectable functions the library needs from a storage
// backend. Adapters return a fully populated Deps; consumers can also
// build one by hand or overwrite individual fields after construction.
type Deps struct {
	Write               func(key string, value []byte) error
	WriteIfKeyNotExists func(key string, value []byte) error
	WriteIfValueEquals  func(key string, value []byte, oldValue []byte) error
	Append              func(key string, value []byte) error
	InsertAt            func(key string, position int64, value []byte) error
	Exists              func(key string) (bool, error)
	Read                func(key string) ([]byte, error)
	ReadAt              func(key string, position int64, size int64) ([]byte, error)
	Delete              func(key string) error
	Lock                func(key string, time int) error
	UnLock              func(key string) error
}
