//go:build ignore

// This file is an illustrative sample, not part of the build.
package main

import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
	"github.com/MateusMoutinhoOrg/Agnos/pkg/lib"
)

func main() {
	// 1. Build deps through an adapter (the opinionated layer).
	deps := standard.New(3)

	// 2. Inject deps into the pure library.
	l := lib.New(deps)

	// 3. Exercise the library — it never knows which adapter is behind it.
	obj := l.NewExampleObject(1, "2")
	println(obj.ExampleObjectMethod())
}
