# Schemas

A database is a value. `api.Props` names the prefix every key is written under and the
collections the database holds; nothing else describes it, so a database can be written,
read, diffed and versioned like any other literal. Building the database writes no key —
the first record does.

```go
var Props = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "username", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
				{Name: "height", Type: api.Float},
				{Name: "bio", Type: api.String},
				{Name: "nickname", Type: api.Key},
				{
					Name: "sessions",
					Type: api.Nested,
					Fields: []api.Field{
						{Name: "token", Type: api.Key, Required: true},
						{Name: "creation", Type: api.Int, Required: true},
					},
				},
			},
		},
	},
}

db, failure := lib.Databases.New(Props)    // InvalidSchema when Props is wrong
users, ok := db.Collection("user")         // ok == false for an unknown name
```

Signatures for every type named here are in [PublicApi](../PublicApi/doc.md). `Databases.New`
copies the `Props`, and every `Fields` and `Prefix` a collection or a record hands back is a
copy too: editing one changes nothing the database does.

## The three units

| Unit | Is | Key field |
|---|---|---|
| `Props` | one database | `Path`, the prefix every key starts with |
| `Schema` | one collection of records | `Name`, what `Database.Collection` takes |
| `Field` | one field of a collection | `Name`, what `Get`, `Update` and the fields map take |

A `Field` of type `api.Link` carries one more: `Target`, the `Name` of the schema it points
at.

`Path` is a prefix, not a directory: it is split on slashes into the leading segments of
every key, so a backend that maps keys to files reads it as one, and an in-memory or remote
backend keeps it as part of the key. Empty segments are dropped, which makes a trailing
slash optional and harmless either way. Where the tree lands is the backend's business —
`filestorage.New(base)` — so a `.` or `..` segment is refused.

## What `Databases.New` refuses

A mistake in a schema is a mistake in the code, so it is refused once, where the database is
declared, as an `InvalidSchema` whose `Field` is the dotted path at fault (`user`,
`user.sessions.token`):

| Refused | Because |
|---|---|
| a schema or field with no `Name`, or one declared twice in the same place | two of them would share their keys |
| a field with no `Type`, or one outside `api.Key` … `api.Bytes` | the constants start at 1, so a forgotten `Type` is not read as a `Key` |
| an `api.Link` with no `Target`, or one naming no schema of the `Props` | the link could never resolve |
| a `Target` on anything but a `Link`, `Fields` on anything but a `Nested` | a misplaced declaration |
| an `api.Nested` field named `position` or `values` | the record keeps its own keys under those names — see [DenseRecordPattern](../DenseRecordPattern/doc.md#key-layout) |
| a `.` or `..` segment in `Path` | the database would land outside the backend's base |

## Field types

| `Field.Type` | Go type written | Go type read back | Indexed |
|---|---|---|---|
| `api.Key` | `string`, or anything with a `String() string` method | `string` | yes — unique across the collection |
| `api.String` | `string`, or anything with a `String() string` method | `string` | no |
| `api.Int` | any Go integer, or a float holding a whole number | `int64` | no |
| `api.Float` | any Go float or integer | `float64` | no |
| `api.Link` | `api.Record` of the `Target` collection, or what `api.Int` takes | `int64` | no |
| `api.Bytes` | `[]byte` | `[]byte` | no |
| `api.Nested` | never written directly | never read directly | — |

`api.String` is `api.Key` without the index: same values in, same values out, but two live
records may hold the same one and `FindByKey` refuses to look one up. `api.Int` takes
`float64(27)` — what `encoding/json` hands back for every number — but refuses `27.5`, and
refuses an unsigned value too large for an `int64`. `api.Float` stores the shortest decimal
form that parses back to the same number, so writing the same value twice writes the same
bytes. `api.Link` stores a record id exactly as `api.Int` does, and adds `Target` — see
[Pointing one record at another](#pointing-one-record-at-another). `api.Bytes` stores the
slice verbatim, with no text encoding, so any content — non-UTF-8 and zero bytes included —
comes back byte for byte; `String()` prints its length rather than its contents.

A value of the wrong Go type is refused with `InvalidField` before anything is written — and
so is a `String()` method that panics; the list of failures is in [Errors](../Errors/doc.md).

`Required: true` makes an insert that leaves the field out, or gives it as `nil`, fail with
`MissingField`. On a field that is not required, `nil` in an insert is the same as leaving it
out, and `Update(name, nil)` clears it. `Required` is ignored on an `api.Nested` field, which
is never provided to an insert.

## Keys

A `Key` field carries a unique index, which is what `Collection.FindByKey` reads:

- **Unique.** Two live records of one collection can never hold the same value for it. An
  insert or an `Update` that would break that fails with `KeyConflict` and writes nothing.
- **Case-insensitive, in every script.** The value is case-folded and canonically
  decomposed before it is hashed, so `MATEUS@GMAIL.COM` and `mateus@gmail.com` are the same
  key, and so are `STRASSE` and `straße`, `ΟΔΟΣ` and `οδος`, or a precomposed `é` and an `e`
  followed by a combining accent. Folding is not confusable detection: a Cyrillic `а` is
  still not a Latin `a`.
- **Scoped to its own collection.** A nested collection has its own index, so the same token
  can live under two different users.
- **Fixed cost.** The lookup is the hash of the value, so it costs the same at any size.

A collection may declare any number of `Key` fields, and each indexes its own values:
`FindByKey("email", …)` and `FindByKey("username", …)` both work on the schema above.

## Nested collections

An `api.Nested` field is a collection rooted at the record that owns it. `Record.Nested`
hands it out as an `api.Collection` of its own — own ids, own unique indexes, same `Insert`,
`FindByKey`, `FindByID`, `List`, `ListAll` and `Repair`:

```go
sessions, failure := user.Nested("sessions")
session, failure := sessions.Insert(map[string]any{"token": "token-1", "creation": 1000})
found, ok, failure := sessions.FindByKey("token", "TOKEN-1")
```

`InsertNested(field, fields)` and `ListNested(field)` are the same calls in one step. `Get`
on such a field fails with `InvalidField`: it is not a value. Removing the owning record
removes every record under it, at any depth, and an insert under a removed owner fails with
`Removed`.

Nesting has no depth limit of its own — a `Field` of an `api.Nested` field may itself be an
`api.Nested` — though a backend may have one; see
[StorageContract](../StorageContract/doc.md#filestorage). A `Link` cannot point into a nested
collection: its `Target` is a top-level schema.

## Pointing one record at another

Every record carries a permanent `ID`, and ids are allocated from a counter that only
grows. An `api.Link` field stores one and names the collection it belongs to, so the record
is followed in one call:

```go
{Name: "author", Type: api.Link, Target: "user", Required: true}

posts.Insert(map[string]any{"slug": "…", "author": author})      // or author.ID
resolved, ok, failure := post.GetLink("author")
```

`GetLink` reports `ok == false` when the field holds no value or when the record the id
names is no longer live, and fails with `InvalidField` on a field that is not a `Link`. `Get`
on the same field still returns the bare id as an `int64`.

A `Link` given a record checks that record belongs to the `Target` collection — a post given
where a user is expected is refused with `InvalidField`, since its id would name some other
user. A bare id is not checked: any number of records may link to the same target, and a
link to a record that was never written is stored like any other. What a link cannot do is
resolve to the wrong record — nothing is ever handed a used id again, so a stale link
resolves to nothing rather than to whatever took its place.

Use an `api.Int` field plus `Collection.FindByID` for the same reference when the
schema cannot name the target — a link across databases, or one whose collection is chosen
at run time:

```go
posts.Insert(map[string]any{"slug": "…", "author": author.ID})   // an api.Int field
resolved, ok, failure := users.FindByID(authorID)
```

Use a nested collection when the children belong to the parent and die with it, and a link
when they do not.

## Reading a collection

`ListAll` reads every record, `List(position, chunk)` the records of `chunk` positions from
`position` on, counted from 1 — `0` for "to the end"; a position below 1 or a negative chunk
is an `InvalidArgument`. Neither is a snapshot. A removal moves the last record into the
freed position, so a caller paging through a collection that is being written may miss the
record that moved, and a record removed while a listing runs may be left out of it.

## Changing a schema

A schema is read on every call and stored nowhere, so changing it is changing the `Props`.
What the records already written make of the change:

| Change | Safe | What to do |
|---|---|---|
| add a field that is not required, a `Nested` field, or a schema | yes | nothing |
| add a required field | for new inserts | old records hold no value: `Get` reports `NoValue` |
| turn a `String` into a `Key` | after `Repair` | `Collection.Repair` indexes every record; two holding the same value fail it with `KeyConflict` |
| turn a `Key` into a `String` | yes | its index entries are left behind, unread |
| change any other `Type` | no | stored values are not converted: `Get` fails with `Internal` on one it cannot parse |
| remove a field | no | its values stay in storage, and a removal no longer deletes them: clear the field on every record with `Update(name, nil)` first |
| rename a field or a schema | no | it is a remove plus an add |

Runnable versions of all of this are in [LibExamples](../LibExamples/doc.md).
