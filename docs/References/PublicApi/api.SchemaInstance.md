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
	FindById  func(id int64) (api.SchemaItem, bool)
	ListAll   func() ([]api.SchemaItem, *api.Error)
	List      func(position int, chunk int) ([]api.SchemaItem, *api.Error)
}
```

## Description

One collection of records, obtained via [`KeepDatabase.GetSchema`](./api.KeepDatabase.md#fields). It is the entry point for creating, finding, and listing records. Its function fields are filled by factories in `sandbox/lib/schemainstance/`. The narrative guide with full examples is [Working with Records](/docs/References/Records.md).

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

### `FindById`

```go
FindById func(id int64) (api.SchemaItem, bool)
```

Looks a record up by its permanent id — the value [`api.SchemaItem.Id`](./api.SchemaItem.md) reports. It is the cheapest lookup the library has: no index entry is read, the id alone names the record's keys. `ok` is `false` when the collection holds no live record under that id, either because it was never allocated or because the record was removed; ids are never reused, so a stale id never resolves to a different record.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `id` | `int64` | The record's permanent id, as reported by `SchemaItem.Id`. |

| Returns | Description |
| :--- | :--- |
| [`api.SchemaItem`](./api.SchemaItem.md) | The record, when `ok` is `true`. |
| `bool` | `false` when no live record carries the id. |

Because ids are stable and never reused, storing one in an `Int` field of another collection makes that field a **pointer to a record** — a foreign key — and `FindById` is what follows it. See [Relations between collections](/docs/References/Records.md#relations-between-collections).

### `ListAll`

```go
ListAll func() ([]api.SchemaItem, *api.Error)
```

Returns every record of the collection. List order is **not stable** across deletions (see [Working with Records](/docs/References/Records.md#list--listall-and-list)).

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
sameUser, ok2 := users.FindById(created.Id)
all, err2 := users.ListAll()
page, err3 := users.List(1, 10)
```
