# Errors

Every operation that can fail returns `*api.Error`, and `nil` means success. It is plain
data: no behaviour, no wrapping, no `errors.Is`.

```go
user, failure := users.Insert(fields)
if failure != nil {
	switch failure.Type {
	case api.KeyConflict:
		// failure.Field names the field, failure.Value the value refused
	case api.MissingField, api.InvalidField:
		// the call was wrong
	case api.Internal:
		// the storage backend failed; failure.Message carries what it said
	}
}
```

Switch on `Type`. `Message` is written for a person and may change between releases; the
constants may not. They start at 1, so a zero `api.Error{}` is none of them.

| Field | Carries |
|---|---|
| `Type` | one of the eight constants below |
| `Field` | the field the failure involves, empty when it involves no particular field. With several bad fields in one insert, the first in schema order — so the same call always names the same one |
| `Value` | the value the failure involves, when there is one |
| `Message` | the human-readable description |

## The eight causes

| `Error.Type` | Raised by | Means |
|---|---|---|
| `api.KeyConflict` | `Insert`, `InsertNested`, `Update`, `Repair` | another live record of the same collection already holds that value for that `Key` field. Nothing was written — or, from `Repair`, two records already do |
| `api.NoValue` | `Get` | the field is declared by the schema and this record has no value stored for it — which is also what a removed record reports |
| `api.MissingField` | `Insert`, `InsertNested`, `Update` | a field declared `Required` was left out of the fields map, given as `nil`, or cleared by an `Update` to `nil` |
| `api.InvalidField` | `Insert`, `InsertNested`, `Get`, `GetLink`, `Update`, `HasValues`, `FindByKey`, `Nested`, `ListNested` | the field is not in the schema, or is not of the kind the call needs (`FindByKey` a `Key`, `GetLink` a `Link`, `Nested` a `Nested`), the value is the wrong Go type for it or its `String()` panicked, or a `Link` was given a record of another collection |
| `api.Internal` | any | the storage backend reported a failure, or the write lock of the collection stayed held for two minutes. `Message` is what it said |
| `api.Removed` | `Update`, `InsertNested`, a nested `Insert` or `Repair` | the record written to is no longer live, or the record owning the nested collection is not. Nothing was written |
| `api.InvalidSchema` | `Databases.New` | the `Props` cannot be built — see [Schemas](../Schemas/doc.md#what-databasesnew-refuses). `Field` is the dotted path at fault |
| `api.InvalidArgument` | `List` | a position below 1 or a negative chunk |

`NoValue` and `InvalidField` are different answers to what looks like one question: a field
the schema does not declare is a mistake in the code, a declared field with no value is a
fact about that record. `HasValues` answers the second for several fields at once
without reading any value.

## What returns no error

A lookup that finds nothing is not a failure — `FindByKey`, `FindByID` and `GetLink` report
it as `ok == false` with a `nil` failure, and `Database.Collection` as `ok == false`, because
an absent record is an ordinary outcome and these types are structs with no nil form:

```go
user, ok, failure := users.FindByKey("email", "nobody@gmail.com")   // ok == false, failure == nil
```

A backend that fails is never folded into that answer: the same lookup then returns `Internal`,
so a caller never reads a storage outage as "no such user". `GetLink` folds the two reasons
a stored link does not resolve — no value, or a record that has been removed — into the same
`ok == false`, because a caller does nothing different for either.

`Remove` on a record that is already gone returns `nil`: absent before and absent after is
the same outcome.

## Where a failure leaves the database

Nothing is half-written on a refusal a caller can see. `KeyConflict`, `MissingField`,
`InvalidField`, `Removed`, `InvalidSchema` and `InvalidArgument` are all decided before the
first write of the operation, so the database is exactly as it was.

`Internal` is the one case where a write sequence can stop part-way through. The orderings
that make the outcome harmless are in [DenseRecordPattern](../DenseRecordPattern/doc.md):

- an insert commits on its last write, and deletes what it wrote when it fails before;
- an update to a `Key` field writes the new index entry before it moves the value;
- a removal run again finishes what it started.

So the answer to an `Internal` is to call again — with the same record, for an `Update` or
a `Remove`. What a crash leaves behind is invisible to every read, and
`Collection.Repair` deletes it.
