# Adapters Specification

## Description
Defines the required shape of an adapter in `adapters/<name>/<name>.go`. This spec describes **how an adapter must be built**, not which concrete dependencies it fills in — those come from the [Deps contract](/docs/Reference/Meta/Deps/Specs.md).

### Rules
- Each adapter lives in its own directory under `adapters/` and uses a package named after that directory.
- Each adapter exposes a single `New(...) deps.Deps` factory as its entry point.
- `New` must return a fully-populated `deps.Deps` — **every** field of the `Deps` struct must be assigned. A partial adapter breaks consumers.
- An adapter may import `pkg/deps` but must never import `pkg/lib` — dependencies flow one way.
- The adapter is the **opinionated** layer: all concrete choices (stdlib, third-party libs, config) live here.

## Structure
1. **Package clause**: `package <name>`.
2. **Imports**: at least `pkg/deps`, plus whatever the implementation needs.
3. **`New(...) deps.Deps` factory**: accepts adapter-specific configuration and returns a `deps.Deps` with every field implemented as a closure.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).