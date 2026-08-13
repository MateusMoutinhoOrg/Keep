# `api.SchemaItem`

**Type:** Struct (struct of function fields)

## Definition

```go
type SchemaItem struct {
	Deps              deps.Deps
	Items             []api.Item
	Prefix            string
	Id                int64
	Get               func(fieldName string) (any, *api.Error)
	Update            func(fieldName string, value any) *api.Error
	Remove            func() *api.Error
	CheckKeysPresence func(keys []string) bool
	ListAll           func(fieldName string) []api.SchemaItem
	NewSubItem        func(fieldName string, fields map[string]any) (api.SchemaItem, *api.Error)
	String            func() string
}
```

## Description

One record of a collection, returned by the [`SchemaInstance`](./api.SchemaInstance.md) fields. `Id` is a plain data field (never reused, even after deletion); every other behavior is a function field filled by factories in `sandbox/lib/schemaitem/`. The narrative guide with full examples is [Working with Records](/docs/References/Records.md).

## Fields

### `Id`

```go
Id int64
```

The record's permanent identifier — read directly, no accessor method needed.

### `Get`

```go
Get func(fieldName string) (any, *api.Error)
```

Reads one field. `Key` fields come back as `string`, `Int` fields as `int64`. Returns a `NotFound` error if the record never stored that field.

### `Update`

```go
Update func(fieldName string, value any) *api.Error
```

Writes a new value for a field. Updating a `Key` field re-indexes it and fails with `KeyConflict` if another record already owns the new value.

### `Remove`

```go
Remove func() *api.Error
```

Deletes the record, its unique index entries, and everything inside its sub-databases. Returns `nil` on success. Removing an already-removed record is a no-op.

### `ListAll`

```go
ListAll func(fieldName string) []api.SchemaItem
```

Returns every record of a sub-database (`Database` type) field of this record.

### `NewSubItem`

```go
NewSubItem func(fieldName string, fields map[string]any) (api.SchemaItem, *api.Error)
```

Inserts a record into a sub-database (`Database` type) field of this record.

### `CheckKeysPresence`

```go
CheckKeysPresence func(keys []string) bool
```

Reports whether every named field has a stored value for this record.

### `String`

```go
String func() string
```

Renders the record's plain fields, so printing a record shows `{id: 1, email: ..., age: ...}`.

## Examples

```go
user, ok := users.FindByKey("email", "a@x.com")

age, err := user.Get("age")            // int64
err = user.Update("age", 31)
session, err := user.NewSubItem("sessions", map[string]any{"token": "t1"})
for _, s := range user.ListAll("sessions") {
	fmt.Println(s.String())
}
e := user.Remove()
if e != nil {
	fmt.Println("error removing:", e.Message)
}
```
