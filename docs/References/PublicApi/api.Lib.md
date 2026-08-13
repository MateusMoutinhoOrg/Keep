# `api.Lib`

**Type:** Struct (struct of function fields)

## Definition

```go
type Lib struct {
	Deps        deps.Deps
	Version     func() string
	NewDatabase func(props api.Props) api.KeepDatabase
}
```

## Description

The library entry point, handed back by [`lib.New`](./lib.New.md). It carries the injected `Deps` and its function fields, each filled by a factory of its own file under `sandbox/lib/publicfunctions/` and assigned by `New` in `sandbox/lib/new.go`. Consumers never build it directly — `lib.New` is the only constructor.

Because it is a struct of function fields rather than an interface, replacing its behavior for a test is a plain field assignment (`l.NewDatabase = func(...) api.KeepDatabase { ... }`) rather than a wrapper type. See [StructContracts.md](/docs/References/StructContracts.md).

## Fields

### `Version`

```go
Version func() string
```

Returns the library's own release — the version the consumer linked against. The number is a compile-time constant in `sandbox/config`, so a release bump is a one-line edit touching no logic, and it is exposed as a field like every other behavior so a consumer can report it without importing anything but this api.

| Returns | Description |
| :--- | :--- |
| `string` | The release tag, e.g. `v0.0.4`. |

### `NewDatabase`

```go
NewDatabase func(props api.Props) api.KeepDatabase
```

Creates a [`KeepDatabase`](./api.KeepDatabase.md) from a [`Props`](./api.Props.md) description, with the lib's `Deps` propagated into it.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `props` | [`api.Props`](./api.Props.md) | The database description: key prefix and schemas. |

| Returns | Description |
| :--- | :--- |
| [`api.KeepDatabase`](./api.KeepDatabase.md) | A database ready to hand out its collections via `GetSchema`. |

## Examples

```go
import (
	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

keep := keeplib.New(keepadapter.New())
fmt.Println("keep", keep.Version())

props := keeptypes.Props{Path: "myDatabase/", Schemas: schemas} // see api.Props.md
db := keep.NewDatabase(props)
```
