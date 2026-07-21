# Handle Library Functions

## Description
Covers adding functions and objects to the library in [pkg/keep/](../../pkg/keep/) and [pkg/database/](../../pkg/database/), wiring them to dependencies, and exposing them in the public API. For adding new dependency requirements, see [HandleDependencies.md](./HandleDependencies.md).

### Rules
- Pure logic only: `pkg/` never imports concrete implementations — all storage access goes through the injected `deps.Deps`.
- Operations must respect the Dense Record Pattern invariants described in [DatabaseSchema.md](../DatabaseSchema.md): single-key reads/writes only, no key listing, and the documented write orderings.

---

## Add a Function

### Workflow
1. Define the function in the layer it belongs to:
   - Entry-point functionality (creating databases) → [pkg/keep/](../../pkg/keep/).
   - Collection or record operations → [pkg/database/](../../pkg/database/), reaching storage through `db.Deps`:
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
2. If the function needs a new dependency, add it following [HandleDependencies.md](./HandleDependencies.md).
3. If the function is public, expose it following [Expose in the Public API](#expose-in-the-public-api).
4. Check whether a sample is needed to demonstrate the function. If so, follow [HandleSamples.md](./HandleSamples.md).
5. Add coverage in [pkg/database/database_test.go](../../pkg/database/database_test.go) — tests run against both built-in adapters via `runWithAdapters`.

---

## Add an Object

### Workflow
1. Declare the type in the relevant file of [pkg/database/](../../pkg/database/) (schema description types in `schema.go`, runtime objects in their own file), keeping deps-carrying fields unexported:
   ```go
   type ExampleObject struct {
       db     *KeepDatabase
       prefix string
   }
   ```
2. Add the constructor and methods in a dedicated file, wiring the deps from the parent object.
3. If the object is public, expose it and its constructor following [Expose in the Public API](#expose-in-the-public-api).
4. Check whether a sample is needed to demonstrate the object. If so, follow [HandleSamples.md](./HandleSamples.md).

---

## Expose in the Public API

### Workflow
1. Open [PublicApi.md](../Usage/PublicApi.md).
2. Document the new function, struct, or method under the relevant section.
3. Add a dedicated detail page under [docs/Usage/PublicApi/](../Usage/PublicApi/) (e.g., `database.SchemaInstance.md`) when it needs extensive explanation, and link to it from the index. Methods of an existing struct are documented in that struct's detail page.
4. If the change affects a narrative guide ([Records.md](../Usage/Records.md), [Schemas.md](../Usage/Schemas.md), [Errors.md](../Usage/Errors.md)), update it in the same commit.
