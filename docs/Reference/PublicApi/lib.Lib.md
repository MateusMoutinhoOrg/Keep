# `lib.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	deps deps.Deps
}
```

## Description

The main library entry point. Holds the injected dependency adapter and creates databases with it wired in. Always constructed via [`lib.New`](./lib.New.md).

## Fields

| Field | Type | Visibility | Description |
| :--- | :--- | :--- | :--- |
| `deps` | [`deps.Deps`](./deps.Deps.md) | Private | The injected storage functions every created database will use, wired in by [`lib.New`](./lib.New.md). |

## Methods

### `NewDatabase`

```go
func (l *Lib) NewDatabase(props lib.Props) *lib.KeepDatabase
```

Creates a [`KeepDatabase`](./lib.KeepDatabase.md) from a [`Props`](./lib.Props.md) description, with the lib's deps wired in.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `props` | [`lib.Props`](./lib.Props.md) | The database description: key prefix and schemas. |

| Returns | Description |
| :--- | :--- |
| [`*lib.KeepDatabase`](./lib.KeepDatabase.md) | A database ready to hand out its collections via `GetSchema`. |

## Examples

```go
keep := lib.New(standard.New())
db := keep.NewDatabase(Props)
```
