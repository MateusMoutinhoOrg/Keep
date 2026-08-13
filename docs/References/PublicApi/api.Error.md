# `api.Error`

**Type:** Struct (plain data, no behavior)

## Definition

```go
type Error struct {
	Type     int    // the cause
	Key      string // the field involved
	KeyValue any    // the value involved, when relevant
	Message  string // human-readable description
}

const (
	KeyConflict = iota
	NotFound
	MissingField
	InvalidField
	Internal
)
```

## Description

The typed failure returned by database operations, as `*api.Error` — `nil` when the operation succeeded. Switch on `Type` to react to each failure; the full guide is [Error Handling](../Errors.md).

It is a **plain struct**, not an interface: `Error` carries no behavior, so callers read `Type`, `Key`, `KeyValue`, and `Message` as fields rather than calling methods, and there is no `Error() string` method. `*Error` stays a pointer so `if err != nil` still works as the success/failure check.

## Error Types

| Type | Meaning |
| :--- | :--- |
| `api.KeyConflict` | A unique key value is already taken. |
| `api.NotFound` | The field has no stored value for this record. |
| `api.MissingField` | A required field was not provided. |
| `api.InvalidField` | The field is not in the schema, or the value has the wrong type. |
| `api.Internal` | The storage backend failed. |

## Examples

```go
import keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"

created, err := users.NewItem(fields)
if err != nil {
	switch err.Type {
	case keeptypes.KeyConflict:
		fmt.Printf("%q %v is already taken\n", err.Key, err.KeyValue)
	case keeptypes.MissingField:
		fmt.Printf("field %q is required\n", err.Key)
	default:
		fmt.Println("unexpected error:", err.Message)
	}
}
```
