# Project Structure

This document maps the project **schema** — the kinds of files this project is built from — rather than listing every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/Reference/Specs.md) to get its description and sample.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview, quick start, and the Doc Index | Readme |
| `AGENTS.md` | Entry point for AI coding agents working on this repository | |
| `CLAUDE.md` | Guidance for Claude Code on this repository's architecture and conventions | |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition and dependencies | |
| `.gitignore` | Intentionally untracked files, including the samples' `testDatabase/` output | |

---

## `/adapters/`
Opinionated implementations of the [`Deps`](#pkg) contract — the concrete storage backends. Each one satisfies the whole contract described in [RequiredApi.md](/docs/Reference/RequiredApi.md).

### `/adapters/<name>/`
One directory per backend, packaged under its own name. `standard` (filesystem) is the default; `native` (in-memory) is the one used by tests and prototypes.

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A `New(...) deps.Deps` factory implementing every `Deps` field | Adapters |

---

## `/docs/`
Documentation of the project, split by the kind of material it holds.

### `/docs/Reference/`
Listable material — structures, rules, specifications, and the public API.

| File | Description | Spec |
|------|-------------|------|
| `RULES.md` | The binding contribution rules and their required companion updates | Rules |
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `PublicApi.md` | Index of the public structs, functions, and methods, with links to their detail pages | ReferenceDocs |
| `RequiredApi.md` | The contract each `Deps` function must honor | ReferenceDocs |
| `Errors.md` | The error types returned by database operations and how to react to them | ReferenceDocs |
| `TemplateFileActions.md` | The action each file takes when the structure is reused for another library | ReferenceDocs |

#### `/docs/Reference/Meta/`
The specifications describing how each kind of file must be shaped. Never browse this directory — locate a specification through `Specs.md`.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/` | One directory per specification, holding its `Specs.md` and its `sample` | |

#### `/docs/Reference/PublicApi/`
One detail page per public API entry.

| File | Description | Spec |
|------|-------------|------|
| `<pkg>.<Symbol>.md` | The fields, methods, and usage of one public entry | ReferenceDocs |

---

### `/docs/Explanation/`
How the project's mechanics and features work.

| File | Description | Spec |
|------|-------------|------|
| `DepsMechanic.md` | How storage dependencies are injected, overwritten, and implemented | ExplanationDocs |
| `Schemas.md` | How collections, field types, and sub-databases are described | ExplanationDocs |
| `Records.md` | How records are created, found, read, updated, deleted, and listed | ExplanationDocs |
| `DenseRecordPattern.md` | The key layout and procedures behind the storage engine | ExplanationDocs |

---

### `/docs/Tutorials/`
Workflow guides, one goal per file.

| File | Description | Spec |
|------|-------------|------|
| `<Goal>.md` | The numbered steps achieving one goal (e.g. `AddSample.md`) | TutorialDocs |

---

## `/examples/`
Runnable samples demonstrating the library, one per database operation. Samples using the standard adapter write their data under `testDatabase/` (gitignored).

### `/examples/<sample>/`

| File | Description | Spec |
|------|-------------|------|
| `<sample>.go` | Self-contained `package main` program showing one use case | Examples |

**Run a sample:**
```sh
go run ./examples/<sample>/<sample>.go
```

---

## `/pkg/`
Core of the project — the **Dependency Injection System** and all **public logic**. Never imports concrete implementations.

### `/pkg/deps/`
The contract every adapter must satisfy.

| File | Description | Spec |
|------|-------------|------|
| `deps.go` | Declares the `Deps` struct of injectable function fields and the sentinel errors | Deps |

### `/pkg/lib/`
The whole library: the entry point plus the database engine — schemas, records, and the key layout described in [DenseRecordPattern.md](/docs/Explanation/DenseRecordPattern.md). Every object it creates carries the injected `Deps` in a private `deps` field.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | Declares `Lib` and the `New` constructor injecting `Deps` into it | LibObjects |
| `types.go` | Declares every type the lib creates: `Props`, `Schema`, `Item`, `ItemType`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`, `Error`, `ErrorType` | LibObjects |
| `database.go` | The `NewDatabase` constructor and `KeepDatabase.GetSchema` | LibObjects |
| `schema_instance.go` | Collection-level operations (`NewItem`, `FindByKey`, `ListAll`, `List`) | LibFunctions |
| `schema_item.go` | Record-level operations (`Get`, `Update`, `Remove`, sub-database methods) | LibFunctions |
| `dense.go` | The key layout and procedures backing every operation | LibFunctions |
| `lib_test.go` | Runs every operation against both built-in adapters | |
