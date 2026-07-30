# `api.SchemaInstance`

**Type:** Interface

## Definition

```go
type SchemaInstance interface {
	NewItem(fields map[string]any) (SchemaItem, api.Error)
	FindByKey(key string, keyValue any) SchemaItem
	ListAll() ([]SchemaItem, api.Error)
	List(position int, chunk int) ([]SchemaItem, api.Error)
}
```

## Description

One collection of records, obtained via [`KeepDatabase.GetSchema`](./api.KeepDatabase.md#methods). It is the entry point for creating, finding, and listing records. The struct implementing it lives in `sandbox/internal/schemainstance/` and is unreachable from outside the sandbox. The narrative guide with full examples is [Working with Records](/docs/Explanations/Records.md).

## Methods

### `NewItem`

```go
func NewItem(fields map[string]any) (api.SchemaItem, api.Error)
```

Inserts a record. Fails with `MissingField` if a required field is absent, `InvalidField` if a field is not in the schema or has the wrong type, and `KeyConflict` if a `Key` value is already taken.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `fields` | `map[string]any` | Field name → value for every plain field to store. |

| Returns | Description |
| :--- | :--- |
| [`api.SchemaItem`](./api.SchemaItem.md) | The created record, with a permanent, never-reused id. |
| [`api.Error`](./api.Error.md) | `nil` on success. |

### `FindByKey`

```go
func FindByKey(key string, keyValue any) api.SchemaItem
```

Looks a record up by any `Key` field, at constant cost and case-insensitively. Returns `nil` when the field is not a `Key` of the schema or no record matches.

### `ListAll`

```go
func ListAll() ([]api.SchemaItem, api.Error)
```

Returns every record of the collection. List order is **not stable** across deletions (see [Working with Records](/docs/Explanations/Records.md#list--listall-and-list)).

### `List`

```go
func List(position int, chunk int) ([]api.SchemaItem, api.Error)
```

Returns up to `chunk` records starting at `position` (1-based). A past-the-end page is empty, not an error.

## Examples

```go
users := db.GetSchema("user")

created, err := users.NewItem(map[string]any{"email": "a@x.com", "username": "alice", "age": 30})
found := users.FindByKey("email", "a@x.com")
all, err2 := users.ListAll()
page, err3 := users.List(1, 10)
```
