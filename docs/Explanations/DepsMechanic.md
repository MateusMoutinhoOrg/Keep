# Dependency Mechanic

## Description
Keep never talks to storage directly. Every read and write goes through the `deps.Deps` interface, so you choose where the bytes live: one of the built-in adapters, your own implementation, or an adapter with some methods overridden.

This guide explains how this injection works, how to write your own backend, and how to override specific dependencies.

---

## Dependency Injection

The `Deps` interface is defined in `sandbox/contracts/deps/deps.go` and its methods are documented in [Required API](/docs/References/RequiredApi.md). It is a plain key-value API: write, read, delete, exists, plus conditional writes and locks.

A backend is any type implementing every method of that interface — the compiler checks the contract is complete, so a half-finished adapter never reaches a consumer.

To inject a built-in adapter:

```go
import (
	"github.com/MateusMoutinhoOrg/Keep/adapters/standard" // filesystem adapter
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

func main() {
	deps := standard.New() // pick a backend...
	keep := lib.New(deps)  // ...and inject it
	_ = keep
}
```

The value flows one way and never comes back out: `lib.New` stores it on the internal `Lib`, which copies it into every `KeepDatabase` it creates, which copies it into every `SchemaInstance`, which copies it into every `SchemaItem`. Every object therefore reaches storage through the exact backend the caller chose, and no object anywhere inside `sandbox/` can reach a different one — see [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).

The shipped adapters and when to use each are listed in [Adapters.md](/docs/References/Adapters.md).

---

## Creating Custom Dependencies

You can run Keep over any storage backend by implementing every method of `deps.Deps` on a type of your own.

### Rules
- **Return the sentinel errors**: The database layer distinguishes "key doesn't exist" from "storage is broken" using the sentinels in `sandbox/contracts/deps` (`ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, `ErrKeyLocked`). Wrap them with `fmt.Errorf("%w: ...", ...)` so `errors.Is` matches.
- **Only per-key operations are needed**: Keep never lists keys or scans prefixes, so any store with `get`/`set`/`delete` can back it.
- **Implement the whole interface**: a missing method is a compile error at the point you pass the value to `lib.New`.

### Workflow

1. Declare a type holding your storage client.
2. Implement every method of `deps.Deps` on it, following [Required API](/docs/References/RequiredApi.md).
3. Expose a constructor returning `deps.Deps`.

```go
package main

import (
	"fmt"

	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// myBackend holds your storage client.
type myBackend struct {
	store map[string][]byte
}

func createMyCustomDeps() deps.Deps {
	return &myBackend{store: map[string][]byte{}}
}

func (b *myBackend) Write(key string, value []byte) error {
	b.store[key] = value
	return nil
}

func (b *myBackend) Read(key string) ([]byte, error) {
	value, found := b.store[key]
	if !found {
		return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
	}
	return value, nil
}

// ...implement the remaining methods...

func main() {
	keep := lib.New(createMyCustomDeps())
	_ = keep
}
```

---

## Overriding Dependencies

You can change one behavior of an existing adapter — logging, a guard, a metric — without writing a full backend from scratch. Since `Deps` is an interface, this is done by **embedding**: the outer type inherits every method of the embedded value, and the methods you declare shadow the originals.

### Workflow

1. Declare a struct embedding `deps.Deps`.
2. Declare only the methods you want to change; everything else is inherited.
3. Wrap an adapter's result in it and inject that.

```go
package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// loggedDeps inherits every method of the embedded Deps...
type loggedDeps struct {
	deps.Deps
}

// ...and shadows just this one.
func (l loggedDeps) Write(key string, value []byte) error {
	fmt.Println("writing key:", key) // custom logging on every write
	return l.Deps.Write(key, value)  // delegate to the embedded adapter
}

func main() {
	myDeps := loggedDeps{Deps: standard.New()}

	keep := lib.New(myDeps)
	_ = keep
}
```
