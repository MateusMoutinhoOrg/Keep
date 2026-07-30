# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Keep is a **storage-independent database** built on plain key-value operations. Schemas with typed fields, unique indexed keys, and nested collections run over any backend that can read, write, and delete a single key — no listing, no prefix scans, no range queries.

## Commands

```bash
go build ./...                                    # build everything
go test ./...                                     # run the engine tests against both built-in adapters
go run ./examples/CreateUser/CreateUser.go        # run a sample
```

Samples using the standard adapter write to `testDatabase/` (gitignored) and never reset it — re-running a sample exercises the "already exists" paths.

## Architecture

Four top-level trees, wired through **interfaces**, with a strict one-way dependency flow:

```
adapters/  ──▶  sandbox/  ◀──  examples/ , tests/
(reaches the OS)  (closed)     (wire the two together)

standard.New()  ──▶  deps.Deps  ──▶  lib.New(deps)  ──▶  api.Lib  ──▶  api.KeepDatabase  ──▶  api.SchemaInstance  ──▶  api.SchemaItem
(opinionated impl)   (contract)      (entry point)       (output interfaces, impls in sandbox/internal/)
```

`sandbox/` is a **closed sandbox**: nothing in it may import `adapters/`, `examples/`, `tests/`, a third-party module, or an OS-bound stdlib package (`os`, `net`, `syscall`, …). Every such effect is a `Deps` method reached through the object's `Deps` field. This is a binding rule — see `docs/References/RULES.md` and `docs/Explanations/SandboxIsolation.md`.

- **`sandbox/new.go`** — package `lib`, the only wiring point consumers touch: `New(deps.Deps) api.Lib`. Never imports `adapters/`. Importers alias it: `lib "github.com/MateusMoutinhoOrg/Keep/sandbox"`.
- **`sandbox/contracts/deps/deps.go`** — the `Deps` **interface** plus the sentinel errors adapters must return. Adding a requirement = adding a method here. Its behavioral contract is `docs/References/RequiredApi.md`.
- **`sandbox/contracts/api/api.go`** — the library's whole public surface, and it holds **nothing but interfaces and constants**: the objects handed back (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`), the description passed in (`Props`, `Schema`, `Item`), and `Error`. A struct in this file is a spec violation. Every method must be exported, or `sandbox/internal/` cannot implement it.
- **`sandbox/description.go`** — package `lib`; the constructors callers need because `Props`/`Schema`/`Item` are interfaces and cannot be built with a composite literal: `NewProps`, `NewSchema`, `NewKeyItem`, `NewIntItem`, `NewDatabaseItem`.
- **`sandbox/internal/description/`, `sandbox/internal/liberror/`** — the structs backing those interfaces.
- **`sandbox/internal/lib/`** — `Lib`, the `api.Lib` implementation. Go's `internal/` rule keeps it unreachable from `adapters/`, `examples/`, and consumers.
- **`sandbox/internal/<object>/`** — one package per object (`database`, `schemainstance`, `schemaitem`), each carrying an exported `Deps` field. No `internal_` prefix — the `internal/` parent already says it.
- **`sandbox/internal/dense/`** — the shared key-layout and encoding helpers. It references **no object type**, which is what keeps the package graph acyclic; the record constructors (`New`, `ListRange`, `ClearCollection`) live in `schemaitem` for the same reason.
- **`adapters/<name>/`** — outside the sandbox; the only place OS-bound and third-party code is allowed. Each declares a struct implementing every `Deps` method and exposes a `New(...) deps.Deps` factory returning the **interface**. `standard` is filesystem-backed; `native` is in-memory. Listed in `docs/References/Adapters.md`.
- **`examples/<name>/<name>.go`**, **`tests/`** — outside the sandbox; the only places an adapter and the sandbox meet.

Every object propagates `Deps` to the objects it creates — see `database.KeepDatabase.GetSchema` for the pattern. A method returning "nothing found" must return a **literal `nil`**, never a typed nil pointer, or the caller's `== nil` check silently fails.

Storage layout follows the Dense Record Pattern (`docs/Explanations/DenseRecordPattern.md`): fixed-cost operations, ids that are never reused, and a dense position list that makes iteration possible without key listing. Changes to `sandbox/internal/dense/` and `sandbox/internal/schemaitem/` must preserve those invariants and the documented write orderings.

## Critical: this repo is documentation-driven

Changes are governed by required-reading docs, and several actions **must** update companion files in the same commit. Each tutorial in `docs/Tutorials/` covers exactly one goal — read the one matching your change:

| If you... | Read | And keep in sync |
|-----------|------|------------------|
| add/rename/delete any file or dir | `docs/References/Structure.md` | `docs/References/Structure.md` |
| add/rename/delete a `.md` file | `docs/Tutorials/AddDocument.md`, `RenameDocument.md`, `DeleteDocument.md` | Doc Index in `README.md` |
| add a lib function/object | `docs/Tutorials/AddLibFunction.md`, `AddLibObject.md`, `AddDatabaseOperation.md` | `sandbox/contracts/api/api.go` + `docs/References/PublicApi.md` (+ detail page in `docs/References/PublicApi/`, see `ExposePublicApi.md`) |
| add a `Deps` method | `docs/Tutorials/AddDependency.md` | **every** adapter in `adapters/`, plus `docs/References/RequiredApi.md` |
| add an adapter | `docs/Tutorials/AddAdapter.md` | `docs/References/Structure.md`, `docs/References/Adapters.md` |
| need an OS/third-party call inside `sandbox/` | `docs/Explanations/SandboxIsolation.md`, `docs/Tutorials/AddDependency.md` | `sandbox/contracts/deps/deps.go` + **every** adapter |
| add/rename/delete a sample | `docs/Tutorials/AddSample.md` | Samples section in `README.md` |
| reuse this structure for another library | `docs/Tutorials/ForkTemplate.md`, `AdaptExistingLib.md` | `docs/References/TemplateFileActions.md` (the per-file copy/create/rewrite/delete list both tutorials follow) |

`docs/References/RULES.md` is the binding rule set and `docs/References/Specs.md` is the index of every file specification; `AGENTS.md` points here. Adding a `Deps` method without implementing it in all adapters breaks every consumer — that's the most common footgun. A new public lib function or object must be declared in `sandbox/contracts/api/api.go` **and** implemented in `sandbox/internal/`, or callers cannot reach it.

## Conventions

- Module path is `github.com/MateusMoutinhoOrg/Keep`; renaming it is a documented procedure — see `docs/Tutorials/RenameModule.md`.
- Public-facing API entries each get a detail page under `docs/References/PublicApi/` named `<pkg>.<Symbol>.md`.
- Every exported signature takes and returns only primitives or interfaces — never a struct. Operations return a typed `api.Error` (nil on success); switch on its `Type()` rather than matching strings — see `docs/References/Errors.md`.
- `docs/References/Meta/` holds the specifications: one directory per kind of file, each pairing a `Specs.md` (how the file must be shaped) with a `sample`. Never browse it — always locate a specification through `docs/References/Specs.md`.
