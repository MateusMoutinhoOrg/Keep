# Working with Records

## Description
Explains every operation a collection and its records support: creating, finding, reading, updating, deleting, listing, and managing sub-databases. To describe the collections themselves, see [Schemas.md](/docs/References/Schemas.md).

---

## Starting point

All operations start from a schema instance:

```go
users, ok := db.GetSchema("user")
```

Each runnable example below has a full version in [examples/libraryExamples/](../../examples/libraryExamples/).

---

## Create — `NewItem`

```go
user, err := users.NewItem(map[string]any{
	"email":    "mateus@gmail.com",
	"username": "mateus",
	"age":      27,
})
```

Fails with `MissingField` if a required field is absent, `InvalidField` if a field is not in the schema or has the wrong type, and `KeyConflict` if a `Key` value is already taken. The returned record has a permanent id (`user.Id`) that is never reused, even after deletion.

Full example: [examples/libraryExamples/CreateUser](../../examples/libraryExamples/CreateUserSample/CreateUserSample.go)

---

## Find — `FindByKey`

Looks a record up by any `Key` field. Constant cost, case-insensitive.

```go
user, ok := users.FindByKey("email", "mateus@gmail.com")
if !ok {
	fmt.Println("not found")
}
```

Full example: [examples/libraryExamples/FindUserByKey](../../examples/libraryExamples/FindUserByKeySample/FindUserByKeySample.go)

---

## Find — `FindById`

Looks a record up by its permanent id — the value `user.Id` reports. No index entry is read: the id alone names the record's keys, which makes it the cheapest lookup the library has.

```go
user, ok := users.FindById(1)
if !ok {
	fmt.Println("no record under that id")
}
```

`ok` is `false` when the collection holds no live record under that id — either it was never allocated, or the record was removed. Ids are **never reused**, so an id kept from a previous run either resolves to the same record it always named or to nothing at all; it can never come back pointing at a different record.

Full example: [examples/libraryExamples/FindUserById](../../examples/libraryExamples/FindUserByIdSample/FindUserByIdSample.go)

---

## Relations between collections

Because ids are permanent and never reused, an id is a safe value to **store**. Putting one in an `Int` field of another collection turns that field into a pointer to a record — a foreign key — and `FindById` is what follows it. This is how one collection references another without the engine needing any join, key listing, or scan.

Describe the pointer as a plain `Int` field:

```go
var Schemas = []keeptypes.Schema{
	{
		Name: "user",
		Itens: []keeptypes.Item{
			{Name: "email", Type: keeptypes.Key, Required: true},
		},
	},
	{
		Name: "post",
		Itens: []keeptypes.Item{
			{Name: "title", Type: keeptypes.Key, Required: true},
			// The foreign key: it holds a user record's Id.
			{Name: "author", Type: keeptypes.Int, Required: true},
		},
	},
}
```

Write the pointer by storing the target's `Id`:

```go
author, _ := users.NewItem(map[string]any{"email": "mateus@gmail.com"})

post, _ := posts.NewItem(map[string]any{
	"title":  "keep-by-id",
	"author": author.Id,
})
```

Follow it by reading the field back (an `Int` field comes back as `int64`) and resolving it:

```go
authorId, _ := post.Get("author")
foundAuthor, ok := users.FindById(authorId.(int64))
```

Two notes on using this:

- **Nothing enforces the reference.** Removing the author leaves the post's `author` field pointing at an id that no longer resolves, and `FindById` reports `ok == false`. Deciding what happens then — deleting the post, clearing the field, showing a placeholder — is the caller's job.
- **A relation is not a sub-database.** Use a nested `Database` field when the sub-records belong to exactly one parent and die with it (see [Sub-databases](#sub-databases--newsubitem-and-listallfield)); use a stored id when the target is an independent record other collections may point at too.

Full example: [examples/libraryExamples/FindUserById](../../examples/libraryExamples/FindUserByIdSample/FindUserByIdSample.go)

---

## Read — `Get`

```go
age, err := user.Get("age") // int64(27)
```

`Key` fields come back as `string`, `Int` fields as `int64`. Returns a `NotFound` error if the record never stored that field (possible for non-required fields).

Full example: [examples/libraryExamples/RetrieveUserInfo](../../examples/libraryExamples/RetrieveUserInfoSample/RetrieveUserInfoSample.go)

---

## Update — `Update`

```go
err := user.Update("age", 28)
```

Works for any plain field. Updating a `Key` field re-indexes it and fails with `KeyConflict` if another record already owns the new value:

```go
err := user.Update("email", "newmail@gmail.com")
```

Full examples: [examples/libraryExamples/UpdateUser](../../examples/libraryExamples/UpdateUserSample/UpdateUserSample.go), [examples/libraryExamples/UpdateUserKey](../../examples/libraryExamples/UpdateUserKeySample/UpdateUserKeySample.go)

---

## Delete — `Remove`

```go
e := user.Remove()
if e != nil {
	fmt.Println("error removing:", e.Message)
}
```

Removes the record, its unique index entries, and everything inside its sub-databases. Removing an already-removed record is a no-op.

Full example: [examples/libraryExamples/DeleteUser](../../examples/libraryExamples/DeleteUserSample/DeleteUserSample.go)

---

## List — `ListAll` and `List`

```go
all, err := users.ListAll()

// Pagination: up to 10 records starting at position 1 (positions are 1-based)
page, err := users.List(1, 10)
```

**List order is not stable.** Deleting a record moves the last record into the freed position (this is what keeps deletion constant-cost). If you need a stable order, store it as a field on the record.

Full examples: [examples/libraryExamples/ListAllUsers](../../examples/libraryExamples/ListAllUsersSample/ListAllUsersSample.go), [examples/libraryExamples/ListUsersPaginated](../../examples/libraryExamples/ListUsersPaginatedSample/ListUsersPaginatedSample.go)

---

## Sub-databases — `NewSubItem` and `ListAll(field)`

For fields of type `keeptypes.Database` (see [Schemas](Schemas.md)):

```go
session, err := user.NewSubItem("sessions", map[string]any{
	"token":      "token-1",
	"creation":   1000,
	"expiration": 2000,
})

for _, s := range user.ListAll("sessions") {
	token, _ := s.Get("token")
	fmt.Println(token)
}
```

Full example: [examples/libraryExamples/SubInfos](../../examples/libraryExamples/SubInfosSample/SubInfosSample.go)

---

## Other helpers

- `user.Id` — the record's permanent identifier (a plain data field, not a function).
- `user.CheckKeysPresence([]string{"email", "age"})` — reports whether every named field has a stored value.
- `user.String()` — renders the record's plain fields as `{id: 1, email: ..., username: ..., age: ...}`.

---

## Concurrency

Keep assumes a **single writer** unless the backend provides atomic operations. Multiple readers are always safe. See [DenseRecordPattern.md](/docs/References/DenseRecordPattern.md#concurrency-and-atomicity).
