package api

// This package is the library's whole public surface, and it declares
// nothing but types: the objects handed back (Lib, KeepDatabase,
// SchemaInstance, SchemaItem), the description passed in (Props, Schema,
// Item), and Error. Every type here is a struct — never an interface —
// following the struct-of-function-fields pattern: a type that carries
// behavior leads with a Deps field and fills the rest of its behavior as
// function fields, each assigned by a factory in sandbox/lib/ (see
// docs/References/StructContracts.md). A type that carries no behavior
// (Item, Schema, Props, Error) is plain data and can be built directly
// with a composite literal — no constructor required.

import "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"

// Field types, reported by Item.Type.
const (
	// Key is a unique, indexed string field.
	Key = iota
	// Int is a plain integer field.
	Int
	// Database is a nested collection of records.
	Database
)

// Failure causes, reported by Error.Type.
const (
	KeyConflict = iota
	NotFound
	MissingField
	InvalidField
	Internal
)

// Item describes one field of a schema.
type Item struct {
	// Name is the field's name, as used in the fields map.
	Name string
	// Type is one of Key, Int, or Database.
	Type int
	// Required reports whether a record must provide this field.
	Required bool
	// Itens are the nested fields, for a Database field; nil otherwise.
	Itens []Item
}

// Schema describes one collection of records and its fields.
type Schema struct {
	// Name is the collection's name, as used in KeepDatabase.GetSchema.
	Name string
	// Itens are the fields each record of the collection can hold.
	Itens []Item
}

// Props is the declarative description of a database.
type Props struct {
	// Path is the prefix every key of the database is written under.
	Path string
	// Schemas are the collections the database holds.
	Schemas []Schema
}

// Error describes one failure reported by a database operation. It
// carries no behavior, so callers switch on Type and read Key, KeyValue
// and Message directly, and a nil *Error means success.
type Error struct {
	// Type is one of KeyConflict, NotFound, MissingField, InvalidField,
	// or Internal.
	Type int
	// Key is the field the failure involves.
	Key string
	// KeyValue is the value the failure involves, when relevant.
	KeyValue any
	// Message is the human-readable description of the failure.
	Message string
}

// SchemaItem is one record of a collection. It is handed back by
// SchemaInstance.NewItem, FindByKey, ListAll and List, and carries the
// Deps it was built with, so every field read or write goes through the
// same injected backend. Its function fields are filled by factories in
// sandbox/lib/schemaitem.
type SchemaItem struct {
	// Deps is the dependency set the record was built with.
	Deps deps.Deps
	// Items are the schema fields this record's collection declares.
	Items []Item
	// Prefix is the collection's key prefix.
	Prefix string
	// Id is the record's permanent, never-reused identifier.
	Id int64
	// Get returns the typed value stored for a field.
	Get func(fieldName string) (any, *Error)
	// Update writes a new value for a field, re-indexing it when the
	// field is a unique key.
	Update func(fieldName string, value any) *Error
	// Remove deletes the record and everything nested under it,
	// returning nil on success.
	Remove func() *Error
	// CheckKeysPresence reports whether every named field has a stored
	// value for this record.
	CheckKeysPresence func(keys []string) bool
	// ListAll returns every record of a nested (Database) field.
	ListAll func(fieldName string) []SchemaItem
	// NewSubItem inserts a record into a nested (Database) field.
	NewSubItem func(fieldName string, fields map[string]any) (SchemaItem, *Error)
	// String renders the record's plain fields.
	String func() string
}

// SchemaInstance is one collection of records, handed back by
// KeepDatabase.GetSchema. Its function fields are filled by factories in
// sandbox/lib/schemainstance.
type SchemaInstance struct {
	// Deps is the dependency set the collection was built with.
	Deps deps.Deps
	// Items are the fields each record of the collection can hold.
	Items []Item
	// Prefix is the collection's key prefix.
	Prefix string
	// NewItem inserts a record, validating the fields against the schema.
	NewItem func(fields map[string]any) (SchemaItem, *Error)
	// FindByKey looks a record up through a unique Key field. ok is
	// false when the field is not an indexed key or no record matches.
	FindByKey func(key string, keyValue any) (SchemaItem, bool)
	// FindById looks a record up through its permanent id — the same
	// value SchemaItem.Id reports. Storing that id in an Int field of
	// another collection is how a record points at another record. ok
	// is false when no live record carries the id.
	FindById func(id int64) (SchemaItem, bool)
	// ListAll returns every record of the collection.
	ListAll func() ([]SchemaItem, *Error)
	// List returns up to chunk records starting at position (1-based).
	List func(position int, chunk int) ([]SchemaItem, *Error)
}

// KeepDatabase is a database bound to a Props description and to the
// injected Deps, handed back by Lib.NewDatabase. Its function fields are
// filled by factories in sandbox/lib/database.
type KeepDatabase struct {
	// Deps is the dependency set the database was built with.
	Deps deps.Deps
	// Props is the description the database was created from.
	Props Props
	// GetSchema returns the collection with the given name. ok is false
	// when the database declares no schema under that name.
	GetSchema func(name string) (SchemaInstance, bool)
}

// Lib is the entry point handed back by lib.New. Its function fields are
// filled by factories in sandbox/lib/publicfunctions.
type Lib struct {
	// Deps is the dependency set injected by lib.New.
	Deps deps.Deps
	// Version returns the library's own version — the release the
	// consumer linked against. It is held as a constant in
	// sandbox/config, so a release bump is a one-line edit touching no
	// logic, and it is exposed as a field like every other behavior so a
	// consumer can report it without importing anything but this api.
	Version func() string
	// NewDatabase creates a database from a Props description.
	NewDatabase func(props Props) KeepDatabase
}
