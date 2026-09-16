# Workflow

Every change this project takes and the command that makes it. `agnos` owns every generated
file; what stays hand-written is listed in [GeneratedFiles](../GeneratedFiles/doc.md), and the
rules each recipe holds to are in [Rules](../Rules/doc.md).

## The loop

```bash
agnos build      # verify + regenerate every generated file + go mod tidy + compile
agnos verify     # the schema check alone, writes nothing
```

`build` is the only thing that regenerates the wiring,
`README.md` and `docs/`, so run it after every hand edit. It is idempotent: a second run leaves
the tree unchanged. Every command below takes `--path <dir>` (default `.`) and `-q`, and runs
`build` for you.

No recipe below asks for a Go file to be created by hand except the cases listed under
[Hand-written code](#hand-written-code).

## Choose what agnos generates

```bash
agnos list-extensions             # every generation mechanic and whether it is on
agnos enable-extension readme     # start generating README.md again
agnos disable-extension doc       # stop generating docs/, keep what is there
```

`AgnosConfig/extensions.yaml` is what `build` reads to decide what to render. Turning a
mechanic off stops the generation and removes nothing: the files stay, and they are yours to
edit. Every key is in [Extensions](../Extensions/doc.md).

## Add the CLI layer

```bash
agnos cli-init     # sandbox/internal/cli, cmd/main, the help and version commands, argvdeps + std
```

From there `agnos add-command <name> --help "..." --category "..."` declares a command and
`agnos add-flag` / `add-arg` its fields. `agnos cli-purge` removes the layer again.


## Add the server layer

```bash
agnos server-init      # serverdeps, sandbox/internal/server, the health route, start-server
Keep start-server  # listens on :8080
```

From there `agnos add-route <name> --trigger /<path> --help "..." --category "..."` declares a
route and `agnos add-segment` / `add-header` / `add-param` / `add-body-field` its fields. A
project with no CLI gets one first: a server needs a command that starts it.
`agnos server-purge` removes the layer again.


## Add the front layer

```bash
agnos front-init                  # pageio, the static route, assets/frontend/
agnos add-page home --trigger /   # a page answering GET /
Keep start-server             # serves it
```

From there `agnos add-page <name>` declares a page and `remove-page` drops it,
html included. A project with no server layer gets one first: a page is answered over http.
`agnos front-purge` removes the layer again, leaving `assets/frontend/` alone.
## Add reusable logic

`sandbox/internal/<pkg>/`, one directory per concern, imported by whatever needs it. No
declaration, no generated counterpart — write the package and run `build`.

## Add a surface to the sandbox api

The api is what a Go caller gets back from `sandbox.New` (see [LibUsage](../LibUsage/doc.md)).
Two hand-written places, then `build` regenerates `sandbox/api/sandbox.go` and
`sandbox/new.go` around them:

1. `sandbox/api/<x>.go` — the contract: `type <X> struct { ... }` of function fields, named
   after the file, every declaration doc-commented (those comments render
   [PublicApi](../PublicApi/doc.md)). It becomes the `api.Sandbox` field `<X>`.
2. `sandbox/internal/<x>/new.go` — `func New<X>(sandbox *api.Sandbox) api.<X>`, assigning each
   field of the contract, with the implementation beside it. `sandbox/new.go` calls it as
   `self.<X> = <x>.New<X>(&self)`.

## Add a dependency

Everything the sandbox is not allowed to do itself — filesystem, clock, network, subprocess —
arrives through `sandbox.Deps`. Install a ready-made one:

```bash
agnos list-deps                 # every installable contract
agnos add-dep <dep>             # sandbox/deps/<dep>/ + its default adapter + the go.mod require
agnos add-dep <dep> --adapter <adapter>
agnos remove-dep <dep> [--with-adapters]
```

[DepList](../DepList/doc.md) is the catalogue. For one of your own, write the two halves:

1. `sandbox/deps/<x>/<x>.go` — `type Sandbox struct { ... }` of function fields, no import at all.
2. `adapters/libs/<x>/<x>.go` — `func Bind(deps *deps.Deps) { deps.<X> = <x>.Sandbox{...} }`, any
   import allowed, beside an `adapter.yaml` saying `dep: <x>`.

Then bind it: add `<x>` to `adapters/availables/standard/available.yaml`, or let
`agnos add-dep` do both for a dep of the catalogue. Reach it as `sandbox.Deps.<X>`
from anywhere inside `sandbox/`.

One contract may have several adapters — see [Adapters](../Adapters/doc.md).

## Add a doc

```bash
agnos add-doc <Name> --theme <id> --description "one line"    # themes: AgnosConfig/themes.yaml
agnos add-doc <Name>/<Sub> --description "one line"           # sub-doc, no theme
agnos remove-doc <Name>
```

Write `docs/<Name>/doc.md`; `README.md`'s index, and the parent `Index.md` of a sub-doc, are
regenerated. Describe any new path worth naming in `AgnosConfig/structure.yaml` — that
file is what renders [Structure](../Structure/doc.md).

## Add an example

```bash
agnos add-lib-example <name>       # examples/lib/<name>/example.go
agnos exec-test                    # run them all, check each against its golden
agnos exec-test --only <name>      # one example, both sides
agnos update-test <name>           # rewrite that one golden with what it produces now
agnos exec-test --update           # rewrite every golden at once
agnos remove-lib-example <name>
```

Write the example itself, ending with the copy out of `TestDir` into `AssertDir` that says what
it asserts: `result.yaml` records `AssertDir`, and an example that copies nothing out fails.
The golden is written by the first `exec-test` and refreshed with `update-test <name>`, which
prints what it changes before writing. Details in [LibExamples](../LibExamples/doc.md).

## Hand-written code

| File | Written when |
| --- | --- |
| `sandbox/internal/<pkg>/*.go` | logic worth reusing |
| `sandbox/api/<x>.go` + `sandbox/internal/<x>/new.go` | a new api surface |
| `sandbox/deps/<x>/<x>.go` + `adapters/libs/<x>/<x>.go` + its `adapter.yaml` | a new dependency |

Everything else is regenerated over. Two more files are yours: `AgnosConfig/docs/ReadmeHeader.md`
is the whole of `README.md` above the documentation index, and `LICENSE` is pasted verbatim into
its License section — put whatever license you want there.

## Ship

This project has no `cmd/main` to compile: it ships as the Go module other programs import
(see [LibUsage](../LibUsage/doc.md)). Bump `version` in `AgnosConfig/project.yaml` and tag
the repository; `agnos cli-init` adds a binary if you want one.
