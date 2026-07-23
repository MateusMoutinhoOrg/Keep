# Add a Nested Collection

## Description
Covers giving a record its own collection of sub-records — a user owning its sessions, an order owning its items. To add a plain field instead, follow [AddSchemaField.md](/docs/Tutorials/AddSchemaField.md). How nesting is stored is explained in [Schemas.md](/docs/Explanation/Schemas.md).

### Rules
- A nested collection is a field of type `lib.Database` whose `Itens` describe the sub-records' own fields.
- Sub-records are reached only through their owner (`NewSubItem`, `ListAll(fieldName)`) — a nested collection is never returned by `GetSchema`.
- Uniqueness of a `lib.Key` field inside a nested collection is scoped to the owning record, not to the whole database.
- Removing the owner removes every sub-record with it; there is no orphan cleanup to write.

---

## Workflow
1. Add the field to the owner's schema, describing the sub-records' fields in its `Itens`:
   ```go
   {
       Name: "sessions",
       Type: lib.Database,
       Itens: []lib.Item{
           {Name: "token", Type: lib.Key, Required: true},
           {Name: "creation", Type: lib.Int, Required: true},
       },
   }
   ```
2. Insert sub-records through the owning record:
   ```go
   session, e := user.NewSubItem("sessions", map[string]any{
       "token":    "abc123",
       "creation": time.Now().Unix(),
   })
   ```
3. Read them back from the owner:
   ```go
   for _, session := range user.ListAll("sessions") {
       token, _ := session.Get("token")
       fmt.Println(token)
   }
   ```
4. Handle the failures of the operation, following [Errors.md](/docs/Reference/Errors.md) — a duplicated `token` for the same user returns `KeyConflict`, and setting the `sessions` field directly returns `InvalidField`.
5. If the nesting demonstrates a use case not yet covered, add a sample following [AddSample.md](/docs/Tutorials/AddSample.md) — [SubInfos](/examples/SubInfos/SubInfos.go) is the existing one.
6. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
