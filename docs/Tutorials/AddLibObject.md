# Add a Library Object

## Description
Covers adding an object created by the library in [pkg/lib/](../../pkg/lib/), with its dependencies wired in by the constructor. To add a plain function on an existing object, follow [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) instead.

### Rules
- An object that needs storage must carry a private `deps deps.Deps` field, filled by the constructor from the parent's `deps`.
- The object is built by a constructor method on its parent — consumers never assemble it by hand.
- Adding a file to [pkg/lib/](../../pkg/lib/) requires updating [Structure.md](/docs/Reference/Structure.md).

---

## Workflow
1. Declare the type in [pkg/lib/types.go](../../pkg/lib/types.go), keeping the wiring fields private:
   ```go
   type SchemaItem struct {
       deps   deps.Deps // the injected storage functions
       items  []Item
       prefix string
       id     int64
   }
   ```
2. Add the constructor as a method on the object's parent, copying the parent's dependencies into the new value:
   ```go
   // GetSchema returns the collection with the given name, or nil when
   // no schema matches.
   func (d *KeepDatabase) GetSchema(name string) *SchemaInstance {
       // ...
       return &SchemaInstance{
           deps:   d.deps, // the injected deps travel with the parent
           items:  schema.Itens,
           prefix: d.Props.Path + schema.Name,
       }
   }
   ```
3. Add the object's methods in a file named after it (e.g. [pkg/lib/schema_item.go](../../pkg/lib/schema_item.go)), reaching storage only through its `deps`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md).
4. If a method needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
5. If the object is public, expose it, its constructor, and its methods following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
6. Register any new file in [Structure.md](/docs/Reference/Structure.md).
7. If the object needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
8. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
