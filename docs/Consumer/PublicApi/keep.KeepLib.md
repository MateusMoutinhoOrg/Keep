# `keep.KeepLib`

**Type:** Struct

## Definition

```go
type KeepLib struct {
	Deps deps.Deps
}
```

## Description

The main library entry point. Holds the injected dependency adapter and creates databases with it wired in. Always constructed via [`keep.New`](./keep.New.md).

## Fields

| Field | Type | Visibility | Description |
| :--- | :--- | :--- | :--- |
| `Deps` | [`deps.Deps`](./deps.Deps.md) | Public | The storage functions every created database will use. |

## Methods

### `NewDatabase`

```go
func (l KeepLib) NewDatabase(props database.Props) *database.KeepDatabase
```

Creates a [`KeepDatabase`](./database.KeepDatabase.md) from a [`Props`](./database.Props.md) description, with the lib's deps wired in.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `props` | [`database.Props`](./database.Props.md) | The database description: key prefix and schemas. |

| Returns | Description |
| :--- | :--- |
| [`*database.KeepDatabase`](./database.KeepDatabase.md) | A database ready to hand out its collections via `GetSchema`. |

## Examples

```go
keep := keep_lib.New(keep_deps.New())
db := keep.NewDatabase(Props)
```
