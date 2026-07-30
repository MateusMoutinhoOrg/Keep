# Add an Adapter

## Description
Covers creating a new opinionated implementation of the `Deps` contract under [adapters/](../../adapters/).

### Rules
- Each adapter lives in its own directory under [adapters/](../../adapters/) and uses a package named after that directory.
- Each adapter declares a struct implementing **every** method of the `Deps` interface, and exposes a single `New(...) deps.Deps` factory returning it as the interface.
- An adapter sits **outside** the sandbox and is the only place OS-bound and third-party code may appear. It may import [sandbox/contracts/deps](../../sandbox/contracts/deps/) but never [sandbox](../../sandbox/) or [sandbox/internal](../../sandbox/internal/) — see [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- Every method must honor the contract described in [RequiredApi.md](/docs/References/RequiredApi.md), including the sentinel errors it has to return.
- The adapter file must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).

---

## Workflow
1. Create the adapter directory and its file, both named after the adapter (e.g. `adapters/redis/redis.go`).
2. Declare the package, a struct holding the adapter's state, and the `New` factory:
   ```go
   package redis

   import (
       "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
   )

   // redis holds the adapter's state — the opinionated part.
   type redis struct {
       backend *connection
   }

   // New creates a deps.Deps backed by a Redis connection.
   func New(addr string) deps.Deps {
       return &redis{backend: connect(addr)} // adapter-specific configuration
   }
   ```
3. Implement every method of the contract on that struct:
   ```go
   func (r *redis) Write(key string, value []byte) error {
       return r.backend.set(key, value)
   }

   func (r *redis) Read(key string) ([]byte, error) {
       value, found := r.backend.get(key)
       if !found {
           return nil, deps.ErrKeyNotFound // expected condition
       }
       return value, nil
   }

   // ... every remaining method of the contract
   ```
   A missing method means the struct does not satisfy `deps.Deps`, and `New` fails to compile — the contract is checked for you.
4. Register the new directory and file in [Structure.md](/docs/References/Structure.md), and add a row for the adapter in [Adapters.md](/docs/References/Adapters.md).
5. Expose the adapter's `New` factory following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
6. If the adapter needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
7. Build the project and run the tests, which exercise every built-in adapter:
   ```bash
   go build ./... && go test ./...
   ```
