package native

type NativeDeps struct{}

func (n *NativeDeps) WriteString(key, value string) error {
	return nil
}
func (n *NativeDeps) WriteBytes(key string, value []byte) error {
	return nil
}
func (n *NativeDeps) AppendString(key, value string) error {
	return nil
}
func (n *NativeDeps) AppendBytes(key string, value []byte) error {
	return nil
}
func (n *NativeDeps) Inc(key string) error {
	return nil
}
func (n *NativeDeps) Dec(key string) error {
	return nil
}

func (n *NativeDeps) WriteIfNotExists(key string, value string) error {
	return nil
}

func (n *NativeDeps) ReadString(key string, pos uint64, length uint64) (string, error) {
	return "", nil
}
func (n *NativeDeps) ReadBytes(key string, pos uint64, length uint64) ([]byte, error) {
	return nil, nil
}

func (n *NativeDeps) Lock(key string) error {
	return nil
}
func (n *NativeDeps) Unlock(key string) error {
	return nil
}

func New() *NativeDeps {
	return &NativeDeps{}
}
