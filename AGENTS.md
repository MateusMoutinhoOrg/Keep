# AGENTS.md

## What this is

Keep is a **storage-independent database**: schemas with typed fields, unique indexed keys
and nested collections, over any backend that can read, write and delete a single key — no
listing, no prefix scans, no range queries.

## This repository is generated

It is an [agnos](https://github.com/MateusMoutinhoOrg/Agnos) project. **agnos owns every
generated file**, and a hand edit to one is lost on the next build. Before changing
anything, read the three pages that govern it — they are generated too, and they are the
authority over anything written here:

| Read | For |
|---|---|
| [docs/Rules/doc.md](docs/Rules/doc.md) | every rule, including the ones `agnos verify` enforces |
| [docs/Workflow/doc.md](docs/Workflow/doc.md) | the agnos command that makes each kind of change |
| [docs/Structure/doc.md](docs/Structure/doc.md) | what lives where |
| [docs/GeneratedFiles/doc.md](docs/GeneratedFiles/doc.md) | which files `build` rewrites |

```bash
agnos build        # verify + regenerate everything + go mod tidy + compile
agnos verify       # the schema check alone, writes nothing
agnos exec-test    # run every example and check it against its golden
```

Run `agnos build` after every hand edit. It is idempotent.

Never create a doc, an example or an adapter by hand: `agnos add-doc`, `agnos
add-lib-example`, `agnos add-dep` / `add-adapter` / `add-available` scaffold them and
rewrite the indexes that list them.

## The four hand-written places

Everything else is regenerated over.

| File | Holds |
|---|---|
| `sandbox/api/<x>.go` | one contract — the file name is the `api.Sandbox` field, and the type of that name must be declared in it |
| `sandbox/internal/<x>/new.go` | `New<X>(sandbox *api.Sandbox) api.<X>`, plus the implementation beside it |
| `sandbox/deps/<x>/<x>.go` | one capability the sandbox needs from outside. **Imports nothing at all** |
| `adapters/libs/<x>/<x>.go` | `Bind(deps *deps.Deps)`, beside its `adapter.yaml`. The only place OS-bound code may live |

## Keep-specific things that are easy to get wrong

- **The sandbox may not import the standard library.** `fmt` is `sandbox.Deps.Std.Sprintf`,
  `strconv`/`strings` are `sandbox.Deps.Stringsdeps`, `crypto/sha256` is
  `sandbox.Deps.Hashdeps`, and storage is `sandbox.Deps.Storagedeps`. A new capability is a
  new contract under `sandbox/deps/`, filled by every available.
- **`sandbox/deps/storagedeps` reports no expected condition as an error.** Absent is
  `found == false`, a conditional write that did not apply is `written == false`. That is why
  it needs no sentinel values and so no import. See
  [docs/StorageContract](docs/StorageContract/doc.md).
- **The write orderings are load-bearing.** An insert commits on its last write; an update
  to a `Key` field writes the new index entry before it moves the value. Changing
  `sandbox/internal/dense` or `sandbox/internal/schemaitem` means preserving every invariant
  of [docs/DenseRecordPattern](docs/DenseRecordPattern/doc.md).
- **No api type carries a `Deps` field.** `Sandbox.Deps` is the only one, and it is the one
  field that does not cross into a consumer. A record reaches storage through the closure it
  was built with, which is what keeps every api type convertible and Keep installable as a
  dep of another agnos repo.
- **Completeness is unchecked by the compiler.** A function field no factory fills is nil
  and panics on first call, so every `<Field>Factory` must be called from the `New` that
  builds its object.
- **A release bump is `version:` in `AgnosConfig/project.yaml`**, then `agnos build`. It
  regenerates `sandbox/internal/config/config.go`, which `api.Info.Version` reports.

## Adding an api surface

Two files, then `agnos build` writes the wiring:

1. `sandbox/api/<x>.go` — `type <X> struct { … }` of function fields, every declaration
   doc-commented (`verify` fails without it, and `docs/PublicApi` is generated from those
   comments).
2. `sandbox/internal/<x>/new.go` — `func New<X>(sandbox *api.Sandbox) api.<X>`, assigning
   each field from its factory.

An example goes with it: `agnos add-lib-example <name>`, write `example.go`, add a
`props.yaml` with a one-line `description:`, then `agnos exec-test` to write the golden. An
example that copies nothing out of `TestDir` into `AssertDir` fails.
