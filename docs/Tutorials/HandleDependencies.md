# Handle Dependencies

## Description
Explains how the library receives its dependencies — the `Deps` contract in [sandbox/contracts/deps/deps.go](/sandbox/contracts/deps/deps.go) — how an injected value propagates through the object graph, how to add a requirement to the contract, and how to build an adapter that fills it. Assumes [SandboxIsolation.md](/docs/References/SandboxIsolation.md) and [StructContracts.md](/docs/References/StructContracts.md). Using the dependency from library code is covered by [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md).

---

## Find the Dependencies You Can Use

`sandbox/contracts/deps` declares what the library needs; `sandbox/contracts/api` declares what it hands back. Nothing else crosses the boundary, and `lib.New(deps.Deps) api.Lib` is the single wiring point.

```go
// sandbox/contracts/deps/deps.go — what the library needs
type Deps struct {
	Write               func(key string, value []byte) error
	WriteIfKeyNotExists func(key string, value []byte) error
	WriteIfValueEquals  func(key string, value []byte, oldValue []byte) error
	Append              func(key string, value []byte) error
	InsertAt            func(key string, position int64, value []byte) error
	Exists              func(key string) (bool, error)
	Read                func(key string) ([]byte, error)
	ReadAt              func(key string, position int64, size int64) ([]byte, error)
	Delete              func(key string) error
	Lock                func(key string, time int) error
	UnLock              func(key string) error
}
```

Every field addresses a **single key**: the library never lists, scans a prefix, or queries a range, which is what lets any key-value backend host it. The behavior each field must honor — including the sentinel error it has to return for an expected condition — is documented field by field in [RequiredApi.md](/docs/References/RequiredApi.md).

`Deps` is the *only* door in the sandbox wall: since nothing under `sandbox/` may import an adapter, a third-party module, or an OS-bound stdlib package, every effect the library performs has to be a field on this struct.

`lib.New` stores the `Deps` on the `api.Lib` struct and runs the factories over it; each closure reads `l.Deps` when the field is *called*, not when the factory ran. Every object the lib creates receives the same `Deps`, passed into the object package's `New` constructor. So a dependency injected once is reachable from anywhere in the object graph — see [DepsMechanic.md](/docs/References/DepsMechanic.md).

---

## Add a Dependency

### Rules
- A requirement is a **function field** declaring behavior, never a concrete implementation.
- A new field must be filled by **every** adapter in [adapters/](/adapters/) in the same commit. The compiler will **not** catch a missing one: the field stays nil and panics on first call.
- The `Deps` struct must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).
- Every field of the contract must have its behavior documented in [RequiredApi.md](/docs/References/RequiredApi.md).
- Adding a field is the **only** way to give the sandbox a capability it lacks; reaching for `os`, `net`, or a third-party module inside `sandbox/` is forbidden by [SandboxIsolation.md](/docs/References/SandboxIsolation.md).

### Workflow
1. Add the field to the `Deps` struct in [sandbox/contracts/deps/deps.go](/sandbox/contracts/deps/deps.go), named after the behavior it provides:
   ```go
   type Deps struct {
       Write      func(key string, value []byte) error
       ReadPrefix func(prefix string) ([]string, error) // new requirement
   }
   ```
2. If the field must distinguish an expected condition (e.g. a missing key) from a real failure, add a sentinel error next to the existing ones in the same file:
   ```go
   var ErrPrefixNotFound = errors.New("keep: prefix not found")
   ```
3. Fill the new field in every adapter under [adapters/](/adapters/) with a `<Field>Factory`, as described in [Create an Adapter in This Repository](#create-an-adapter-in-this-repository).
4. Document the field's contract — arguments, return values, and the sentinel it must return — in [RequiredApi.md](/docs/References/RequiredApi.md), and update the [deps.Deps](/docs/References/PublicApi/deps.Deps.md) detail page.
5. Use the dependency from the library through the injected `Deps`, following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#add-a-library-function).
6. If the requirement changes how dependencies behave for consumers, update [DepsMechanic.md](/docs/References/DepsMechanic.md).
7. Build the project and confirm no adapter, sample, or test breaks:
   ```bash
   go build ./... && go test ./...
   ```

---

## Overwrite an Adapter Function

Take the `deps.Deps` an adapter returns and reassign the field you want; every other field keeps the adapter's implementation:

```go
myDeps := keepadapter.New()

// Replace only the reader — every write still goes to the real filesystem
realRead := myDeps.Read
myDeps.Read = func(key string) ([]byte, error) {
	fmt.Println("reading", key)
	return realRead(key)
}

keep := keeplib.New(myDeps)
```

> **Careful:** patch the `deps.Deps` value **before** calling `lib.New`. The factories close over the `api.Lib` they ran on, so assigning to `keep.Deps.Read` afterwards changes nothing — see [StructContracts.md](/docs/References/StructContracts.md#what-it-costs).

---

## Create an Adapter in This Repository

Covers creating a new opinionated implementation of the `Deps` contract under [adapters/](/adapters/). The shipped adapters are listed in [Adapters.md](/docs/References/Adapters.md).

### Rules
- Each adapter lives in its own directory under [adapters/](/adapters/) and uses a package named after that directory.
- Each adapter declares a struct carrying a `Deps deps.Deps` field, fills **every** field of that contract with a `<Field>Factory`, and exposes a single `New(...) deps.Deps` constructor returning the filled contract struct — never the concrete adapter type.
- An adapter sits **outside** the sandbox and is the only place OS-bound and third-party code may appear. It may import [sandbox/contracts/deps](/sandbox/contracts/deps/) but never [sandbox](/sandbox/) or [sandbox/lib](/sandbox/lib/) — see [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- Every field must honor the contract described in [RequiredApi.md](/docs/References/RequiredApi.md), including the sentinel errors it has to return.
- The adapter file must follow its specification — locate it in [Specs.md](/docs/References/Specs.md), and follow the factory shape in [Factories](/docs/References/Specs/Factories/Specs.md).

### Workflow
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
6. Publish the adapter's `New` factory following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#expose-in-the-public-api).
7. If the adapter needs a runnable demonstration, add one following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample).
8. Build the project and run the tests, which exercise every built-in adapter:
   ```bash
   go build ./... && go test ./...
   ```

---

## Create an Adapter in Your Project

An adapter does not have to live in this repository. For complete control, build the `deps.Deps` as a struct literal in your own module — there is no type to declare and no method set to satisfy:

```go
import (
	keepdeps "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

myDeps := keepdeps.Deps{
	Read:  func(key string) ([]byte, error) { return myStore.Get(key) },
	Write: func(key string, value []byte) error { return myStore.Put(key, value) },
	// ... one literal per remaining field
}
keep := keeplib.New(myDeps)
```

Every field your implementation fills must honor [RequiredApi.md](/docs/References/RequiredApi.md), returning the `keepdeps` sentinel errors for the expected conditions so the library can tell them apart from real failures.

> **Careful:** the compiler cannot tell you a field is missing — an unfilled field panics on first call. In practice, start from a shipped adapter and patch what you need, as shown in [Overwrite an Adapter Function](#overwrite-an-adapter-function).
