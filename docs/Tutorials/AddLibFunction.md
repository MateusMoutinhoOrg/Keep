# Add a Library Function

## Description
Covers adding a function field to an object in the closed sandbox in [sandbox/internal/](../../sandbox/internal/), and wiring it to the injected dependencies through a factory. To add a whole new object, follow [AddLibObject.md](/docs/Tutorials/AddLibObject.md) instead.

### Rules
- Sandbox code must never import [adapters/](../../adapters/), [examples/](../../examples/), [tests/](../../tests/), a third-party module, or an OS-bound stdlib package — reach storage only through the `Deps` field the object carries. See [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- A function field is filled by a `<Field>Factory(carrier *api.<Object>) <FieldType>` that returns a closure reading `carrier.Deps` (and any other carrier field) at call time — never captured at factory-run time. See [StructContracts.md](/docs/Explanations/StructContracts.md).
- Storage access must respect the invariants of the [Dense Record Pattern](/docs/Explanations/DenseRecordPattern.md): single-key reads and writes, no key listing.
- Logic shared by more than one object goes in [sandbox/internal/dense/](../../sandbox/internal/dense/) as a package-level function taking `deps.Deps` first, and must reference no object type — that is what keeps the package graph acyclic.
- Adding a package or file to [sandbox/internal/](../../sandbox/internal/) requires updating [Structure.md](/docs/References/Structure.md).

---

## Workflow
1. Add the function field to the object's struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go):
   ```go
   type SchemaInstance struct {
       // ...
       Count func() (int64, *Error)
   }
   ```
2. Write the factory in the package of the object it belongs to — entry-point behavior for `Lib` in [sandbox/internal/lib/](../../sandbox/internal/lib/), collection behavior in [schemainstance/](../../sandbox/internal/schemainstance/), record behavior in [schemaitem/](../../sandbox/internal/schemaitem/):
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
3. Call the new factory from the object's constructor, assigning its return value to the field — a field no factory fills stays `nil` and panics on first call:
   ```go
   func New(d deps.Deps, items []api.Item, prefix string) api.SchemaInstance {
       si := api.SchemaInstance{Deps: d, Items: items, Prefix: prefix}
       si.Count = CountFactory(&si)
       // ... every other field factory
       return si
   }
   ```
4. If the function needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
5. If the function is public, expose it following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md). A field missing from `api.go` is unreachable by consumers.
6. If a new package or file was created, register it in [Structure.md](/docs/References/Structure.md).
7. If the function needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
8. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
