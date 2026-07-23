# Add a Library Object

## Description
Covers adding an object created by the library in [pkg/keep/](../../pkg/keep/) and [pkg/database/](../../pkg/database/), with its dependencies wired in by the constructor. To add a plain function on an existing object, follow [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) instead.

### Rules
- An object that needs storage must carry the injected `Deps`, either directly or through a back-pointer to the object that owns them.
- The object is built by a constructor method on its parent — consumers never assemble it by hand.
- Adding a file to [pkg/](../../pkg/) requires updating [Structure.md](/docs/Reference/Structure.md).

---

## Workflow
1. Declare the type in a file named after it (e.g. `pkg/database/schema_item.go`), keeping the wiring fields private:
   ```go
   type SchemaItem struct {
       db     *KeepDatabase // carries the injected Deps
       items  []Item
       prefix string
       id     int64
   }
   ```
2. Add the constructor as a method on the object's parent, wiring the parent's dependencies into the new value:
   ```go
   // GetSchema returns the collection with the given name, or nil when
   // no schema matches.
   func (d *KeepDatabase) GetSchema(name string) *SchemaInstance {
       // ...
       return &SchemaInstance{
           db:     d, // the injected Deps travel with the parent
           schema: schema,
           items:  schema.Itens,
           prefix: d.Props.Path + schema.Name,
       }
   }
   ```
3. Add the object's methods, reaching storage only through the injected `Deps`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md).
4. If a method needs a dependency that is not yet in the contract, add it following [AddDependency.md](/docs/Tutorials/AddDependency.md).
5. If the object is public, expose it, its constructor, and its methods following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
6. Register the new file in [Structure.md](/docs/Reference/Structure.md).
7. If the object needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
8. Build the project and run the tests:
   ```bash
   go build ./... && go test ./...
   ```
