# LibraryExamples Specification

## Description
Defines the required shape of a runnable library example in `examples/libraryExamples/<example>/<example>.go`. A library example is a self-contained `package main` program that wires an adapter into the lib to demonstrate real usage from Go code.

### Rules
- Each example lives in its own directory under `examples/libraryExamples/` named after the feature it demonstrates, suffixed `Sample` (e.g. `examples/libraryExamples/CreateUserSample/`).
- The file is named after its directory (`<example>/<example>.go`) and declares `package main` with a `main` function.
- An example wires the two layers together: it builds a `deps.Deps` through an adapter's `New(...)` factory, then passes it to `lib.New(...)`, which returns an `api.Lib`.
- An example may import `adapters/<name>` (aliased `keepadapter`), `sandbox` (aliased `keeplib`), and `sandbox/contracts/api` (aliased `keeptypes`); it must never import `sandbox/lib/` — that tree is private to the sandbox — nor reconstruct dependencies by hand, which is the adapter's job.
- Examples live outside the sandbox and are the only place an adapter and the library are named in the same file.
- The database description is built as package-level `Schemas` and `Props` values, using composite literals directly — `Props`, `Schema`, and `Item` are plain structs with no behavior, so no constructor call is needed.
- Keep examples minimal and runnable via `go run ./examples/libraryExamples/<example>/<example>.go`; add explanatory comments on the key wiring steps.
- Adding, renaming, or deleting an example requires updating [ApiSamplesList.md](/docs/References/ApiSamplesList.md) — see [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md).

## Structure
1. **Package clause**: `package main`.
2. **Imports**: every import of this module is aliased with the `keep` prefix, so a reader sees at a glance which layer a call belongs to — the adapter as `keepadapter` (e.g. `keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"`), the sandbox entry point as `keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"`, and the output types as `keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"`.
3. **`Schemas` and `Props` package-level values**: the `keeptypes.Props` description, built directly with composite literals.
4. **`main` function**: build deps via `keepadapter.New(...)`, inject them with `keeplib.New`, then exercise the returned `keeptypes.Lib` against `Props`.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
