# Sandbox Isolation

## Description
Explains why the database engine lives in `sandbox/` and what "closed sandbox" means in practice: the engine reaches nothing outside itself — no adapter, no third-party module, no OS-bound standard-library package — so everything it can do is exactly what the injected `Deps` allows.

---

## The Four Trees

The project is split into four top-level directories, and the arrows only point one way:

```
adapters/  ──▶  sandbox/  ◀──  examples/libraryExamples/ , tests/
(reaches the OS)  (closed)     (wire the two together)
```

- `sandbox/` is the engine. It is closed: it imports only itself and OS-independent standard-library packages (`errors`, `fmt`, `strconv`, `strings`, `crypto/sha256`, …).
- `adapters/` is outside the wall. It is the only place `os`, `net`, a database driver, or any third-party module may appear.
- `examples/libraryExamples/` and `tests/` are outside the wall too, and are the only places an adapter and the sandbox are named in the same file.

This split is what makes Keep **storage-independent**. The engine cannot be affected by which operating system, filesystem, or network the program runs on, because it has no way to reach any of them — it only knows how to read, write, and delete a single key.

---

## What the Wall Forbids

A file under `sandbox/` may not import:

| Forbidden | Why |
|-----------|-----|
| `adapters/…` | The engine would bind itself to one concrete backend, and injection would be pointless. |
| `examples/libraryExamples/…`, `tests/…` | Consumers of the library are never part of it. |
| Any third-party module | A dependency the caller cannot replace is a dependency the caller cannot test around. |
| OS-bound stdlib (`os`, `net`, `os/exec`, `syscall`, …) | The effect belongs in an adapter, reached through a `Deps` field. |

Everything the engine needs from the outside world is declared as a function field on `Deps`:

```go
// sandbox/contracts/deps/deps.go — the only door in the wall
type Deps struct {
	Write  func(key string, value []byte) error // instead of os.WriteFile
	Read   func(key string) ([]byte, error)      // instead of os.ReadFile
	Exists func(key string) (bool, error)        // instead of os.Stat
	Delete func(key string) error                // instead of os.Remove
	// … one field per remaining requirement
}
```

Inside the sandbox, the same behaviors are reached only through the object's `Deps` field, from inside a factory's closure:

```go
// sandbox/lib/schemaitem/schemaitem.go — no os, no net, no third party
func CheckKeysPresenceFactory(s *api.SchemaItem) func(keys []string) bool {
	return func(keys []string) bool {
		for _, key := range keys {
			exists, err := s.Deps.Exists(dense.ValueKey(s.Prefix, s.Id, key))
			if err != nil || !exists {
				return false
			}
		}
		return true
	}
}
```

The wall does not reach the library's *words*: the version string, and any other fixed text the engine reports, are compile-time constants in `sandbox/config`, which is inside the sandbox. Text is data, not an effect — writing it out is the effect, and that goes through a `Deps` field like everything else.

To add a new door, follow [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency).

---

## What the Wall Forbids in the Other Direction

The wall is not only about what the sandbox imports — it also limits what the outside may reach into. `sandbox/lib/` holds the factories that fill the contract structs, and it is private to the sandbox: nothing in `adapters/`, `examples/libraryExamples/`, `tests/`, or a consuming project may import it. Unlike the import rules above, this one is a convention the reviewer enforces, not the compiler — the tree carries no `internal/` path element, so a stray import compiles. Treat it as binding all the same, and see [RULES.md](/docs/References/RULES.md#sandbox-isolation).

So the outside world sees exactly three packages:

| Package | Who imports it | For what |
|---------|----------------|----------|
| `sandbox` (package `lib`) | consumers, examples, tests | `lib.New(deps) api.Lib` — the single wiring point |
| `sandbox/contracts/deps` | adapters, consumers | the contract to fill |
| `sandbox/contracts/api` | consumers, examples, tests | every type and constant the library exchanges |

Everything else in `sandbox/` is off limits, which is why `KeepDatabase`, `SchemaInstance`, `SchemaItem`, and the dense-record helpers can be renamed or restructured without breaking a single consumer.

---

## Why the Entry Point Lives Inside

`sandbox/new.go` is the one file in the sandbox that consumers import directly, and it stays inside the wall because it obeys the same rule — it names no adapter:

```go
// sandbox/new.go
func New(d deps.Deps) api.Lib {
	return lib.New(d)
}
```

It accepts a filled `deps.Deps` struct and returns the filled `api.Lib` struct. The caller decides which implementation flows in, so the engine never learns what is behind the contract:

```go
import (
	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

// This line is in examples/libraryExamples/, outside the wall — the only place
// an adapter and the sandbox meet.
keep := keeplib.New(keepadapter.New())
```

For how the injected value then travels through the object graph, see [DepsMechanic.md](/docs/References/DepsMechanic.md).
