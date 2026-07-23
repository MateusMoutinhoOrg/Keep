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

The DI flow is a strict one-way dependency chain — the core never imports concrete implementations:

```
adapters/standard  ──constructs──▶  deps.Deps  ──injected into──▶  lib.New(deps)
  (filesystem impl)                (the contract)                  (pure logic)
```

- **`pkg/deps/deps.go`** — the `Deps` struct: a bag of injectable **function fields** (not an interface), plus the sentinel errors adapters must return. Adding a requirement = adding a field here. Its behavioral contract is `docs/Reference/RequiredApi.md`.
- **`pkg/lib/`** — the whole library, one package: `new.go` declares `Lib` (the entry point, holding a private `deps`) and `types.go` declares every object it creates — `KeepDatabase` → `SchemaInstance` (a collection) → `SchemaItem` (a record). Each object carries its own private copy of the injected `deps`. Never imports `adapters/`.
- **`adapters/<name>/`** — opinionated concrete backends. Each exposes a `New(...) deps.Deps` factory filling in **every** `Deps` field. `standard` is filesystem-backed; `native` is in-memory.
- **`examples/<name>/<name>.go`** — self-contained `package main` programs wiring an adapter into the lib.

Storage layout follows the Dense Record Pattern (`docs/Explanation/DenseRecordPattern.md`): fixed-cost operations, ids that are never reused, and a dense position list that makes iteration possible without key listing. Changes to `pkg/lib/dense.go` must preserve those invariants and the documented write orderings.

## Critical: this repo is documentation-driven

Changes are governed by required-reading docs, and several actions **must** update companion files in the same commit. Each tutorial in `docs/Tutorials/` covers exactly one goal — read the one matching your change:

| If you... | Read | And keep in sync |
|-----------|------|------------------|
| add/rename/delete any file or dir | `docs/Reference/Structure.md` | `docs/Reference/Structure.md` |
| add/rename/delete a `.md` file | `docs/Tutorials/AddDocument.md`, `RenameDocument.md`, `DeleteDocument.md` | Doc Index in `README.md` |
| add a lib function/object | `docs/Tutorials/AddLibFunction.md`, `AddLibObject.md`, `AddDatabaseOperation.md` | `docs/Reference/PublicApi.md` (+ detail page in `docs/Reference/PublicApi/`, see `ExposePublicApi.md`) |
| add a `Deps` field | `docs/Tutorials/AddDependency.md` | **every** adapter in `adapters/`, plus `docs/Reference/RequiredApi.md` |
| add an adapter | `docs/Tutorials/AddAdapter.md` | `docs/Reference/Structure.md` |
| add/rename/delete a sample | `docs/Tutorials/AddSample.md` | Samples section in `README.md` |
| reuse this structure for another library | `docs/Tutorials/ForkTemplate.md`, `AdaptExistingLib.md` | `docs/Reference/TemplateFileActions.md` (the per-file copy/create/rewrite/delete list both tutorials follow) |

`docs/Reference/RULES.md` is the binding rule set and `docs/Reference/Specs.md` is the index of every file specification; `AGENTS.md` points here. Adding a `Deps` field without updating all adapters breaks every consumer — that's the most common footgun.

## Conventions

- Module path is `github.com/MateusMoutinhoOrg/Keep`; renaming it is a documented procedure — see `docs/Tutorials/RenameModule.md`.
- Public-facing API entries each get a detail page under `docs/Reference/PublicApi/` named `<pkg>.<Symbol>.md`.
- Operations return a typed `*lib.Error`; switch on its `Type` rather than matching strings — see `docs/Reference/Errors.md`.
- `docs/Reference/Meta/` holds the specifications: one directory per kind of file, each pairing a `Specs.md` (how the file must be shaped) with a `sample`. Never browse it — always locate a specification through `docs/Reference/Specs.md`.
