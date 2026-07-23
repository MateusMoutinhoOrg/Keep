# Add an Adapter

## Description
Covers creating a new opinionated implementation of the `Deps` contract under [adapters/](../../adapters/).

### Rules
- Each adapter lives in its own directory under [adapters/](../../adapters/) and uses a package named after that directory.
- Each adapter exposes a single `New(...) deps.Deps` factory that assigns **every** field of the `Deps` struct.
- An adapter may import [pkg/deps](../../pkg/deps/) but never [pkg/lib](../../pkg/lib/).
- Every function must honor the contract described in [RequiredApi.md](/docs/Reference/RequiredApi.md), including the sentinel errors it has to return.
- The adapter file must follow its specification — locate it in [Specs.md](/docs/Reference/Specs.md).

---

## Workflow
1. Create the adapter directory and its file, both named after the adapter (e.g. `adapters/redis/redis.go`).
2. Declare the package and implement the `New` factory, filling in every field of the `Deps` struct:
   ```go
   package redis

   import (
       "github.com/MateusMoutinhoOrg/Keep/pkg/deps"
   )

   // New creates a deps.Deps backed by a Redis connection.
   func New(addr string) deps.Deps {
       backend := connect(addr) // adapter-specific configuration

       return deps.Deps{
           Write: func(key string, value []byte) error {
               return backend.set(key, value)
           },
           Read: func(key string) ([]byte, error) {
               value, found := backend.get(key)
               if !found {
                   return nil, deps.ErrKeyNotFound // expected condition
               }
               return value, nil
           },
           // ... every remaining field of the contract
       }
   }
   ```
3. Register the new directory and file in [Structure.md](/docs/Reference/Structure.md).
4. Expose the adapter's `New` factory following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
5. If the adapter needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
6. Build the project and run the tests, which exercise every built-in adapter:
   ```bash
   go build ./... && go test ./...
   ```
