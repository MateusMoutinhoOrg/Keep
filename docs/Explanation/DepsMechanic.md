# Dependency Mechanic

## Description
Keep never talks to storage directly. Every read and write goes through the `deps.Deps` struct of injectable functions, so you choose where the bytes live: one of the built-in adapters, your own implementation, or an adapter with some functions overwritten.

This guide explains how this injection works, how to write your own backend, and how to overwrite specific dependencies.

---

## Dependency Injection

The `Deps` struct is defined in `pkg/deps/deps.go` and its fields are documented in [Required API](/docs/Reference/RequiredApi.md). It is a plain key-value API: write, read, delete, exists, plus conditional writes and locks.

Because the fields are plain function values, a backend is just a populated struct — there is no interface to implement.

To inject a built-in adapter:

```go
import (
	"github.com/MateusMoutinhoOrg/Keep/adapters/standard" // filesystem adapter
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

func main() {
	deps := standard.New()     // pick a backend...
	keep := lib.New(deps)  // ...and inject it
	_ = keep
}
```

---

## Creating Custom Dependencies

You can run Keep over any storage backend by providing your own implementation for each function field of `deps.Deps`.

### Rules
- **Return the sentinel errors**: The database layer distinguishes "key doesn't exist" from "storage is broken" using the sentinels in `pkg/deps` (`ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, `ErrKeyLocked`). Wrap them with `fmt.Errorf("%w: ...", ...)` so `errors.Is` matches.
- **Only per-key operations are needed**: Keep never lists keys or scans prefixes, so any store with `get`/`set`/`delete` can back it.

### Workflow

1. Create a function that constructs a `deps.Deps` object.
2. Provide your own implementation for each function field, following [Required API](/docs/Reference/RequiredApi.md).

```go
package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

func createMyCustomDeps() deps.Deps {
	store := map[string][]byte{} // your storage client goes here
	return deps.Deps{
		Write: func(key string, value []byte) error {
			store[key] = value
			return nil
		},
		Read: func(key string) ([]byte, error) {
			value, found := store[key]
			if !found {
				return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
			}
			return value, nil
		},
		// ...implement the remaining fields...
	}
}

func main() {
	myDeps := createMyCustomDeps()
	keep := lib.New(myDeps)
	_ = keep
}
```

---

## Overwriting Dependencies

You can overwrite specific function fields of an already populated `Deps` struct to alter behavior (like logging) without creating a full adapter from scratch.

### Workflow

1. Get a fully populated `deps.Deps` struct from an existing adapter.
2. Overwrite specific function fields in the struct to inject your custom behavior.

```go
package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

func main() {
	myDeps := standard.New()

	defaultWrite := myDeps.Write
	myDeps.Write = func(key string, value []byte) error {
		fmt.Println("writing key:", key) // custom logging on every write
		return defaultWrite(key, value)
	}

	keep := lib.New(myDeps)
	_ = keep
}
```
