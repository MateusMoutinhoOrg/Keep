# LibFunctions Specification

## Description
Defines the required shape of a library function inside the sandbox — the factories filling `api.Lib`'s function fields in `sandbox/lib/publicfunctions/`, the `New` constructor that runs them in `sandbox/lib/new.go`, and the shared engine helpers in `sandbox/lib/dense/`. A library function is pure logic that reaches storage only through the injected `Deps`.

### Rules
- One factory per function field of `api.Lib`, named `<Field>Factory` and taking a single `*api.Lib` parameter: `func NewDatabaseFactory(l *api.Lib) func(props api.Props) api.KeepDatabase`. Its only job is to build and return the closure — see [Factories](/docs/References/Specs/Factories/Specs.md).
- Each factory lives in **its own file** under `sandbox/lib/publicfunctions/`, named after the field it fills (`NewDatabase.go` fills `NewDatabase`). Behavior belonging to a created object is a factory in that object's package instead — see [LibObjects](/docs/References/Specs/LibObjects/Specs.md).
- Every factory must be called from `New(d deps.Deps) api.Lib` in `sandbox/lib/new.go`, which assigns its return value into the matching field and doubles as the factory aggregate. A field whose factory's return value is never assigned stays nil and panics on first call; the compiler does not check this. `sandbox/new.go` does nothing but delegate to it.
- Compile-time constants the library reports — its version, and any other fixed text — live in `sandbox/config/` and are read from there, never hard-coded into a closure.
- Shared logic used by more than one object goes in `sandbox/lib/dense/` as an exported package-level function taking the deps it needs as its first parameter — a plain function, not a factory, since `dense` fills no `api` struct's fields. It must not reference any object type, so the dependency graph stays acyclic.
- Storage is touched **only** through the `Deps` the carrier struct holds (`l.Deps.<Field>()`, `s.Deps.<Field>()`) or the one passed in as a parameter. Never construct or import a concrete implementation. Reading `l.Deps` inside the closure rather than capturing it at factory time is what lets the injected value stay authoritative.
- No package under `sandbox/` may import `adapters/`, `examples/libraryExamples/`, `tests/`, a third-party module, or an OS-bound stdlib package — see [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- Every storage access must respect the invariants of the [Dense Record Pattern](/docs/References/DenseRecordPattern.md): single-key reads and writes only, no key listing, and the documented write orderings.
- Expected failures are returned as a typed `*api.Error`; storage failures are wrapped as `api.Internal` through `dense.InternalError`.
- A function field that is part of the public surface must be declared on an [`api`](/docs/References/Specs/Outputs/Specs.md) struct and listed in [PublicApi.md](/docs/References/PublicApi.md).

## Structure
1. **Package clause**: `package publicfunctions`, `package lib` (for `new.go`), or `package dense`.
2. **Doc comment**: one sentence describing what the factory or helper does, naming the field it fills.
3. **Field factory**: `func <Field>Factory(l *api.Lib) <FieldType>` returning a closure for `l.<Field>`, calling dependencies via `l.Deps.<Field>(...)` — or a `dense` package-level helper taking `d deps.Deps` first, reaching storage only through it and returning the composed result plus a `*api.Error` where it can fail.
4. **`New` constructor** (in `sandbox/lib/new.go`): `func New(d deps.Deps) api.Lib` building `api.Lib{Deps: d}`, calling every field factory in `publicfunctions` exactly once and assigning its return value into the matching field, and returning the struct.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
