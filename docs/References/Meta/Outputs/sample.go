//go:build ignore

// This file is an illustrative sample, not part of the build.
package api

import "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"

// Classification is exported as plain int constants, reported through a
// plain Type field — an int is a primitive, a defined type is not.
const (
	ExampleKeyItem = iota
	ExampleIntItem
)

// ExampleItem is data the caller passes in. It carries no behavior, so
// it is plain data, built directly with a composite literal.
type ExampleItem struct {
	Name string
	Type int
}

// ExampleError is the typed failure handed back. It carries no
// behavior — no methods, not even Error() string — so callers read its
// fields directly. *ExampleError stays a pointer so nil still means
// success.
type ExampleError struct {
	Type int
	Key  string
}

// ExampleLibObject is an object handed back by the library, created
// through Lib.NewExampleObject with the deps already wired in. It leads
// with Deps and fills the rest as function fields, each closing over
// this struct.
type ExampleLibObject struct {
	Deps                deps.Deps
	ExampleObjectMethod func(fieldName string) (any, *ExampleError)
}

// Lib is the entry point handed back by lib.New. Every object it
// creates carries the same Deps it was built with.
type Lib struct {
	Deps             deps.Deps
	NewExampleObject func(itens ...ExampleItem) ExampleLibObject
}
