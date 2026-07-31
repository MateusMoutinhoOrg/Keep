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

Four top-level trees, wired through **structs of function fields** (never interfaces), with a strict one-way dependency flow:

```
adapters/  ──▶  sandbox/  ◀──  examples/ , tests/
(reaches the OS)  (closed)     (wire the two together)

standard.New()  ──▶  deps.Deps  ──▶  lib.New(deps)  ──▶  api.Lib  ──▶  api.KeepDatabase  ──▶  api.SchemaInstance  ──▶  api.SchemaItem
(opinionated impl)   (contract)      (entry point)       (output structs, filled by sandbox/internal/ factories)
```

Contracts are structs whose fields hold functions, and **every** one of them is filled by **factories** — `func <Field>Factory(carrier *T) <FieldType>` bodies that return one closure reading the carrier at call time, with the assignment made explicitly by the caller. Inside the sandbox the carrier is the `api` struct, which carries its own `Deps` field, and each object's `New` assigns the result (`s.Get = GetFactory(&s)`, reading `s.Deps` inside the closure); inside `adapters/` the carrier is the adapter struct, which declares a `Deps deps.Deps` field its `New` assigns into from each factory's return value (`s.Deps.Read = ReadFactory(s)`). No methods bound into fields, no internal mirror type. This is a binding rule — see `docs/References/RULES.md#factory-pattern` and `docs/Explanations/StructContracts.md`.

Two trade-offs, neither caught by the compiler: **completeness is unchecked** — a field no factory fills is nil and panics on first call, so every factory must be called from its package's `New` constructor; and **`Deps` is read-only after construction** — the closures captured the struct the factories ran over, so patch `deps.Deps` before calling `lib.New`, never on the returned struct.

`Props`, `Schema`, and `Item` carry no behavior and no `Deps` — they are plain data, buildable directly with a composite literal (see the Conventions section below for the shape). `Error` is the same: a plain data struct returned as `*api.Error`, nil on success.

`sandbox/` is a **closed sandbox**: nothing in it may import `adapters/`, `examples/`, `tests/`, a third-party module, or an OS-bound stdlib package (`os`, `net`, `syscall`, …). Every such effect is a `Deps` field reached through the object's `Deps` field. This is a binding rule — see `docs/References/RULES.md` and `docs/Explanations/SandboxIsolation.md`.

- **`sandbox/new.go`** — package `lib`, the only wiring point consumers touch: `New(deps.Deps) api.Lib`. Never imports `adapters/`. Importers alias it: `keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"`.
- **`sandbox/contracts/deps/deps.go`** — the `Deps` **struct of function fields** plus the sentinel errors adapters must return. Adding a requirement = adding a field here. Its behavioral contract is `docs/References/RequiredApi.md`.
- **`sandbox/contracts/api/api.go`** — the library's whole public surface, and it holds **nothing but types**: the objects handed back (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`), the description passed in (`Props`, `Schema`, `Item`), and `Error`. Every type in the project is declared here — `sandbox/internal/` declares no types at all, only factories. Every field must be exported, or `sandbox/internal/` cannot fill it.
- **`sandbox/internal/liberror/`** — plain constructors for `*api.Error` (not factories — `Error` has no `Deps` and nothing to fill after construction).
- **`sandbox/internal/lib/`** — one `<Field>Factory(l *api.Lib)` per function field of `api.Lib`, plus `New(d deps.Deps) api.Lib` assigning every factory's return value into the matching field. Go's `internal/` rule keeps it unreachable from `adapters/`, `examples/`, and consumers.
- **`sandbox/internal/<object>/`** — one package per object (`database`, `schemainstance`, `schemaitem`), holding that object's `<Field>Factory` functions plus a `New(...)` constructor that runs them all against the `api.<Object>` it builds. No `internal_` prefix — the `internal/` parent already says it.
- **`sandbox/internal/dense/`** — the shared key-layout and encoding helpers. It references **no object type**, which is what keeps the package graph acyclic; the record constructors (`New`, `ListRange`, `ClearCollection`) live in `schemaitem` for the same reason.
- **`adapters/<name>/`** — outside the sandbox; the only place OS-bound and third-party code is allowed. Each declares an exported struct carrying a `Deps deps.Deps` field, one `<Field>Factory(a *<Name>Adapter)` per `Deps` field, and a `New(...) deps.Deps` constructor assigning each factory's return value into `a.Deps` and returning it — the populated **contract struct**, never the adapter type. `standard` is filesystem-backed; `native` is in-memory. Listed in `docs/References/Adapters.md`.
- **`examples/<name>/<name>.go`**, **`tests/`** — outside the sandbox; the only places an adapter and the sandbox meet.

Every object propagates `Deps` to the objects it creates — see `database.GetSchemaFactory` for the pattern. Lookups that can fail (`KeepDatabase.GetSchema`, `SchemaInstance.FindByKey`) return `(value, ok bool)` instead of a nilable object, since `SchemaInstance`/`SchemaItem` are structs, not interfaces, and have no nil form.

Storage layout follows the Dense Record Pattern (`docs/Explanations/DenseRecordPattern.md`): fixed-cost operations, ids that are never reused, and a dense position list that makes iteration possible without key listing. Changes to `sandbox/internal/dense/` and `sandbox/internal/schemaitem/` must preserve those invariants and the documented write orderings.

## Critical: this repo is documentation-driven

Changes are governed by required-reading docs, and several actions **must** update companion files in the same commit. Each tutorial in `docs/Tutorials/` covers exactly one goal — read the one matching your change:

| If you... | Read | And keep in sync |
|-----------|------|------------------|
| add/rename/delete any file or dir | `docs/References/Structure.md` | `docs/References/Structure.md` |
| add/rename/delete a `.md` file | `docs/Tutorials/AddDocument.md`, `RenameDocument.md`, `DeleteDocument.md` | Doc Index in `README.md` |
| add a lib function/object | `docs/Tutorials/AddLibFunction.md`, `AddLibObject.md`, `AddDatabaseOperation.md` | `sandbox/contracts/api/api.go` + `docs/References/PublicApi.md` (+ detail page in `docs/References/PublicApi/`, see `ExposePublicApi.md`) |
| add a `Deps` field | `docs/Tutorials/AddDependency.md` | **every** adapter in `adapters/`, plus `docs/References/RequiredApi.md` |
| add an adapter | `docs/Tutorials/AddAdapter.md` | `docs/References/Structure.md`, `docs/References/Adapters.md` |
| need an OS/third-party call inside `sandbox/` | `docs/Explanations/SandboxIsolation.md`, `docs/Tutorials/AddDependency.md` | `sandbox/contracts/deps/deps.go` + **every** adapter |
| write or edit any `<Field>Factory` (sandbox or adapter) | `docs/Explanations/StructContracts.md` | the `New` constructor that must call it |
| add/rename/delete a sample | `docs/Tutorials/AddSample.md` | Samples section in `README.md` |
| reuse this structure for another library | `docs/Tutorials/ForkTemplate.md`, `AdaptExistingLib.md` | `docs/References/TemplateFileActions.md` (the per-file copy/create/rewrite/delete list both tutorials follow) |

`docs/References/RULES.md` is the binding rule set and `docs/References/Specs.md` is the index of every file specification; `AGENTS.md` points here. Adding a `Deps` field without implementing it in all adapters breaks every consumer — that's the most common footgun. A new public lib function or object must be declared in `sandbox/contracts/api/api.go` **and** implemented (as factories) in `sandbox/internal/`, or callers cannot reach it.

## Conventions

- Module path is `github.com/MateusMoutinhoOrg/Keep`; renaming it is a documented procedure — see `docs/Tutorials/RenameModule.md`.
- Public-facing API entries each get a detail page under `docs/References/PublicApi/` named `<pkg>.<Symbol>.md`.
- `Item`, `Schema`, and `Props` are plain structs built with a composite literal — no constructor required:
  ```go
  var Schemas = []database.Schema{
      {
          Name: "User",
          Itens: []database.Item{
              {Type: database.Key, Required: true, Name: "Email"},
              {Type: database.Key, Required: true, Name: "UserName"},
              {Name: "Age", Required: true, Type: database.Int},
          },
      },
  }

  var Props = database.Props{
      Path:    "testDatabase/",
      Schemas: Schemas,
  }
  ```
  where `database` aliases `github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api` — see the import-alias table below.
- Operations return a typed `*api.Error` (nil on success); switch on its `Type` field rather than matching strings — see `docs/References/Errors.md`.
- Import aliases: consumers (`examples/`, `tests/`, third-party code) alias the library's packages so each call site says which layer it belongs to — `adapters/<name>` as `keepadapter`, `sandbox` as `keeplib`, and `sandbox/contracts/api` as `database` (the domain name reads better here than a generic `keeptypes`). Files inside the library itself (`sandbox/` and its own `adapters/`) keep the plain package names.
- `docs/References/Meta/` holds the specifications: one directory per kind of file, each pairing a `Specs.md` (how the file must be shaped) with a `sample`. Never browse it — always locate a specification through `docs/References/Specs.md`.
