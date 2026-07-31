# Dependency Mechanic

## Description
Keep never talks to storage directly. Every read and write goes through the `deps.Deps` struct of function fields, so you choose where the bytes live: one of the built-in adapters, your own implementation, or an adapter with some fields overridden.

This guide explains how this injection works, how to write your own backend, and how to override specific dependencies. The struct-of-function-fields shape itself is explained in [StructContracts.md](/docs/Explanations/StructContracts.md).

---

## Dependency Injection

The `Deps` struct is defined in `sandbox/contracts/deps/deps.go` and its fields are documented in [Required API](/docs/References/RequiredApi.md). It is a plain key-value API: write, read, delete, exists, plus conditional writes and locks.

A backend is any value with every field of that struct filled — nothing checks completeness at compile time, unlike an interface, so a half-finished adapter builds fine and panics on first call to the missing field.

To inject a built-in adapter:

```go
import (
	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard" // filesystem adapter
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

func main() {
	deps := keepadapter.New() // pick a backend...
	keep := keeplib.New(deps) // ...and inject it
	_ = keep
}
```

The value flows one way and never comes back out: `lib.New` stores it on the internal `Lib`, which copies it into every `KeepDatabase` it creates, which copies it into every `SchemaInstance`, which copies it into every `SchemaItem`. Every object therefore reaches storage through the exact backend the caller chose, and no object anywhere inside `sandbox/` can reach a different one — see [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).

The shipped adapters and when to use each are listed in [Adapters.md](/docs/References/Adapters.md).

---

## Creating Custom Dependencies

You can run Keep over any storage backend by filling every field of `deps.Deps` with a factory of your own.

### Rules
- **Return the sentinel errors**: The database layer distinguishes "key doesn't exist" from "storage is broken" using the sentinels in `sandbox/contracts/deps` (`ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, `ErrKeyLocked`). Wrap them with `fmt.Errorf("%w: ...", ...)` so `errors.Is` matches.
- **Only per-key operations are needed**: Keep never lists keys or scans prefixes, so any store with `get`/`set`/`delete` can back it.
- **Fill the whole struct**: a field no factory assigns stays `nil` and panics on first call — nothing checks this for you.

### Workflow

1. Declare a struct holding your storage client, leading with a `Deps deps.Deps` field.
2. Write a `<Field>Factory` per field of `deps.Deps`, following [Required API](/docs/References/RequiredApi.md).
3. Expose a `New` factory aggregate that assigns every factory's return value and returns the filled `deps.Deps`.

```go
package main

import (
	"fmt"

	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
)

// MyBackend holds your storage client — the carrier every factory closes over.
type MyBackend struct {
	Deps  deps.Deps
	store map[string][]byte
}

func WriteFactory(b *MyBackend) func(key string, value []byte) error {
	return func(key string, value []byte) error {
		b.store[key] = value
		return nil
	}
}

func ReadFactory(b *MyBackend) func(key string) ([]byte, error) {
	return func(key string) ([]byte, error) {
		value, found := b.store[key]
		if !found {
			return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
		}
		return value, nil
	}
}

// ...one factory per remaining field...

func NewMyCustomDeps() deps.Deps {
	b := &MyBackend{store: map[string][]byte{}}
	b.Deps.Write = WriteFactory(b)
	b.Deps.Read = ReadFactory(b)
	// ...one assignment per remaining field...
	return b.Deps
}

func main() {
	keep := keeplib.New(NewMyCustomDeps())
	_ = keep
}
```

---

## Overriding Dependencies

You can change one behavior of an existing adapter — logging, a guard, a metric — without writing a full backend from scratch. Since `Deps` is a struct of function fields, this is a plain **field assignment**: take the filled contract and replace the one field you want to change.

### Workflow

1. Build the `deps.Deps` you want to start from — usually an adapter's `New()`.
2. Assign a new closure to the field you want to change.
3. Inject the patched value, **before** calling `lib.New` — the field is read-only once the struct is returned (see [StructContracts.md](/docs/Explanations/StructContracts.md#what-it-costs)).

```go
package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

func main() {
	myDeps := keepadapter.New()

	// Keep every other field of the adapter; add logging to Write.
	inner := myDeps.Write
	myDeps.Write = func(key string, value []byte) error {
		fmt.Println("writing key:", key) // custom logging on every write
		return inner(key, value)         // delegate to the original closure
	}

	keep := keeplib.New(myDeps)
	_ = keep
}
```
