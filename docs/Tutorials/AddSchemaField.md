# Add a Field to a Schema

## Description
Covers adding a field to a collection that already holds records, without breaking the records written before the change. To describe a database from scratch, follow [DefineDatabase.md](/docs/Tutorials/DefineDatabase.md); to add a nested collection instead of a plain field, follow [AddNestedCollection.md](/docs/Tutorials/AddNestedCollection.md).

### Rules
- Records are stored field by field, so old records simply have no value for the new field — reading it returns a `NotFound` [Error](/docs/References/Errors.md), never a corrupted record.
- Marking a new field `Required: true` only affects `NewItem` calls made after the change; it never rewrites existing records.
- Adding a `keeptypes.Key` field to a collection with existing records leaves those records unindexed for it — they will not be found by `FindByKey` until the value is written.

---

## Workflow
1. Add the new `keeptypes.Item` to the collection's `Itens` in the `Schemas` value every program touching this database shares:
   ```go
   {Name: "nickname", Type: keeptypes.Key, Required: false}
   ```
2. Decide how old records get the value:
   - leave it absent and treat `NotFound` as "not set", or
   - backfill it by iterating the collection and writing the value:
     ```go
     all, err := users.ListAll()
     if err != nil {
         return err
     }
     for _, user := range all {
         if user.CheckKeysPresence([]string{"nickname"}) {
             continue
         }
         if e := user.Update("nickname", defaultNicknameFor(user)); e != nil {
             return e
         }
     }
     ```
3. Handle the field's absence wherever it is read, following [Errors.md](/docs/References/Errors.md):
   ```go
   nickname, e := user.Get("nickname")
   if e != nil && e.Type == keeptypes.NotFound {
       nickname = ""
   }
   ```
4. If the field demonstrates a use case not yet covered, add a sample following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample).
5. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
