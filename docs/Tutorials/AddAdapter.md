# Add an Adapter

## Description
Covers creating a new opinionated implementation of the `Deps` contract under [adapters/](../../adapters/).

### Rules
- Each adapter lives in its own directory under [adapters/](../../adapters/) and uses a package named after that directory.
- Each adapter declares a struct carrying a `Deps deps.Deps` field, fills **every** field of that contract with a `<Field>Factory`, and exposes a single `New(...) deps.Deps` factory aggregate returning the filled contract struct — never the concrete adapter type.
- An adapter sits **outside** the sandbox and is the only place OS-bound and third-party code may appear. It may import [sandbox/contracts/deps](../../sandbox/contracts/deps/) but never [sandbox](../../sandbox/) or [sandbox/internal](../../sandbox/internal/) — see [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).
- Every field must honor the contract described in [RequiredApi.md](/docs/References/RequiredApi.md), including the sentinel errors it has to return.
- The adapter file must follow its specification — locate it in [Specs.md](/docs/References/Specs.md), and follow the factory shape in [StructContracts.md](/docs/Explanations/StructContracts.md).

---

## Workflow
1. Create the adapter directory and its file, both named after the adapter (e.g. `adapters/redis/redis.go`).
2. Declare the package and the struct holding the adapter's state — the carrier every factory closes over:
   ```go
   package redis

   import (
       "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
   )

   // RedisAdapter fills deps.Deps against a Redis connection.
   type RedisAdapter struct {
       Deps    deps.Deps // the contract this adapter fills
       backend *connection
   }
   ```
3. Write a `<Field>Factory` per requirement, each taking a pointer to the adapter and returning that field's closure:
   ```go
   func WriteFactory(r *RedisAdapter) func(key string, value []byte) error {
       return func(key string, value []byte) error {
           return r.backend.set(key, value)
       }
   }

   func ReadFactory(r *RedisAdapter) func(key string) ([]byte, error) {
       return func(key string) ([]byte, error) {
           value, found := r.backend.get(key)
           if !found {
               return nil, deps.ErrKeyNotFound // expected condition
           }
           return value, nil
       }
   }

   // ... one factory per remaining field of the contract
   ```
4. Write `New`, the factory aggregate: it builds the adapter, assigns every factory's return value to the matching field, and returns the filled `deps.Deps`:
   ```go
   // New creates a deps.Deps backed by a Redis connection.
   func New(addr string) deps.Deps {
       r := &RedisAdapter{backend: connect(addr)} // adapter-specific configuration
       r.Deps.Write = WriteFactory(r)
       r.Deps.Read = ReadFactory(r)
       // ... one assignment per remaining field
       return r.Deps
   }
   ```
   A field no factory fills stays `nil` and panics on first call — nothing checks completeness for you, unlike an interface. Visit every field before considering the adapter done.
5. Register the new directory and file in [Structure.md](/docs/References/Structure.md), and add a row for the adapter in [Adapters.md](/docs/References/Adapters.md).
6. Expose the adapter's `New` factory following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
7. If the adapter needs a runnable demonstration, add one following [AddSample.md](/docs/Tutorials/AddSample.md).
8. Build the project and run the tests, which exercise every built-in adapter:
   ```bash
   go build ./... && go test ./...
   ```
