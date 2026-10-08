# PublicApi

Every exported symbol of `github.com/MateusMoutinhoOrg/Keep`, read straight from the contract sources on
every build: `sandbox/api/` is the surface `sandbox.New` returns, `sandbox/deps/`
the contracts an adapter fills and a caller may replace. Each description is the
doc comment of the declaration itself — change the comment, run `build`, and the page
follows.

One page per contract: the tables below say which page declares a symbol, so open that page
rather than reading the whole surface.

## Entry points

| Symbol | Signature |
| --- | --- |
| `sandbox.New` | `func(deps *deps.Deps) *api.Sandbox` |
| `standard.New` | `func() deps.Deps` (`adapters/bindings/standard`) |

Implementations live under `sandbox/internal` and are unreachable: every contract is a
struct of function fields, filled by a binder.

## The sandbox api

| Page | Declares |
| --- | --- |
| [`sandbox/api/sandbox.go`](api.sandbox.md) | `Sandbox` |
| [`sandbox/api/config.go`](api.config.md) | `Config` |
| [`sandbox/api/databases.go`](api.databases.md) | `Key`, `Int`, `Database`, `Float`, `String`, `Link`, `Bytes`, `KeyConflict`, `NotFound`, `MissingField`, `InvalidField`, `Internal`, `Item`, `Schema`, `Props`, `Error`, `SchemaItem`, `SchemaInstance`, `DatabaseHandle`, `Databases` |
| [`sandbox/api/info.go`](api.info.md) | `Info` |
| [`sandbox/api/projectconfig.go`](api.projectconfig.md) | `ProjectConfig` |
| [`sandbox/api/projectsandbox.go`](api.projectsandbox.md) | `ProjectSandbox` |

## Dependency contracts

`deps.Deps` has one field per directory of `sandbox/deps/`, named by title-casing it. Each
field is that package's `Sandbox` struct, filled by `adapters/impls/<name>.Bind(&deps)`.

| Page | Declares |
| --- | --- |
| [`deps.HashDeps`](deps.hashdeps.md) | `Contract` |
| [`deps.StdDeps`](deps.stddeps.md) | `Contract` |
| [`deps.StorageDeps`](deps.storagedeps.md) | `Contract` |
| [`deps.StringsDeps`](deps.stringsdeps.md) | `Contract` |
