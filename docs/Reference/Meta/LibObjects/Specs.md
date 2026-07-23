# LibObjects Specification

## Description
Defines the required shape of an object created by the library in `pkg/lib/`. A library object carries a private copy of the injected `Deps` so its own methods can reach storage.

### Rules
- The object type is declared in [pkg/lib/types.go](../../../../pkg/lib/types.go), holding a private `deps deps.Deps` field plus its own properties. `Lib` itself is declared in [pkg/lib/new.go](../../../../pkg/lib/new.go).
- The object is created by a constructor method on its parent (e.g. `Lib.NewDatabase`, `KeepDatabase.GetSchema`, `SchemaInstance.NewItem`) that copies the parent's `deps` into the new value — callers never build it directly.
- Methods reach storage only through that private `deps` field; `pkg/lib/` never imports `adapters/`.
- Fields that must not be mutated by consumers (deps, ids, key prefixes) stay private and are exposed through methods.
- Exported objects, constructors, and methods must have doc comments and be listed in [PublicApi.md](/docs/Reference/PublicApi.md).

## Structure
1. **Type declaration** (in `types.go`): a struct with a private `deps deps.Deps` field and its own properties.
2. **Constructor**: a method on the parent object that copies the parent's `deps` into the returned object.
3. **Methods**: `func (o *<Object>) ...` operating on the object's properties and, when needed, on storage through its `deps`.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
