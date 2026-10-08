# {{.ProjectName}}

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Keep.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Keep)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Keep)](https://github.com/MateusMoutinhoOrg/Keep/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)](go.mod)

**A storage-independent database.** Schemas with typed fields, unique indexed keys, nested
collections and pagination — over any backend that can read, write and delete a single key.
No listing, no prefix scans, no range queries.

```bash
go get github.com/MateusMoutinhoOrg/Keep@latest
```

```go
package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

var Props = api.Props{
	Path: "database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
			},
		},
	},
}

func main() {
	deps := standard.New()    // one file per key, under the working directory
	lib := sandbox.New(&deps) // *api.Sandbox

	db, failure := lib.Databases.New(Props) // InvalidSchema when Props is wrong
	if failure != nil {
		panic(failure.Message)
	}
	users, ok := db.Collection("user") // ok == false for a name Props does not declare
	if !ok {
		panic("no user schema")
	}

	if _, failure := users.Insert(map[string]any{"email": "mateus@gmail.com", "age": 27}); failure != nil {
		panic(failure.Message)
	}

	found, ok, failure := users.FindByKey("email", "MATEUS@gmail.com")
	if failure != nil || !ok {
		panic("not found")
	}
	fmt.Println(found.String()) // {id: 1, email: "mateus@gmail.com", age: 27}
}
```

Swap `adapters/bindings/standard` for `adapters/bindings/memory` and the same program
runs entirely in memory. Nothing else changes: the library only reaches storage through six
single-key functions of `sandbox/deps/storagedeps` — `Write`, `Read`, `Exists`, `Delete`,
and the `Lock` and `Unlock` of a lease — and which implementation stands behind them is
decided by the one import a program picks.

Every object it hands back is safe to share between goroutines: the writes to one collection
take its write lock, which holds across processes sharing the same files. A write that a
crash cuts short leaves nothing a read mistakes for data, and `Collection.Repair` clears
what it left.

| Start here | For |
|---|---|
| [Schemas](docs/Schemas/doc.md) | declaring a database — field types, keys, nested collections |
| [PublicApi](docs/PublicApi/doc.md) | every exported symbol, generated from the contracts |
| [LibExamples](docs/LibExamples/doc.md) | runnable programs, each checked against a golden |
| [StorageContract](docs/StorageContract/doc.md) | writing a backend of your own |
| [DenseRecordPattern](docs/DenseRecordPattern/doc.md) | how iteration works without listing, and what survives a crash |
| [Errors](docs/Errors/doc.md) | every failure an operation reports |

This repository is generated and checked by [agnos](https://github.com/MateusMoutinhoOrg/Agnos):
`agnos build` rewrites every generated file, `agnos verify` checks the schema, and
`agnos run-examples` runs every example against its golden. See
[Requirements](docs/Requirements/doc.md) and [Workflow](docs/Workflow/doc.md).
