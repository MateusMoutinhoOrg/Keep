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

type Deps interface {
	Write(key string, value []byte) error
	WriteIfKeyNotExists(key string, value []byte) error
	WriteIfValueEquals(key string, value []byte, oldValue []byte) error
	Append(key string, value []byte) error
	InsertAt(key string, position int64, value []byte) error
	Exists(key string) (bool, error)
	Read(key string) ([]byte, error)
	ReadAt(key string, position int64, size int64) ([]byte, error)
	Delete(key string) error
	Lock(key string, time int) error
	UnLock(key string) error
}
