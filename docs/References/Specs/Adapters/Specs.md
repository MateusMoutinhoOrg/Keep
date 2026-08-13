# Adapters Specification

## Description
Defines the required shape of an adapter in `adapters/<name>/<name>.go`. This spec describes **how an adapter must be built**, not which concrete dependencies it fills in — those come from the [Deps contract](/docs/References/Specs/Deps/Specs.md). The page listing the shipped adapters is governed by [AdaptersDoc](/docs/References/Specs/AdaptersDoc/Specs.md).

### Rules
- Each adapter lives in its own directory under `adapters/` and uses a package named after that directory.
- Each adapter declares a struct carrying an exported `Deps deps.Deps` field — the contract its factories fill — plus whatever state its closures read.
- Each function field of `deps.Deps` is filled by a `<Field>Factory(carrier *<Name>Adapter) <FieldType>` that returns one closure, reading the carrier's state at call time.
- Each adapter exposes a single `New(...) deps.Deps` factory aggregate as its entry point: it builds the carrier, assigns every factory's return value to the matching field, and returns the filled **contract struct** — never the concrete adapter type.
- An adapter may import `sandbox/contracts/deps` but must never import `sandbox`, `sandbox/contracts/api`, or `sandbox/lib/` — dependencies flow one way.
- The adapter sits **outside the sandbox** and is the only place OS-bound and third-party code is allowed: all concrete choices (stdlib, third-party libs, config) live here.
- Expected conditions must be reported with the sentinel errors of the `deps` package, wrapped so `errors.Is` still matches.
- A field no factory fills stays `nil` and panics on first call — nothing checks this for you, unlike a missing method on an interface. Every adapter must have a row in [Adapters.md](/docs/References/Adapters.md).

## Structure
1. **Package clause**: `package <name>`.
2. **Imports**: at least `sandbox/contracts/deps`, plus whatever the implementation needs.
3. **State struct**: an exported type leading with `Deps deps.Deps`, plus its backing store and configuration.
4. **One `<Field>Factory` per contract field**: `func <Field>Factory(carrier *<Name>Adapter) <FieldType>`, returning that field's closure.
5. **`New(...) deps.Deps` factory aggregate**: accepts adapter-specific configuration, builds the carrier, assigns every factory's return value to `carrier.Deps.<Field>`, and returns `carrier.Deps`.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
