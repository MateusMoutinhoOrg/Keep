package api

// This file is the whole of Keep's database surface: the Databases contract
// that reaches the Sandbox as the field of the same name, and the types that
// cross it in either direction.
//
// Every type here is a struct — never an interface. A type that carries
// behaviour holds it as function fields, filled by sandbox/internal/databases
// at the moment the object is built; a type that carries none (Item, Schema,
// Props, Error) is plain data a caller writes as a composite literal. None of
// them names the dependencies the library was wired with: a record reaches
// storage through the closure it was built with, not through a field a caller
// could read or replace.

// Field types, reported by Item.Type.
const (
	// Key is a unique, indexed string field. Two live records of one
	// collection can never hold the same value for it, and it is the only
	// kind of field SchemaInstance.FindByKey looks a record up by.
	Key = iota
	// Int is a plain integer field, stored in its decimal form and handed
	// back as an int64.
	Int
	// Database is a nested collection of records, reached through
	// SchemaItem.NewSubItem and SchemaItem.ListAll rather than read as a
	// value.
	Database
)

// Failure causes, reported by Error.Type. Switch on the constant rather
// than matching Error.Message: the message is written for a person and may
// change between releases, the constant may not.
const (
	// KeyConflict is a value another live record already holds for a Key
	// field.
	KeyConflict = iota
	// NotFound is a field of an existing record that has no stored value.
	NotFound
	// MissingField is a Required field absent from an insert.
	MissingField
	// InvalidField is a field the schema does not declare, a value of the
	// wrong Go type for the field it is written to, or a Database field
	// used where a plain value was expected.
	InvalidField
	// Internal is a failure the storage backend reported. Error.Message
	// carries what the backend said.
	Internal
)

// Item describes one field of a schema.
type Item struct {
	// Name is the field's name, as used in the fields map of an insert and
	// in SchemaItem.Get.
	Name string
	// Type is one of Key, Int or Database.
	Type int
	// Required reports whether an insert must provide this field. It is
	// ignored on a Database field, which is never provided directly.
	Required bool
	// Itens are the nested fields, for a Database field; nil otherwise.
	Itens []Item
}

// Schema describes one collection of records and the fields each of them
// can hold.
type Schema struct {
	// Name is the collection's name, as used in DatabaseHandle.GetSchema.
	Name string
	// Itens are the fields each record of the collection can hold.
	Itens []Item
}

// Props is the declarative description a database is created from. It is
// the only thing Database.New takes, so a database is fully described by a
// value a caller can write, read and version.
type Props struct {
	// Path is the prefix every key of the database is written under. It is
	// split on slashes into the leading segments of every key, so a backend
	// that maps keys to files reads it as a directory; empty segments are
	// dropped, which makes a trailing slash optional.
	Path string
	// Schemas are the collections the database holds.
	Schemas []Schema
}

// Error describes one failure reported by a database operation. It carries
// no behaviour, so a caller switches on Type and reads Key, KeyValue and
// Message directly. A nil *Error means success.
type Error struct {
	// Type is one of KeyConflict, NotFound, MissingField, InvalidField or
	// Internal.
	Type int
	// Key is the name of the field the failure involves, empty when the
	// failure involves no particular field.
	Key string
	// KeyValue is the value the failure involves, when there is one.
	KeyValue any
	// Message is the human-readable description of the failure.
	Message string
}

// SchemaItem is one record of a collection, handed back by
// SchemaInstance.NewItem, FindByKey, FindById, ListAll and List, and by
// SchemaItem.ListAll and NewSubItem for a nested collection.
type SchemaItem struct {
	// Items are the fields the record's own collection declares.
	Items []Item
	// Prefix is the key prefix of the collection the record belongs to,
	// held as the list of segments every key under it is built from.
	Prefix []string
	// Id is the record's permanent identifier. It is never reused, so an
	// id stored in an Int field of another collection stays a reference to
	// this record or to nothing at all — never to a different record.
	Id int64
	// Get returns the typed value stored for a field: a string for a Key
	// field, an int64 for an Int field. It fails with NotFound when the
	// field has no stored value and with InvalidField when the schema
	// declares no such field or the field is a nested collection.
	Get func(fieldName string) (any, *Error)
	// Update writes a new value for a field, re-indexing it when the field
	// is a Key. It fails with KeyConflict when another live record already
	// holds the new value for that Key.
	Update func(fieldName string, value any) *Error
	// Remove deletes the record, its index entries and every record of
	// every collection nested under it. Removing a record that is already
	// gone is not an error.
	Remove func() *Error
	// CheckKeysPresence reports whether every named field has a stored
	// value for this record.
	CheckKeysPresence func(keys []string) bool
	// ListAll returns every record of a nested (Database) field, and nil
	// when the schema declares no such nested field.
	ListAll func(fieldName string) []SchemaItem
	// NewSubItem inserts a record into a nested (Database) field, taking
	// the same fields map SchemaInstance.NewItem takes.
	NewSubItem func(fieldName string, fields map[string]any) (SchemaItem, *Error)
	// String renders the record's id and its plain fields, for printing.
	String func() string
}

// SchemaInstance is one collection of records, handed back by
// DatabaseHandle.GetSchema.
type SchemaInstance struct {
	// Items are the fields each record of the collection can hold.
	Items []Item
	// Prefix is the collection's key prefix, held as the list of segments
	// every key under it is built from.
	Prefix []string
	// NewItem inserts a record, validating the fields against the schema
	// and against the unique index of every Key field.
	NewItem func(fields map[string]any) (SchemaItem, *Error)
	// FindByKey looks a record up through a unique Key field. ok is false
	// when the schema declares no such Key field, or when no live record
	// holds that value.
	FindByKey func(key string, keyValue any) (SchemaItem, bool)
	// FindById looks a record up through its permanent id — the value
	// SchemaItem.Id reports. ok is false when the collection holds no live
	// record under that id.
	FindById func(id int64) (SchemaItem, bool)
	// ListAll returns every record of the collection, in the order the
	// dense position list holds them.
	ListAll func() ([]SchemaItem, *Error)
	// List returns up to chunk records starting at position, counted from
	// 1. A chunk of 0 means "to the end of the collection".
	List func(position int, chunk int) ([]SchemaItem, *Error)
}

// DatabaseHandle is one database, bound to the Props it was described with.
type DatabaseHandle struct {
	// Props is the description the database was created from.
	Props Props
	// GetSchema returns the collection with the given name. ok is false
	// when the Props declares no schema under that name.
	GetSchema func(name string) (SchemaInstance, bool)
}

// Databases is the contract every database is reached through, carried by
// the Sandbox as the field of the same name.
type Databases struct {
	// New builds a database from a Props description. It touches no key:
	// a database is a handle over a prefix, so building one is free and
	// creates nothing until the first record is written.
	New func(props Props) DatabaseHandle
}
