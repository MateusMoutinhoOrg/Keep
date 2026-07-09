package deps

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
