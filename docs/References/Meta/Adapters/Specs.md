# Adapters Specification

## Description
Defines the required shape of an adapter in `adapters/<name>/<name>.go`. This spec describes **how an adapter must be built**, not which concrete dependencies it fills in — those come from the [Deps contract](/docs/References/Meta/Deps/Specs.md). The page listing the shipped adapters is governed by [AdaptersDoc](/docs/References/Meta/AdaptersDoc/Specs.md).

### Rules
- Each adapter lives in its own directory under `adapters/` and uses a package named after that directory.
- Each adapter declares an **unexported struct** holding its state, and implements **every** method of the `Deps` interface on it. A partial adapter does not satisfy the contract and fails to compile at the factory.
- Each adapter exposes a single `New(...) deps.Deps` factory as its entry point, returning the **interface** rather than the concrete struct.
- An adapter may import `sandbox/contracts/deps` but must never import `sandbox`, `sandbox/contracts/api`, or `sandbox/internal/` — dependencies flow one way.
- The adapter sits **outside the sandbox** and is the only place OS-bound and third-party code is allowed: all concrete choices (stdlib, third-party libs, config) live here.
- Expected conditions must be reported with the sentinel errors of the `deps` package, wrapped so `errors.Is` still matches.
- Every adapter must have a row in [Adapters.md](/docs/References/Adapters.md).

## Structure
1. **Package clause**: `package <name>`.
2. **Imports**: at least `sandbox/contracts/deps`, plus whatever the implementation needs.
3. **State struct**: an unexported type holding the adapter's backing store and configuration.
4. **`New(...) deps.Deps` factory**: accepts adapter-specific configuration and returns the state struct as a `deps.Deps`.
5. **Method set**: one method per contract requirement, implemented on the state struct.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
