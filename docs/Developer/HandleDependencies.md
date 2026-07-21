# Handle Dependencies

## Description
Covers adding new requirements to the `Deps` contract in [pkg/deps/deps.go](../../pkg/deps/deps.go) and creating or updating the adapters in [adapters/](../../adapters/) that implement them.

---

## Add a Dependency Requirement

### Workflow
1. Open [pkg/deps/deps.go](../../pkg/deps/deps.go).
2. Add a new descriptive field to the `Deps` struct for the required function:
   ```go
   type Deps struct {
       Write func(key string, value []byte) error
       // New dependency added:
       ReadPrefix func(prefix string) ([]string, error)
   }
   ```
3. If the function must distinguish an expected condition (e.g., a missing key) from a real failure, add a sentinel error next to the existing ones in the same file, and document when to return it in [RequiredApi.md](../Usage/RequiredApi.md).
4. Update every existing adapter in [adapters/](../../adapters/) to implement the new dependency, following [Create or Update an Adapter](#create-or-update-an-adapter). This prevents adapters from breaking when consumers upgrade.
5. Document the new function's contract in [RequiredApi.md](../Usage/RequiredApi.md) and update [deps.Deps.md](../Usage/PublicApi/deps.Deps.md).
6. Compile and run the project's tests to verify no adapters or samples are broken.

---

## Create or Update an Adapter

### Workflow
1. For a new adapter, create a directory inside [adapters/](../../adapters/) (e.g., `adapters/myadapter/`) and an adapter file (e.g., `myadapter.go`).
2. Define a `New()` factory that constructs and returns a `deps.Deps`, implementing every field required by the `Deps` struct in [pkg/deps/deps.go](../../pkg/deps/deps.go) and honoring the contract in [RequiredApi.md](../Usage/RequiredApi.md):
   ```go
   package myadapter

   import (
       "fmt"

       "github.com/MateusMoutinhoOrg/Keep/pkg/deps"
   )

   // New creates a deps.Deps with custom behavior for this adapter.
   func New() deps.Deps {
       return deps.Deps{
           Read: func(key string) ([]byte, error) {
               value, found := /* fetch from your store */
               if !found {
                   // Sentinel errors are part of the contract.
                   return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
               }
               return value, nil
           },
           // ...implement the remaining fields...
       }
   }
   ```
3. Keep external dependencies documented or self-contained.
4. Add the new adapter to the cross-adapter test matrix in [pkg/database/database_test.go](../../pkg/database/database_test.go) (`runWithAdapters`) so every operation is verified against it.
5. For a major adapter, update [STRUCTURE.md](../STRUCTURE.md) to reference it, following [HandleDocumentation.md](./HandleDocumentation.md).
