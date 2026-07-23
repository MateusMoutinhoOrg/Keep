# Project Structure

This document maps the project **schema** — the kinds of files a project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/Reference/Meta/Structure/Specs.md) to get its description and sample. Use them to build a project from scratch, even before any `pkg/` code exists.

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and quick-start guide | Readme |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition | |

---

## `/docs/`
Documentation of the project.

### `/docs/Reference/`

| File | Description | Spec |
|------|-------------|------|
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |

#### `/docs/Reference/Meta/`
The specifications describing how each kind of file in the project must be shaped.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/` | One directory per specification, holding its `Specs.md` and `sample` | |

---

## `/adapters/`
Opinionated implementations of the [`Deps`](#pkg) contract.

### `/adapters/<name>/`

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A `New(...) deps.Deps` factory implementing every `Deps` field | Adapters |

---

## `/pkg/`
Core dependency-injection system and public library logic.

| File | Description | Spec |
|------|-------------|------|
| `deps/deps.go` | The `Deps` contract of injectable function fields | Deps |
| `lib/*.go` | A library function reaching deps via `l.deps.<Field>()` | LibFunctions |
