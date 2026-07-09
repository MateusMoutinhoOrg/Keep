# Getting Started

## Installation

```sh
go get github.com/MateusMoutinhoOrg/Keep@v0.0.1
```

Requires Go 1.22 or newer.

## The three building blocks

Every Keep program follows the same three steps:

1. **Pick a backend (`deps`)** — where the bytes are physically stored.
2. **Describe your data (`Props` + `Schemas`)** — the collections and their fields.
3. **Create the database and use it.**

```go
package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

// 2. Describe your data: one "user" collection with three fields.
var Props = database.Props{
	Path: "myDatabase/", // prefix for every stored key (a folder, with the standard adapter)
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
	// 1. Pick a backend: the standard adapter stores each key as a file.
	deps := keep_deps.New()

	// 3. Create the database and use it.
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
	fmt.Println("created:", created) // {id: 1, email: mateus@gmail.com, username: mateus, age: 27}
}
```

Run it twice and the second run fails with a `KeyConflict` error — `email` and `username` are `Key` fields, so they are unique. See [Error Handling](Errors.md).

## Choosing a backend

| Adapter | Import | Storage |
|---|---|---|
| `standard` | `adapters/standard` | Filesystem — one file per key, survives restarts |
| `native` | `adapters/native` | In-memory — gone when the process exits, great for tests |

```go
deps := standard.New()               // files relative to the working directory
deps := standard.NewWithBase("/srv") // files under /srv
deps := native.New()                 // in-memory
```

You can also write your own backend — see [Dependency Injection](DependencyInjection.md).

## Next steps

- [Schemas](Schemas.md) — all the field types, including nested sub-databases.
- [Working with Records](Records.md) — find, update, delete, and list records.
- [samples/](../samples/) — one runnable program per operation.
