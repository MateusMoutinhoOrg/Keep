package liberror

// Constructors for api.Error. Error carries no behavior, so these are
// plain builders rather than factories — there is no Deps field to
// close over and nothing to fill after construction.

import "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"

// New builds a failure of the given type. Callers must return the
// result directly as *api.Error, never through a typed nil hidden in an
// interface.
func New(failureType int, key string, msg string) *api.Error {
	return &api.Error{Type: failureType, Key: key, Message: msg}
}

// NewWithValue builds a failure that also carries the offending value.
func NewWithValue(failureType int, key string, value any, msg string) *api.Error {
	return &api.Error{Type: failureType, Key: key, KeyValue: value, Message: msg}
}
