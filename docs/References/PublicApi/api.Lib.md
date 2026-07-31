# `api.Lib`

**Type:** Struct (struct of function fields)

## Definition

```go
type Lib struct {
	Deps        deps.Deps
	NewDatabase func(props api.Props) api.KeepDatabase
}
```

## Description

The library entry point, handed back by [`lib.New`](./lib.New.md). It carries the injected `Deps` and one function field, `NewDatabase`, filled by a factory in `sandbox/internal/lib/` (`NewDatabaseFactory`). Consumers never build it directly — `lib.New` is the only constructor.

Because it is a struct of function fields rather than an interface, replacing its behavior for a test is a plain field assignment (`l.NewDatabase = func(...) api.KeepDatabase { ... }`) rather than a wrapper type. See [StructContracts.md](/docs/Explanations/StructContracts.md).

## Fields

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
	database "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

keep := keeplib.New(keepadapter.New())
props := database.Props{Path: "myDatabase/", Schemas: schemas} // see api.Props.md
db := keep.NewDatabase(props)
```
