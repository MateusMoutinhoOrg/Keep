# `lib.SchemaItem`

**Type:** Struct

## Definition

```go
type SchemaItem struct { /* internal */ }
```

## Description

One record of a collection, returned by the [`SchemaInstance`](./lib.SchemaInstance.md) methods. The narrative guide with full examples is [Working with Records](/docs/Explanation/Records.md).

## Methods

### `Id`

```go
func (s *SchemaItem) Id() int64
```

Returns the record's permanent identifier. Ids are never reused, even after deletion.

### `Get`

```go
func (s *SchemaItem) Get(fieldName string) (any, *Error)
```

Reads one field. `Key` fields come back as `string`, `Int` fields as `int64`. Returns a `NotFound` error if the record never stored that field.

### `Update`

```go
func (s *SchemaItem) Update(fieldName string, value any) *Error
```

Writes a new value for a field. Updating a `Key` field re-indexes it and fails with `KeyConflict` if another record already owns the new value.

### `Remove`

```go
func (s *SchemaItem) Remove() Error
```

Deletes the record, its unique index entries, and everything inside its sub-databases. Removing an already-removed record is a no-op. Returns a value, not a pointer — check `e.Msg != ""` to detect failure.

### `ListAll`

```go
func (s *SchemaItem) ListAll(fieldName string) []*SchemaItem
```

Returns every record of a sub-database (`Database` type) field of this record.

### `NewSubItem`

```go
func (s *SchemaItem) NewSubItem(fieldName string, fields map[string]any) (*SchemaItem, *Error)
```

Inserts a record into a sub-database (`Database` type) field of this record.

### `CheckKeysPresence`

```go
func (s *SchemaItem) CheckKeysPresence(keys []string) bool
```

Reports whether every named field has a stored value for this record.

### `String`

```go
func (s *SchemaItem) String() string
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
if e.Msg != "" {
	fmt.Println("error removing:", e)
}
```
