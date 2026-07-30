# `api.KeepDatabase`

**Type:** Interface

## Definition

```go
type KeepDatabase interface {
	GetSchema(name string) SchemaInstance
	Props() api.Props
}
```

## Description

A database bound to a storage backend ([`deps.Deps`](./deps.Deps.md)) and a schema description ([`Props`](./api.Props.md)). Always constructed via [`Lib.NewDatabase`](./api.Lib.md#methods). The struct implementing it lives in `sandbox/internal/database/` and is unreachable from outside the sandbox.

## Methods

### `GetSchema`

```go
func GetSchema(name string) api.SchemaInstance
```

Returns the collection whose schema has the given name.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | The `Schema.Name` to look up. |

| Returns | Description |
| :--- | :--- |
| [`api.SchemaInstance`](./api.SchemaInstance.md) | The collection, or `nil` if no schema has that name. |

### `Props`

```go
func Props() api.Props
```

Returns the [`Props`](./api.Props.md) description the database was created from.

## Examples

```go
db := keep.NewDatabase(Props)
users := db.GetSchema("user")
if users == nil {
	panic("schema not declared in Props")
}
```
