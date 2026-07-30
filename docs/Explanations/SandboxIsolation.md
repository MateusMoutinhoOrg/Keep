# Sandbox Isolation

## Description
Explains why the database engine lives in `sandbox/` and what "closed sandbox" means in practice: the engine reaches nothing outside itself — no adapter, no third-party module, no OS-bound standard-library package — so everything it can do is exactly what the injected `Deps` allows.

---

## The Four Trees

The project is split into four top-level directories, and the arrows only point one way:

```
adapters/  ──▶  sandbox/  ◀──  examples/ , tests/
(reaches the OS)  (closed)     (wire the two together)
```

- `sandbox/` is the engine. It is closed: it imports only itself and OS-independent standard-library packages (`errors`, `fmt`, `strconv`, `strings`, `crypto/sha256`, …).
- `adapters/` is outside the wall. It is the only place `os`, `net`, a database driver, or any third-party module may appear.
- `examples/` and `tests/` are outside the wall too, and are the only places an adapter and the sandbox are named in the same file.

This split is what makes Keep **storage-independent**. The engine cannot be affected by which operating system, filesystem, or network the program runs on, because it has no way to reach any of them — it only knows how to read, write, and delete a single key.

---

## What the Wall Forbids

A file under `sandbox/` may not import:

| Forbidden | Why |
|-----------|-----|
| `adapters/…` | The engine would bind itself to one concrete backend, and injection would be pointless. |
| `examples/…`, `tests/…` | Consumers of the library are never part of it. |
| Any third-party module | A dependency the caller cannot replace is a dependency the caller cannot test around. |
| OS-bound stdlib (`os`, `net`, `os/exec`, `syscall`, …) | The effect belongs in an adapter, reached through a `Deps` method. |

Everything the engine needs from the outside world is declared as a method on `Deps`:

```go
// sandbox/contracts/deps/deps.go — the only door in the wall
type Deps interface {
	Write(key string, value []byte) error   // instead of os.WriteFile
	Read(key string) ([]byte, error)        // instead of os.ReadFile
	Exists(key string) (bool, error)        // instead of os.Stat
	Delete(key string) error                // instead of os.Remove
	// … one method per remaining requirement
}
```

Inside the sandbox, the same behaviors are reached only through the object's `Deps` field:

```go
// sandbox/internal/schemaitem/schemaitem.go — no os, no net, no third party
func (s *SchemaItem) CheckKeysPresence(keys []string) bool {
	for _, key := range keys {
		exists, err := s.Deps.Exists(dense.ValueKey(s.Prefix, s.RecordID, key))
		if err != nil || !exists {
			return false
		}
	}
	return true
}
```

To add a new door, follow [AddDependency.md](/docs/Tutorials/AddDependency.md).

---

## What the Wall Forbids in the Other Direction

The wall is not only about what the sandbox imports — it also limits what the outside may reach into. `sandbox/internal/` is protected by Go's `internal/` rule: only packages rooted at `sandbox/` can import it. An adapter or a consumer that tries gets a compile error, not a convention warning.

So the outside world sees exactly three packages:

| Package | Who imports it | For what |
|---------|----------------|----------|
| `sandbox` (package `lib`) | consumers, examples, tests | `lib.New(deps) api.Lib` — the single wiring point |
| `sandbox/contracts/deps` | adapters, consumers | the contract to implement |
| `sandbox/contracts/api` | consumers, examples, tests | every interface and constant the library exchanges |

Everything else in `sandbox/` is unreachable, which is why `KeepDatabase`, `SchemaInstance`, `SchemaItem`, and the dense-record helpers can be renamed or restructured without breaking a single consumer.

---

## Why the Entry Point Lives Inside

`sandbox/new.go` is the one file in the sandbox that consumers import directly, and it stays inside the wall because it obeys the same rule — it names no adapter:

```go
// sandbox/new.go
func New(d deps.Deps) api.Lib {
	return &internallib.Lib{Deps: d}
}
```

It accepts an interface and returns an interface. The caller decides which implementation flows in, so the engine never learns what is behind the contract:

```go
import (
	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

// This line is in examples/, outside the wall — the only place
// an adapter and the sandbox meet.
keep := lib.New(standard.New())
```

For how the injected value then travels through the object graph, see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).
