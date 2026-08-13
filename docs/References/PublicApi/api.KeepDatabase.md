# `api.KeepDatabase`

**Type:** Struct (struct of function fields)

## Definition

```go
type KeepDatabase struct {
	Deps      deps.Deps
	Props     api.Props
	GetSchema func(name string) (api.SchemaInstance, bool)
}
```

## Description

A database bound to a storage backend ([`deps.Deps`](./deps.Deps.md)) and a schema description ([`Props`](./api.Props.md), carried directly as a field). Always constructed via [`Lib.NewDatabase`](./api.Lib.md#fields). `GetSchema` is a function field filled by `GetSchemaFactory` in `sandbox/lib/database/`.

## Fields

### `Props`

```go
Props api.Props
```

The description the database was created from — read directly, no accessor method needed since `Props` is a plain data struct.

### `GetSchema`

```go
GetSchema func(name string) (api.SchemaInstance, bool)
```

Returns the collection whose schema has the given name. `ok` is `false` when no schema has that name — there is no nil form to check, because `SchemaInstance` is a struct.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | The `Schema.Name` to look up. |

| Returns | Description |
| :--- | :--- |
| [`api.SchemaInstance`](./api.SchemaInstance.md) | The collection. |
| `bool` | `false` when no schema has that name. |

## Examples

```go
props := keeptypes.Props{Path: "myDatabase/", Schemas: schemas} // see api.Props.md
db := keep.NewDatabase(props)
users, ok := db.GetSchema("user")
if !ok {
	panic("schema not declared in Props")
}
```
