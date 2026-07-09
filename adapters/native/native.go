package native

type native struct{}

func (n *native) Write(key string, value []byte) error {
	return nil
}

func (n *native) WriteIfKeyNotExists(key string, value []byte) error {
	return nil
}

func (n *native) WriteIfValueEquals(key string, value []byte, oldValue []byte) error {
	return nil
}

func (n *native) Append(key string, value []byte) error {
	return nil
}

func (n *native) InsertAt(key string, position int64, value []byte) error {
	return nil
}

func (n *native) Exists(key string) (bool, error) {
	return false, nil
}

func (n *native) Read(key string) ([]byte, error) {
	return nil, nil
}

func (n *native) ReadAt(key string, position int64, size int64) ([]byte, error) {
	return nil, nil
}

func (n *native) Delete(key string) error {
	return nil
}

func (n *native) Lock(key string, time int) error {
	return nil
}

func (n *native) UnLock(key string) error {
	return nil
}

func New() *native {
	return &native{}
}
