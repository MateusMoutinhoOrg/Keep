# Add a Database Operation

## Description
Covers adding a new operation to the storage engine in [pkg/lib/](../../pkg/lib/) — a method on a collection or a record — without breaking the key layout it shares with every other operation. For a function that needs no storage access, follow [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md); for a new object, follow [AddLibObject.md](/docs/Tutorials/AddLibObject.md).

### Rules
- Read [DenseRecordPattern.md](/docs/Explanation/DenseRecordPattern.md) first: the operation must preserve its invariants — single-key reads and writes, no key listing, ids never reused, and the dense position list left without holes.
- Keys are never assembled inline: reuse the builders in [dense.go](../../pkg/lib/dense.go) (`valueKey`, `indexKey`, `listKey`, `positionKey`, `sizeKey`).
- Storage is reached only through the object's private `deps` field — `si.deps.<Field>()`, `s.deps.<Field>()`. Never import `adapters/`.
- Expected failures return a typed `*Error` with the matching [ErrorType](/docs/Reference/Errors.md); storage failures are wrapped with `internalError(err)`.
- Every write ordering must leave the database readable if the process dies between two writes — publish the commit point last.

---

## Workflow
1. Pick the object the operation belongs to: collection-wide behavior goes in [schema_instance.go](../../pkg/lib/schema_instance.go), record behavior in [schema_item.go](../../pkg/lib/schema_item.go), shared key procedures in [dense.go](../../pkg/lib/dense.go).
2. Write the method, reaching storage only through the injected deps and the key builders:
   ```go
   // Count returns the number of live records in the collection.
   func (si *SchemaInstance) Count() (int64, *Error) {
       size, err := readCount(si.deps, sizeKey(si.prefix))
       if err != nil {
           return 0, internalError(err)
       }
       return size, nil
   }
   ```
3. If the operation writes more than one key, order the writes so the last one commits the change, and document the ordering in a step comment — as `newItem` and `Remove` do.
4. If it needs a storage call the contract does not have yet, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
5. Cover the operation in [lib_test.go](../../pkg/lib/lib_test.go), which runs every test against both built-in adapters.
6. Expose it following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md), and document any new failure in [Errors.md](/docs/Reference/Errors.md).
7. If it deserves a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
8. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
