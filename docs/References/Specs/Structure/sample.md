# Project Structure

This document maps the project **schema** — the kinds of files a project is built from — not every concrete file. A slot with a **Spec** name is governed by a specification; resolve the name through [Specs.md](/docs/References/Specs.md) to get its description and sample. Use them to build a project from scratch, even before any `sandbox/` code exists.

The project is split into top-level trees, and the dependency flow between them is one-way:

```
adapters/  ──▶  sandbox/  ◀──  examples/libraryExamples/
(reaches the OS)  (closed)     (wires the two together)
```

## Root

| File | Description | Spec |
|------|-------------|------|
| `README.md` | Project overview and the Doc Index pointing at each theme index | Readme |
| `LICENSE` | License terms for the project | |
| `go.mod` | Go module definition | |

---

## `/docs/`
Documentation of the project, split by kind of page: `Index/` holds one entry point per theme, `Tutorials/` holds every workflow, `References/` holds every lookup and explanation.

### `/docs/Index/`

| File | Description | Spec |
|------|-------------|------|
| `<Theme>.md` | The theme's entry point: its Tutorials and its References | Index |

### `/docs/References/`

| File | Description | Spec |
|------|-------------|------|
| `Structure.md` | The project's schema and the purpose of each component | Structure |
| `Specs.md` | Index of every specification and the files each one governs | |
| `Adapters.md` | Lists every shipped adapter and when to use each one | AdaptersDoc |

#### `/docs/References/Specs/`
The specifications describing how each kind of file in the project must be shaped.

| File | Description | Spec |
|------|-------------|------|
| `<Spec>/` | One directory per specification, holding its `Specs.md` and `sample` | |

---

## `/adapters/`
Outside the sandbox. Opinionated implementations of the [`Deps`](#sandboxcontracts) contract, and the only place OS-bound and third-party code is allowed.

### `/adapters/<name>/`

| File | Description | Spec |
|------|-------------|------|
| `<name>.go` | A struct carrying a `Deps` field, one `<Field>Factory` per field, and the `New(...) deps.Deps` constructor that assigns every return value | Adapters |

---

## `/sandbox/`
The closed sandbox — the pure library. Every effect it needs arrives through the injected `Deps`.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | The `New` constructor delegating to the internal lib constructor | |

### `/sandbox/contracts/`
The contracts everything is wired through — the only part of the sandbox anything outside it may import. Contracts are structs of function fields, never interfaces.

| File | Description | Spec |
|------|-------------|------|
| `deps/deps.go` | The `Deps` struct, one function field per injectable behavior | Deps |
| `api/api.go` | Every struct and constant the library exchanges — structs only | Outputs |

### `/sandbox/config/`
Compile-time constants the library reports.

| File | Description | Spec |
|------|-------------|------|
| `version.go` | The `Version` constant reported by the `api.Lib.Version` field | |

### `/sandbox/lib/`
The factories filling the `api` structs' function fields, plus the constructors that run them. It declares no types; the whole tree is private to the sandbox.

| File | Description | Spec |
|------|-------------|------|
| `new.go` | `New(d deps.Deps) api.Lib`, assigning every lib factory's return value into its field | LibFunctions |
| `publicfunctions/<Function>.go` | One file per function field of `api.Lib`, holding its `<Field>Factory` | LibFunctions |
| `<object>/<object>.go` | One package per object the lib creates, carrying the propagated `Deps` | LibObjects |
