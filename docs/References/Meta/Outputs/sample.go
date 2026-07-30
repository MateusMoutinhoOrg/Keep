//go:build ignore

// This file is an illustrative sample, not part of the build.
package api

// Classification is exported as plain int constants, reported by an
// int-returning method — an int is a primitive, a defined type is not.
const (
	ExampleKeyItem = iota
	ExampleIntItem
)

// ExampleItem is data the caller passes in. It is an interface like
// everything else here, so it is built through the matching constructor
// in sandbox/description.go rather than a composite literal.
type ExampleItem interface {
	Name() string
	Type() int
}

// ExampleError is the typed failure handed back. Embedding the standard
// error lets callers treat it as one.
type ExampleError interface {
	error
	Type() int
	Key() string
}

// ExampleLibObject is an object handed back by the library, created
// through Lib.NewExampleObject with the deps already wired in. Every
// method takes and returns only primitives or interfaces.
type ExampleLibObject interface {
	ExampleObjectMethod(fieldName string) (any, ExampleError)
}

// Lib is the entry point handed back by lib.New. Every object it
// creates carries the same deps it was built with.
type Lib interface {
	NewExampleObject(itens ...ExampleItem) ExampleLibObject
}
