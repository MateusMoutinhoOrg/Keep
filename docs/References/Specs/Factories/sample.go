//go:build ignore

// This file is an illustrative sample, not part of the build.
// It shows the same factory pattern on both sides of the sandbox wall:
// inside sandbox/, where the carrier is an api struct, and inside
// adapters/, where the carrier is the adapter struct.
package example_factories

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// ---------------------------------------------------------------
// Inside the sandbox: the carrier is the api struct being filled,
// and the state the closure reads is its propagated Deps.
// ---------------------------------------------------------------

// ExampleLibFunctionFactory returns the closure that fills
// api.Lib.ExampleLibFunction, closed over l, so the injected dependency is
// read at call time through l.Deps.
func ExampleLibFunctionFactory(l *api.Lib) func(key string) bool {
	return func(key string) bool {
		exists, err := l.Deps.Exists(key)
		return err == nil && exists
	}
}

// New builds the api.Lib entry point and runs every lib factory over it,
// assigning each return value into its matching field. It is the factory
// aggregate — a field left unassigned here stays nil and panics on first
// call.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.ExampleLibFunction = ExampleLibFunctionFactory(&l)
	return l
}

// ---------------------------------------------------------------
// Outside the sandbox: the carrier is the adapter struct, whose Deps
// field is the contract the factories fill. The state the closure reads
// is the adapter's own configuration.
// ---------------------------------------------------------------

// ExampleAdapter fills deps.Deps against an in-process map.
type ExampleAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps deps.Deps
	// values is adapter-specific state, read by the closure at call time.
	values map[string][]byte
}

// ExistsFactory returns the closure that fills deps.Deps.Exists, reading the
// adapter's backing map through a.
func ExistsFactory(a *ExampleAdapter) func(key string) (bool, error) {
	return func(key string) (bool, error) {
		_, ok := a.values[key]
		return ok, nil
	}
}

// NewAdapter creates a deps.Deps backed by the example adapter. It builds the
// adapter instance, runs every field factory over it and assigns its return
// value into the matching field, and returns the **contract struct** — never
// the concrete adapter type.
func NewAdapter() deps.Deps {
	adapter := &ExampleAdapter{values: map[string][]byte{}}
	adapter.Deps.Exists = ExistsFactory(adapter)
	return adapter.Deps
}
