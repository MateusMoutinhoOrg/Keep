package liberror

// The struct backing every failure the library reports. It implements
// api.Error, so the contracts package can stay interfaces-only. The
// package is named liberror rather than error to avoid shadowing the
// predeclared type.

import "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"

// Error implements api.Error.
type Error struct {
	FailureType  int
	FailureKey   string
	FailureValue any
	Msg          string
}

func (e *Error) Error() string { return e.Msg }
func (e *Error) Type() int     { return e.FailureType }
func (e *Error) Key() string   { return e.FailureKey }
func (e *Error) KeyValue() any { return e.FailureValue }

// New builds a failure of the given type. Callers must return the
// result directly as an api.Error, never through a typed nil.
func New(failureType int, key string, msg string) api.Error {
	return &Error{FailureType: failureType, FailureKey: key, Msg: msg}
}

// NewWithValue builds a failure that also carries the offending value.
func NewWithValue(failureType int, key string, value any, msg string) api.Error {
	return &Error{FailureType: failureType, FailureKey: key, FailureValue: value, Msg: msg}
}
