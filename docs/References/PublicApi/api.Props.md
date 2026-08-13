# `api.Props`, `api.Schema`, `api.Item`

**Type:** Structs (plain data, no behavior)

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
	Type     int // one of Key, Int, Database
	Required bool
	Itens    []Item // only when Type == Database
}

const (
	Key = iota
	Int
	Database
)
```

## Description

The declarative description of a database, passed to [`Lib.NewDatabase`](./api.Lib.md#fields). `Path` is a prefix added to every stored key (a folder, with the standard adapter); each `Schema` is one collection; each `Item` is one field.

They are **plain structs**, not interfaces: `Item`, `Schema`, `Props`, and `Error` carry no behavior and no `Deps` field, so they are built directly with a composite literal — there is no `lib.NewSchema` / `lib.NewKeyItem` / `lib.NewProps` constructor to call (`sandbox/description.go` was removed along with them). The full guide, including field types and nested sub-databases, is in [Schemas](/docs/References/Schemas.md).

## Field Types

| `Type` value | Holds | Notes |
| :--- | :--- | :--- |
| `keeptypes.Key` | `string` | Unique and indexed, case-insensitive; usable with `FindByKey`. |
| `keeptypes.Int` | `int`, `int32`, or `int64` | Always read back as `int64`. |
| `keeptypes.Database` | a nested collection | The field is a sub-database with its own `Itens`. |

## Examples

```go
import keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"

var Schemas = []keeptypes.Schema{
	{
		Name: "user",
		Itens: []keeptypes.Item{
			{Name: "email", Type: keeptypes.Key, Required: true},
			{Name: "age", Type: keeptypes.Int, Required: true},
			{
				Name: "sessions",
				Type: keeptypes.Database,
				Itens: []keeptypes.Item{
					{Name: "token", Type: keeptypes.Key, Required: true},
					{Name: "creation", Type: keeptypes.Int, Required: true},
				},
			},
		},
	},
}

var Props = keeptypes.Props{
	Path:    "myDatabase/",
	Schemas: Schemas,
}
```

Every sample under `examples/libraryExamples/` builds `Props` this way: a package-level `Schemas` value listing each collection, and a `Props` value wrapping it with the key prefix.
