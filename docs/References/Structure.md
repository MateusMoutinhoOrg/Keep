# Project Structure

This document maps the project **schema** — the kinds of files this project is built from — rather than listing every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/References/Specs.md) to get its description and sample.

The project is a **library** and nothing else: it ships no binary and no interface of its own. It is split into four top-level trees, and the dependency flow between them is one-way:

```
adapters/  ──▶  sandbox/  ◀──  examples/libraryExamples/ , tests/
(reaches the OS)  (closed)     (wire the two together)
```

- **`/sandbox/`** is a **closed sandbox**: the whole database engine. Nothing inside it may import `adapters/`, `examples/libraryExamples/`, `tests/`, a third-party module, or any OS-bound standard-library package. Every effect it needs arrives through the injected `Deps`. See [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- **`/adapters/`** sits outside the sandbox and is the only place OS-bound and third-party code is allowed. Each adapter imports `sandbox/contracts/deps` and nothing else from the sandbox.
- **`/examples/libraryExamples/`** and **`/tests/`** sit outside the sandbox too, and are the only places where an adapter and the sandbox meet.

Because every behavior is a field of `api.Lib`, a consumer never calls into the sandbox's internals: it builds a `Deps`, hands it to `lib.New`, and reads the struct it gets back.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and the Doc Index pointing at each theme index under `docs/Index/` | Readme |
| `AGENTS.md` | Entry point for AI coding agents working on this repository | |
| `CLAUDE.md` | Guidance for Claude Code on this repository's architecture and conventions | |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition and dependencies | |
| `.gitignore` | Intentionally untracked files, including the samples' `testDatabase/` output | |

---

## `/sandbox/`
The closed sandbox — the whole database engine. It holds its own entry point, the contracts everything is wired through, the configuration constants, and the internal implementation. It reaches nothing outside itself: every OS-bound or third-party effect arrives through the injected `Deps`. Its package is named `lib`, so consumers import it as `lib "…/sandbox"` and call `lib.New`.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | The `New` constructor delegating to the internal lib constructor, which stores `Deps` on `api.Lib` and runs the factories over it | |

### `/sandbox/contracts/`
The struct contracts the rest of the project is wired through — the only part of the sandbox anything outside it may import. Contracts hold the project's **public types** and are structs of function fields, never interfaces; see [StructContracts.md](/docs/References/StructContracts.md). Contracts import nothing from `adapters/` or `sandbox/lib/`.

#### `/sandbox/contracts/deps/`
The contract every adapter must fill. Its behavioral requirements are spelled out in [RequiredApi.md](/docs/References/RequiredApi.md).

| File | Description | Spec |
|------|-------------|------|
| `deps.go` | The `Deps` struct, one function field per injectable behavior, plus the sentinel errors adapters must return | Deps |

#### `/sandbox/contracts/api/`
The library's whole public surface. **Every type in the project is declared here** — never in `sandbox/lib/`. A struct that carries behavior (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`) leads with a `Deps` field and fills the rest as function fields, each assigned by a factory in `sandbox/lib/`; a struct with no behavior (`Props`, `Schema`, `Item`, `Error`) is plain data, buildable with a composite literal. It depends only on `sandbox/contracts/deps`, so it sits near the bottom of the dependency graph.

| File | Description | Spec |
|------|-------------|------|
| `api.go` | The objects handed back (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`), the description passed in (`Props`, `Schema`, `Item`), `Error`, and the field-type and failure-cause constants | Outputs |

### `/sandbox/config/`
Static configuration the sandbox reads at compile time: the fixed text and numbers the library reports. Holding them as Go constants rather than as files keeps every reference under the compiler's eye — a renamed constant is a build failure rather than a blank line at runtime — and costs no read at all. Nothing outside the sandbox imports this package.

| File | Description | Spec |
|------|-------------|------|
| `version.go` | The `Version` constant reported by the `api.Lib.Version` field | |

### `/sandbox/lib/`
The factories filling the [`api`](#sandboxcontractsapi) structs' function fields, plus the entry point's constructor. `sandbox/lib/` declares no types at all — only `<Field>Factory` functions and the constructors that run them. The whole tree is private to the sandbox: nothing in `adapters/`, `examples/libraryExamples/`, `tests/`, or a consuming project may import it.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | `New(d deps.Deps) api.Lib`, the factory aggregate building the entry point and assigning every `publicfunctions` factory's return value into its field | LibFunctions |

#### `/sandbox/lib/publicfunctions/`
One file per public function field of `api.Lib`. Each file holds its `<Field>Factory(l *api.Lib)` that returns a closure; the `New` constructor in `sandbox/lib/new.go` calls every factory here.

| File | Description | Spec |
|------|-------------|------|
| `<Function>.go` | One file per lib function, holding its `<Field>Factory(l *api.Lib)` — `NewDatabase.go`, `Version.go` | LibFunctions |

#### `/sandbox/lib/<object>/`
One package per object the library creates, named after the object itself. Each factory closes over the `api` struct being built, reading its `Deps` field at call time.

| File | Description | Spec |
|------|-------------|------|
| `database/database.go` | Factories for `api.KeepDatabase` (`GetSchemaFactory`) plus `New`, resolving collections through `GetSchema` | LibObjects |
| `schemainstance/schemainstance.go` | Factories for `api.SchemaInstance` (`NewItemFactory`, `FindByKeyFactory`, `ListAllFactory`, `ListFactory`) plus `New` | LibObjects |
| `schemaitem/schemaitem.go` | Factories for `api.SchemaItem` (`GetFactory`, `UpdateFactory`, `RemoveFactory`, sub-database factories) plus the record constructors (`New`, `ResolveLive`, `ListRange`, `ClearCollection`) the other packages call | LibObjects |
| `liberror/liberror.go` | Plain builders for `*api.Error` (`New`, `NewWithValue`) — not factories, since `Error` carries no `Deps` field and nothing to fill after construction | LibObjects |

#### `/sandbox/lib/dense/`
The shared engine helpers every object is built on — not an object itself, so it fills no `api` struct's fields.

| File | Description | Spec |
|------|-------------|------|
| `dense.go` | The key layout, value encoding, and counters described in [DenseRecordPattern.md](/docs/References/DenseRecordPattern.md) | LibFunctions |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontractsdeps) contract — the concrete storage backends, each satisfying every behavior described in [RequiredApi.md](/docs/References/RequiredApi.md). This is where OS-bound and third-party code lives; an adapter imports `sandbox/contracts/deps` and nothing else from `sandbox/`. An adapter fills its contract with the same **factories** [`sandbox/lib/`](#sandboxlib) uses — the carrier is the adapter struct, which declares the `Deps` field the factories' return values are assigned into. Which adapter to use when is listed in [Adapters.md](/docs/References/Adapters.md).

### `/adapters/<name>/`
One directory per backend, packaged under its own name. `standard` (filesystem) is the default; `native` (in-memory) is the one used by tests and prototypes.

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct carrying a `Deps` field, one `<Field>Factory` per `Deps` field returning a closure, plus the `New(...) deps.Deps` constructor that assigns every factory's return value and returns the filled contract struct | Adapters |

---

## `/examples/libraryExamples/`
Outside the sandbox. Runnable Go samples demonstrating the library, one per database operation — the only place, together with `tests/`, where an adapter and the sandbox are wired together. Samples using the standard adapter write their data under `testDatabase/` (gitignored). Every one of them is listed in [ApiSamplesList.md](/docs/References/ApiSamplesList.md).

### `/examples/libraryExamples/<sample>/`

| File | Description | Spec |
|------|-------------|------|
| `<sample>.go` | Self-contained `package main` program showing one use case | LibraryExamples |

**Run a sample:**
```sh
go run ./examples/libraryExamples/<sample>/<sample>.go
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
Documentation of the project, split by **kind of page**: `Index/` holds one entry point per theme, `Tutorials/` holds every workflow, `References/` holds every lookup and explanation. A **theme** — what the reader wants to accomplish — is not a directory: it is the index that lists a page. `Tutorials/` and `References/` are flat, so a page's file name must be unique inside the directory it lands in. The [README](/README.md) links to the three indexes and to nothing else inside `docs/`.

| Directory | Description |
|-----------|-------------|
| `Index/` | One entry point per theme, each listing the pages of its theme |
| `Tutorials/` | Every workflow page of the project, whatever theme it belongs to |
| `References/` | Every lookup and explanation page, whatever theme it belongs to |

### `/docs/Index/`
One page per theme. The three themes are `LibUsage` — installing the module, consuming the library from Go, and the public API; `Development` — contributing: the rules, the mechanics, the workflows, the specifications; and `Templating` — turning this repository into another library.

| File | Description | Spec |
|------|-------------|------|
| `<Theme>.md` | The theme's entry point: its Tutorials and its References, each entry listing that page's sections | Index |

### `/docs/Tutorials/`
One page per subject, its workflows written as numbered steps. A page can belong to one or more themes, and the theme indexes in [`/docs/Index/`](#docsindex) are what say which.

| File | Description | Spec |
|------|-------------|------|
| `<Subject>.md` | The numbered steps achieving the subject's workflows (e.g. `HandleLibrarySamples.md`) | TutorialDocs |

### `/docs/References/`
One page per lookup table or explained mechanic, plus the two directories the project's biggest listings live in.

| File | Description | Spec |
|------|-------------|------|
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `PublicApi.md` | Index of the public structs, fields, and functions, with links to their detail pages | ReferenceDocs |
| `Adapters.md` | Lists every shipped adapter and when to use each one | AdaptersDoc |
| `ApiSamplesList.md` | Every Go sample shipped in `examples/libraryExamples/` | ReferenceDocs |
| `RequiredApi.md` | The contract each `Deps` field must honor | ReferenceDocs |
| `Errors.md` | The error types returned by database operations and how to react to them | ReferenceDocs |
| `TemplateFileActions.md` | The action each file takes when the structure is reused for another library | ReferenceDocs |
| `SandboxIsolation.md` | Why the engine lives in a closed sandbox and what the wall forbids | ExplanationDocs |
| `StructContracts.md` | Why every contract is a struct of function fields filled by factories, not an interface | ExplanationDocs |
| `DepsMechanic.md` | How storage dependencies are injected, propagated, and implemented | ExplanationDocs |
| `Schemas.md` | How collections, field types, and sub-databases are described | ExplanationDocs |
| `Records.md` | How records are created, found, read, updated, deleted, and listed | ExplanationDocs |
| `DenseRecordPattern.md` | The key layout and procedures behind the storage engine | ExplanationDocs |

#### `/docs/References/Specs/`
The specifications describing how each kind of file must be shaped. Never browse this directory — locate a specification through `Specs.md`.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/Specs.md` | The required shape of the artifact the specification governs | |
| `<Spec>/sample.<ext>` | Concrete reference implementation of the specification | |

#### `/docs/References/PublicApi/`
One detail page per public API entry, indexed by `PublicApi.md`. Reach a page through that index rather than by browsing the directory.

| File | Description | Spec |
|------|-------------|------|
| `<pkg>.<Symbol>.md` | The fields and usage of one public entry | ReferenceDocs |
