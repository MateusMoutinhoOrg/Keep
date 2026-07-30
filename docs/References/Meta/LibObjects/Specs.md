# LibObjects Specification

## Description
Defines the required shape of an object created by the library in `sandbox/internal/<object>/`. A library object is the concrete struct implementing one of the [`api`](/docs/References/Meta/Outputs/Specs.md) interfaces, and it carries the injected `Deps` so its own methods can reach storage.

### Rules
- Each object lives in **its own package** under `sandbox/internal/`, named after the object itself (e.g. `schemaitem` for `SchemaItem`). The `internal/` parent already marks it private, so the package takes no `internal_` prefix.
- The struct must declare an exported `Deps deps.Deps` field. Its fields are exported because sibling internal packages construct it; Go's `internal/` rule is what keeps it unreachable from consumers, so nothing needs to be hidden by naming.
- Every method of the matching `api` interface must be implemented, with **identical signatures**. A method returning another library object returns that object's **`api` interface**, never the concrete struct.
- When a method returns "nothing found", it must return a **literal `nil`**, never a typed nil pointer — a typed nil wrapped in an interface is not `== nil` and would break every caller's check.
- The object is created by a constructor on its parent (e.g. `Lib.NewDatabase`, `KeepDatabase.GetSchema`, `SchemaInstance.NewItem`) that copies the parent's `Deps` into the new value — callers never build it directly.
- Methods reach storage only through that `Deps` field; no package under `sandbox/` may import `adapters/`, as required by [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- A field name that would collide with a method of the interface must be renamed rather than the method (e.g. field `RecordID` for method `Id()`).
- Exported objects, constructors, and methods must have doc comments, and the interface they implement must be listed in [PublicApi.md](/docs/References/PublicApi.md).

## Structure
1. **Package clause**: `package <object>`.
2. **Type declaration**: a struct with an exported `Deps deps.Deps` field plus its own properties, and a doc comment naming the `api` interface it implements.
3. **Methods**: `func (o *<Object>) ...` implementing the interface, operating on the object's properties and, when needed, on storage through its `Deps`.
4. **Constructors** (when the object creates its own kind): package-level functions returning the `api` interface, taking the deps to propagate as their first parameter.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
