# Error Handling

## Description
Lists the errors returned by database operations. Every operation returns `api.Error` — `nil` on success — carrying a machine-checkable `Type()` plus context about which field caused it.

---

## The Error interface

`Error` is an interface, not a struct: every value crossing the library boundary is a primitive or an interface. It embeds the standard `error`, so it can be returned and printed as one.

```go
type Error interface {
	error          // Error() string — human-readable description
	Type() int     // what kind of failure
	Key() string   // the field involved
	KeyValue() any // the value involved (when relevant)
}
```

---

## Error types

| Type | Meaning | Typical cause |
|---|---|---|
| `api.KeyConflict` | A unique key value is already taken | Creating or updating with an email/username that exists |
| `api.NotFound` | The field has no stored value for this record | Reading a non-required field that was never set |
| `api.MissingField` | A required field was not provided | `NewItem` without all `Required` fields |
| `api.InvalidField` | The field is not in the schema, or the value has the wrong type | Typo in a field name, passing a string to an `IntItem` field |
| `api.Internal` | The storage backend failed | I/O error, permissions, corrupted data |

---

## Reacting to an error

Switch on `Type()` to decide what to do:

```go
created, err := users.NewItem(fields)
if err != nil {
	switch err.Type() {
	case api.KeyConflict:
		fmt.Printf("%q %v is already taken\n", err.Key(), err.KeyValue())
	case api.MissingField:
		fmt.Printf("field %q is required\n", err.Key())
	default:
		fmt.Println("unexpected error:", err)
	}
	return
}
```

A common pattern from the samples — treat "already exists" as fine and reuse the record:

```go
_, err := users.NewItem(fields)
if err != nil && err.Type() != api.KeyConflict {
	return err
}
user := users.FindByKey("email", email)
```

---

## Special cases

- `FindByKey` and `GetSchema` do not return errors — they return `nil` when nothing matches.
- `Remove` returns `nil` on success; check `e != nil` to detect failure.
