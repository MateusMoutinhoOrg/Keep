# Project Structure

## Root

| File | Description |
|------|-------------|
| `README.md` | Project overview and quick-start guide |
| `AGENTS.md` | Instructions for AI coding agents working on this repository |
| `LICENSE` | License terms for the project |
| `go.mod` | Go module definition and dependencies |
| `.gitignore` | Specifies intentionally untracked files to ignore |

---

## `/adapters/`
Implementations of the [`/pkg/deps`](#pkgdeps) contract, providing different **opinionated** storage backends with distinct behaviors.

### `/adapters/standard/`
The default adapter, storing each key as a file on the filesystem.

| File | Description |
|------|-------------|
| `defaultImplementation.go` | Implements the filesystem-backed `Deps` (constructors `New` and `NewWithBase`) |

### `/adapters/native/`
The in-memory adapter, ideal for tests and prototypes.

| File | Description |
|------|-------------|
| `native.go` | Implements the in-memory `Deps` (constructor `New`) |

---

## `/docs/`
Documentation of the project.

### `/docs/Developer/`
Developer documentation, required only to contribute to this project.

| File | Description |
|------|-------------|
| `RULES.md` | Rules to follow when contributing to this project |
| `STRUCTURE.md` | Structure of this project |
| `DocumentationStandards.md` | Documentation standards and conventions for writing project documentation |
| `DatabaseSchema.md` | The Dense Record Pattern design document behind the storage layer |
| `HandleDocumentation.md` | How to add, update, rename, or delete documentation |
| `HandleLibFunctions.md` | How to add functions and objects to the library and expose them in the public API |
| `HandleDependencies.md` | How to add dependency requirements and create or update adapters |
| `HandleSamples.md` | How to add and run executable samples |
| `HandleModule.md` | How to rename the Go module |

---

### `/docs/Usage/`
End-user documentation that explains how to integrate and use the library.

| File | Description |
|------|-------------|
| `Schemas.md` | Defining collections, field types, and sub-databases |
| `Records.md` | Create, find, read, update, delete, and list records |
| `Errors.md` | The error types and how to react to them |
| `DepsMechanic.md` | Choosing a backend, overwriting deps, or writing your own |
| `RequiredApi.md` | The contract each `Deps` function must honor |
| `LibInitialization.md` | How to initialize the library |
| `RunSample.md` | How to run the provided examples |

#### `/docs/Usage/PublicApi.md`
Index of all public-facing components (structs, functions, and methods), with links to their respective detail files.

#### `/docs/Usage/PublicApi/`
Detailed documentation for each individual public API entry.

| File | Description |
|------|-------------|
| `keep.KeepLib.md` | Documents the main `KeepLib` struct and its `NewDatabase` method |
| `keep.New.md` | Documents the `keep.New` constructor |
| `database.Props.md` | Documents the schema description structs (`Props`, `Schema`, `Item`, `ItemType`) |
| `database.KeepDatabase.md` | Documents the `KeepDatabase` struct and its `GetSchema` method |
| `database.SchemaInstance.md` | Documents the `SchemaInstance` struct and its record operations |
| `database.SchemaItem.md` | Documents the `SchemaItem` struct and its record methods |
| `database.Error.md` | Documents the typed `Error` struct and its `ErrorType` values |
| `deps.Deps.md` | Documents the `Deps` struct of injectable functions |
| `standard.New.md` | Documents the `standard.New` and `standard.NewWithBase` constructors |
| `native.New.md` | Documents the `native.New` constructor |


---

## `/examples/`
Runnable samples demonstrating how to use the library, one per database operation.

### `/examples/<sample>/`

| File | Description |
|------|-------------|
| `<sample>.go` | Self-contained sample showing a specific use case |

**Run a sample:**
```sh
go run ./examples/<sample>/<sample>.go
```

Samples using the filesystem adapter write their data under `testDatabase/` (gitignored).

---

## `/pkg/`
Core of the project — contains the **Dependency Injection System** and all **public interfaces**. Never imports concrete implementations.

### `/pkg/deps/`
Defines the `Deps` struct that all adapters must satisfy.

| File | Description |
|------|-------------|
| `deps.go` | Declares the `Deps` struct with injectable function fields and the sentinel errors |

### `/pkg/keep/`
The main library entry point. Accepts a `Deps` implementation and creates databases with it wired in.

| File | Description |
|------|-------------|
| `new.go` | Constructor (`New`) for creating a `KeepLib` instance with injected deps, and its `NewDatabase` method |

### `/pkg/database/`
The database engine: schemas, records, and the Dense Record Pattern (see [DatabaseSchema.md](./DatabaseSchema.md)).

| File | Description |
|------|-------------|
| `database.go` | Declares `Props` and `KeepDatabase` with its `GetSchema` method |
| `schema.go` | Declares the schema description types (`Schema`, `Item`, `ItemType`) |
| `schema_instance.go` | Collection-level operations (`NewItem`, `FindByKey`, `ListAll`, `List`) |
| `schema_item.go` | Record-level operations (`Get`, `Update`, `Remove`, sub-database methods) |
| `dense.go` | Implementation of the Dense Record Pattern key layout and procedures |
| `error.go` | Declares the typed `Error` struct and its `ErrorType` values |
| `database_test.go` | Tests running every operation against both built-in adapters |
