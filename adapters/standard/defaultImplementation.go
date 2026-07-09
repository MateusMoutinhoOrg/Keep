package standard

type standard struct{}

func (n *standard) WriteString(key, value string) error {
	return nil
}
func (n *standard) WriteBytes(key string, value []byte) error {
	return nil
}
func (n *standard) AppendString(key, value string) error {
	return nil
}
func (n *standard) AppendBytes(key string, value []byte) error {
	return nil
}
func (n *standard) Inc(key string) error {
	return nil
}
func (n *standard) Dec(key string) error {
	return nil
}

func (n *standard) WriteIfNotExists(key string, value string) error {
	return nil
}

func (n *standard) ReadString(key string, pos uint64, length uint64) (string, error) {
	return "", nil
}
func (n *standard) ReadBytes(key string, pos uint64, length uint64) ([]byte, error) {
	return nil, nil
}

func (n *standard) Lock(key string) error {
	return nil
}
func (n *standard) Unlock(key string) error {
	return nil
}

func New() *standard {
	return &standard{}
}
