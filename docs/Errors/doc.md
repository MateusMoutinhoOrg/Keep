# Errors

Every operation that can fail returns `*api.Error`, and `nil` means success. It is plain
data: no behaviour, no wrapping, no `errors.Is`.

```go
user, failure := users.NewItem(fields)
if failure != nil {
	switch failure.Type {
	case api.KeyConflict:
		// failure.Key names the field, failure.KeyValue the value refused
	case api.MissingField, api.InvalidField:
		// the call was wrong
	case api.Internal:
		// the storage backend failed; failure.Message carries what it said
	}
}
```

Switch on `Type`. `Message` is written for a person and may change between releases; the
constants may not.

| Field | Carries |
|---|---|
| `Type` | one of the five constants below |
| `Key` | the field the failure involves, empty when it involves no particular field |
| `KeyValue` | the value the failure involves, when there is one |
| `Message` | the human-readable description |

## The five causes

| `Error.Type` | Raised by | Means |
|---|---|---|
| `api.KeyConflict` | `NewItem`, `NewSubItem`, `Update` | another live record of the same collection already holds that value for that `Key` field. Nothing was written |
| `api.NotFound` | `Get` | the field is declared by the schema and this record has no value stored for it |
| `api.MissingField` | `NewItem`, `NewSubItem` | a field declared `Required` was left out of the fields map |
| `api.InvalidField` | `NewItem`, `NewSubItem`, `Get`, `Update` | the field is not in the schema, the value is the wrong Go type for it, or a nested (`api.Database`) field was used where a plain value was expected |
| `api.Internal` | any | the storage backend reported a failure. `Message` is what it said |

`NotFound` and `InvalidField` are different answers to what looks like one question: a field
the schema does not declare is a mistake in the code, a declared field with no value is a
fact about that record. `CheckKeysPresence` answers the second for several fields at once
without reading any value.

## What returns no error

A lookup that finds nothing is not a failure — `FindByKey`, `FindById` and
`DatabaseHandle.GetSchema` report it as `ok == false`, because an absent record is an
ordinary outcome and these types are structs with no nil form:

```go
user, ok := users.FindByKey("email", "nobody@gmail.com")   // ok == false
```

`SchemaItem.ListAll(fieldName)` returns `nil` when the schema declares no nested field of
that name, and `Remove` on a record that is already gone returns `nil`: absent before and
absent after is the same outcome.

## Where a failure leaves the database

Nothing is half-written on a refusal a caller can see. `KeyConflict`, `MissingField` and
`InvalidField` are all decided before the first write of the operation, so the database is
exactly as it was.

`Internal` is the one case where a write sequence can stop part-way through. The orderings
that make the outcome harmless are in [DenseRecordPattern](../DenseRecordPattern/doc.md):
an insert commits on its last write, and an update to a `Key` field writes the new index
entry before it moves the value.
