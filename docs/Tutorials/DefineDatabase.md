# Define a Database

## Description
Covers describing a database — its key prefix and its collections — and opening it in a program. To add a field to a collection that already exists, follow [AddSchemaField.md](/docs/Tutorials/AddSchemaField.md); to nest a collection inside a record, follow [AddNestedCollection.md](/docs/Tutorials/AddNestedCollection.md). The background on field types is in [Schemas.md](/docs/Explanations/Schemas.md).

### Rules
- A database is a value, not a migration: the `Props` description is the only source of truth and is passed on every run.
- `Path` is prefixed to every stored key, so two databases sharing a backend must not share a `Path`.
- Fields of type `database.Key` are unique and indexed case-insensitively; every collection meant to be looked up needs at least one.

---

## Workflow
1. Declare a package-level `Schemas` value listing every collection and its fields. `Props`, `Schema`, and `Item` are plain structs, so they are built directly with composite literals:
   ```go
   var Schemas = []database.Schema{
       {
           Name: "user",
           Itens: []database.Item{
               {Name: "email", Type: database.Key, Required: true},
               {Name: "username", Type: database.Key, Required: true},
               {Name: "age", Type: database.Int, Required: true},
           },
       },
   }
   ```
2. Wrap the schemas in a `Props` value, choosing the prefix every key of this database is written under:
   ```go
   var Props = database.Props{
       Path:    "myDatabase/",
       Schemas: Schemas,
   }
   ```
3. Build the dependencies with an adapter and inject them into the lib, following [LibInitialization.md](/docs/Tutorials/LibInitialization.md):
   ```go
   keep := keeplib.New(keepadapter.New())
   ```
4. Open the database and take the collection to operate on:
   ```go
   db := keep.NewDatabase(Props)
   users, ok := db.GetSchema("user")
   if !ok {
       panic("schema not declared in Props")
   }
   ```
5. Operate on the collection — see [Records.md](/docs/Explanations/Records.md) for the available operations, and react to failures with [Errors.md](/docs/References/Errors.md).
6. Run the program:
   ```bash
   go run main.go
   ```
