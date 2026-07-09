# Keep

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Keep.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Keep)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Keep)](https://github.com/MateusMoutinhoOrg/Keep/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A storage-independent database built on top of plain key-value operations.

Keep lets you define schemas with typed fields, unique indexed keys, and nested collections — and runs them over **any** backend that can read, write, and delete a key. It needs no key listing, no prefix scans, and no range queries, so it works the same over the local filesystem, memory, S3-like blob stores, or anything you can wrap in a small interface.

- **Storage independent** — bring your own backend by implementing a small interface, or use the built-in ones ([filesystem](adapters/standard/), [in-memory](adapters/native/)).
- **Constant-time operations** — create, lookup by key, and delete each touch a fixed number of keys, no matter how many records exist.
- **Unique keys** — fields of type `Key` are indexed and enforced unique (case-insensitive).
- **Nested collections** — a record can own sub-databases (e.g. a user owning its sessions).

## Quick Start

Install:

```sh
go get github.com/MateusMoutinhoOrg/Keep@v0.0.1
```

Create a database, insert a record, and find it back:

```go
package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

var Props = database.Props{
	Path: "myDatabase/",
	Schemas: []database.Schema{
		{
			Name: "user",
			Itens: []database.Item{
				{Name: "email", Type: database.Key, Required: true},
				{Name: "username", Type: database.Key, Required: true},
				{Name: "age", Type: database.Int, Required: true},
			},
		},
	},
}

func main() {
	deps := keep_deps.New() // filesystem backend
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("user")

	created, err := users.NewItem(map[string]any{
		"email":    "mateus@gmail.com",
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		fmt.Println("error creating user:", err)
		return
	}
	fmt.Println("created:", created)

	found := users.FindByKey("email", "mateus@gmail.com")
	fmt.Println("found:", found)
}
```

Runnable examples for every operation live in [samples/](samples/).

## Use Cases

- **Apps without a database server** — persist structured data straight to the filesystem with the standard adapter; no external service to install or run.
- **Exotic or minimal storage backends** — run a schema-based database over any store that only offers `get`/`set`/`delete` (blob storage, embedded KV stores, remote APIs) by writing a small adapter.
- **Tests and prototypes** — swap the backend for the in-memory adapter and get a zero-setup database with the exact same API.
- **Large collections on dumb storage** — lookups and deletes stay constant-cost even with millions of records, because nothing ever scans or lists keys.

## Documentation

- [Getting Started](docs/GettingStarted.md) — installation and your first database.
- [Schemas](docs/Schemas.md) — defining collections, field types, and sub-databases.
- [Working with Records](docs/Records.md) — create, find, read, update, delete, and list records.
- [Error Handling](docs/Errors.md) — the error types and how to react to them.
- [Dependency Injection](docs/DependencyInjection.md) — choosing a backend or writing your own.
- [Required API](docs/RequiredApi.md) — the interface a custom backend must implement.
- [Database Internals](docs/DatabaseSchema.md) — the Dense Record Pattern that powers Keep.
