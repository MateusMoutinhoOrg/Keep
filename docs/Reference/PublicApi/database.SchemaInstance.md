# `database.SchemaInstance`

**Type:** Struct

## Definition

```go
type SchemaInstance struct { /* internal */ }
```

## Description

One collection of records, obtained via [`KeepDatabase.GetSchema`](./database.KeepDatabase.md#methods). It is the entry point for creating, finding, and listing records. The narrative guide with full examples is [Working with Records](/docs/Explanation/Records.md).

## Methods

### `NewItem`

```go
func (si *SchemaInstance) NewItem(fields map[string]any) (*SchemaItem, *Error)
```

Inserts a record. Fails with `MissingField` if a required field is absent, `InvalidField` if a field is not in the schema or has the wrong type, and `KeyConflict` if a `Key` value is already taken.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `fields` | `map[string]any` | Field name → value for every plain field to store. |

| Returns | Description |
| :--- | :--- |
| [`*SchemaItem`](./database.SchemaItem.md) | The created record, with a permanent, never-reused id. |
| [`*Error`](./database.Error.md) | `nil` on success. |

### `FindByKey`

```go
func (si *SchemaInstance) FindByKey(key string, keyValue any) *SchemaItem
```

Looks a record up by any `Key` field, at constant cost and case-insensitively. Returns `nil` when the field is not a `Key` of the schema or no record matches.

### `ListAll`

```go
func (si *SchemaInstance) ListAll() ([]*SchemaItem, *Error)
```

Returns every record of the collection. List order is **not stable** across deletions (see [Working with Records](/docs/Explanation/Records.md#list--listall-and-list)).

### `List`

```go
func (si *SchemaInstance) List(position int, chunk int) ([]*SchemaItem, *Error)
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
