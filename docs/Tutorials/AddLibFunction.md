# Add a Library Function

## Description
Covers adding a function to the pure library in [pkg/keep/](../../pkg/keep/) and [pkg/database/](../../pkg/database/), and wiring it to the injected dependencies. To add an object created by the library, follow [AddLibObject.md](/docs/Tutorials/AddLibObject.md) instead.

### Rules
- Library code must never import anything from [adapters/](../../adapters/) — reach storage only through the injected `Deps`.
- Storage access must respect the invariants of the [Dense Record Pattern](/docs/Explanation/DenseRecordPattern.md): single-key reads and writes, no key listing.
- Adding a file to [pkg/](../../pkg/) requires updating [Structure.md](/docs/Reference/Structure.md).

---

## Workflow
1. Define the function in the layer it belongs to — entry-point behavior in [pkg/keep/](../../pkg/keep/), collection or record operations in [pkg/database/](../../pkg/database/):
   ```go
   // Count returns the number of live records in the collection.
   func (si *SchemaInstance) Count() (int64, *Error) {
       size, err := readCount(si.db.Deps, sizeKey(si.prefix))
       if err != nil {
           return 0, internalError(err)
       }
       return size, nil
   }
   ```
2. If the function needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
3. If the function is public, expose it following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
4. If a new file was created, register it in [Structure.md](/docs/Reference/Structure.md).
5. If the function needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
6. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
