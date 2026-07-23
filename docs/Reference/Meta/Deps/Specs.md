# Deps Specification

## Description
Defines the required shape of the dependency contract in `pkg/deps/deps.go`. This spec describes **how the contract must be declared**, not which concrete dependencies a real library needs.

### Rules
- `deps.go` must declare a single `Deps` struct — the one contract every adapter satisfies.
- Every dependency must be a **function field** (not a plain value and not an interface), so adapters inject behavior rather than data.
- Field names must be descriptive and exported.
- `deps.go` must not import anything from `adapters/` or `pkg/lib/` — the contract stays free of implementations.

## Structure
1. **Package clause**: `package deps`.
2. **`Deps` struct**: a bag of exported function fields, one per requirement.
3. Each field is a `func(...) ...` signature describing a single injectable behavior.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).