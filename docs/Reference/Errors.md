# Error Handling

## Description
Lists the errors returned by database operations. Every operation returns `*lib.Error` (or a plain `lib.Error` from `Remove`), carrying a machine-checkable `Type` plus context about which field caused it.

---

## The Error struct

```go
type Error struct {
	Type     ErrorType // what kind of failure
	Key      string    // the field involved
	KeyValue any       // the value involved (when relevant)
	Msg      string    // human-readable description
}
```

---

## Error types

| Type | Meaning | Typical cause |
|---|---|---|
| `lib.KeyConflict` | A unique `Key` value is already taken | Creating or updating with an email/username that exists |
| `lib.NotFound` | The field has no stored value for this record | Reading a non-required field that was never set |
| `lib.MissingField` | A required field was not provided | `NewItem` without all `Required` fields |
| `lib.InvalidField` | The field is not in the schema, or the value has the wrong type | Typo in a field name, passing a string to an `Int` field |
| `lib.Internal` | The storage backend failed | I/O error, permissions, corrupted data |

---

## Reacting to an error

Switch on `Type` to decide what to do:

```go
created, err := users.NewItem(fields)
if err != nil {
	switch err.Type {
	case lib.KeyConflict:
		fmt.Printf("%q %v is already taken\n", err.Key, err.KeyValue)
	case lib.MissingField:
		fmt.Printf("field %q is required\n", err.Key)
	default:
		fmt.Println("unexpected error:", err)
	}
	return
}
```

A common pattern from the samples — treat "already exists" as fine and reuse the record:

```go
_, err := users.NewItem(fields)
if err != nil && err.Type != lib.KeyConflict {
	return err
}
user := users.FindByKey("email", email)
```

---

## Special cases

- `FindByKey` and `GetSchema` do not return errors — they return `nil` when nothing matches.
- `Remove` returns a value, not a pointer; check `e.Msg != ""` to detect failure.
