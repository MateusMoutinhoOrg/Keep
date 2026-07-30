# Outputs Specification

## Description
Defines the required shape of `sandbox/contracts/api/api.go` — the library's whole public surface. It is the file that guarantees the central rule of this project: **every value crossing the library boundary is either a primitive or an interface**, never a concrete type of the library. The structs implementing these interfaces live in `sandbox/internal/` and are unreachable from outside the sandbox.

### Rules
- `api.go` must declare **interfaces and constants only**. A `struct` in this file is a specification violation, no matter how plain the data it holds.
- Every method of every interface must take and return only primitives (`string`, `int`, `int64`, `bool`, `any`, and slices or maps of those) or interfaces declared here. A method exposing a struct breaks the rule even if the struct is internal.
- One **interface** per object the library hands back, including the `Lib` entry point returned by `lib.New`.
- Data the caller passes *in* is also an interface (e.g. `Props`, `Schema`, `Item`). Because an interface cannot be built with a composite literal, each one needs a matching constructor in `sandbox/description.go` returning that interface.
- Classification is expressed as **exported constants** (`KeyItem`, `KeyConflict`, …) reported by an `int`-returning method, not as a named type — an `int` is a primitive, a defined type is not.
- Every method must be **exported**: an interface with unexported methods cannot be implemented by the `sandbox/internal/` packages.
- A method returning "nothing found" or "no failure" must be documented as returning `nil`, and its implementation must return a **literal nil**, never a typed nil pointer.
- The error interface must embed the standard `error`, so callers can treat it as one.
- `api.go` must not import anything — not `adapters/`, `examples/`, `tests/`, `sandbox/internal/`, `sandbox`, nor any other package. It sits at the bottom of the dependency graph.
- Exported interfaces must have a doc comment and be listed in [PublicApi.md](/docs/References/PublicApi.md).

## Structure
1. **Package clause**: `package api`.
2. **Constant blocks**: the exported classification constants, each with a doc comment naming the method that reports it.
3. **Input interfaces**: the description the caller supplies, innermost first.
4. **Error interface**: embeds `error` and exposes the failure's cause and subject.
5. **Output interfaces**: one per object the library hands back, innermost object first.
6. **`Lib` interface**: the entry point, declaring the constructors and functions the library exposes.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
