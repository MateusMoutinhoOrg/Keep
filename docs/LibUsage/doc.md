# LibUsage

`Keep` is a Go module before it is anything else: every feature lives in `sandbox/`
and is reachable from any Go program that imports it.

```bash
go get github.com/MateusMoutinhoOrg/Keep@latest
```

## Wiring

`sandbox/` performs no OS effects of its own — filesystem, clock, stdout, processes all
arrive through a `deps.Deps` struct. `adapters/bindings/standard` builds the ready-made
assembly, and `sandbox.New` turns it into the API object, which carries the deps on
`Sandbox.Deps` — so everything inside reaches them through the api it was handed.

```go
package main

import (
	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
)

func main() {
	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox

	_ = lib
}
```

## What the sandbox exposes

`*api.Sandbox` is a flat struct, one field per contract declared in `sandbox/api/`.
Everything callable from Go is behind one of them.

| Field | Type |
| --- | --- |
| `lib.Config` | `api.Config` |
| `lib.Databases` | `api.Databases` |
| `lib.Info` | `api.Info` |

[PublicApi](../PublicApi/doc.md) lists every one of them — signatures, props structs and
dependency contracts — generated from `sandbox/api/` itself on every build.

## Custom deps

Every sub-contract is a struct of function fields, so any of them can be swapped for a
test double, an in-memory implementation or an instrumented wrapper. Patch fields **before**
`sandbox.New(&deps)`: the constructors capture the pointer.

```go
deps := standard.New()

var out bytes.Buffer
deps.StdDeps.Printf = func(f string, a ...any) (int, error) {
	return fmt.Fprintf(&out, f, a...)
}

lib := sandbox.New(&deps)
```

The contracts available to patch:

| Field | Contract package |
| --- | --- |
| `deps.HashDeps` | `sandbox/deps/hashdeps` |
| `deps.StdDeps` | `sandbox/deps/stddeps` |
| `deps.StorageDeps` | `sandbox/deps/storagedeps` |
| `deps.StringsDeps` | `sandbox/deps/stringsdeps` |

Each one is filled by a matching implementation under `adapters/impls/`, every package
exposing the same `Bind(deps *deps.Deps)` entry point:

| Adapter lib | Binder |
| --- | --- |
| `adapters/impls/filestorage` | `filestorage.Bind(&deps)` |
| `adapters/impls/memstorage` | `memstorage.Bind(&deps)` |
| `adapters/impls/osstd` | `osstd.Bind(&deps)` |
| `adapters/impls/sha256hash` | `sha256hash.Bind(&deps)` |
| `adapters/impls/stdstrings` | `stdstrings.Bind(&deps)` |

Starting from `standard.New()` is the safe default: an unfilled field is a nil func that
panics on first call. For a permanent mix, write your own
`adapters/bindings/<name>/new.go` binding only the libs you want — `standard/new.go` is
regenerated on every build, while other directories under `bindings/` are left alone.

`sandbox/api` is pure contract and `sandbox/` never touches the OS, so both are safe to import
anywhere; the rest of the rules a caller can count on are in [Rules](../Rules/doc.md#layers),
and [DepList](../DepList/doc.md) lists every contract that can be added.
