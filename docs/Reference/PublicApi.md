# Public API

## Description
Index of all public-facing components (structs, functions, and methods), with links to their respective detail files.

---

## Structs

### [keep.KeepLib](./PublicApi/keep.KeepLib.md)
The main library entry point. Constructed via `keep.New`; creates databases with the injected deps wired in.

### [database.Props / database.Schema / database.Item](./PublicApi/database.Props.md)
The declarative description of a database: its key prefix and its collections with typed fields.

### [database.KeepDatabase](./PublicApi/database.KeepDatabase.md)
A database bound to a `Deps` backend and a `Props` description.

### [database.SchemaInstance](./PublicApi/database.SchemaInstance.md)
One collection of records; the entry point for creating, finding, and listing records.

### [database.SchemaItem](./PublicApi/database.SchemaItem.md)
One record; reads, updates, removes, and manages sub-database records.

### [database.Error](./PublicApi/database.Error.md)
The typed error returned by database operations.

### [deps.Deps](./PublicApi/deps.Deps.md)
The struct of injectable storage functions every backend must populate.

---

## Functions

### [keep.New](./PublicApi/keep.New.md)
Initializes and returns a new `KeepLib` instance configured with the provided dependency adapter.

### [standard.New / standard.NewWithBase](./PublicApi/standard.New.md)
Creates a `deps.Deps` backed by the filesystem.

### [native.New](./PublicApi/native.New.md)
Creates a `deps.Deps` backed by process memory.

---

## Methods

### [keep.KeepLib.NewDatabase](./PublicApi/keep.KeepLib.md#methods)
Creates a `KeepDatabase` from a `Props` description with the lib's deps wired in.

### [database.KeepDatabase.GetSchema](./PublicApi/database.KeepDatabase.md#methods)
Returns the collection with the given name, or `nil` if no schema has that name.

### [database.SchemaInstance methods](./PublicApi/database.SchemaInstance.md#methods)
`NewItem`, `FindByKey`, `ListAll`, and `List`.

### [database.SchemaItem methods](./PublicApi/database.SchemaItem.md#methods)
`Id`, `Get`, `Update`, `Remove`, `ListAll`, `NewSubItem`, `CheckKeysPresence`, and `String`.
