# PublicApi

Every exported symbol of `github.com/MateusMoutinhoOrg/Keep`, read straight from the contract sources on
every build: `sandbox/api/` is the surface `sandbox.New` returns, `sandbox/deps/`
the contracts an adapter fills and a caller may replace. Each description below is the
doc comment of the declaration itself — change the comment, run `build`, and this page
follows.

## Entry points

| Symbol | Signature |
| --- | --- |
| `sandbox.New` | `func(deps *deps.Deps) *api.Sandbox` |
| `standard.New` | `func() deps.Deps` (`adapters/availables/standard`) |

Implementations live under `sandbox/internal` and are unreachable: every contract is a
struct of function fields, filled by a binder.

# The sandbox api

## `sandbox/api/sandbox.go`

### `Sandbox`

Sandbox is the whole library: one field per contract declared in sandbox/api/, each built by the New<Contract> of its own package under sandbox/internal/. sandbox.New returns it, and nothing callable lives outside of it.

| Field | Type | Description |
| --- | --- | --- |
| `Deps` | `*deps.Deps` | Deps is every capability the sandbox reaches the outside world through. It rides on the api so that a function handed the Sandbox holds the whole of what it needs, and can call another field of the api besides — which is what makes a field a caller replaced take effect everywhere. It is also the one field that does not cross into a consumer: an installed copy of this contract carries the api, never the wiring behind it. |
| `Databases` | `Databases` |  |
| `Info` | `Info` |  |

## `sandbox/api/databases.go`

| Constant | Value | Description |
| --- | --- | --- |
| `Key` | `iota` | Key is a unique, indexed string field. Two live records of one collection can never hold the same value for it, and it is the only kind of field SchemaInstance.FindByKey looks a record up by. |
| `Int` |  | Int is a plain integer field, stored in its decimal form and handed back as an int64. |
| `Database` |  | Database is a nested collection of records, reached through SchemaItem.NewSubItem and SchemaItem.ListAll rather than read as a value. |
| `KeyConflict` | `iota` | KeyConflict is a value another live record already holds for a Key field. |
| `NotFound` |  | NotFound is a field of an existing record that has no stored value. |
| `MissingField` |  | MissingField is a Required field absent from an insert. |
| `InvalidField` |  | InvalidField is a field the schema does not declare, a value of the wrong Go type for the field it is written to, or a Database field used where a plain value was expected. |
| `Internal` |  | Internal is a failure the storage backend reported. Error.Message carries what the backend said. |

### `Item`

Item describes one field of a schema.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the field's name, as used in the fields map of an insert and in SchemaItem.Get. |
| `Type` | `int` | Type is one of Key, Int or Database. |
| `Required` | `bool` | Required reports whether an insert must provide this field. It is ignored on a Database field, which is never provided directly. |
| `Itens` | `[]Item` | Itens are the nested fields, for a Database field; nil otherwise. |

### `Schema`

Schema describes one collection of records and the fields each of them can hold.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `string` | Name is the collection's name, as used in DatabaseHandle.GetSchema. |
| `Itens` | `[]Item` | Itens are the fields each record of the collection can hold. |

### `Props`

Props is the declarative description a database is created from. It is the only thing Database.New takes, so a database is fully described by a value a caller can write, read and version.

| Field | Type | Description |
| --- | --- | --- |
| `Path` | `string` | Path is the prefix every key of the database is written under. A backend that maps keys to files reads it as a directory, so it usually ends with a slash. |
| `Schemas` | `[]Schema` | Schemas are the collections the database holds. |

### `Error`

Error describes one failure reported by a database operation. It carries no behaviour, so a caller switches on Type and reads Key, KeyValue and Message directly. A nil *Error means success.

| Field | Type | Description |
| --- | --- | --- |
| `Type` | `int` | Type is one of KeyConflict, NotFound, MissingField, InvalidField or Internal. |
| `Key` | `string` | Key is the name of the field the failure involves, empty when the failure involves no particular field. |
| `KeyValue` | `any` | KeyValue is the value the failure involves, when there is one. |
| `Message` | `string` | Message is the human-readable description of the failure. |

### `SchemaItem`

SchemaItem is one record of a collection, handed back by SchemaInstance.NewItem, FindByKey, FindById, ListAll and List, and by SchemaItem.ListAll and NewSubItem for a nested collection.

| Field | Type | Description |
| --- | --- | --- |
| `Items` | `[]Item` | Items are the fields the record's own collection declares. |
| `Prefix` | `string` | Prefix is the key prefix of the collection the record belongs to. |
| `Id` | `int64` | Id is the record's permanent identifier. It is never reused, so an id stored in an Int field of another collection stays a reference to this record or to nothing at all — never to a different record. |
| `Get` | `func(fieldName string) (any, *Error)` | Get returns the typed value stored for a field: a string for a Key field, an int64 for an Int field. It fails with NotFound when the field has no stored value and with InvalidField when the schema declares no such field or the field is a nested collection. |
| `Update` | `func(fieldName string, value any) *Error` | Update writes a new value for a field, re-indexing it when the field is a Key. It fails with KeyConflict when another live record already holds the new value for that Key. |
| `Remove` | `func() *Error` | Remove deletes the record, its index entries and every record of every collection nested under it. Removing a record that is already gone is not an error. |
| `CheckKeysPresence` | `func(keys []string) bool` | CheckKeysPresence reports whether every named field has a stored value for this record. |
| `ListAll` | `func(fieldName string) []SchemaItem` | ListAll returns every record of a nested (Database) field, and nil when the schema declares no such nested field. |
| `NewSubItem` | `func(fieldName string, fields map[string]any) (SchemaItem, *Error)` | NewSubItem inserts a record into a nested (Database) field, taking the same fields map SchemaInstance.NewItem takes. |
| `String` | `func() string` | String renders the record's id and its plain fields, for printing. |

### `SchemaInstance`

SchemaInstance is one collection of records, handed back by DatabaseHandle.GetSchema.

| Field | Type | Description |
| --- | --- | --- |
| `Items` | `[]Item` | Items are the fields each record of the collection can hold. |
| `Prefix` | `string` | Prefix is the collection's key prefix. |
| `NewItem` | `func(fields map[string]any) (SchemaItem, *Error)` | NewItem inserts a record, validating the fields against the schema and against the unique index of every Key field. |
| `FindByKey` | `func(key string, keyValue any) (SchemaItem, bool)` | FindByKey looks a record up through a unique Key field. ok is false when the schema declares no such Key field, or when no live record holds that value. |
| `FindById` | `func(id int64) (SchemaItem, bool)` | FindById looks a record up through its permanent id — the value SchemaItem.Id reports. ok is false when the collection holds no live record under that id. |
| `ListAll` | `func() ([]SchemaItem, *Error)` | ListAll returns every record of the collection, in the order the dense position list holds them. |
| `List` | `func(position int, chunk int) ([]SchemaItem, *Error)` | List returns up to chunk records starting at position, counted from 1. A chunk of 0 means "to the end of the collection". |

### `DatabaseHandle`

DatabaseHandle is one database, bound to the Props it was described with.

| Field | Type | Description |
| --- | --- | --- |
| `Props` | `Props` | Props is the description the database was created from. |
| `GetSchema` | `func(name string) (SchemaInstance, bool)` | GetSchema returns the collection with the given name. ok is false when the Props declares no schema under that name. |

### `Databases`

Databases is the contract every database is reached through, carried by the Sandbox as the field of the same name.

| Field | Type | Description |
| --- | --- | --- |
| `New` | `func(props Props) DatabaseHandle` | New builds a database from a Props description. It touches no key: a database is a handle over a prefix, so building one is free and creates nothing until the first record is written. |

## `sandbox/api/info.go`

### `Info`

Info is the contract reporting the library's own identity, carried by the Sandbox as the field of the same name. Both values are compile-time constants of sandbox/internal/config, generated from AgnosConfig/project.yaml, so a release bump is a one-line edit touching no logic — and a caller can report which Keep it linked against without importing anything but this package.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `func() string` | Name is the library's name, "Keep". |
| `Version` | `func() string` | Version is the release the caller linked against, in the "v0.0.0" spelling the repository tags with. |

# Dependency contracts

`deps.Deps` has one field per directory of `sandbox/deps/`, named by title-casing it. Each
field is that package's `Sandbox` struct, filled by `adapters/libs/<name>.Bind(&deps)`.

## `deps.Hashdeps`

`sandbox/deps/hashdeps`

### `Sandbox`

Sandbox is the hashing library injected whole as the Deps.Hashdeps field.

| Field | Type | Description |
| --- | --- | --- |
| `Sha256Hex` | `func(content []byte) string` | Sha256Hex returns the SHA-256 digest of content, lower-case hexadecimal. It is what every recorded example tree is compared by, so the encoding is part of the golden and may not change. |

## `deps.Std`

`sandbox/deps/std`

### `Sandbox`

Sandbox is the runtime library injected whole as the Deps.Std field.

| Field | Type | Description |
| --- | --- | --- |
| `Now` | `func() int64` | Now returns the current wall-clock time as nanoseconds since the Unix epoch, UTC. The sandbox may not name a `time.Time`, so an instant crosses this boundary as a plain integer. |
| `Printf` | `func(format string, a ...any) (n int, err error)` | Printf writes one formatted message to standard output. It carries the command's result — the data a script would read — so it is never silenced. |
| `Log` | `func(format string, a ...any) (n int, err error)` | Log writes one formatted progress message to standard error. It is the channel every "… started with path …" notice goes through, so a caller can keep stdout free of log noise, and it is what --quiet turns off. |
| `Error` | `func(format string, a ...any) (n int, err error)` | Error writes one formatted message to standard error. |
| `Errorf` | `func(format string, a ...any) error` | Errorf formats an error message and returns it as an error. |
| `Sprintf` | `func(format string, a ...any) string` | Sprintf formats a message and returns it as a string. It is the one formatting entry point the sandbox has: every string it builds out of values rather than out of concatenation goes through here. |
| `Goos` | `func() string` | Goos is the name of the operating system the process runs on, in the spelling the Go toolchain uses ("darwin", "linux", "windows", …). |

## `deps.Storagedeps`

`sandbox/deps/storagedeps`

### `Sandbox`

Sandbox is the storage library injected whole as the Deps.Storagedeps field. The first group of fields writes, the second reads, and the last two are the optional advisory lease a multi-writer backend can offer.

| Field | Type | Description |
| --- | --- | --- |
| `Write` | `func(key string, value []byte) error` | Write stores value under key, overwriting any current value and creating the key when it is absent. |
| `WriteIfKeyNotExists` | `func(key string, value []byte) (written bool, err error)` | WriteIfKeyNotExists stores value only when key holds nothing. written is false when the key already exists, which is not an error. |
| `WriteIfValueEquals` | `func(key string, value []byte, old_value []byte) (written bool, err error)` | WriteIfValueEquals stores value only when the current value of key is exactly old_value. written is false when the key is absent or holds something else, which is not an error. |
| `Append` | `func(key string, value []byte) error` | Append adds value to the end of the current value of key, creating the key when it is absent. |
| `InsertAt` | `func(key string, position int64, value []byte) error` | InsertAt splices value into the current value of key at position, counted in bytes from the start. A position beyond the current length is an error: it would leave a hole. |
| `Exists` | `func(key string) (exists bool, err error)` | Exists reports whether key currently holds a value. |
| `Read` | `func(key string) (value []byte, found bool, err error)` | Read returns the whole value of key. found is false when the key holds nothing, and value is then nil. |
| `ReadAt` | `func(key string, position int64, size int64) (value []byte, found bool, err error)` | ReadAt returns at most size bytes of the value of key, starting at position, counted in bytes from the start. A range reaching past the end is truncated rather than refused. found is false when the key holds nothing. |
| `Delete` | `func(key string) error` | Delete removes key. Removing a key that holds nothing is not an error: absent before and absent after is the same outcome. |
| `Lock` | `func(key string, seconds int) (locked bool, err error)` | Lock takes an advisory lease on key for seconds seconds. locked is false when someone else already holds a lease that has not expired, which is not an error. Keep never calls it itself — a database is written by one writer at a time — so a backend with no leases may report locked == true and do nothing. |
| `UnLock` | `func(key string) error` | UnLock releases a lease taken by Lock. Releasing a lease nobody holds is not an error. |

## `deps.Stringsdeps`

`sandbox/deps/stringsdeps`

### `Sandbox`

Sandbox is the text library injected whole as the Deps.Stringsdeps field. The first group of fields is string manipulation, the second is conversion between strings and numbers.

| Field | Type | Description |
| --- | --- | --- |
| `TrimSpace` | `func(s string) string` | TrimSpace returns s with leading and trailing white space removed. |
| `Trim` | `func(s string, cutset string) string` | Trim returns s with every leading and trailing character contained in cutset removed. |
| `TrimLeft` | `func(s string, cutset string) string` | TrimLeft returns s with every leading character contained in cutset removed. |
| `TrimRight` | `func(s string, cutset string) string` | TrimRight returns s with every trailing character contained in cutset removed. |
| `TrimPrefix` | `func(s string, prefix string) string` | TrimPrefix returns s without the given leading prefix. When s does not start with prefix, s is returned unchanged. |
| `TrimSuffix` | `func(s string, suffix string) string` | TrimSuffix returns s without the given trailing suffix. When s does not end with suffix, s is returned unchanged. |
| `HasPrefix` | `func(s string, prefix string) bool` | HasPrefix reports whether s begins with prefix. |
| `HasSuffix` | `func(s string, suffix string) bool` | HasSuffix reports whether s ends with suffix. |
| `Contains` | `func(s string, substr string) bool` | Contains reports whether substr is within s. |
| `ContainsAny` | `func(s string, chars string) bool` | ContainsAny reports whether any character of chars is within s. |
| `LastIndex` | `func(s string, substr string) int` | LastIndex returns the index of the last instance of substr in s, or -1 when substr is absent. |
| `Count` | `func(s string, substr string) int` | Count returns the number of non-overlapping instances of substr in s. When substr is empty it returns one plus the number of runes in s. |
| `Split` | `func(s string, sep string) []string` | Split slices s into every substring separated by sep. |
| `Join` | `func(elems []string, sep string) string` | Join concatenates elems, placing sep between consecutive elements. |
| `Fields` | `func(s string) []string` | Fields slices s around each run of white space, returning the substrings between them. |
| `FieldsFunc` | `func(s string, f func(rune) bool) []string` | FieldsFunc slices s at each run of runes satisfying f, returning the substrings between them. |
| `Repeat` | `func(s string, count int) string` | Repeat returns count copies of s concatenated. |
| `ReplaceAll` | `func(s string, old string, new string) string` | ReplaceAll returns s with every non-overlapping instance of old replaced by new. |
| `ToUpper` | `func(s string) string` | ToUpper returns s with every letter mapped to its upper case. |
| `ToLower` | `func(s string) string` | ToLower returns s with every letter mapped to its lower case. |
| `Quote` | `func(s string) string` | Quote returns s as a double-quoted Go string literal, escaping what the Go syntax requires. |
| `MatchPattern` | `func(pattern string, s string) (bool, error)` | MatchPattern reports whether s is matched by the regular expression pattern, and errors when the pattern itself does not compile. It is the one matching primitive the sandbox has: `regexp` lives on the adapter side like every other standard package. |
| `Atoi` | `func(s string) (int, error)` | Atoi parses s as a decimal integer. The error reports a string that is not one. |
| `ParseInt` | `func(s string, base int, bit_size int) (int64, error)` | ParseInt parses s as an integer in the given base with the given bit size. The error reports a string that is not one. |
| `ParseFloat` | `func(s string, bit_size int) (float64, error)` | ParseFloat parses s as a floating-point number of the given bit size. The error reports a string that is not one. |
| `FormatInt` | `func(value int64, base int) string` | FormatInt returns the string representation of value in the given base. |
| `FormatFloat` | `func(value float64, format byte, precision int, bit_size int) string` | FormatFloat returns the string representation of value, formatted according to the format byte, the precision and the bit size — the same three controls the standard library takes. |
