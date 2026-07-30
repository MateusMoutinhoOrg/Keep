# Add a Library Function

## Description
Covers adding a function to the closed sandbox in [sandbox/internal/](../../sandbox/internal/), and wiring it to the injected dependencies. To add an object created by the library, follow [AddLibObject.md](/docs/Tutorials/AddLibObject.md) instead.

### Rules
- Sandbox code must never import [adapters/](../../adapters/), [examples/](../../examples/), [tests/](../../tests/), a third-party module, or an OS-bound stdlib package — reach storage only through the `Deps` field the object carries. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- Storage access must respect the invariants of the [Dense Record Pattern](/docs/Explanations/DenseRecordPattern.md): single-key reads and writes, no key listing.
- Logic shared by more than one object goes in [sandbox/internal/dense/](../../sandbox/internal/dense/) as a package-level function taking `deps.Deps` first, and must reference no object type — that is what keeps the package graph acyclic.
- Adding a package or file to [sandbox/internal/](../../sandbox/internal/) requires updating [Structure.md](/docs/References/Structure.md).

---

## Workflow
1. Define the function in the package of the object it hangs off — entry-point behavior on `Lib` in [sandbox/internal/lib/](../../sandbox/internal/lib/), collection operations in [schemainstance/](../../sandbox/internal/schemainstance/), record operations in [schemaitem/](../../sandbox/internal/schemaitem/):
   ```go
   // Count returns the number of live records in the collection.
   func (si *SchemaInstance) Count() (int64, api.Error) {
       size, err := dense.ReadCount(si.Deps, dense.SizeKey(si.Prefix))
       if err != nil {
           return 0, dense.InternalError(err)
       }
       return size, nil
   }
   ```
2. If the function needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
3. If the function is public, declare it on the object's interface in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go) and expose it following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md). A method missing from the interface is unreachable by consumers.
4. If a new package or file was created, register it in [Structure.md](/docs/References/Structure.md).
5. If the function needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
6. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
