package database

type ErrorType int

const (
	KeyConflict ErrorType = iota
)

type Error struct {
	Type     ErrorType
	Key      string
	KeyValue any
	Msg      string
}

func (e Error) Error() string { return e.Msg }
