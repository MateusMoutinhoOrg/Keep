# Outputs Specification

## Description
Defines the required shape of `sandbox/contracts/api/api.go` — the library's whole public surface. It is the file that guarantees the central rule of this project: **every type crossing the library boundary is declared here**, as a struct, never an interface — a behaviorless struct for plain data, a struct of function fields for anything with behavior. The factories filling those function fields live in `sandbox/lib/` and are unreachable from outside the sandbox.

### Rules
- `api.go` must declare **structs and constants only**. Every type in the project lives here; `sandbox/lib/` declares no types at all.
- A struct that carries behavior (`Lib`, `KeepDatabase`, `SchemaInstance`, `SchemaItem`) must lead with an exported `Deps deps.Deps` field, followed by any plain data fields, followed by its function fields. A struct with no behavior (`Props`, `Schema`, `Item`, `Error`) carries no `Deps` field and no function fields — it is buildable directly with a composite literal.
- Every function field must take and return only primitives (`string`, `int`, `int64`, `bool`, `any`, and slices or maps of those) or structs declared here. A field exposing an unlisted type breaks the rule even if that type is internal.
- One **struct** per object the library hands back, including the `Lib` entry point returned by `lib.New`.
- Data the caller passes *in* (e.g. `Props`, `Schema`, `Item`) is a plain struct too, built directly with a composite literal — there is no constructor package, because none of them are interfaces.
- Classification is expressed as **exported constants** (`Key`, `KeyConflict`, …) reported through a plain `Type int` field, not as a named type — an `int` is a primitive, a defined type is not.
- Every field must be **exported**: a struct with unexported fields cannot be filled by a factory in the `sandbox/lib/` packages.
- A lookup that can fail must return `(value, ok bool)`, never a typed nil — a struct has no nil form, unlike an interface. `bool` is `false` when nothing matches.
- `*Error` stays a **pointer**, so `err != nil` is still the success/failure check, even though `Error` itself carries no methods.
- `api.go` must not import anything except `sandbox/contracts/deps` (needed for the `Deps deps.Deps` field). It never imports `adapters/`, `examples/libraryExamples/`, `tests/`, `sandbox/lib/`, or `sandbox`.
- Exported structs and fields must have a doc comment and be listed in [PublicApi.md](/docs/References/PublicApi.md).

## Structure
1. **Package clause**: `package api`.
2. **Constant blocks**: the exported classification constants, each with a doc comment naming the field that reports it.
3. **Input structs**: the description the caller supplies, innermost first — plain data, no `Deps` field.
4. **Error struct**: plain data (`Type`, `Key`, `KeyValue`, `Message`); no methods, no embedded `error`.
5. **Output structs**: one per object the library hands back, innermost object first, each leading with `Deps deps.Deps`.
6. **`Lib` struct**: the entry point, leading with `Deps` and declaring the function fields the library exposes.

> **Note**: For a concrete example, refer to [sample.go](./sample.go).
