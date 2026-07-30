# `api.Error`

**Type:** Interface

## Definition

```go
type Error interface {
	error          // Error() string
	Type() int     // the cause
	Key() string   // the field involved
	KeyValue() any // the value involved, when relevant
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

The typed error returned by database operations, and `nil` when the operation succeeded. Switch on `Type()` to react to each failure; the full guide is [Error Handling](../Errors.md).

It is an **interface**, not a struct, because every value crossing the library boundary must be a primitive or an interface. It embeds the standard `error`, so it can be returned, printed, and compared as one.

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
import "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"

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
}
```
