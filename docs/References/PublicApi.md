# Public API

## Description
Index of all public-facing components (structs, fields, and functions), with links to their respective detail files.

The library's whole surface is **structs of function fields and plain data structs, never behaviorless interfaces**. `sandbox/contracts/api` declares every type — a struct type that carries behavior leads with a `Deps` field and fills the rest as function fields, each assigned by a factory in `sandbox/lib/`; a struct type with no behavior (`Item`, `Schema`, `Props`, `Error`) is plain data. Because none of them are interfaces, every one is buildable directly with a composite literal — there is no separate constructor package. See [StructContracts.md](/docs/References/StructContracts.md).

---

## Structs

### [api.Lib](./PublicApi/api.Lib.md)
The main library entry point. Obtained via `lib.New`; creates databases with the injected deps wired in.

### [api.KeepDatabase](./PublicApi/api.KeepDatabase.md)
A database bound to a `Deps` backend and a `Props` description.

### [api.SchemaInstance](./PublicApi/api.SchemaInstance.md)
One collection of records; the entry point for creating, finding, and listing records.

### [api.SchemaItem](./PublicApi/api.SchemaItem.md)
One record; reads, updates, removes, and manages sub-database records.

### [api.Props / api.Schema / api.Item](./PublicApi/api.Props.md)
The declarative description of a database: its key prefix and its collections with typed fields. Plain data — build with a composite literal.

### [api.Error](./PublicApi/api.Error.md)
The typed failure returned by database operations, as `*api.Error` (`nil` on success). Plain data — no methods, read `Type`/`Key`/`KeyValue`/`Message` as fields.

### [deps.Deps](./PublicApi/deps.Deps.md)
The contract of injectable storage operations every backend must fill.

---

## Functions

### [lib.New](./PublicApi/lib.New.md)
Injects a filled `Deps` into the library and returns the `api.Lib` entry point.

### [standard.New / standard.NewWithBase](./PublicApi/standard.New.md)
Creates a `deps.Deps` backed by the filesystem.

### [native.New](./PublicApi/native.New.md)
Creates a `deps.Deps` backed by process memory.

---

## Fields

### [api.Lib.Version](./PublicApi/api.Lib.md#fields)
Returns the library's own release, held as a compile-time constant in `sandbox/config`.

### [api.Lib.NewDatabase](./PublicApi/api.Lib.md#fields)
Creates a `KeepDatabase` from a `Props` description with the lib's deps wired in.

### [api.KeepDatabase fields](./PublicApi/api.KeepDatabase.md#fields)
`Props` (plain data) and `GetSchema` (function field, `(SchemaInstance, bool)`).

### [api.SchemaInstance fields](./PublicApi/api.SchemaInstance.md#fields)
`NewItem`, `FindByKey`, `ListAll`, and `List`.

### [api.SchemaItem fields](./PublicApi/api.SchemaItem.md#fields)
`Id` (plain data), `Get`, `Update`, `Remove`, `ListAll`, `NewSubItem`, `CheckKeysPresence`, and `String`.

### [api.Error fields](./PublicApi/api.Error.md)
`Type`, `Key`, `KeyValue`, and `Message` — plain data, no methods.
