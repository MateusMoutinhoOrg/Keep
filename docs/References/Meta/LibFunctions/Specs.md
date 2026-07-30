# LibFunctions Specification

## Description
Defines the required shape of a library function inside the sandbox — the methods of `Lib` in `sandbox/internal/lib/` and the shared engine helpers in `sandbox/internal/dense/`. A library function is pure logic that reaches storage only through the injected `Deps`.

### Rules
- Entry-point behavior hangs off `Lib`: `func (l *Lib) ...` in `sandbox/internal/lib/`. Behavior belonging to a created object is a method on that object instead — see [LibObjects](/docs/References/Meta/LibObjects/Specs.md).
- Shared logic used by more than one object goes in `sandbox/internal/dense/` as an exported package-level function taking the deps it needs as its first parameter. It must not reference any object type, so the dependency graph stays acyclic.
- Storage is touched **only** through the `Deps` the object carries (`l.Deps.<Method>()`, `s.Deps.<Method>()`) or the one passed in as a parameter. Never construct or import a concrete implementation.
- No package under `sandbox/` may import `adapters/`, `examples/`, `tests/`, a third-party module, or an OS-bound stdlib package — see [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- Every storage access must respect the invariants of the [Dense Record Pattern](/docs/Explanations/DenseRecordPattern.md): single-key reads and writes only, no key listing, and the documented write orderings.
- Expected failures are returned as a typed `api.Error`; storage failures are wrapped as `api.Internal` through `dense.InternalError`.
- A function that is part of the public surface must be declared on an [`api`](/docs/References/Meta/Outputs/Specs.md) interface and listed in [PublicApi.md](/docs/References/PublicApi.md).

## Structure
1. **Package clause**: `package lib` or `package dense`.
2. **Doc comment**: one sentence describing what the function does.
3. **The function**: a method on `Lib`, or a package-level helper taking `d deps.Deps` first, reaching storage only through it and returning the composed result plus a `api.Error` where it can fail.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
