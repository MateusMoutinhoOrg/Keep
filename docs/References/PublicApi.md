# Public API

## Description
Index of all public-facing components (interfaces, functions, and methods), with links to their respective detail files.

The library's whole surface is **interfaces and primitives**. `sandbox/contracts/api` declares every interface and constant, and holds no struct at all; the structs implementing them live in `sandbox/internal/` and are unreachable from outside the sandbox. Because the description types are interfaces too, they are built with the constructors the `sandbox` package exports rather than composite literals.

---

## Interfaces

### [api.Lib](./PublicApi/api.Lib.md)
The main library entry point. Obtained via `lib.New`; creates databases with the injected deps wired in.

### [api.KeepDatabase](./PublicApi/api.KeepDatabase.md)
A database bound to a `Deps` backend and a `Props` description.

### [api.SchemaInstance](./PublicApi/api.SchemaInstance.md)
One collection of records; the entry point for creating, finding, and listing records.

### [api.SchemaItem](./PublicApi/api.SchemaItem.md)
One record; reads, updates, removes, and manages sub-database records.

### [api.Props / api.Schema / api.Item](./PublicApi/api.Props.md)
The declarative description of a database: its key prefix and its collections with typed fields.

### [api.Error](./PublicApi/api.Error.md)
The typed error returned by database operations, and `nil` on success.

### [deps.Deps](./PublicApi/deps.Deps.md)
The contract of injectable storage operations every backend must implement.

---

## Functions

### [lib.New](./PublicApi/lib.New.md)
Injects a `Deps` implementation into the library and returns the `api.Lib` entry point.

### [lib.NewProps / lib.NewSchema / lib.NewKeyItem / lib.NewIntItem / lib.NewDatabaseItem](./PublicApi/api.Props.md#constructors)
The constructors building a database description, since `Props`, `Schema`, and `Item` are interfaces.

### [standard.New / standard.NewWithBase](./PublicApi/standard.New.md)
Creates a `deps.Deps` backed by the filesystem.

### [native.New](./PublicApi/native.New.md)
Creates a `deps.Deps` backed by process memory.

---

## Methods

### [api.Lib.NewDatabase](./PublicApi/api.Lib.md#methods)
Creates a `KeepDatabase` from a `Props` description with the lib's deps wired in.

### [api.KeepDatabase methods](./PublicApi/api.KeepDatabase.md#methods)
`GetSchema` and `Props`.

### [api.SchemaInstance methods](./PublicApi/api.SchemaInstance.md#methods)
`NewItem`, `FindByKey`, `ListAll`, and `List`.

### [api.SchemaItem methods](./PublicApi/api.SchemaItem.md#methods)
`Id`, `Get`, `Update`, `Remove`, `ListAll`, `NewSubItem`, `CheckKeysPresence`, and `String`.

### [api.Error methods](./PublicApi/api.Error.md)
`Error`, `Type`, `Key`, and `KeyValue`.
