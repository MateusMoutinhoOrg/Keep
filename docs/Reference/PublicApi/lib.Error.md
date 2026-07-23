# `lib.Error`

**Type:** Struct

## Definition

```go
type Error struct {
	Type     ErrorType // what kind of failure
	Key      string    // the field involved
	KeyValue any       // the value involved (when relevant)
	Msg      string    // human-readable description
}

type ErrorType int

const (
	KeyConflict ErrorType = iota
	NotFound
	MissingField
	InvalidField
	Internal
)
```

## Description

The typed error returned by database operations. Switch on `Type` to react to each failure; the full guide is [Error Handling](../Errors.md).

## Error Types

| Type | Meaning |
| :--- | :--- |
| `KeyConflict` | A unique `Key` value is already taken. |
| `NotFound` | The field has no stored value for this record. |
| `MissingField` | A required field was not provided. |
| `InvalidField` | The field is not in the schema, or the value has the wrong type. |
| `Internal` | The storage backend failed. |

## Examples

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
}
```
