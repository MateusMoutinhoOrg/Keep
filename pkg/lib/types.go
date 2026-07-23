package lib

import "github.com/MateusMoutinhoOrg/Keep/pkg/deps"

// ItemType is the type of a schema field.
type ItemType int

const (
	// Key is a unique, indexed string field.
	Key ItemType = iota
	// Int is a plain integer field.
	Int
	// Database is a nested collection of records.
	Database
)

// Item describes one field of a schema. Nested collections (Database
// type) describe their own fields in Itens.
type Item struct {
	Name     string
	Type     ItemType
	Required bool
	Itens    []Item
}

// Schema describes one collection of records and its fields.
type Schema struct {
	Name  string
	Itens []Item
}

// Props is the declarative description of a database: the prefix every
// key is written under, and the collections it holds.
type Props struct {
	Path    string
	Schemas []Schema
}

// KeepDatabase is a database bound to the injected dependencies and to
// a Props description. Created by Lib.NewDatabase.
type KeepDatabase struct {
	deps  deps.Deps
	Props Props
}

// SchemaInstance is one collection of records, created by
// KeepDatabase.GetSchema.
type SchemaInstance struct {
	deps   deps.Deps
	items  []Item
	prefix string
}

// SchemaItem is one record of a collection, created by
// SchemaInstance.NewItem, FindByKey, ListAll, or List.
type SchemaItem struct {
	deps   deps.Deps
	items  []Item
	prefix string
	id     int64
}

// ErrorType classifies an Error so callers can switch on the cause
// instead of matching messages.
type ErrorType int

const (
	KeyConflict ErrorType = iota
	NotFound
	MissingField
	InvalidField
	Internal
)

// Error is the typed error returned by database operations.
type Error struct {
	Type     ErrorType
	Key      string
	KeyValue any
	Msg      string
}

func (e Error) Error() string { return e.Msg }
