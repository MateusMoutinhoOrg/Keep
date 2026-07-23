# `lib.KeepDatabase`

**Type:** Struct

## Definition

```go
type KeepDatabase struct {
	deps  deps.Deps
	Props Props
}
```

## Description

A database bound to a storage backend ([`deps.Deps`](./deps.Deps.md)) and a schema description ([`Props`](./lib.Props.md)). Always constructed via [`Lib.NewDatabase`](./lib.Lib.md#methods).

## Methods

### `GetSchema`

```go
func (d *KeepDatabase) GetSchema(name string) *SchemaInstance
```

Returns the collection whose schema has the given name.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | The `Schema.Name` to look up. |

| Returns | Description |
| :--- | :--- |
| [`*SchemaInstance`](./lib.SchemaInstance.md) | The collection, or `nil` if no schema has that name. |

## Examples

```go
db := keep.NewDatabase(Props)
users := db.GetSchema("user")
if users == nil {
	panic("schema not declared in Props")
}
```
