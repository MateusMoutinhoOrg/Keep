# `api.SchemaItem`

**Type:** Interface

## Definition

```go
type SchemaItem interface {
	Id() int64
	Get(fieldName string) (any, api.Error)
	Update(fieldName string, value any) api.Error
	Remove() api.Error
	CheckKeysPresence(keys []string) bool
	ListAll(fieldName string) []SchemaItem
	NewSubItem(fieldName string, fields map[string]any) (SchemaItem, api.Error)
	String() string
}
```

## Description

One record of a collection, returned by the [`SchemaInstance`](./api.SchemaInstance.md) methods. The struct implementing it lives in `sandbox/internal/schemaitem/` and is unreachable from outside the sandbox. The narrative guide with full examples is [Working with Records](/docs/Explanations/Records.md).

## Methods

### `Id`

```go
func Id() int64
```

Returns the record's permanent identifier. Ids are never reused, even after deletion.

### `Get`

```go
func Get(fieldName string) (any, api.Error)
```

Reads one field. `Key` fields come back as `string`, `Int` fields as `int64`. Returns a `NotFound` error if the record never stored that field.

### `Update`

```go
func Update(fieldName string, value any) api.Error
```

Writes a new value for a field. Updating a `Key` field re-indexes it and fails with `KeyConflict` if another record already owns the new value.

### `Remove`

```go
func Remove() api.Error
```

Deletes the record, its unique index entries, and everything inside its sub-databases. Returns `nil` on success. Removing an already-removed record is a no-op.

### `ListAll`

```go
func ListAll(fieldName string) []api.SchemaItem
```

Returns every record of a sub-database (`Database` type) field of this record.

### `NewSubItem`

```go
func NewSubItem(fieldName string, fields map[string]any) (api.SchemaItem, api.Error)
```

Inserts a record into a sub-database (`Database` type) field of this record.

### `CheckKeysPresence`

```go
func CheckKeysPresence(keys []string) bool
```

Reports whether every named field has a stored value for this record.

### `String`

```go
func String() string
```

Renders the record's plain fields, so printing a record shows `{id: 1, email: ..., age: ...}`.

## Examples

```go
user := users.FindByKey("email", "a@x.com")

age, err := user.Get("age")            // int64
err = user.Update("age", 31)
session, err := user.NewSubItem("sessions", map[string]any{"token": "t1"})
for _, s := range user.ListAll("sessions") {
	fmt.Println(s)
}
e := user.Remove()
if e != nil {
	fmt.Println("error removing:", e)
}
```
