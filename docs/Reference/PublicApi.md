# Public API

## Description
Index of all public-facing components (structs, functions, and methods), with links to their respective detail files.

---

## Structs

### [lib.Lib](./PublicApi/lib.Lib.md)
The main library entry point. Constructed via `lib.New`; creates databases with the injected deps wired in.

### [lib.Props / lib.Schema / lib.Item](./PublicApi/lib.Props.md)
The declarative description of a database: its key prefix and its collections with typed fields.

### [lib.KeepDatabase](./PublicApi/lib.KeepDatabase.md)
A database bound to a `Deps` backend and a `Props` description.

### [lib.SchemaInstance](./PublicApi/lib.SchemaInstance.md)
One collection of records; the entry point for creating, finding, and listing records.

### [lib.SchemaItem](./PublicApi/lib.SchemaItem.md)
One record; reads, updates, removes, and manages sub-database records.

### [lib.Error](./PublicApi/lib.Error.md)
The typed error returned by database operations.

### [deps.Deps](./PublicApi/deps.Deps.md)
The struct of injectable storage functions every backend must populate.

---

## Functions

### [lib.New](./PublicApi/lib.New.md)
Initializes and returns a new `Lib` instance configured with the provided dependency adapter.

### [standard.New / standard.NewWithBase](./PublicApi/standard.New.md)
Creates a `deps.Deps` backed by the filesystem.

### [native.New](./PublicApi/native.New.md)
Creates a `deps.Deps` backed by process memory.

---

## Methods

### [lib.Lib.NewDatabase](./PublicApi/lib.Lib.md#methods)
Creates a `KeepDatabase` from a `Props` description with the lib's deps wired in.

### [lib.KeepDatabase.GetSchema](./PublicApi/lib.KeepDatabase.md#methods)
Returns the collection with the given name, or `nil` if no schema has that name.

### [lib.SchemaInstance methods](./PublicApi/lib.SchemaInstance.md#methods)
`NewItem`, `FindByKey`, `ListAll`, and `List`.

### [lib.SchemaItem methods](./PublicApi/lib.SchemaItem.md#methods)
`Id`, `Get`, `Update`, `Remove`, `ListAll`, `NewSubItem`, `CheckKeysPresence`, and `String`.
