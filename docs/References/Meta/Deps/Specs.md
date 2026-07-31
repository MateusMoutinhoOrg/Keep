# Deps Specification

## Description
Defines the required shape of the dependency contract in `sandbox/contracts/deps/deps.go`. This spec describes **how the contract must be declared**, not which concrete dependencies a real library needs — the behavior each field must honor is documented in [RequiredApi.md](/docs/References/RequiredApi.md).

### Rules
- `deps.go` must declare a single `Deps` **struct of function fields** — the one contract every adapter fills. It is never an interface.
- Every dependency must be a **function field** on that struct, so adapters inject behavior by assigning a closure, not by implementing a method set.
- Field names must be descriptive and exported; an unexported field cannot be assigned by an adapter's factory in another package.
- Expected conditions must be reported through exported sentinel errors declared alongside the struct, so the library can tell them apart from real failures with `errors.Is`.
- `deps.go` must not import anything from `adapters/`, `examples/`, `tests/`, `sandbox/internal/`, `sandbox/contracts/api`, or `sandbox` (the entry point) — the contract stays free of implementations.

## Structure
1. **Package clause**: `package deps`.
2. **Sentinel errors**: the exported `Err*` values every adapter must return for expected conditions.
3. **`Deps` struct**: a set of exported function fields, one per requirement.
4. Each field is a `Name func(...) ...` signature describing a single injectable behavior.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
