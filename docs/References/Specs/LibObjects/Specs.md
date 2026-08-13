# LibObjects Specification

## Description
Defines the required shape of an object created by the library in `sandbox/lib/<object>/`. A library object's package holds the **factories** filling one of the [`api`](/docs/References/Specs/Outputs/Specs.md) structs' function fields, plus the constructor that runs them all — it declares no types of its own.

### Rules
- Each object's factories live in **their own package** under `sandbox/lib/`, named after the object itself (e.g. `schemaitem` for `SchemaItem`). The whole `sandbox/lib/` tree is private to the sandbox, so the package takes no `internal_` prefix.
- `sandbox/lib/` holds **only factories and constructors**. The object's type is declared in `sandbox/contracts/api`; a package under `sandbox/lib/` may import contracts, never the reverse.
- The object's struct — declared in `sandbox/contracts/api/api.go`, never here — must lead with an exported `Deps deps.Deps` field.
- Every function field must be filled by a `<Field>Factory(carrier *api.<Object>) <FieldType>` that returns one closure, reading `carrier.Deps` (and any other carrier field) at call time — never captured when the factory itself runs.
- A factory whose field returns another library object returns that object's **`api` struct**, never a pointer to it.
- A lookup that can fail must return `(value, ok bool)` from the closure — a struct has no nil form, unlike an interface, so there is no "return nil" to fall back on.
- The object is built by a constructor (e.g. `Lib.NewDatabase`'s factory, `database.New`, `schemainstance.New`) that copies the parent's `Deps` into the new struct before running every field factory over it — callers never build it directly.
- Factories reach storage only through the carrier's `Deps` field; no package under `sandbox/` may import `adapters/`, as required by [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- Exported factories and constructors must have doc comments naming the field they fill, and the struct they belong to must be listed in [PublicApi.md](/docs/References/PublicApi.md).

## Structure
1. **Package clause**: `package <object>`.
2. **One `<Field>Factory` per function field**: `func <Field>Factory(carrier *api.<Object>) <FieldType>`, returning the field's closure.
3. **Constructor**: `func New(d deps.Deps, ...) api.<Object>` (or a shared aggregate like `build`), building the struct with `Deps` set, assigning every factory's return value to its field, and returning the struct by value.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
