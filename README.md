# Keep

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


## Documentation

### LibUsage

Using the project as a Go module - wiring the deps, calling the sandbox

| Doc | Description |
| --- | --- |
| [LibUsage](docs/LibUsage/doc.md) | Use Keep as a Go module: wire the deps, build the sandbox, call its API |
| [PublicApi](docs/PublicApi/doc.md) | Every exported symbol of Keep, generated from the contract sources and their doc comments |
| [LibExamples](docs/LibExamples/doc.md) | Index of every runnable example of Keep as a Go module |
| [Schemas](docs/Schemas/doc.md) | Declaring a database: Props, schemas, field types and nested collections |

### ApiUsage

The exported surface - contracts, deps and what each one promises

| Doc | Description |
| --- | --- |
| [Errors](docs/Errors/doc.md) | Every failure a database operation reports, and what raises it |

### Architecture

How the project is put together - layers, boundaries, data flow

| Doc | Description |
| --- | --- |
| [Adapters](docs/Adapters/doc.md) | Contract, adapter and available: three units, one field of Deps, and who fills it |
| [Dense Record Pattern](docs/DenseRecordPattern/doc.md) | The key layout that makes a schema database run over single-key storage |

### Development

Changing this repository - schema, build mechanics, recipes

| Doc | Description |
| --- | --- |
| [Requirements](docs/Requirements/doc.md) | The two tools this project needs — Go and agnos — installed per platform |
| [Workflow](docs/Workflow/doc.md) | Every change this project takes and the agnos command that makes it |
| [Rules](docs/Rules/doc.md) | Every rule the generators, `verify` and the hand-written files must hold to |
| [Structure](docs/Structure/doc.md) | The project schema: what lives where, what is generated, what verify enforces |

### Reference

Lookup tables - schemas, file formats, generated file listings

| Doc | Description |
| --- | --- |
| [EntriesYaml](docs/EntriesYaml/doc.md) | Every key of a command's entries.yaml and what the generated code does with it |
| [Extensions](docs/Extensions/doc.md) | The generation mechanics this project turns on, and what each one writes |
| [DepList](docs/DepList/doc.md) | Every dep `agnos add-dep` can add, the adapters that fill it, and what backs each one |
| [GeneratedFiles](docs/GeneratedFiles/doc.md) | Every file agnos writes into this project and whether build overwrites it |
| [LibExamples](docs/LibExamples/doc.md) | Index of every runnable example of Keep as a Go module |
| [Storage Contract](docs/StorageContract/doc.md) | What an adapter filling Deps.Storagedeps has to guarantee, field by field |

## License

MIT License

Copyright (c) 2026 MateusMoutinho

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

