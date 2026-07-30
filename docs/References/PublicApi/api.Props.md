# `api.Props`, `api.Schema`, `api.Item`

**Type:** Interfaces (schema description)

## Definition

```go
type Props interface {
	Path() string
	Schemas() []Schema
}

type Schema interface {
	Name() string
	Itens() []Item
}

type Item interface {
	Name() string
	Type() int // one of KeyItem, IntItem, DatabaseItem
	Required() bool
	Itens() []Item // only for Type() == DatabaseItem
}

const (
	KeyItem = iota
	IntItem
	DatabaseItem
)
```

## Description

The declarative description of a database, passed to [`Lib.NewDatabase`](./api.Lib.md#methods). `Path` is a prefix added to every stored key (a folder, with the standard adapter); each `Schema` is one collection; each `Item` is one field.

They are **interfaces**, not structs, because every value crossing the library boundary must be a primitive or an interface. An interface cannot be built with a composite literal, so you build one with the constructors below. The full guide, including field types and nested sub-databases, is in [Schemas](/docs/Explanations/Schemas.md).

## Constructors

```go
func NewProps(path string, schemas ...api.Schema) api.Props
func NewSchema(name string, itens ...api.Item) api.Schema
func NewKeyItem(name string, required bool) api.Item
func NewIntItem(name string, required bool) api.Item
func NewDatabaseItem(name string, itens ...api.Item) api.Item
```

All five live in the `sandbox` package, imported as `lib`.

## Field Types

| Constructor | `Type()` | Holds | Notes |
| :--- | :--- | :--- | :--- |
| `lib.NewKeyItem` | `api.KeyItem` | `string` | Unique and indexed, case-insensitive; usable with `FindByKey`. |
| `lib.NewIntItem` | `api.IntItem` | `int`, `int32`, or `int64` | Always read back as `int64`. |
| `lib.NewDatabaseItem` | `api.DatabaseItem` | a nested collection | The field is a sub-database with its own `Itens()`. |

## Examples

```go
import lib "github.com/MateusMoutinhoOrg/Keep/sandbox"

var Props = lib.NewProps("myDatabase/",
	lib.NewSchema("user",
		lib.NewKeyItem("email", true),
		lib.NewIntItem("age", true),
		lib.NewDatabaseItem("sessions",
			lib.NewKeyItem("token", true),
			lib.NewIntItem("creation", true),
		),
	),
)
```
