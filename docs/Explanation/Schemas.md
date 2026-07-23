# Schemas

## Description
Explains how a database is described: its collections, their typed fields, and the sub-databases a record can own. To operate on the records themselves, see [Records.md](/docs/Explanation/Records.md).

---

## Collections

A schema describes one collection of records: its name and the fields (`Itens`) each record can hold. Schemas are passed to the database through `lib.Props`.

```go
var Props = lib.Props{
	Path:    "myDatabase/",
	Schemas: []lib.Schema{ /* ... */ },
}
```

- `Path` is a prefix added to every key the database stores. With the standard (filesystem) adapter it behaves like a folder.
- `GetSchema(name)` returns the collection, or `nil` if no schema has that name.

---

## Fields (`lib.Item`)

Every field has a `Name`, a `Type`, and an optional `Required` flag. Required fields must be present when creating a record.

| Type | Holds | Notes |
|---|---|---|
| `lib.Key` | `string` | Unique and indexed: two records can never share the same value, and `FindByKey` can look records up by it. Uniqueness is case-insensitive (`User@x.com` and `user@x.com` conflict). |
| `lib.Int` | `int`, `int32`, or `int64` | Plain integer field, always read back as `int64`. |
| `lib.Database` | a nested collection | The field is itself a sub-database with its own `Itens`. See below. |

---

## Example

A `user` collection where each user owns a nested `sessions` collection:

```go
var Schemas = []lib.Schema{
	{
		Name: "user",
		Itens: []lib.Item{
			{Name: "email", Type: lib.Key, Required: true},
			{Name: "username", Type: lib.Key, Required: true},
			{Name: "age", Type: lib.Int, Required: true},
			{
				Name: "sessions",
				Type: lib.Database,
				Itens: []lib.Item{
					{Name: "token", Type: lib.Key, Required: true},
					{Name: "creation", Type: lib.Int, Required: true},
					{Name: "expiration", Type: lib.Int, Required: true},
				},
			},
		},
	},
}
```

---

## Sub-databases

A `Database` field gives each record its own private collection. In the example above, every user has its own list of sessions, isolated from other users' sessions.

Sub-database fields cannot be set with `NewItem` or `Update` — they are managed through the record itself:

```go
user := users.FindByKey("email", "mateus@gmail.com")

// Insert into the user's sessions
session, err := user.NewSubItem("sessions", map[string]any{
	"token":      "token-1",
	"creation":   1000,
	"expiration": 2000,
})

// List them
for _, s := range user.ListAll("sessions") {
	fmt.Println(s)
}
```

Sub-databases nest to any depth: an item inside `sessions` could itself have a `Database` field. When a record is removed, all its sub-database records are removed with it.

---

## Naming rules

Collection and field names must not contain the `-` character — it is the separator used internally in the key layout (see [DenseRecordPattern.md](/docs/Explanation/DenseRecordPattern.md)).
