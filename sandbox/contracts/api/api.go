package api

// This package is the library's whole public surface, and it declares
// nothing but interfaces and constants. Every value crossing the
// boundary — in or out — is either a primitive or one of the interfaces
// below, so a consumer never depends on a concrete type of the library.
// The structs implementing these interfaces live in sandbox/internal/
// and are unreachable from outside the sandbox.

// Field types, reported by Item.Type.
const (
	// KeyItem is a unique, indexed string field.
	KeyItem = iota
	// IntItem is a plain integer field.
	IntItem
	// DatabaseItem is a nested collection of records.
	DatabaseItem
)

// Failure causes, reported by Error.Type. Switch on them instead of
// matching messages.
const (
	KeyConflict = iota
	NotFound
	MissingField
	InvalidField
	Internal
)

// Item describes one field of a schema. Build one with lib.NewKeyItem,
// lib.NewIntItem, or lib.NewDatabaseItem.
type Item interface {
	// Name is the field's name, as used in the fields map.
	Name() string
	// Type is one of KeyItem, IntItem, or DatabaseItem.
	Type() int
	// Required reports whether a record must provide this field.
	Required() bool
	// Itens are the nested fields, for a DatabaseItem; nil otherwise.
	Itens() []Item
}

// Schema describes one collection of records and its fields. Build one
// with lib.NewSchema.
type Schema interface {
	// Name is the collection's name, as used in GetSchema.
	Name() string
	// Itens are the fields each record of the collection can hold.
	Itens() []Item
}

// Props is the declarative description of a database. Build one with
// lib.NewProps.
type Props interface {
	// Path is the prefix every key of the database is written under.
	Path() string
	// Schemas are the collections the database holds.
	Schemas() []Schema
}

// Error is the typed error returned by database operations. It
// satisfies the standard error interface, so it can be returned and
// compared as one.
type Error interface {
	error
	// Type is one of KeyConflict, NotFound, MissingField, InvalidField,
	// or Internal.
	Type() int
	// Key is the field the failure involves.
	Key() string
	// KeyValue is the value the failure involves, when relevant.
	KeyValue() any
}

// SchemaItem is one record of a collection. It is handed back by
// SchemaInstance.NewItem, FindByKey, ListAll and List, and carries the
// deps it was built with, so every field read or write goes through the
// same injected backend.
type SchemaItem interface {
	// Id returns the record's permanent, never-reused identifier.
	Id() int64
	// Get returns the typed value stored for a field.
	Get(fieldName string) (any, Error)
	// Update writes a new value for a field, re-indexing it when the
	// field is a unique key.
	Update(fieldName string, value any) Error
	// Remove deletes the record and everything nested under it,
	// returning nil on success.
	Remove() Error
	// CheckKeysPresence reports whether every named field has a stored
	// value for this record.
	CheckKeysPresence(keys []string) bool
	// ListAll returns every record of a nested (DatabaseItem) field.
	ListAll(fieldName string) []SchemaItem
	// NewSubItem inserts a record into a nested (DatabaseItem) field.
	NewSubItem(fieldName string, fields map[string]any) (SchemaItem, Error)
	// String renders the record's plain fields.
	String() string
}

// SchemaInstance is one collection of records, handed back by
// KeepDatabase.GetSchema.
type SchemaInstance interface {
	// NewItem inserts a record, validating the fields against the schema.
	NewItem(fields map[string]any) (SchemaItem, Error)
	// FindByKey looks a record up through a unique KeyItem field. It
	// returns nil when the field is not an indexed key or no record
	// matches.
	FindByKey(key string, keyValue any) SchemaItem
	// ListAll returns every record of the collection.
	ListAll() ([]SchemaItem, Error)
	// List returns up to chunk records starting at position (1-based).
	List(position int, chunk int) ([]SchemaItem, Error)
}

// KeepDatabase is a database bound to a Props description and to the
// injected deps, handed back by Lib.NewDatabase.
type KeepDatabase interface {
	// GetSchema returns the collection with the given name, or nil when
	// the database declares no schema under that name.
	GetSchema(name string) SchemaInstance
	// Props returns the description the database was created from.
	Props() Props
}

// Lib is the entry point handed back by lib.New. Every object it
// creates carries the same deps it was built with.
type Lib interface {
	// NewDatabase creates a database from a Props description.
	NewDatabase(props Props) KeepDatabase
}
