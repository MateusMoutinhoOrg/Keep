# Handle Library Elements

## Description
Covers the three moves that grow the library itself: adding a function field to an existing object, adding a whole new object, and publishing either one in the public API index. All of them happen inside the closed sandbox in [sandbox/lib/](/sandbox/lib/), reaching every effect through the injected `Deps`. Adding the dependency a new behavior needs is [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md); demonstrating it is [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md).

### Rules
- Sandbox code must never import [adapters/](/adapters/), [examples/libraryExamples/](/examples/libraryExamples/), `tests/`, a third-party module, or an OS-bound stdlib package — reach storage only through the `Deps` field the object carries. See [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- Every function field is filled by a `<Field>Factory(carrier *api.<Object>) <FieldType>` that returns a closure reading `carrier.Deps` (and any other carrier field) at call time — never captured at factory-run time. See [Factories](/docs/References/Specs/Factories/Specs.md) and [StructContracts.md](/docs/References/StructContracts.md).
- Storage access must respect the invariants of the [Dense Record Pattern](/docs/References/DenseRecordPattern.md): single-key reads and writes, no key listing.
- Logic shared by more than one object goes in [sandbox/lib/dense/](/sandbox/lib/dense/) as a package-level function taking `deps.Deps` first, and must reference no object type — that is what keeps the package graph acyclic.
- Adding a package or file to [sandbox/lib/](/sandbox/lib/) requires updating [Structure.md](/docs/References/Structure.md).

---

## Add a Library Function
1. Add the function field to the object's struct in [sandbox/contracts/api/api.go](/sandbox/contracts/api/api.go):
   ```go
   type SchemaInstance struct {
       // ...
       Count func() (int64, *Error)
   }
   ```
2. Write the factory in the package of the object it belongs to — entry-point behavior for `Lib` in its own file under [sandbox/lib/publicfunctions/](/sandbox/lib/publicfunctions/), collection behavior in [schemainstance/](/sandbox/lib/schemainstance/), record behavior in [schemaitem/](/sandbox/lib/schemaitem/):
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
   A field of `api.Lib` is assigned in `New` in [sandbox/lib/new.go](/sandbox/lib/new.go) instead, from the factory in `publicfunctions`.
4. If the function needs a dependency that is not yet in the contract, add it following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency).
5. If the function is public, publish it as described in [Expose in the Public API](#expose-in-the-public-api). A field missing from `api.go` is unreachable by consumers.
6. If a new package or file was created, register it in [Structure.md](/docs/References/Structure.md).
7. If the function needs a runnable demonstration, add one following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample).
8. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```

---

## Add a Library Object
1. Declare the object's struct in [sandbox/contracts/api/api.go](/sandbox/contracts/api/api.go), leading with `Deps` and declaring one function field per behavior:
   ```go
   type SchemaInstance struct {
       Deps      deps.Deps
       Items     []Item
       Prefix    string
       NewItem   func(fields map[string]any) (SchemaItem, *Error)
       FindByKey func(key string, keyValue any) (SchemaItem, bool)
   }
   ```
   Data it exchanges is declared in the same file too — `api.go` holds every type in the project. A type with no behavior (like `Item`, `Schema`, `Props`) is plain data, buildable directly with a composite literal; no constructor package is needed for it.
2. Create the package under `sandbox/lib/`, named after the object itself, and write one `<Field>Factory` per function field, each taking a pointer to the `api` struct being built (the carrier) and returning that field's closure:
   ```go
   // sandbox/lib/schemainstance/schemainstance.go
   package schemainstance

   // NewItemFactory fills api.SchemaInstance.NewItem.
   func NewItemFactory(si *api.SchemaInstance) func(fields map[string]any) (api.SchemaItem, *api.Error) {
       return func(fields map[string]any) (api.SchemaItem, *api.Error) {
           return schemaitem.New(si.Deps, si.Items, si.Prefix, fields)
       }
   }
   ```
3. Write the package's `New` constructor — the factory aggregate — propagating the parent's `Deps` and running every field factory:
   ```go
   // New builds an api.SchemaInstance, storing the injected Deps and the
   // schema's fields and prefix, and runs every factory over it.
   func New(d deps.Deps, items []api.Item, prefix string) api.SchemaInstance {
       si := api.SchemaInstance{Deps: d, Items: items, Prefix: prefix}
       si.NewItem = NewItemFactory(&si)
       si.FindByKey = FindByKeyFactory(&si)
       return si
   }
   ```
   The parent calls it from its own factory:
   ```go
   // GetSchemaFactory fills api.KeepDatabase.GetSchema.
   func GetSchemaFactory(kd *api.KeepDatabase) func(name string) (api.SchemaInstance, bool) {
       return func(name string) (api.SchemaInstance, bool) {
           for _, schema := range kd.Props.Schemas {
               if schema.Name == name {
                   return schemainstance.New(kd.Deps, schema.Itens, kd.Props.Path+schema.Name), true
               }
           }
           return api.SchemaInstance{}, false
       }
   }
   ```
   A lookup that can fail returns `(value, ok bool)`, never a typed nil pointer — the object is a struct, not an interface, and has no nil form.
4. Add the object's remaining factories in its own package file, reaching storage only through `si.Deps`, as described in [Add a Library Function](#add-a-library-function). Logic shared with other objects goes in [sandbox/lib/dense/](/sandbox/lib/dense/) so no import cycle forms.
5. If a factory needs a dependency that is not yet in the contract, add it following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency). Never import `os`, `net`, or a third-party module inside the sandbox.
6. If the object is public, publish it and its fields as described in [Expose in the Public API](#expose-in-the-public-api).
7. Register any new package or file in [Structure.md](/docs/References/Structure.md).
8. If the object needs a runnable demonstration, add one following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample).
9. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```

---

## Expose in the Public API
1. Open [PublicApi.md](/docs/References/PublicApi.md).
2. Add the function, struct, or field to the section matching its kind, with a one-line description.
3. Create the detail page under [docs/References/PublicApi/](/docs/References/PublicApi/), named `<pkg>.<Symbol>.md` (e.g., `api.SchemaInstance.md`), following [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md#add-a-document). Document a struct's function fields under a `## Fields` section, and any plain data field directly.
4. Link the new detail page from its entry in [PublicApi.md](/docs/References/PublicApi.md).
5. Register the detail page in [Structure.md](/docs/References/Structure.md).
