package liberror

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Constructors for *api.Error. Every failure the library reports is built
// here, so the shape of an error is decided in one place and a caller can
// rely on Type being one of the api constants and never on the wording of
// Message.

// New builds a failure of the given type, naming the field it involves.
// Callers return the result directly as *api.Error: a typed nil hidden in
// an interface would read as a failure to a caller comparing against nil.
func New(sandbox *api.Sandbox, failureType int, key string, message string) *api.Error {
	return &api.Error{Type: failureType, Key: key, Message: message}
}

// NewWithValue builds a failure that also carries the offending value, so a
// caller handling a KeyConflict or an InvalidField can report which value
// was refused without re-reading it.
func NewWithValue(sandbox *api.Sandbox, failureType int, key string, value any, message string) *api.Error {
	return &api.Error{Type: failureType, Key: key, KeyValue: value, Message: message}
}
