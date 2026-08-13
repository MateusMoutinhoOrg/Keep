# Contribution Rules

Rules to follow when contributing to this project. Every file must also be shaped by the specification that governs it — locate it in [Specs.md](/docs/References/Specs.md).

---

## Tutorials Guide
Before making anything, read the [README.md](/README.md), open the theme index matching your goal (`docs/Index/<Theme>.md`) and search for a tutorial about what you want to do. If there is one, follow it; if there isn't, you need to create one following the spec defined in [TutorialDocs](./Specs/TutorialDocs/).


## Specification Compliance

Before creating or editing any file, read [Specs.md](/docs/References/Specs.md) and check whether the file matches an **Applies To** entry. If it does, create or edit it following the specification that entry points to — reproduce the shape it requires, using its `sample` as reference.

---

## Sandbox Isolation

[sandbox/](/sandbox/) is a closed sandbox. No file inside it may import [adapters/](/adapters/), [examples/libraryExamples/](/examples/libraryExamples/), `tests/`, a third-party module, or an OS-bound standard-library package (`os`, `net`, `os/exec`, `syscall`, …). Every such effect must be declared as a function field on the `Deps` contract and reached through the object's `Deps` field, following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency). The mechanic is explained in [SandboxIsolation.md](/docs/References/SandboxIsolation.md).

Contracts are **structs of function fields**, never interfaces — in `sandbox/contracts/deps` and `sandbox/contracts/api` alike. Every type in the project is declared in `sandbox/contracts/`; `sandbox/lib/` declares no types at all, only factories. See [StructContracts.md](/docs/References/StructContracts.md).

Conversely, nothing outside the sandbox may reach into it beyond its three public packages: `sandbox` (package `lib`), `sandbox/contracts/deps`, and `sandbox/contracts/api`. `sandbox/lib/` and `sandbox/config/` are private to the sandbox. The compiler does not enforce that half — the tree carries no `internal/` path element — so it is on the reviewer: an import of `sandbox/lib/...` from `adapters/`, `examples/libraryExamples/`, `tests/`, or a consuming project is a rejected change, not a warning.

---

## Factory Pattern

Every struct of function fields in this project — an `api` struct inside the sandbox, a `deps.Deps` filled by an adapter outside it — is filled by **factories**, never by methods bound into fields and never by an internal mirror type.

A factory takes a pointer to the **carrier** — the struct holding the state the closure reads — and returns exactly one field's value; the caller assigns it:

```go
// sandbox/lib/schemaitem/ — the carrier is the api struct being filled
func CheckKeysPresenceFactory(s *api.SchemaItem) func(keys []string) bool {
	return func(keys []string) bool {
		exists, err := s.Deps.Exists(dense.ValueKey(s.Prefix, s.Id, "email"))
		return err == nil && exists
	}
}

func build(d deps.Deps, items []api.Item, prefix string, id int64) api.SchemaItem {
	s := api.SchemaItem{Deps: d, Items: items, Prefix: prefix, Id: id}
	s.CheckKeysPresence = CheckKeysPresenceFactory(&s)
	return s
}

// adapters/<name>/ — the carrier is the adapter, whose Deps field is the contract
func ReadFactory(s *StandardAdapter) func(key string) ([]byte, error) {
	return func(key string) ([]byte, error) { /* ... */ }
}

func New() deps.Deps {
	s := &StandardAdapter{}
	s.Deps.Read = ReadFactory(s)
	return s.Deps
}
```

Four rules follow, and none of them is checked by the compiler:

- Every `api` struct whose behavior needs dependencies declares a `Deps deps.Deps` field, and closures reach dependencies through it — `s.Deps.<Field>(...)`, read inside the closure, never captured at factory time. Every adapter struct declares the same field, as the contract its factories fill.
- Every field factory must be called and its return value assigned from its package's constructor (`New`, or a shared aggregate like `build`), which is the factory aggregate — there is no separate `Factory` function. A field no factory fills stays nil and panics on first call.
- A constructor returns the filled **contract struct** by value — `api.Lib`, `api.<Object>`, or `adapter.Deps` — never the carrier type of an adapter.
- The `Deps` field is **read-only once the struct is returned**: closures capture the struct the factories ran over, so a caller patching `Deps` on a copy changes nothing. Patch the `deps.Deps` value before calling `lib.New`.

See [StructContracts.md](/docs/References/StructContracts.md) for the full explanation, including why lookups that can fail return `(value, bool)` instead of a nil struct.

---

## Import Aliases

Any file that **consumes** the library from outside it — [examples/libraryExamples/](/examples/libraryExamples/), `tests/`, and third-party consumers — imports it under `keep`-prefixed aliases, so each call site says which layer it belongs to:

| Import | Alias |
|--------|-------|
| `adapters/<name>` | `keepadapter` |
| `sandbox` | `keeplib` |
| `sandbox/contracts/api` | `keeptypes` |
| `sandbox/contracts/deps` | `keepdeps` |

```go
import (
	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)
```

Files that belong to the library — everything under `sandbox/` and its own [adapters/](/adapters/) — keep the plain package names (`api`, `deps`): there the prefix would be noise, since the import is already local.

---

## Adapter Changes

When you create, delete, or rename an adapter inside [adapters/](/adapters/), update the table in [Adapters.md](/docs/References/Adapters.md) in the same commit. When you add a field to the `Deps` contract, every adapter must fill it with a factory — a missing field compiles but panics at runtime on first call, since the contract is a struct, not an interface.

---

## File Changes

Before creating, deleting, or renaming any file or directory, read [Structure.md](/docs/References/Structure.md) and check whether the change affects the project structure. If it does, update [Structure.md](/docs/References/Structure.md) in the same commit.

---

## Specification Changes

When you create, delete, or rename a specification inside [Specs/](./Specs), you MUST adapt all the files that match the spec's Applies To rule, and update the index in [Specs.md](/docs/References/Specs.md).

---

## Documentation Changes

When you create, delete, or rename a `.md` file, update the theme index that lists it under [docs/Index/](/docs/Index/) and, when it is a new structural component, [Structure.md](/docs/References/Structure.md) — following [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md). The [README.md](/README.md) links to theme indexes only, so it changes when a **theme** is added, renamed, or removed, never for a single page.

---

## Sample Changes

When you create, delete, or rename a sample (any file inside [examples/libraryExamples/](/examples/libraryExamples/)), update [ApiSamplesList.md](/docs/References/ApiSamplesList.md) — following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample).
