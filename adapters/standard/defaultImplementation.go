package standard

type standard struct{}

func (n *standard) Write(key string, value []byte) error {
	return nil
}

func (n *standard) WriteIfKeyNotExists(key string, value []byte) error {
	return nil
}

func (n *standard) WriteIfValueEquals(key string, value []byte, oldValue []byte) error {
	return nil
}

func (n *standard) Append(key string, value []byte) error {
	return nil
}

func (n *standard) InsertAt(key string, position int64, value []byte) error {
	return nil
}

func (n *standard) Exists(key string) (bool, error) {
	return false, nil
}

func (n *standard) Read(key string) ([]byte, error) {
	return nil, nil
}

func (n *standard) ReadAt(key string, position int64, size int64) ([]byte, error) {
	return nil, nil
}

func (n *standard) Delete(key string) error {
	return nil
}

func (n *standard) Lock(key string, time int) error {
	return nil
}

func (n *standard) UnLock(key string) error {
	return nil
}

func New() *standard {
	return &standard{}
}
