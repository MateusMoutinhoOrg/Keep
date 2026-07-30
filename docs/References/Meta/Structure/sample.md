# Project Structure

This document maps the project **schema** — the kinds of files a project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/References/Meta/Structure/Specs.md) to get its description and sample. Use them to build a project from scratch, even before any `sandbox/` code exists.

The project is split into top-level trees, and the dependency flow between them is one-way:

```
adapters/  ──▶  sandbox/  ◀──  examples/
(reaches the OS)  (closed)     (wires the two together)
```

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and quick-start guide | Readme |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition | |

---

## `/docs/`
Documentation of the project.

### `/docs/References/`

| File | Description | Spec |
|------|-------------|------|
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `Adapters.md` | Lists every shipped adapter and when to use each one | AdaptersDoc |

#### `/docs/References/Meta/`
The specifications describing how each kind of file in the project must be shaped.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/` | One directory per specification, holding its `Specs.md` and `sample` | |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontractsdeps) contract, and the only place OS-bound and third-party code is allowed.

### `/adapters/<name>/`

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct implementing every `Deps` method, exposed by a `New(...) deps.Deps` factory | Adapters |

---

## `/sandbox/`
The closed sandbox — the pure library. Every effect it needs arrives through the injected `Deps`.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | The `New` constructor wiring `Deps` into the internal `Lib` | |

### `/sandbox/contracts/`
The contracts everything is wired through — the only part of the sandbox anything outside it may import.

| File | Description | Spec |
|------|-------------|------|
| `deps/deps.go` | The `Deps` interface, one method per injectable behavior | Deps |
| `api/api.go` | Every interface and constant the library exchanges — interfaces only | Outputs |

### `/sandbox/internal/`
The concrete structs implementing the `api` interfaces, kept unreachable from outside by Go's `internal/` rule.

| File | Description | Spec |
|------|-------------|------|
| `lib/lib.go` | `Lib`: holds `Deps` and constructs the lib's objects | LibFunctions |
| `<object>/<object>.go` | One package per object the lib creates, carrying the propagated `Deps` | LibObjects |
