# LibObjects Specification

## Description
Defines the required shape of an object created by the library in `pkg/keep/` and `pkg/database/`. A library object carries the injected `Deps` — directly, or through a back-pointer to the object that owns them — so its own methods can reach storage.

### Rules
- The object type is declared in the file named after it (e.g. `pkg/database/schema_item.go`), holding either an exported `Deps deps.Deps` field or a private back-pointer to the object that carries them, plus its own properties.
- The object is created by a constructor method on its parent (e.g. `KeepLib.NewDatabase`, `KeepDatabase.GetSchema`, `SchemaInstance.NewItem`) that wires the parent's `Deps` into the new value — callers never build it directly.
- Methods reach storage only through those injected `Deps`; `pkg/` never imports `adapters/`.
- Fields that must not be mutated by consumers (back-pointers, ids, key prefixes) stay private and are exposed through methods.
- Exported objects, constructors, and methods must have doc comments and be listed in [PublicApi.md](/docs/Reference/PublicApi.md).

## Structure
1. **Type declaration**: a struct carrying the injected `Deps` (or a back-pointer to their owner) and its own properties.
2. **Constructor**: a method on the parent object that wires the parent's `Deps` into the returned object.
3. **Methods**: `func (o *<Object>) ...` operating on the object's properties and, when needed, on storage through the injected `Deps`.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
