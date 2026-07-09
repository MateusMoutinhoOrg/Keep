package database

type ErrorType int

const (
	KeyConflict ErrorType = iota
	NotFound
	MissingField
	InvalidField
	Internal
)

type Error struct {
	Type     ErrorType
	Key      string
	KeyValue any
	Msg      string
}

func (e Error) Error() string { return e.Msg }
