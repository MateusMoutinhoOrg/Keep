# Add a Database Operation

## Description
Covers adding a new operation to the storage engine in [sandbox/lib/](../../sandbox/lib/) — a function field on a collection or a record — without breaking the key layout it shares with every other operation. For a function that needs no storage access, follow [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#add-a-library-function); for a new object, follow [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#add-a-library-object).

### Rules
- Read [DenseRecordPattern.md](/docs/References/DenseRecordPattern.md) first: the operation must preserve its invariants — single-key reads and writes, no key listing, ids never reused, and the dense position list left without holes.
- Keys are never assembled inline: reuse the builders in [dense/](../../sandbox/lib/dense/) (`ValueKey`, `IndexKey`, `ListKey`, `PositionKey`, `SizeKey`).
- Storage is reached only through the object's `Deps` field, from inside a factory's closure — `si.Deps.<Field>()`, `s.Deps.<Field>()`. Never import `adapters/` inside the sandbox.
- Expected failures return a typed `*api.Error` with the matching [ErrorType](/docs/References/Errors.md); storage failures are wrapped with `dense.InternalError(err)`.
- Every write ordering must leave the database readable if the process dies between two writes — publish the commit point last.
- A public operation must be declared as a function field on the object's struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go), filled by a factory, or consumers cannot call it.

---

## Workflow
1. Pick the object the operation belongs to: collection-wide behavior goes in [schemainstance/](../../sandbox/lib/schemainstance/), record behavior in [schemaitem/](../../sandbox/lib/schemaitem/), shared key procedures in [dense/](../../sandbox/lib/dense/).
2. Write the factory, reaching storage only through the injected deps and the key builders:
   ```go
   // CountFactory fills api.SchemaInstance.Count.
   func CountFactory(si *api.SchemaInstance) func() (int64, *api.Error) {
       return func() (int64, *api.Error) {
           size, err := dense.ReadCount(si.Deps, dense.SizeKey(si.Prefix))
           if err != nil {
               return 0, dense.InternalError(err)
           }
           return size, nil
       }
   }
   ```
3. If the operation writes more than one key, order the writes so the last one commits the change, and document the ordering in a step comment — as `schemaitem.New` and `RemoveFactory` do.
4. If it needs a storage call the contract does not have yet, add it following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency).
5. Declare the field on the matching struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go), call the factory from the object's constructor to fill it, and add any new data type it exchanges to [sandbox/contracts/api/](../../sandbox/contracts/api/).
6. Cover the operation with a test exercising both built-in adapters (`adapters/standard`, `adapters/native`).
7. Expose it following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#expose-in-the-public-api), and document any new failure in [Errors.md](/docs/References/Errors.md).
8. If it deserves a runnable demonstration, add one following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample).
9. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
