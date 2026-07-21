# Dependency Injection

Keep never talks to storage directly. Every read and write goes through the `deps.Deps` struct of injectable functions, so you choose where the bytes live: one of the built-in adapters, your own implementation, or an adapter with some functions overwritten.

```go
deps := standard.New()      // pick a backend...
keep := keep_lib.New(deps)  // ...and inject it
```

---

## The contract

`Deps` is a struct whose fields are functions, defined in [pkg/deps/deps.go](../../pkg/deps/deps.go) and documented field-by-field in [Required API](RequiredApi.md). It is a plain key-value API: write, read, delete, exists, plus conditional writes and locks.

Because the fields are plain function values, a backend is just a populated struct — there is no interface to implement, and any individual function can be replaced after construction (see [Overwriting Dependencies](OverwritingDeps.md)).

---

## Built-in adapters

Adapters live in [adapters/](../../adapters/):

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

---

## Writing your own backend

Build a `deps.Deps` with your own functions and pass it to `keep_lib.New` — the full workflow is in [Creating Custom Deps](CreatingCustomDeps.md). Two rules matter:

1. **Return the sentinel errors.** The database layer distinguishes "key doesn't exist" from "storage is broken" using the sentinels in `pkg/deps` (`ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, `ErrKeyLocked`). Wrap them with `fmt.Errorf("%w: ...", ...)` so `errors.Is` matches.
2. **Only per-key operations are needed.** Keep never lists keys or scans prefixes, so any store with `get`/`set`/`delete` can back it.

The [native adapter](../../adapters/native/native.go) (~150 lines) is the best starting point to copy from.
