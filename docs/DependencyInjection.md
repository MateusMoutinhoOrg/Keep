# Dependency Injection

Keep never talks to storage directly. Every read and write goes through the `deps.Deps` interface, so you choose where the bytes live: one of the built-in adapters, or your own implementation.

```go
deps := standard.New()      // pick a backend...
keep := keep_lib.New(deps)  // ...and inject it
```

## The interface

The full contract is defined in [pkg/deps/deps.go](../pkg/deps/deps.go) and documented method-by-method in [Required API](RequiredApi.md). It is a plain key-value API: write, read, delete, exists, plus conditional writes and locks.

## Built-in adapters

Adapters live in [adapters/](../adapters/):

### `standard` — filesystem (the default choice)

Each key becomes a file. Data survives restarts. Works in most scenarios.

```go
import keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"

deps := keep_deps.New()               // relative to the working directory
deps := keep_deps.NewWithBase("/srv") // under a specific directory
```

### `native` — in-memory

Data lives only for the lifetime of the process. Ideal for tests and prototypes.

```go
import keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"

deps := keep_deps.New()
```

## Writing your own backend

Implement every method of `deps.Deps` and pass your type to `keep_lib.New`. Two rules matter:

1. **Return the sentinel errors.** The database layer distinguishes "key doesn't exist" from "storage is broken" using the sentinels in `pkg/deps` (`ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, `ErrKeyLocked`). Wrap them with `fmt.Errorf("%w: ...", ...)` so `errors.Is` matches.
2. **Only per-key operations are needed.** Keep never lists keys or scans prefixes, so any store with `get`/`set`/`delete` can back it.

Minimal sketch:

```go
type myBackend struct{ /* your client */ }

func (m *myBackend) Read(key string) ([]byte, error) {
	value, found := /* fetch from your store */
	if !found {
		return nil, fmt.Errorf("%w: %s", deps.ErrKeyNotFound, key)
	}
	return value, nil
}

// ...implement the remaining methods of deps.Deps...

keep := keep_lib.New(&myBackend{})
```

The [native adapter](../adapters/native/native.go) (~130 lines) is the best starting point to copy from.
