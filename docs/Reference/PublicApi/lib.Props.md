# `lib.Props`, `lib.Schema`, `lib.Item`

**Type:** Structs (schema description)

## Definition

```go
type Props struct {
	Path    string
	Schemas []Schema
}

type Schema struct {
	Name  string
	Itens []Item
}

type Item struct {
	Name     string
	Type     ItemType
	Required bool
	Itens    []Item // only for Type == Database
}

type ItemType int

const (
	Key ItemType = iota
	Int
	Database
)
```

## Description

The declarative description of a database, passed to [`Lib.NewDatabase`](./lib.Lib.md#methods). `Path` is a prefix added to every stored key (a folder, with the standard adapter); each `Schema` is one collection; each `Item` is one field. The full guide, including field types and nested sub-databases, is in [Schemas](/docs/Explanation/Schemas.md).

## Field Types

| Type | Holds | Notes |
| :--- | :--- | :--- |
| `lib.Key` | `string` | Unique and indexed, case-insensitive; usable with `FindByKey`. |
| `lib.Int` | `int`, `int32`, or `int64` | Always read back as `int64`. |
| `lib.Database` | a nested collection | The field is a sub-database with its own `Itens`. |

## Examples

```go
var Props = lib.Props{
	Path: "myDatabase/",
	Schemas: []lib.Schema{
		{
			Name: "user",
			Itens: []lib.Item{
				{Name: "email", Type: lib.Key, Required: true},
				{Name: "age", Type: lib.Int, Required: true},
			},
		},
	},
}
```
