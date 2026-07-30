# Deps Specification

## Description
Defines the required shape of the dependency contract in `sandbox/contracts/deps/deps.go`. This spec describes **how the contract must be declared**, not which concrete dependencies a real library needs — the behavior each method must honor is documented in [RequiredApi.md](/docs/References/RequiredApi.md).

### Rules
- `deps.go` must declare a single `Deps` **interface** — the one contract every adapter satisfies.
- Every dependency must be a **method** on that interface, so adapters inject behavior rather than data.
- Method names must be descriptive and exported; an unexported method cannot be implemented by an adapter in another package.
- Expected conditions must be reported through exported sentinel errors declared alongside the interface, so the library can tell them apart from real failures with `errors.Is`.
- `deps.go` must not import anything from `adapters/`, `examples/`, `tests/`, `sandbox/internal/`, `sandbox/contracts/api`, or `sandbox` (the entry point) — the contract stays free of implementations.

## Structure
1. **Package clause**: `package deps`.
2. **Sentinel errors**: the exported `Err*` values every implementation must return for expected conditions.
3. **`Deps` interface**: a set of exported methods, one per requirement.
4. Each method is a `Name(...) ...` signature describing a single injectable behavior.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
