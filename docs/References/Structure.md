# Project Structure

This document maps the project **schema** — the kinds of files this project is built from — rather than listing every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/References/Specs.md) to get its description and sample.

The project is split into four top-level trees, and the dependency flow between them is one-way:

```
adapters/  ──▶  sandbox/  ◀──  examples/ , tests/
(reaches the OS)  (closed)     (wire the two together)
```

- **`/sandbox/`** is a **closed sandbox**: the whole database engine. Nothing inside it may import `adapters/`, `examples/`, `tests/`, a third-party module, or any OS-bound standard-library package. Every effect it needs arrives through the injected `Deps`. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- **`/adapters/`** sits outside the sandbox and is the only place OS-bound and third-party code is allowed. Each adapter imports `sandbox/contracts/deps` and nothing else from the sandbox.
- **`/examples/`** and **`/tests/`** sit outside the sandbox too, and are the only places where an adapter and the sandbox meet.

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

## `/sandbox/`
The closed sandbox — the whole database engine. It holds its own entry point, the contracts everything is wired through, and the internal implementation. It reaches nothing outside itself: every OS-bound or third-party effect arrives through the injected `Deps`. Its package is named `lib`, so consumers import it as `lib "…/sandbox"` and call `lib.New`.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | The `New` constructor wiring `Deps` into the internal `Lib` | |

### `/sandbox/contracts/`
The struct contracts the rest of the project is wired through — the only part of the sandbox anything outside it may import. Contracts import nothing from `adapters/` or `sandbox/internal/`.

#### `/sandbox/contracts/deps/`
The contract every adapter must fill. Its behavioral requirements are spelled out in [RequiredApi.md](/docs/References/RequiredApi.md).

| File | Description | Spec |
|------|-------------|------|
| `deps.go` | The `Deps` struct, one function field per injectable behavior, plus the sentinel errors adapters must return | Deps |

#### `/sandbox/contracts/api/`
The library's whole public surface. **Every type in the project is declared here** — never in `sandbox/internal/`. A struct that carries behavior (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`) leads with a `Deps` field and fills the rest as function fields, each assigned by a factory in `sandbox/internal/`; a struct with no behavior (`Props`, `Schema`, `Item`, `Error`) is plain data, buildable with a composite literal. It depends only on `sandbox/contracts/deps`, so it sits near the bottom of the dependency graph.

| File | Description | Spec |
|------|-------------|------|
| `api.go` | The objects handed back (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`), the description passed in (`Props`, `Schema`, `Item`), `Error`, and the field-type and failure-cause constants | Outputs |

### `/sandbox/internal/`
The factories filling the [`api`](#sandboxcontractsapi) structs' function fields. `sandbox/internal/` declares no types at all — only `<Field>Factory` functions and the constructors that run them. Go's `internal/` rule makes it unreachable from outside `sandbox/`, so neither consumers nor `adapters/` can reach in — the sandbox wall is enforced by the compiler, not by convention alone.

#### `/sandbox/internal/lib/`
The entry-point factories. The `internal/` parent already marks it private, so the package carries no `internal_` prefix.

| File | Description | Spec |
|------|-------------|------|
| `lib.go` | `NewDatabaseFactory` plus `New`, the factory aggregate building `api.Lib` and running every factory over it | LibFunctions |

#### `/sandbox/internal/<object>/`
One package per object the library creates, named after the object itself. Each factory closes over the `api` struct being built, reading its `Deps` field at call time.

| File | Description | Spec |
|------|-------------|------|
| `database/database.go` | Factories for `api.KeepDatabase` (`GetSchemaFactory`) plus `New`, resolving collections through `GetSchema` | LibObjects |
| `schemainstance/schemainstance.go` | Factories for `api.SchemaInstance` (`NewItemFactory`, `FindByKeyFactory`, `ListAllFactory`, `ListFactory`) plus `New` | LibObjects |
| `schemaitem/schemaitem.go` | Factories for `api.SchemaItem` (`GetFactory`, `UpdateFactory`, `RemoveFactory`, sub-database factories) plus the record constructors (`New`, `ResolveLive`, `ListRange`, `ClearCollection`) the other packages call | LibObjects |
| `liberror/liberror.go` | Plain builders for `*api.Error` (`New`, `NewWithValue`) — not factories, since `Error` carries no `Deps` field and nothing to fill after construction | LibObjects |

#### `/sandbox/internal/dense/`
The shared engine helpers every object is built on — not an object itself, so it fills no `api` struct's fields.

| File | Description | Spec |
|------|-------------|------|
| `dense.go` | The key layout, value encoding, and counters described in [DenseRecordPattern.md](/docs/Explanations/DenseRecordPattern.md) | LibFunctions |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontractsdeps) contract — the concrete storage backends, each satisfying every behavior described in [RequiredApi.md](/docs/References/RequiredApi.md). This is where OS-bound and third-party code lives; an adapter imports `sandbox/contracts/deps` and nothing else from `sandbox/`. Which adapter to use when is listed in [Adapters.md](/docs/References/Adapters.md).

### `/adapters/<name>/`
One directory per backend, packaged under its own name. `standard` (filesystem) is the default; `native` (in-memory) is the one used by tests and prototypes.

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct filling every `Deps` field via factories, exposed by a `New(...) deps.Deps` factory aggregate | Adapters |

---

## `/examples/`
Outside the sandbox. Runnable samples demonstrating the library, one per database operation — the only place, together with `tests/`, where an adapter and the sandbox are wired together. Samples using the standard adapter write their data under `testDatabase/` (gitignored).

### `/examples/<sample>/`

| File | Description | Spec |
|------|-------------|------|
| `<sample>.go` | Self-contained `package main` program showing one use case | Examples |

**Run a sample:**
```sh
go run ./examples/<sample>/<sample>.go
```

---

## `/tests/`
Outside the sandbox. The engine's test suite, kept here rather than beside the code because it wires in real adapters — something no file under `sandbox/` is allowed to do.

| File | Description | Spec |
|------|-------------|------|
| `<name>_test.go` | Runs every operation against both built-in adapters | |

**Run the tests:**
```sh
go test ./...
```

---

## `/docs/`
Documentation of the project, split by the kind of material it holds.

### `/docs/References/`
Listable material — structures, rules, specifications, and the public API.

| File | Description | Spec |
|------|-------------|------|
| `RULES.md` | The binding contribution rules and their required companion updates | Rules |
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `PublicApi.md` | Index of the public structs, fields, and functions, with links to their detail pages | ReferenceDocs |
| `Adapters.md` | Lists every shipped adapter and when to use each one | AdaptersDoc |
| `RequiredApi.md` | The contract each `Deps` field must honor | ReferenceDocs |
| `Errors.md` | The error types returned by database operations and how to react to them | ReferenceDocs |
| `TemplateFileActions.md` | The action each file takes when the structure is reused for another library | ReferenceDocs |

#### `/docs/References/Meta/`
The specifications describing how each kind of file must be shaped. Never browse this directory — locate a specification through `Specs.md`.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/Specs.md` | The required shape of the artifact the specification governs | |
| `<Spec>/sample.<ext>` | Concrete reference implementation of the specification | |

#### `/docs/References/PublicApi/`
One detail page per public API entry.

| File | Description | Spec |
|------|-------------|------|
| `<pkg>.<Symbol>.md` | The fields and usage of one public entry | ReferenceDocs |

---

### `/docs/Explanations/`
How the project's mechanics and features work.

| File | Description | Spec |
|------|-------------|------|
| `SandboxIsolation.md` | Why the engine lives in a closed sandbox and what the wall forbids | ExplanationDocs |
| `StructContracts.md` | Why every contract is a struct of function fields filled by factories, not an interface | ExplanationDocs |
| `DepsMechanic.md` | How storage dependencies are injected, propagated, and implemented | ExplanationDocs |
| `Schemas.md` | How collections, field types, and sub-databases are described | ExplanationDocs |
| `Records.md` | How records are created, found, read, updated, deleted, and listed | ExplanationDocs |
| `DenseRecordPattern.md` | The key layout and procedures behind the storage engine | ExplanationDocs |

---

### `/docs/Tutorials/`
Workflow guides, one goal per file.

| File | Description | Spec |
|------|-------------|------|
| `<Goal>.md` | The numbered steps achieving one goal (e.g. `AddSample.md`) | TutorialDocs |
