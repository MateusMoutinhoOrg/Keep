# Schemas

## Description
Explains how a database is described: its collections, their typed fields, and the sub-databases a record can own. To operate on the records themselves, see [Records.md](/docs/Explanations/Records.md).

---

## Collections

A schema describes one collection of records: its name and the fields (`Itens`) each record can hold. Schemas are passed to the database through `api.Props`.

`Props`, `Schema`, and `Item` are **interfaces**, so they are built with the constructors the `lib` package exports rather than composite literals:

```go
func createProps() api.Props {
	user := lib.NewSchema("user" /* , fields... */)
	return lib.NewProps("myDatabase/", user)
}
```

Naming each part before passing it on keeps the description readable as it grows — the nested collection below is built the same way, one variable at a time.

- The first argument is a prefix added to every key the database stores. With the standard (filesystem) adapter it behaves like a folder.
- `GetSchema(name)` returns the collection, or `nil` if no schema has that name.

---

## Fields (`api.Item`)

Every field has a name, a type, and — for plain fields — a `required` flag. Required fields must be present when creating a record.

| Constructor | `Type()` | Holds | Notes |
|---|---|---|---|
| `lib.NewKeyItem(name, required)` | `api.KeyItem` | `string` | Unique and indexed: two records can never share the same value, and `FindByKey` can look records up by it. Uniqueness is case-insensitive (`User@x.com` and `user@x.com` conflict). |
| `lib.NewIntItem(name, required)` | `api.IntItem` | `int`, `int32`, or `int64` | Plain integer field, always read back as `int64`. |
| `lib.NewDatabaseItem(name, itens...)` | `api.DatabaseItem` | a nested collection | The field is itself a sub-database with its own `Itens()`. See below. |

---

## Example

A `user` collection where each user owns a nested `sessions` collection:

Build the innermost collection first, then the schema that owns it:

```go
func createProps() api.Props {

	//========================Sessions==========================
	token := lib.NewKeyItem("token", true)
	creation := lib.NewIntItem("creation", true)
	expiration := lib.NewIntItem("expiration", true)
	sessions := lib.NewDatabaseItem("sessions", token, creation, expiration)

	//========================User==========================
	email := lib.NewKeyItem("email", true)
	username := lib.NewKeyItem("username", true)
	age := lib.NewIntItem("age", true)
	user := lib.NewSchema("user", email, username, age, sessions)

	//========================Props==========================
	return lib.NewProps("myDatabase/", user)
}
```

---

## Sub-databases

A `DatabaseItem` field gives each record its own private collection. In the example above, every user has its own list of sessions, isolated from other users' sessions.

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

Collection and field names must not contain the `-` character — it is the separator used internally in the key layout (see [DenseRecordPattern.md](/docs/Explanations/DenseRecordPattern.md)).
