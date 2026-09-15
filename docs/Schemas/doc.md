# Schemas

A database is a value. `api.Props` names the prefix every key is written under and the
collections the database holds; nothing else describes it, so a database can be written,
read, diffed and versioned like any other literal. Building the handle writes no key —
the first record does.

```go
var Props = api.Props{
	Path: "TestDir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Itens: []api.Item{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "username", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
				{Name: "nickname", Type: api.Key},
				{
					Name: "sessions",
					Type: api.Database,
					Itens: []api.Item{
						{Name: "token", Type: api.Key, Required: true},
						{Name: "creation", Type: api.Int, Required: true},
					},
				},
			},
		},
	},
}

db := lib.Databases.New(Props)
users, ok := db.GetSchema("user")
```

`ok` is false when the `Props` declares no schema under that name. Signatures for every
type named here are in [PublicApi](../PublicApi/doc.md).

## The three units

| Unit | Is | Key field |
|---|---|---|
| `Props` | one database | `Path`, the prefix every key starts with |
| `Schema` | one collection of records | `Name`, what `GetSchema` takes |
| `Item` | one field of a collection | `Name`, what `Get`, `Update` and the fields map take |

`Path` is a prefix, not a directory: it is split on slashes into the leading segments of
every key, so a backend that maps keys to files reads it as one, and an in-memory or remote
backend keeps it as part of the key. Empty segments are dropped, which makes a trailing
slash optional and harmless either way.

`Itens` is spelled that way in `Schema` and in `Item`. It is part of the api.

## Field types

| `Item.Type` | Go type written | Go type read back | Indexed |
|---|---|---|---|
| `api.Key` | `string`, or anything with a `String() string` method | `string` | yes — unique across the collection |
| `api.Int` | `int`, `int32`, `int64` | `int64` | no |
| `api.Database` | never written directly | never read directly | — |

A value of the wrong Go type is refused with `InvalidField` before anything is written; the
list of failures is in [Errors](../Errors/doc.md).

`Required: true` makes an insert that leaves the field out fail with `MissingField`. It is
ignored on an `api.Database` field, which is never provided to an insert.

## Keys

A `Key` field carries a unique index, which is what `SchemaInstance.FindByKey` reads:

- **Unique.** Two live records of one collection can never hold the same value for it. An
  insert or an `Update` that would break that fails with `KeyConflict` and writes nothing.
- **Case-insensitive.** The value is lower-cased before it is hashed, so `MATEUS@GMAIL.COM`
  and `mateus@gmail.com` are the same key.
- **Scoped to its own collection.** A nested collection has its own index, so the same token
  can live under two different users.
- **One read.** The lookup is the hash of the value, so it costs the same at any size.

A collection may declare any number of `Key` fields, and each indexes its own values:
`FindByKey("email", …)` and `FindByKey("username", …)` both work on the schema above.

## Nested collections

An `api.Database` field is a collection rooted at the record that owns it. It behaves like
a top-level collection in every way — own ids, own unique indexes, same listing — and it is
reached through the record rather than through `GetSchema`:

```go
session, failure := user.NewSubItem("sessions", map[string]any{
	"token": "token-1", "creation": 1000,
})
for _, session := range user.ListAll("sessions") { … }
```

`Get` on such a field fails with `InvalidField`: it is not a value. Removing the owning
record removes every record under it, at any depth.

Nesting has no depth limit — an `Item` of an `api.Database` field may itself be an
`api.Database`.

## Pointing one record at another

There is no reference type, and there does not need to be one. Every record carries a
permanent `Id`, ids are allocated from a counter that only grows, and `FindById` resolves
one in a single read:

```go
posts.NewItem(map[string]any{"slug": "…", "author": author.Id})   // an api.Int field
resolved, ok := users.FindById(authorId)
```

`ok` is false once the record is removed, and stays false forever: nothing is ever handed
that id again, so a stale reference resolves to nothing rather than to whatever took its
place. Use a nested collection when the children belong to the parent and die with it, and
an id field when they do not.

Runnable versions of all of this are in [LibExamples](../LibExamples/doc.md).
