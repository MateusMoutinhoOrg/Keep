# LibUsage

`Keep` is a Go module before it is anything else: every feature lives in `sandbox/`
and is reachable from any Go program that imports it.

```bash
go get github.com/MateusMoutinhoOrg/Keep@latest
```

## Wiring

`sandbox/` performs no OS effects of its own — filesystem, clock, stdout, processes all
arrive through a `deps.Deps` struct. `adapters/availables/standard` builds the ready-made
assembly, and `sandbox.New` turns it into the API object, which carries the deps on
`Sandbox.Deps` — so everything inside reaches them through the api it was handed.

```go
package main

import (
	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
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
deps.Std.Printf = func(f string, a ...any) (int, error) {
	return fmt.Fprintf(&out, f, a...)
}

lib := sandbox.New(&deps)
```

The contracts available to patch:

| Field | Contract package |
| --- | --- |
| `deps.Hashdeps` | `sandbox/deps/hashdeps` |
| `deps.Std` | `sandbox/deps/std` |
| `deps.Storagedeps` | `sandbox/deps/storagedeps` |
| `deps.Stringsdeps` | `sandbox/deps/stringsdeps` |

Each one is filled by a matching implementation under `adapters/libs/`, every package
exposing the same `Bind(deps *deps.Deps)` entry point:

| Adapter lib | Binder |
| --- | --- |
| `adapters/libs/filestorage` | `filestorage.Bind(&deps)` |
| `adapters/libs/hashdeps` | `hashdeps.Bind(&deps)` |
| `adapters/libs/memstorage` | `memstorage.Bind(&deps)` |
| `adapters/libs/std` | `std.Bind(&deps)` |
| `adapters/libs/stringsdeps` | `stringsdeps.Bind(&deps)` |

Starting from `standard.New()` is the safe default: an unfilled field is a nil func that
panics on first call. For a permanent mix, write your own
`adapters/availables/<name>/new.go` binding only the libs you want — `standard/new.go` is
regenerated on every build, while other directories under `availables/` are left alone.

`sandbox/api` is pure contract and `sandbox/` never touches the OS, so both are safe to import
anywhere; the rest of the rules a caller can count on are in [Rules](../Rules/doc.md#layers),
and [DepList](../DepList/doc.md) lists every contract that can be added.
