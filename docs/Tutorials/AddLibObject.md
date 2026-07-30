# Add a Library Object

## Description
Covers adding an object created by the library in [sandbox/internal/](../../sandbox/internal/), with its dependencies wired in by the constructor. To add a plain function on an existing object, follow [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) instead.

### Rules
- Every object gets **its own package** under `sandbox/internal/`, named after the object itself.
- The object is declared as an interface in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go) and implemented by a struct in its internal package. Consumers only ever see the interface.
- An object that needs storage must carry an exported `Deps deps.Deps` field, filled by the constructor from the parent's `Deps`.
- The object is built by a constructor method on its parent — consumers never assemble it by hand.
- A method returning "nothing found" must return a **literal `nil`**, never a typed nil pointer, or the caller's `== nil` check silently fails.
- Adding a package or file to [sandbox/internal/](../../sandbox/internal/) requires updating [Structure.md](/docs/References/Structure.md).

---

## Workflow
1. Declare the object's interface in [sandbox/contracts/api/api.go](../../sandbox/contracts/api/api.go), returning other objects as their interfaces:
   ```go
   type SchemaInstance interface {
       NewItem(fields map[string]any) (SchemaItem, api.Error)
       FindByKey(key string, keyValue any) SchemaItem
   }
   ```
   Data it exchanges is declared as an interface in the same file — `api.go` holds interfaces and constants only, never a struct. An input interface also needs a constructor in [sandbox/description.go](../../sandbox/description.go), since an interface cannot be built with a composite literal.
2. Create the package and declare the struct implementing it, keeping the wiring fields exported so sibling internal packages can build it:
   ```go
   // sandbox/internal/schemaitem/schemaitem.go
   package schemaitem

   // SchemaItem implements api.SchemaItem.
   type SchemaItem struct {
       Deps     deps.Deps // the injected storage backend
       Items    []api.Item
       Prefix   string
       RecordID int64 // renamed to avoid colliding with the Id() method
   }
   ```
3. Add the constructor as a method on the object's parent, copying the parent's dependencies into the new value and returning the **interface**:
   ```go
   // GetSchema returns the collection with the given name, or nil when
   // no schema matches.
   func (d *KeepDatabase) GetSchema(name string) api.SchemaInstance {
       // ...
       return &schemainstance.SchemaInstance{
           Deps:   d.Deps, // the injected deps travel with the parent
           Items:  schema.Itens,
           Prefix: d.Description.Path + schema.Name,
       }
   }
   ```
4. Add the object's methods in its own package file, reaching storage only through its `Deps`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md). Logic shared with other objects goes in [sandbox/internal/dense/](../../sandbox/internal/dense/) so no import cycle forms.
5. If a method needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md). Never import `os`, `net`, or a third-party module inside the sandbox.
6. If the object is public, expose it, its constructor, and its methods following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
7. Register any new package or file in [Structure.md](/docs/References/Structure.md).
8. If the object needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
9. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
