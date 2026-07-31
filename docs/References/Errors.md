# Error Handling

## Description
Lists the errors returned by database operations. Every operation returns `*api.Error` — `nil` on success — carrying a machine-checkable `Type` field plus context about which field caused it.

---

## The Error struct

`Error` is a plain data struct, not an interface: every field carrying behavior in `sandbox/contracts/api` leads with a `Deps` field, but `Error` needs no live computation, so it is pure data with no methods at all — not even `Error() string`. `*Error` stays a pointer so `err != nil` is still the success/failure check.

```go
type Error struct {
	Type     int    // what kind of failure
	Key      string // the field involved
	KeyValue any    // the value involved (when relevant)
	Message  string // human-readable description
}
```

---

## Error types

| Type | Meaning | Typical cause |
|---|---|---|
| `api.KeyConflict` | A unique key value is already taken | Creating or updating with an email/username that exists |
| `api.NotFound` | The field has no stored value for this record | Reading a non-required field that was never set |
| `api.MissingField` | A required field was not provided | `NewItem` without all `Required` fields |
| `api.InvalidField` | The field is not in the schema, or the value has the wrong type | Typo in a field name, passing a string to an `Int` field |
| `api.Internal` | The storage backend failed | I/O error, permissions, corrupted data |

---

## Reacting to an error

Switch on `Type` to decide what to do:

```go
created, err := users.NewItem(fields)
if err != nil {
	switch err.Type {
	case database.KeyConflict:
		fmt.Printf("%q %v is already taken\n", err.Key, err.KeyValue)
	case database.MissingField:
		fmt.Printf("field %q is required\n", err.Key)
	default:
		fmt.Println("unexpected error:", err.Message)
	}
	return
}
```

A common pattern from the samples — treat "already exists" as fine and reuse the record:

```go
_, err := users.NewItem(fields)
if err != nil && err.Type != database.KeyConflict {
	fmt.Println("Error creating user", err.Message)
	return
}
user, ok := users.FindByKey("email", email)
```

---

## Special cases

- `FindByKey` and `GetSchema` return `(value, ok bool)` — `ok` is `false` when nothing matches, since `SchemaItem`/`SchemaInstance` are structs with no nil form.
- `Remove` returns `nil` on success; check `e != nil` to detect failure.
