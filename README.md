# Keep

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Keep.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Keep)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Keep)](https://github.com/MateusMoutinhoOrg/Keep/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A storage-independent database built on top of plain key-value operations.

---

## Overview

Keep lets you define schemas with typed fields, unique indexed keys, and nested collections — and runs them over **any** backend that can read, write, and delete a key. It needs no key listing, no prefix scans, and no range queries, so it works the same over the local filesystem, memory, S3-like blob stores, or anything you can wrap in a small struct of functions.

It uses a **Dependency Injection** pattern in which:

- **`/pkg/`** contains the pure library logic — it never imports concrete implementations.
- **`/adapters/`** contains opinionated, concrete implementations of the dependency contract.
- **`/pkg/deps/`** defines the `Deps` struct of injectable functions that all adapters must satisfy.

What you get on top of that:

- **Storage independent** — bring your own backend by populating a small struct, or use the built-in ones ([filesystem](adapters/standard/), [in-memory](adapters/native/)).
- **Constant-time operations** — create, lookup by key, and delete each touch a fixed number of keys, no matter how many records exist.
- **Unique keys** — fields of type `Key` are indexed and enforced unique (case-insensitive).
- **Nested collections** — a record can own sub-databases (e.g. a user owning its sessions).

---

## Quick Start

**1. Install the library:**

```sh
go get github.com/MateusMoutinhoOrg/Keep@v0.0.1
```

**2. Create a `main.go` file:**

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
	// 1. Create deps via an adapter (the "opinionated" part)
	deps := keep_deps.New() // filesystem backend

	// 2. Inject deps into the pure library
	keep := keep_lib.New(deps)

	// 3. Use the library — it never knows which adapter is behind the scenes
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

**3. Run:**

```sh
go run main.go
```

---

### Consumer Documentation

> For developers who want to **use** this project as a library.

#### Guides

| Guide | Description |
|-------|-------------|
| [Schemas](./docs/Consumer/Schemas.md) | Defining collections, field types, and sub-databases |
| [Records](./docs/Consumer/Records.md) | Create, find, read, update, delete, and list records |
| [Errors](./docs/Consumer/Errors.md) | The error types and how to react to them |
| [DependencyInjection](./docs/Consumer/DependencyInjection.md) | Choosing a backend or writing your own |
| [RequiredApi](./docs/Consumer/RequiredApi.md) | The contract each `Deps` function must honor |

#### Use Cases

| Use Case | Description |
|----------|-------------|
| [LibInitialization](./docs/Consumer/UseCases/LibInitialization.md) | How to initialize the library |
| [CreatingCustomDeps](./docs/Consumer/UseCases/CreatingCustomDeps.md) | How to create custom dependency adapters |
| [OverwritingDeps](./docs/Consumer/UseCases/OverwritingDeps.md) | How to overwrite specific dependencies |
| [RunSample](./docs/Consumer/UseCases/RunSample.md) | How to run the provided examples |

#### [Public API](./docs/Consumer/PublicApi.md)
Index of all public-facing components (structs, functions, and methods), with links to their respective detail files.

#### Samples

| Sample | Description |
|--------|-------------|
| [CreateUser](./examples/CreateUser/CreateUser.go) | Insert a record with unique keys |
| [FindUserByKey](./examples/FindUserByKey/FindUserByKey.go) | Look a record up by a unique field |
| [RetrieveUserInfo](./examples/RetrieveUserInfo/RetrieveUserInfo.go) | Read individual fields of a record |
| [UpdateUser](./examples/UpdateUser/UpdateUser.go) | Update a plain field |
| [UpdateUserKey](./examples/UpdateUserKey/UpdateUserKey.go) | Update a unique indexed field (re-index) |
| [DeleteUser](./examples/DeleteUser/DeleteUser.go) | Remove a record and its index entries |
| [ListAllUsers](./examples/ListAllUsers/ListAllUsers.go) | Iterate every record of a collection |
| [ListUsersPaginated](./examples/ListUsersPaginated/ListUsersPaginated.go) | Paginate through a collection |
| [SubInfos](./examples/SubInfos/SubInfos.go) | Manage nested sub-database records |

---

### Developer Documentation

> For developers who want to **contribute** to or extend this project.

> [!IMPORTANT]
> **Must Read before contributing.** The following documents are **required reading** for every developer. Do not open a pull request or make changes without first reading both:
>
> | Document | Why it's required |
> |----------|-------------------|
> | [Rules](./docs/Developer/RULES.md) | The contribution rules and guidelines that **must** be followed for any change to be accepted. |
> | [Structure](./docs/Developer/STRUCTURE.md) | The project's directory layout and the purpose of each component — needed to know **where** changes belong. |

#### Internals

| Document | Description |
|----------|-------------|
| [DatabaseSchema](./docs/Developer/DatabaseSchema.md) | The Dense Record Pattern that powers Keep's storage layer |
| [DocumentationStandards](./docs/Developer/DocumentationStandards.md) | Documentation standards and conventions |

#### Use Cases

| Use Case | Description |
|----------|-------------|
| [HandleDocumentation](./docs/Developer/UseCases/HandleDocumentation.md) | How to add, update, rename, or delete documentation and use cases |
| [HandleLibFunctions](./docs/Developer/UseCases/HandleLibFunctions.md) | How to add functions and objects to the library and expose them in the public API |
| [HandleDependencies](./docs/Developer/UseCases/HandleDependencies.md) | How to add dependency requirements and create or update adapters |
| [HandleSamples](./docs/Developer/UseCases/HandleSamples.md) | How to add and run executable samples |
| [HandleModule](./docs/Developer/UseCases/HandleModule.md) | How to rename the Go module |

---

## License

This project is licensed under the [MIT License](./LICENSE).
