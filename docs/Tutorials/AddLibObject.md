# Add a Library Object

## Description
Covers adding an object created by the library in [sandbox/internal/](../../sandbox/internal/), with its dependencies wired in by the object's constructor. To add a plain function field on an existing object, follow [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) instead.

### Rules
- Every object gets **its own package** under `sandbox/internal/`, named after the object itself, holding only `<Field>Factory` functions and the constructor that runs them — no type declarations.
- The object is declared as a **struct of function fields** in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go) — never an interface, and never declared in `sandbox/internal/`.
- An object that needs storage must lead with an exported `Deps deps.Deps` field, propagated by the constructor from the parent's `Deps`.
- The object is built by a constructor function on its parent (e.g. `GetSchemaFactory`'s closure calling `schemainstance.New`) — consumers never assemble it by hand.
- A lookup that can fail returns `(value, ok bool)`, never a typed nil pointer — the object is a struct, not an interface, and has no nil form. See [StructContracts.md](/docs/Explanations/StructContracts.md).
- Adding a package or file to [sandbox/internal/](../../sandbox/internal/) requires updating [Structure.md](/docs/References/Structure.md).

---

## Workflow
1. Declare the object's struct in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go), leading with `Deps` and declaring one function field per behavior:
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
2. Create the package and write one `<Field>Factory` per function field, each taking a pointer to the `api` struct being built (the carrier) and returning that field's closure:
   ```go
   // sandbox/internal/schemainstance/schemainstance.go
   package schemainstance

   // NewItemFactory fills api.SchemaInstance.NewItem.
   func NewItemFactory(si *api.SchemaInstance) func(fields map[string]any) (api.SchemaItem, *api.Error) {
       return func(fields map[string]any) (api.SchemaItem, *api.Error) {
           return schemaitem.New(si.Deps, si.Items, si.Prefix, fields)
       }
   }
   ```
3. Write the package's `New` constructor — the factory aggregate — as a method on the object's parent (or a plain function the parent's factory calls), propagating the parent's `Deps` and running every field factory:
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
4. Add the object's remaining factories in its own package file, reaching storage only through `si.Deps`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md). Logic shared with other objects goes in [sandbox/internal/dense/](../../sandbox/internal/dense/) so no import cycle forms.
5. If a factory needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md). Never import `os`, `net`, or a third-party module inside the sandbox.
6. If the object is public, expose it and its fields following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
7. Register any new package or file in [Structure.md](/docs/References/Structure.md).
8. If the object needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
9. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
