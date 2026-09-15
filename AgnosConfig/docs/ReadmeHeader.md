# {{.Name}}

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

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

var Props = api.Props{
	Path: "database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Itens: []api.Item{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
			},
		},
	},
}

func main() {
	deps := standard.New()    // one file per key
	lib := sandbox.New(&deps) // *api.Sandbox

	users, _ := lib.Databases.New(Props).GetSchema("user")

	users.NewItem(map[string]any{"email": "mateus@gmail.com", "age": 27})

	found, _ := users.FindByKey("email", "mateus@gmail.com")
	fmt.Println(found.String())
}
```

Swap `adapters/availables/standard` for `adapters/availables/native` and the same program
runs entirely in memory. Nothing else changes: the library only ever calls the eleven
single-key functions of `sandbox/deps/storagedeps`, and which implementation stands behind
them is decided by the one import a program picks.

| Start here | For |
|---|---|
| [Schemas](docs/Schemas/doc.md) | declaring a database — field types, keys, nested collections |
| [PublicApi](docs/PublicApi/doc.md) | every exported symbol, generated from the contracts |
| [LibExamples](docs/LibExamples/doc.md) | eleven runnable programs, each checked against a golden |
| [StorageContract](docs/StorageContract/doc.md) | writing a backend of your own |
| [DenseRecordPattern](docs/DenseRecordPattern/doc.md) | how iteration works without listing |

This repository is generated and checked by [agnos](https://github.com/MateusMoutinhoOrg/Agnos):
`agnos build` rewrites every generated file, `agnos verify` checks the schema, and
`agnos exec-test` runs every example against its golden. See
[Requirements](docs/Requirements/doc.md) and [Workflow](docs/Workflow/doc.md).
