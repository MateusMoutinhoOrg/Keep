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

- **`/pkg/lib/`** contains the pure library logic — it never imports concrete implementations.
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

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

var Props = lib.Props{
	Path: "myDatabase/",
	Schemas: []lib.Schema{
		{
			Name: "user",
			Itens: []lib.Item{
				{Name: "email", Type: lib.Key, Required: true},
				{Name: "username", Type: lib.Key, Required: true},
				{Name: "age", Type: lib.Int, Required: true},
			},
		},
	},
}

func main() {
	// 1. Create deps via an adapter (the "opinionated" part)
	deps := standard.New() // filesystem backend

	// 2. Inject deps into the pure library
	keep := lib.New(deps)

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

> [!IMPORTANT]
> **Must Read before contributing.** The following documents are **required reading** for every developer. Do not open a pull request or make changes without first reading them:
>
> | Document | Why it's required |
> |----------|-------------------|
> | [Rules](/docs/Reference/RULES.md) | The contribution rules and guidelines that **must** be followed for any change to be accepted. |
> | [Structure](/docs/Reference/Structure.md) | The project's directory layout and the purpose of each component — needed to know **where** changes belong. |
> | [Specs](/docs/Reference/Specs.md) | The index of every specification — needed to know **how** the file you are about to touch must be shaped. |

### Reference Documentation

> Listable material — structures, rules, specifications, and the public API.

| Name | Description |
|:-|:-|
| <a id="reference-structure"></a>[Structure.md](/docs/Reference/Structure.md) | **Reference** — The project's directory layout and the purpose of each component. |
| <a id="reference-rules"></a>[RULES.md](/docs/Reference/RULES.md) | **Reference** — The binding contribution rules and their required companion updates. |
| <a id="reference-specs"></a>[Specs.md](/docs/Reference/Specs.md) | **Reference** — Lists every specification and the files each one governs. |
| <a id="reference-public-api"></a>[PublicApi.md](/docs/Reference/PublicApi.md) | **Reference** — Index of all public structs, functions, and methods with detail links. |
| <a id="reference-required-api"></a>[RequiredApi.md](/docs/Reference/RequiredApi.md) | **Reference** — The contract each `Deps` function must honor to power the library. |
| <a id="reference-errors"></a>[Errors.md](/docs/Reference/Errors.md) | **Reference** — The error types returned by operations and how to react to them. |
| <a id="reference-template-file-actions"></a>[TemplateFileActions.md](/docs/Reference/TemplateFileActions.md) | **Reference** — The action each file takes when forking or adapting: copy, create, rewrite, delete. |

---

### Explanation Documentation

> How the project's mechanics and features work.

| Name | Description |
|:-|:-|
| <a id="explanation-deps-mechanic"></a>[DepsMechanic.md](/docs/Explanation/DepsMechanic.md) | **Explanation** — Choosing a backend, overwriting deps, or writing your own. |
| <a id="explanation-schemas"></a>[Schemas.md](/docs/Explanation/Schemas.md) | **Explanation** — Defining collections, field types, and nested sub-databases. |
| <a id="explanation-records"></a>[Records.md](/docs/Explanation/Records.md) | **Explanation** — Creating, finding, reading, updating, deleting, and listing records. |
| <a id="explanation-dense-record-pattern"></a>[DenseRecordPattern.md](/docs/Explanation/DenseRecordPattern.md) | **Explanation** — The key layout and procedures behind the storage engine. |

---

### Tutorials

> Workflow guides, grouped by context. Each tutorial covers a single goal.

#### Getting Started

| Name | Description |
|:-|:-|
| <a id="tutorial-lib-initialization"></a>[LibInitialization.md](/docs/Tutorials/LibInitialization.md) | **Tutorial** — Install the lib, create deps via an adapter, and run a first program. |
| <a id="tutorial-run-sample"></a>[RunSample.md](/docs/Tutorials/RunSample.md) | **Tutorial** — Browse and run the executable samples in the examples/ directory. |

#### Using the Database

| Name | Description |
|:-|:-|
| <a id="tutorial-define-database"></a>[DefineDatabase.md](/docs/Tutorials/DefineDatabase.md) | **Tutorial** — Describe a database with its collections and open it in a program. |
| <a id="tutorial-add-schema-field"></a>[AddSchemaField.md](/docs/Tutorials/AddSchemaField.md) | **Tutorial** — Add a field to a collection that already holds records. |
| <a id="tutorial-add-nested-collection"></a>[AddNestedCollection.md](/docs/Tutorials/AddNestedCollection.md) | **Tutorial** — Give a record its own nested collection of sub-records. |

#### Library Development

| Name | Description |
|:-|:-|
| <a id="tutorial-add-lib-function"></a>[AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) | **Tutorial** — Add a function to pkg/lib/ and wire it to the injected deps. |
| <a id="tutorial-add-lib-object"></a>[AddLibObject.md](/docs/Tutorials/AddLibObject.md) | **Tutorial** — Add an object created by the lib, with its deps wired in by the constructor. |
| <a id="tutorial-add-database-operation"></a>[AddDatabaseOperation.md](/docs/Tutorials/AddDatabaseOperation.md) | **Tutorial** — Add an engine operation without breaking the dense key layout. |
| <a id="tutorial-add-dependency"></a>[AddDependency.md](/docs/Tutorials/AddDependency.md) | **Tutorial** — Add a field to the Deps contract and implement it in every adapter. |
| <a id="tutorial-add-adapter"></a>[AddAdapter.md](/docs/Tutorials/AddAdapter.md) | **Tutorial** — Create a new opinionated storage backend for the Deps contract. |
| <a id="tutorial-add-sample"></a>[AddSample.md](/docs/Tutorials/AddSample.md) | **Tutorial** — Create a runnable sample in examples/ and register it in the README. |

#### Documentation

| Name | Description |
|:-|:-|
| <a id="tutorial-add-document"></a>[AddDocument.md](/docs/Tutorials/AddDocument.md) | **Tutorial** — Create or update a .md file and register it in README and Structure. |
| <a id="tutorial-rename-document"></a>[RenameDocument.md](/docs/Tutorials/RenameDocument.md) | **Tutorial** — Rename or move a .md file without leaving broken references behind. |
| <a id="tutorial-delete-document"></a>[DeleteDocument.md](/docs/Tutorials/DeleteDocument.md) | **Tutorial** — Remove a .md file and clear every reference pointing to it. |
| <a id="tutorial-expose-public-api"></a>[ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md) | **Tutorial** — Publish a lib function, object, or method in the public API index. |

#### Templating

| Name | Description |
|:-|:-|
| <a id="tutorial-rename-module"></a>[RenameModule.md](/docs/Tutorials/RenameModule.md) | **Tutorial** — Rename the Go module path and update all internal imports. |
| <a id="tutorial-fork-template"></a>[ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) | **Tutorial** — Use this repo as a template to start a new DI library. |
| <a id="tutorial-adapt-existing-lib"></a>[AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md) | **Tutorial** — Convert a pre-existing library to this DI structure. |

---

#### Samples

| Sample | Description |
|--------|-------------|
| [CreateUser](/examples/CreateUser/CreateUser.go) | Insert a record with unique keys |
| [FindUserByKey](/examples/FindUserByKey/FindUserByKey.go) | Look a record up by a unique field |
| [RetrieveUserInfo](/examples/RetrieveUserInfo/RetrieveUserInfo.go) | Read individual fields of a record |
| [UpdateUser](/examples/UpdateUser/UpdateUser.go) | Update a plain field |
| [UpdateUserKey](/examples/UpdateUserKey/UpdateUserKey.go) | Update a unique indexed field (re-index) |
| [DeleteUser](/examples/DeleteUser/DeleteUser.go) | Remove a record and its index entries |
| [ListAllUsers](/examples/ListAllUsers/ListAllUsers.go) | Iterate every record of a collection |
| [ListUsersPaginated](/examples/ListUsersPaginated/ListUsersPaginated.go) | Paginate through a collection |
| [SubInfos](/examples/SubInfos/SubInfos.go) | Manage nested sub-database records |

---

## License

This project is licensed under the [MIT License](./LICENSE).
