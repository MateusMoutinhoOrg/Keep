# `sandbox/api/databases.go`

| Constant | Value | Description |
| --- | --- | --- |
| `Key` | `iota` | Key is a unique, indexed string field. Two live records of one collection can never hold the same value for it, and it is the only kind of field SchemaInstance.FindByKey looks a record up by. |
| `Int` |  | Int is a plain integer field, stored in its decimal form and handed back as an int64. |
| `Database` |  | Database is a nested collection of records, reached through SchemaItem.NewSubItem and SchemaItem.ListAll rather than read as a value. |
| `Float` |  | Float is a plain floating-point field, stored in the shortest decimal form that parses back to the same number and handed back as a float64. |
| `String` |  | String is a plain text field. It is written like a Key and read back as a string, but it carries no index: two live records of one collection may hold the same value for it, and SchemaInstance.FindByKey never looks a record up by one. |
| `Link` |  | Link is a reference to a record of another collection, named by the field's Target. It is stored as that record's id and read back as an int64, and SchemaItem.GetLink resolves it to the record itself. Like any id it is never reused, so a link to a removed record resolves to nothing rather than to whatever took its place. |
| `Bytes` |  | Bytes is a plain binary field. It is written as a []byte and stored exactly as given, with no text encoding, so any content — a file, an image, a hash — comes back byte for byte as a []byte. Like a String it carries no index, and SchemaItem.String prints its length rather than its contents. |
| `KeyConflict` | `iota` | KeyConflict is a value another live record already holds for a Key field. |
| `NotFound` |  | NotFound is a field of an existing record that has no stored value. |
| `MissingField` |  | MissingField is a Required field absent from an insert. |
| `InvalidField` |  | InvalidField is a field the schema does not declare, a value of the wrong Go type for the field it is written to, a Database field used where a plain value was expected, or a Link field declaring no Target. |
| `Internal` |  | Internal is a failure the storage backend reported. Error.Message carries what the backend said. |

## `Item`

Item describes one field of a schema.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the field's name, as used in the fields map of an insert and in SchemaItem.Get. |
| `Type` | `int` | Type is one of Key, Int, Float, String, Link, Bytes or Database. |
| `Required` | `bool` | Required reports whether an insert must provide this field. It is ignored on a Database field, which is never provided directly. |
| `Target` | `string` | Target is the name of the schema a Link field points at, as used in DatabaseHandle.GetSchema. It is empty on every other kind of field. |
| `Itens` | `[]Item` | Itens are the nested fields, for a Database field; nil otherwise. |

## `Schema`

Schema describes one collection of records and the fields each of them can hold.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the collection's name, as used in DatabaseHandle.GetSchema. |
| `Itens` | `[]Item` | Itens are the fields each record of the collection can hold. |

## `Props`

Props is the declarative description a database is created from. It is the only thing Database.New takes, so a database is fully described by a value a caller can write, read and version.

| Field | Type | Description |
| --- | --- | --- |
| `Path` | `string` | Path is the prefix every key of the database is written under. It is split on slashes into the leading segments of every key, so a backend that maps keys to files reads it as a directory; empty segments are dropped, which makes a trailing slash optional. |
| `Schemas` | `[]Schema` | Schemas are the collections the database holds. |

## `Error`

Error describes one failure reported by a database operation. It carries no behaviour, so a caller switches on Type and reads Key, KeyValue and Message directly. A nil *Error means success.

| Field | Type | Description |
| --- | --- | --- |
| `Type` | `int` | Type is one of KeyConflict, NotFound, MissingField, InvalidField or Internal. |
| `Key` | `string` | Key is the name of the field the failure involves, empty when the failure involves no particular field. |
| `KeyValue` | `any` | KeyValue is the value the failure involves, when there is one. |
| `Message` | `string` | Message is the human-readable description of the failure. |

## `SchemaItem`

SchemaItem is one record of a collection, handed back by SchemaInstance.NewItem, FindByKey, FindById, ListAll and List, and by SchemaItem.ListAll and NewSubItem for a nested collection.

| Field | Type | Description |
| --- | --- | --- |
| `Items` | `[]Item` | Items are the fields the record's own collection declares. |
| `Prefix` | `[]string` | Prefix is the key prefix of the collection the record belongs to, held as the list of segments every key under it is built from. |
| `Id` | `int64` | Id is the record's permanent identifier. It is never reused, so an id stored in a Link or Int field of another collection stays a reference to this record or to nothing at all — never to a different record. |
| `Get` | `func(fieldName string) (any, *Error)` | Get returns the typed value stored for a field: a string for a Key or String field, an int64 for an Int or Link field, a float64 for a Float field, a []byte for a Bytes field. It fails with NotFound when the field has no stored value and with InvalidField when the schema declares no such field or the field is a nested collection. |
| `GetLink` | `func(fieldName string) (SchemaItem, bool)` | GetLink resolves a Link field to the record it points at, looked up by id in the collection the field's Target names. ok is false when the schema declares no such Link field, when the field holds no stored value, when the database declares no schema under that Target, or when the record the stored id names is no longer live. |
| `Update` | `func(fieldName string, value any) *Error` | Update writes a new value for a field, re-indexing it when the field is a Key. It fails with KeyConflict when another live record already holds the new value for that Key. |
| `Remove` | `func() *Error` | Remove deletes the record, its index entries and every record of every collection nested under it. Removing a record that is already gone is not an error. |
| `CheckKeysPresence` | `func(keys []string) bool` | CheckKeysPresence reports whether every named field has a stored value for this record. |
| `ListAll` | `func(fieldName string) []SchemaItem` | ListAll returns every record of a nested (Database) field, and nil when the schema declares no such nested field. |
| `NewSubItem` | `func(fieldName string, fields map[string]any) (SchemaItem, *Error)` | NewSubItem inserts a record into a nested (Database) field, taking the same fields map SchemaInstance.NewItem takes. |
| `String` | `func() string` | String renders the record's id and its plain fields, for printing. |

## `SchemaInstance`

SchemaInstance is one collection of records, handed back by DatabaseHandle.GetSchema.

| Field | Type | Description |
| --- | --- | --- |
| `Items` | `[]Item` | Items are the fields each record of the collection can hold. |
| `Prefix` | `[]string` | Prefix is the collection's key prefix, held as the list of segments every key under it is built from. |
| `NewItem` | `func(fields map[string]any) (SchemaItem, *Error)` | NewItem inserts a record, validating the fields against the schema and against the unique index of every Key field. |
| `FindByKey` | `func(key string, keyValue any) (SchemaItem, bool)` | FindByKey looks a record up through a unique Key field. ok is false when the schema declares no such Key field, or when no live record holds that value. |
| `FindById` | `func(id int64) (SchemaItem, bool)` | FindById looks a record up through its permanent id — the value SchemaItem.Id reports. ok is false when the collection holds no live record under that id. |
| `ListAll` | `func() ([]SchemaItem, *Error)` | ListAll returns every record of the collection, in the order the dense position list holds them. |
| `List` | `func(position int, chunk int) ([]SchemaItem, *Error)` | List returns up to chunk records starting at position, counted from 1. A chunk of 0 means "to the end of the collection". |

## `DatabaseHandle`

DatabaseHandle is one database, bound to the Props it was described with.

| Field | Type | Description |
| --- | --- | --- |
| `Props` | `Props` | Props is the description the database was created from. |
| `GetSchema` | `func(name string) (SchemaInstance, bool)` | GetSchema returns the collection with the given name. ok is false when the Props declares no schema under that name. |

## `Databases`

Databases is the contract every database is reached through, carried by the Sandbox as the field of the same name.

| Field | Type | Description |
| --- | --- | --- |
| `New` | `func(props Props) DatabaseHandle` | New builds a database from a Props description. It touches no key: a database is a handle over a prefix, so building one is free and creates nothing until the first record is written. |

[every contract](doc.md)
