package deps

type Deps interface {
	WriteString(key, value string) error
	WriteBytes(key string, value []byte) error
	AppendString(key, value string) error
	AppendBytes(key string, value []byte) error
	//locks the key,them increase the value
	Inc(key string) error
	//locks the key,them decrease the value
	Dec(key string) error
	//locks the key,them writes the value if the key is equal to the expected value
	WriteIfNotExists(key string, value string, expectedValue string) error

	ReadString(key string, pos uint64, length uint64) (string, error)
	ReadBytes(key string, pos uint64, length uint64) ([]byte, error)
	Lock(key string) error
	Unlock(key string) error
}
