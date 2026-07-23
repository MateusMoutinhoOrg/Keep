# Examples Specification

## Description
Defines the required shape of a runnable example in `examples/<example>/<example>.go`. An example is a self-contained `package main` program that wires an adapter into the lib to demonstrate real usage.

### Rules
- Each example lives in its own directory under `examples/` named after the feature it demonstrates (e.g. `examples/ExampleSample/`).
- The file is named after its directory (`<example>/<example>.go`) and declares `package main` with a `main` function.
- An example wires the two layers together: it constructs `deps.Deps` through an adapter's `New(...)` factory, then passes it to `lib.New(...)`.
- An example may import `adapters/<name>` and `pkg/lib`; it must never reconstruct dependencies by hand — that is the adapter's job.
- Keep examples minimal and runnable via `go run ./examples/<example>/<example>.go`; add explanatory comments on the key wiring steps.
- Adding, renaming, or deleting an example requires updating the Samples section of [README.md](/README.md) — see [AddSample.md](/docs/Tutorials/AddSample.md).

## Structure
1. **Package clause**: `package main`.
2. **Imports**: an adapter (e.g. `adapters/standard`) and `pkg/lib`.
3. **`main` function**: build deps via the adapter, inject them with `lib.New`, then exercise the library.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
