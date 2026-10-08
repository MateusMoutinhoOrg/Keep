# `sandbox/api/databases.go`

| Constant | Value | Description |
| --- | --- | --- |
| `Key` | `iota + 1` | Key is a unique, indexed string field. Two live records of one collection can never hold the same value for it, compared without regard to case, and it is the only kind of field Collection.FindByKey looks a record up by. |
| `Int` |  | Int is a plain integer field, stored in its decimal form and handed back as an int64. |
| `Nested` |  | Nested is a nested collection of records, reached through Record.Nested, Record.InsertNested and Record.ListNested rather than read as a value. |
| `Float` |  | Float is a plain floating-point field, stored in the shortest decimal form that parses back to the same number and handed back as a float64. |
| `String` |  | String is a plain text field. It is written like a Key and read back as a string, but it carries no index: two live records of one collection may hold the same value for it, and Collection.FindByKey never looks a record up by one. |
| `Link` |  | Link is a reference to a record of another collection, named by the field's Target. It is stored as that record's id and read back as an int64, and Record.GetLink resolves it to the record itself. Like any id it is never reused, so a link to a removed record resolves to nothing rather than to whatever took its place. |
| `Bytes` |  | Bytes is a plain binary field. It is written as a []byte and stored exactly as given, with no text encoding, so any content — a file, an image, a hash — comes back byte for byte as a []byte. Like a String it carries no index, and Record.String prints its length rather than its contents. |
| `KeyConflict` | `iota + 1` | KeyConflict is a value another live record already holds for a Key field. |
| `NoValue` |  | NoValue is a field of an existing record that has no stored value. |
| `MissingField` |  | MissingField is a Required field absent from an insert, given as nil to an insert, or cleared by an Update to nil. |
| `InvalidField` |  | InvalidField is a field the schema does not declare, a value of the wrong Go type for the field it is written to, a field of the wrong kind for the call (a Nested field where a plain value was expected, a plain field where a Key, a Link or a Nested one was), or a Link value that is a record of another collection. |
| `Internal` |  | Internal is a failure the storage backend reported, or a write lock that could not be taken in time. Error.Message carries what the backend said. |
| `Removed` |  | Removed is a write to a record that is no longer live — it was removed — or an insert into a nested collection whose owning record is no longer live. Nothing was written. |
| `InvalidSchema` |  | InvalidSchema is a Props that Databases.New refuses: an empty or repeated name, a missing or unknown Type, a Link with no Target or an unknown one, a Nested field under a reserved name, or a Path holding a "." or ".." segment. Error.Field holds the dotted path of the offending schema or field. |
| `InvalidArgument` |  | InvalidArgument is a call argument out of its range: a List position below 1 or a negative chunk. |

## `Field`

Field describes one field of a schema.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the field's name, as used in the fields map of an insert and in Record.Get. A Nested field may not be named "position" or "values": those names are the record's own keys. |
| `Type` | `int` | Type is one of Key, Int, Float, String, Link, Bytes or Nested. It has no default: a Field that leaves it out is refused. |
| `Required` | `bool` | Required reports whether an insert must provide this field, and whether an Update may clear it. It is ignored on a Nested field, which is never provided directly. |
| `Target` | `string` | Target is the name of the schema a Link field points at, as used in Database.Collection. It must name a schema of the same Props, and is empty on every other kind of field. |
| `Fields` | `[]Field` | Fields are the nested fields, for a Nested field; nil otherwise. |

## `Schema`

Schema describes one collection of records and the fields each of them can hold.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the collection's name, as used in Database.Collection. |
| `Fields` | `[]Field` | Fields are the fields each record of the collection can hold. |

## `Props`

Props is the declarative description a database is created from. It is the only thing Databases.New takes, so a database is fully described by a value a caller can write, read and version.

| Field | Type | Description |
| --- | --- | --- |
| `Path` | `string` | Path is the prefix every key of the database is written under. It is split on slashes into the leading segments of every key, so a backend that maps keys to files reads it as a directory; empty segments are dropped, which makes a trailing slash optional. A "." or ".." segment is refused: where the tree lands is the backend's base, not the Path. |
| `Schemas` | `[]Schema` | Schemas are the collections the database holds. |

## `Error`

Error describes one failure reported by a database operation. It carries no behaviour, so a caller switches on Type and reads Field, Value and Message directly. A nil *Error means success.

| Field | Type | Description |
| --- | --- | --- |
| `Type` | `int` | Type is one of KeyConflict, NoValue, MissingField, InvalidField, Internal, Removed, InvalidSchema or InvalidArgument. |
| `Field` | `string` | Field is the name of the field the failure involves, empty when the failure involves no particular field. When an insert has several bad fields, it is the first one in schema order. |
| `Value` | `any` | Value is the value the failure involves, when there is one. |
| `Message` | `string` | Message is the human-readable description of the failure. |

## `Record`

Record is one record of a collection, handed back by Collection.Insert, FindByKey, FindByID, ListAll and List, and by Record.ListNested and InsertNested for a nested collection. Its Fields and Prefix are copies: changing them changes nothing the record does.

| Field | Type | Description |
| --- | --- | --- |
| `Fields` | `[]Field` | Fields are the fields the record's own collection declares. |
| `Prefix` | `[]string` | Prefix is the key prefix of the collection the record belongs to, held as the list of segments every key under it is built from. |
| `ID` | `int64` | ID is the record's permanent identifier. It is never reused, so an id stored in a Link or Int field of another collection stays a reference to this record or to nothing at all — never to a different record. |
| `Get` | `func(fieldName string) (any, *Error)` | Get returns the typed value stored for a field: a string for a Key or String field, an int64 for an Int or Link field, a float64 for a Float field, a []byte for a Bytes field. It fails with NoValue when the field has no stored value — which is also what a removed record reports — and with InvalidField when the schema declares no such field or the field is a nested collection. |
| `GetLink` | `func(fieldName string) (Record, bool, *Error)` | GetLink resolves a Link field to the record it points at, looked up by id in the collection the field's Target names. ok is false when the field holds no stored value or when the record the stored id names is no longer live. It fails with InvalidField when the schema declares no such Link field, and with Internal when the backend fails. |
| `Update` | `func(fieldName string, value any) *Error` | Update writes a new value for a field, re-indexing it when the field is a Key. A nil value clears an optional field, and fails with MissingField on a Required one. It fails with KeyConflict when another live record already holds the new value for that Key, and with Removed when this record is no longer live. |
| `Remove` | `func() *Error` | Remove deletes the record, its index entries and every record of every collection nested under it. Removing a record that is already gone is not an error, and calling it again after it failed with Internal finishes what the failed call started. |
| `HasValues` | `func(fields []string) (bool, *Error)` | HasValues reports whether every named field has a stored value for this record; an empty list is true. It fails with InvalidField when a name is not a plain field of the schema. |
| `Nested` | `func(fieldName string) (Collection, *Error)` | Nested returns a Nested field as a Collection of its own, rooted at this record: FindByKey, FindByID, List and Repair work on it exactly as on a top-level collection. It fails with InvalidField when the schema declares no such nested field. |
| `ListNested` | `func(fieldName string) ([]Record, *Error)` | ListNested returns every record of a Nested field — the same as Nested(fieldName) followed by ListAll. |
| `InsertNested` | `func(fieldName string, fields map[string]any) (Record, *Error)` | InsertNested inserts a record into a Nested field, taking the same fields map Collection.Insert takes — the same as Nested(fieldName) followed by Insert. |
| `String` | `func() string` | String renders the record's id and its plain fields, for printing. Text values are quoted, so a value holding a comma cannot read as two fields. |

## `Collection`

Collection is one collection of records, handed back by Database.Collection, or by Record.Nested for a nested one. Its Fields and Prefix are copies: changing them changes nothing the collection does.

| Field | Type | Description |
| --- | --- | --- |
| `Fields` | `[]Field` | Fields are the fields each record of the collection can hold. |
| `Prefix` | `[]string` | Prefix is the collection's key prefix, held as the list of segments every key under it is built from. |
| `Insert` | `func(fields map[string]any) (Record, *Error)` | Insert writes a record, validating the fields against the schema and against the unique index of every Key field. A nil value is the same as leaving the field out. Inserting into a nested collection whose owning record is no longer live fails with Removed. |
| `FindByKey` | `func(field string, value any) (Record, bool, *Error)` | FindByKey looks a record up through a unique Key field. ok is false when no live record holds that value. It fails with InvalidField when the schema declares no such Key field or the value is the wrong Go type, and with Internal when the backend fails — never reporting a failing backend as an absent record. |
| `FindByID` | `func(id int64) (Record, bool, *Error)` | FindByID looks a record up through its permanent id — the value Record.ID reports. ok is false when the collection holds no live record under that id. It fails with Internal when the backend fails. |
| `ListAll` | `func() ([]Record, *Error)` | ListAll returns every record of the collection, in the order the dense position list holds them. It is not a snapshot: a record removed while it runs may be left out, and another one moved by that removal may be too. |
| `List` | `func(position int, chunk int) ([]Record, *Error)` | List returns the records of up to chunk positions starting at position, counted from 1. A chunk of 0 means "to the end of the collection". A removal between two pages moves the last record into the freed position, so a caller paging through a collection that is being written may miss that record. It fails with InvalidArgument for a position below 1 or a negative chunk. |
| `Repair` | `func() *Error` | Repair restores every invariant of the collection after a crash or a failed write: it closes holes in the position list, deletes what records that are not live left behind, and writes every missing index entry of every live record — which is also how a field turned from String into Key gets indexed. It recurses into every nested collection. It fails with KeyConflict when two live records hold the same value for a Key field; everything else it repaired stays repaired. |

## `Database`

Database is one database, bound to the Props it was described with.

| Field | Type | Description |
| --- | --- | --- |
| `Props` | `Props` | Props is a copy of the description the database was created from. |
| `Collection` | `func(name string) (Collection, bool)` | Collection returns the collection with the given name. ok is false when the Props declares no schema under that name. |

## `Databases`

Databases is the contract every database is reached through, carried by the Sandbox as the field of the same name.

| Field | Type | Description |
| --- | --- | --- |
| `New` | `func(props Props) (Database, *Error)` | New builds a database from a Props description, after checking it: it fails with InvalidSchema when the Props is not one a database can be built from. It touches no key: a database is a value over a prefix, so building one is free and creates nothing until the first record is written. The Props is copied, so changing it afterwards changes nothing the database does. |

[every contract](doc.md)
