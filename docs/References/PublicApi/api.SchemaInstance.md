# `api.SchemaInstance`

**Type:** Struct (struct of function fields)

## Definition

```go
type SchemaInstance struct {
	Deps      deps.Deps
	Items     []api.Item
	Prefix    string
	NewItem   func(fields map[string]any) (api.SchemaItem, *api.Error)
	FindByKey func(key string, keyValue any) (api.SchemaItem, bool)
	ListAll   func() ([]api.SchemaItem, *api.Error)
	List      func(position int, chunk int) ([]api.SchemaItem, *api.Error)
}
```

## Description

One collection of records, obtained via [`KeepDatabase.GetSchema`](./api.KeepDatabase.md#fields). It is the entry point for creating, finding, and listing records. Its function fields are filled by factories in `sandbox/internal/schemainstance/`. The narrative guide with full examples is [Working with Records](/docs/Explanations/Records.md).

## Fields

### `NewItem`

```go
NewItem func(fields map[string]any) (api.SchemaItem, *api.Error)
```

Inserts a record. Fails with `MissingField` if a required field is absent, `InvalidField` if a field is not in the schema or has the wrong type, and `KeyConflict` if a `Key` value is already taken.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `fields` | `map[string]any` | Field name → value for every plain field to store. |

| Returns | Description |
| :--- | :--- |
| [`api.SchemaItem`](./api.SchemaItem.md) | The created record, with a permanent, never-reused id. |
| `*api.Error` | `nil` on success. |

### `FindByKey`

```go
FindByKey func(key string, keyValue any) (api.SchemaItem, bool)
```

Looks a record up by any `Key` field, at constant cost and case-insensitively. `ok` is `false` when the field is not a `Key` of the schema or no record matches.

### `ListAll`

```go
ListAll func() ([]api.SchemaItem, *api.Error)
```

Returns every record of the collection. List order is **not stable** across deletions (see [Working with Records](/docs/Explanations/Records.md#list--listall-and-list)).

### `List`

```go
List func(position int, chunk int) ([]api.SchemaItem, *api.Error)
```

Returns up to `chunk` records starting at `position` (1-based). A past-the-end page is empty, not an error.

## Examples

```go
users, _ := db.GetSchema("user")

created, err := users.NewItem(map[string]any{"email": "a@x.com", "username": "alice", "age": 30})
found, ok := users.FindByKey("email", "a@x.com")
all, err2 := users.ListAll()
page, err3 := users.List(1, 10)
```
