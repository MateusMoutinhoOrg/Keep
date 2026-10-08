package dense

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
	writelock "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/writelock"
)

// Key layout and value encoding of the Dense Record Pattern, documented in
// docs/DenseRecordPattern. Everything here is expressed as single-key reads
// and writes against sandbox.Deps.StorageDeps. The record operations built on
// top of these helpers live in the record package; what makes a record, a
// list slot or an index entry live is in live.go.
//
// This package names no object type of the api beyond the plain-data ones,
// with the single exception of api.Record, which EncodeValue accepts as
// the value of a Link field and reads nothing but the ID off — and api.Collection,
// which a Scope builds without naming the package that implements it. It
// names the types, never the record or collection packages, which is what
// keeps the package graph acyclic: record may import dense, dense may never
// import record.

// stringer is what a caller-supplied value may implement to be stored in a
// Key field without being a string. The sandbox may not import `fmt`, so
// the one method of fmt.Stringer is restated here.
type stringer interface {
	String() string
}

// maxInt64 and minInt64 bound the integers an Int or Link field stores.
const (
	maxInt64 = 1<<63 - 1
	minInt64 = -1 << 63
)

// reservedNestedNames are the names a Nested field may not take: a nested
// collection's prefix is {c}/{id}/{field}, which shares its parent with the
// record's own {c}/{id}/position and {c}/{id}/values keys.
var reservedNestedNames = []string{"position", "values"}

// Key builds a storage key out of a collection prefix and the segments that
// follow it. A key is a list, never a joined string: the separator lives in
// the adapter, so no character of a field name or of a prefix can be read
// back as a boundary and collide with another key. Every call allocates a
// slice of its own, so a key never aliases the prefix it was derived from.
func Key(sandbox *api.Sandbox, prefix []string, segments ...string) []string {
	key := make([]string, 0, len(prefix)+len(segments))
	key = append(key, prefix...)
	return append(key, segments...)
}

// RootPrefix is the prefix of a top-level collection: the slash-separated
// segments of Props.Path, empty ones dropped, followed by the schema name.
// Splitting the path here is what keeps a caller free to write it as the
// directory-looking string it has always been.
func RootPrefix(sandbox *api.Sandbox, path string, name string) []string {
	prefix := make([]string, 0, 4)
	for _, segment := range PathSegments(sandbox, path) {
		prefix = append(prefix, segment)
	}
	return append(prefix, name)
}

// PathSegments splits Props.Path into the segments every key of the database
// starts with, dropping the empty ones.
func PathSegments(sandbox *api.Sandbox, path string) []string {
	segments := make([]string, 0, 4)
	for _, segment := range sandbox.Deps.StringsDeps.Split(path, "/") {
		if segment == "" {
			continue
		}
		segments = append(segments, segment)
	}
	return segments
}

// SizeKey holds the number of positions of a collection's dense list — the
// highest occupied position.
func SizeKey(sandbox *api.Sandbox, prefix []string) []string {
	return Key(sandbox, prefix, "size")
}

// LastIDKey holds the highest id ever allocated in a collection. It only
// grows, which is what makes an id never reused.
func LastIDKey(sandbox *api.Sandbox, prefix []string) []string {
	return Key(sandbox, prefix, "last-id")
}

// ListKey holds the id living at one position of a collection's dense list.
// Positions run from 1 to the value of SizeKey, which is what makes
// iteration possible without listing keys.
func ListKey(sandbox *api.Sandbox, prefix []string, position int64) []string {
	return Key(sandbox, prefix, "list", formatID(sandbox, position))
}

// PositionKey holds the position a record currently occupies in the dense
// list. It is the back-pointer that makes a removal cost the same whatever
// the size of the collection; whether the record is live is decided by Live.
func PositionKey(sandbox *api.Sandbox, prefix []string, id int64) []string {
	return Key(sandbox, prefix, formatID(sandbox, id), "position")
}

// ValueKey holds one field value of one record.
func ValueKey(sandbox *api.Sandbox, prefix []string, id int64, field string) []string {
	return Key(sandbox, prefix, formatID(sandbox, id), "values", field)
}

// IndexKey holds the id owning one value of one Key field — the unique
// index, addressed by the hash of the value so a lookup costs the same at
// any size.
func IndexKey(sandbox *api.Sandbox, prefix []string, field string, hash string) []string {
	return Key(sandbox, prefix, "keys", field, hash)
}

// SubPrefix is the collection prefix of a Nested field of one record. A
// nested collection is a collection like any other, which is why every
// helper here works on it unchanged.
func SubPrefix(sandbox *api.Sandbox, prefix []string, id int64, field string) []string {
	return Key(sandbox, prefix, formatID(sandbox, id), field)
}

// ReservedNestedName reports whether name is one a Nested field may not take.
func ReservedNestedName(sandbox *api.Sandbox, name string) bool {
	for _, reserved := range reservedNestedNames {
		if name == reserved {
			return true
		}
	}
	return false
}

// formatID renders an id or a position as the decimal segment every key
// above carries it as.
func formatID(sandbox *api.Sandbox, value int64) string {
	return sandbox.Deps.StringsDeps.FormatInt(value, 10)
}

// HashIndexValue folds and hashes an encoded value, so index lookups ignore
// case in every script and a key never grows with the value it indexes.
// Folding is Unicode's canonical caseless match, which lower-cases ASCII, so
// an ASCII value hashes exactly as it did when the index was lower-cased.
func HashIndexValue(sandbox *api.Sandbox, encoded string) string {
	folded := sandbox.Deps.FoldDeps.Fold(encoded)
	return sandbox.Deps.HashDeps.Sha256Hex([]byte(folded))
}

// FindField returns the schema field with the given name. ok is false when
// the schema declares no such field.
func FindField(sandbox *api.Sandbox, fields []api.Field, name string) (field api.Field, ok bool) {
	for _, candidate := range fields {
		if candidate.Name == name {
			return candidate, true
		}
	}
	return api.Field{}, false
}

// CopyFields returns a deep copy of a field list, nested fields included, so
// a value handed to a caller shares no array with the one the library reads.
func CopyFields(sandbox *api.Sandbox, fields []api.Field) []api.Field {
	if fields == nil {
		return nil
	}
	copied := make([]api.Field, len(fields))
	for index, field := range fields {
		copied[index] = field
		copied[index].Fields = CopyFields(sandbox, field.Fields)
	}
	return copied
}

// SamePrefix reports whether two prefixes name the same collection.
func SamePrefix(sandbox *api.Sandbox, left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

// LinkResolver returns the fields and the key prefix of the collection a
// Link field targets, by the schema name Field.Target carries. ok is false
// when the database declares no schema under that name. It is built once
// from the Props a database was created with and carried down every nesting
// level, so a record of a nested collection follows a link exactly the way
// a top-level one does.
type LinkResolver func(target string) (fields []api.Field, prefix []string, ok bool)

// NewLinkResolver builds the resolver a database hands to every collection
// it creates, closing over its Props. It is the link half of what
// Database.Collection does: a record reaches the collection it points at
// through the closure it was built with, never through a field a caller
// could read or replace.
func NewLinkResolver(sandbox *api.Sandbox, props api.Props) LinkResolver {
	return func(target string) ([]api.Field, []string, bool) {
		for _, schema := range props.Schemas {
			if schema.Name == target {
				return schema.Fields, RootPrefix(sandbox, props.Path, schema.Name), true
			}
		}
		return nil, nil, false
	}
}

// Scope is everything a collection and its records carry beside their own
// fields and prefix: the database they belong to, the write lock they take,
// and the record owning them when they are nested. It travels from the
// database that built a collection down through every record and every
// nested collection, as a parameter and never as a field of an api type, for
// the same reason storage does: an api type carries no wiring a caller could
// read or replace.
type Scope struct {
	// Resolve finds the collection a Link field targets.
	Resolve LinkResolver
	// Locks is the registry of the sandbox's in-process write locks.
	Locks *writelock.Registry
	// Root is the prefix of the top-level collection this one is, or is
	// nested under. It is the key of the write lock every write takes, so a
	// write to a nested collection and the removal of the record owning it
	// never interleave.
	Root []string
	// Owner is the prefix of the collection holding the record that owns
	// this one, and OwnerID that record's id; Owner is nil at the top level.
	// An insert checks the owner is still live, so nothing is written under
	// a record that has been removed.
	Owner   []string
	OwnerID int64
	// NewCollection builds an api.Collection over a prefix. The database
	// fills it with collection.New, which is how a record hands out one of
	// its nested collections without the record package importing the
	// collection package that imports it.
	NewCollection func(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope Scope) api.Collection
}

// NestedScope is the scope of the collection a Nested field of record id
// holds: the same database and write lock, owned by that record.
func NestedScope(sandbox *api.Sandbox, scope Scope, prefix []string, id int64) Scope {
	nested := scope
	nested.Owner = prefix
	nested.OwnerID = id
	return nested
}

// TargetScope is the scope of the top-level collection at prefix, in the
// same database as scope: the one a Link field resolves into.
func TargetScope(sandbox *api.Sandbox, scope Scope, prefix []string) Scope {
	target := scope
	target.Root = prefix
	target.Owner = nil
	target.OwnerID = 0
	return target
}

// InternalError wraps a storage failure as a typed *api.Error, which is the
// only way a backend error ever reaches a caller.
func InternalError(sandbox *api.Sandbox, err error) *api.Error {
	return liberror.New(sandbox, api.Internal, "", err.Error())
}

// ParseID reads an id back from the decimal form WriteInt stores it in.
func ParseID(sandbox *api.Sandbox, raw []byte) (int64, error) {
	id, err := sandbox.Deps.StringsDeps.ParseInt(string(raw), 10, 64)
	if err != nil {
		return 0, sandbox.Deps.StdDeps.Errorf("keep: invalid id: %q", string(raw))
	}
	return id, nil
}

// EncodeFieldValue is EncodeValue plus the one check that needs the
// database: a record given as the value of a Link field has to belong to the
// collection the field's Target names, or its id would resolve to another
// record of that collection.
func EncodeFieldValue(sandbox *api.Sandbox, scope Scope, field api.Field, value any) (string, *api.Error) {
	if field.Type == api.Link {
		if linked, ok := value.(api.Record); ok {
			_, targetPrefix, found := scope.Resolve(field.Target)
			if !found || !SamePrefix(sandbox, linked.Prefix, targetPrefix) {
				return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, linked.ID,
					sandbox.Deps.StdDeps.Sprintf("field %q links to %q, and the record given belongs to %q",
						field.Name, field.Target, sandbox.Deps.StringsDeps.Join(linked.Prefix, "/")))
			}
		}
	}
	return EncodeValue(sandbox, field, value)
}

// EncodeValue converts a caller-provided value to the canonical string form
// it is stored in, validating it against the field's type on the way.
func EncodeValue(sandbox *api.Sandbox, field api.Field, value any) (string, *api.Error) {
	switch field.Type {
	case api.Key, api.String:
		switch typed := value.(type) {
		case string:
			return typed, nil
		case stringer:
			text, ok := stringOf(sandbox, typed)
			if !ok {
				return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
					sandbox.Deps.StdDeps.Sprintf("field %q: the String method of %T panicked", field.Name, value))
			}
			return text, nil
		default:
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects a string value, got %T", field.Name, value))
		}
	case api.Int:
		number, ok := integerOf(sandbox, value)
		if !ok {
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects an integer value, got %T", field.Name, value))
		}
		return sandbox.Deps.StringsDeps.FormatInt(number, 10), nil
	case api.Float:
		number, ok := floatOf(sandbox, value)
		if !ok {
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects a floating-point value, got %T", field.Name, value))
		}
		return formatFloat(sandbox, number), nil
	case api.Link:
		// A Link naming no collection is a mistake in the schema, which
		// Databases.New refuses; it is checked again here because nothing
		// about the link can be stored without it. Whether the id still
		// names a live record is a read, and is left to GetLink.
		if field.Target == "" {
			return "", liberror.New(sandbox, api.InvalidField, field.Name,
				sandbox.Deps.StdDeps.Sprintf("link field %q declares no Target schema", field.Name))
		}
		if linked, ok := value.(api.Record); ok {
			return sandbox.Deps.StringsDeps.FormatInt(linked.ID, 10), nil
		}
		number, ok := integerOf(sandbox, value)
		if !ok {
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects a record or a record id, got %T", field.Name, value))
		}
		return sandbox.Deps.StringsDeps.FormatInt(number, 10), nil
	case api.Bytes:
		// A Go string holds any bytes at all, so the conversion is lossless
		// and the value reaches storage exactly as the caller gave it.
		typed, ok := value.([]byte)
		if !ok {
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects a byte slice value, got %T", field.Name, value))
		}
		return string(typed), nil
	default:
		return "", liberror.New(sandbox, api.InvalidField, field.Name,
			sandbox.Deps.StdDeps.Sprintf("field %q cannot be encoded as a plain value", field.Name))
	}
}

// stringOf calls a caller's String method, reporting ok == false rather than
// letting a panic inside it — a nil pointer receiver, typically — unwind
// through the library into the caller.
func stringOf(sandbox *api.Sandbox, value stringer) (text string, ok bool) {
	defer func() {
		if recover() != nil {
			text, ok = "", false
		}
	}()
	return value.String(), true
}

// integerOf converts any Go integer, and any float holding a whole number,
// to the int64 an Int or Link field stores. ok is false for any other type
// and for a value outside the range of an int64 — a float64 is what
// encoding/json hands back for every number, so a whole one is accepted.
func integerOf(sandbox *api.Sandbox, value any) (int64, bool) {
	switch typed := value.(type) {
	case int:
		return int64(typed), true
	case int8:
		return int64(typed), true
	case int16:
		return int64(typed), true
	case int32:
		return int64(typed), true
	case int64:
		return typed, true
	case uint:
		return unsignedOf(sandbox, uint64(typed))
	case uint8:
		return int64(typed), true
	case uint16:
		return int64(typed), true
	case uint32:
		return int64(typed), true
	case uint64:
		return unsignedOf(sandbox, typed)
	case float32:
		return wholeOf(sandbox, float64(typed))
	case float64:
		return wholeOf(sandbox, typed)
	default:
		return 0, false
	}
}

// unsignedOf converts an unsigned integer, refusing one an int64 cannot hold.
func unsignedOf(sandbox *api.Sandbox, value uint64) (int64, bool) {
	if value > maxInt64 {
		return 0, false
	}
	return int64(value), true
}

// wholeOf converts a float holding a whole number in the range of an int64.
// The range is checked before converting, since converting an out-of-range
// float to an integer is implementation-defined in Go; NaN fails both
// comparisons.
func wholeOf(sandbox *api.Sandbox, value float64) (int64, bool) {
	if !(value >= minInt64 && value < maxInt64) {
		return 0, false
	}
	whole := int64(value)
	if float64(whole) != value {
		return 0, false
	}
	return whole, true
}

// floatOf converts any Go float or integer to the float64 a Float field
// stores.
func floatOf(sandbox *api.Sandbox, value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	default:
		return 0, false
	}
}

// formatFloat renders a float in the shortest decimal form that parses back
// to the same number, which is what keeps a stored value stable byte for
// byte across writes of the same value.
func formatFloat(sandbox *api.Sandbox, value float64) string {
	return sandbox.Deps.StringsDeps.FormatFloat(value, 'g', -1, 64)
}

// DecodeValue converts a stored value back to the typed form a caller of
// Record.Get receives: an int64 for an Int or Link field, a float64 for
// a Float field, a []byte for a Bytes field, a string otherwise. A Bytes
// value is handed back as a copy of its own, so a caller writing into it
// never reaches the bytes the backend holds. A stored value carries no type
// tag, so the schema is the only thing that says how to read its bytes back
// — a field type with no case here falls through and is handed back as a
// string.
func DecodeValue(sandbox *api.Sandbox, field api.Field, raw []byte) (any, *api.Error) {
	switch field.Type {
	case api.Int, api.Link:
		number, err := sandbox.Deps.StringsDeps.ParseInt(string(raw), 10, 64)
		if err != nil {
			return nil, InternalError(sandbox, err)
		}
		return number, nil
	case api.Float:
		number, err := sandbox.Deps.StringsDeps.ParseFloat(string(raw), 64)
		if err != nil {
			return nil, InternalError(sandbox, err)
		}
		return number, nil
	case api.Bytes:
		copied := make([]byte, len(raw))
		copy(copied, raw)
		return copied, nil
	default:
		return string(raw), nil
	}
}

// ReadInt reads an integer key, treating a key that holds nothing as
// zero: a collection nothing was ever written to has no size key, and its
// size is zero.
func ReadInt(sandbox *api.Sandbox, key []string) (int64, error) {
	raw, found, err := sandbox.Deps.StorageDeps.Read(key)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, nil
	}
	return sandbox.Deps.StringsDeps.ParseInt(string(raw), 10, 64)
}

// WriteInt stores an integer under key in the canonical decimal form every
// reader here expects.
func WriteInt(sandbox *api.Sandbox, key []string, value int64) error {
	encoded := sandbox.Deps.StringsDeps.FormatInt(value, 10)
	return sandbox.Deps.StorageDeps.Write(key, []byte(encoded))
}
