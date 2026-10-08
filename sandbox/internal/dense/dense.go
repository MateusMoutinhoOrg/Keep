package dense

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
)

// Key layout and value encoding of the Dense Record Pattern, documented in
// docs/DenseRecordPattern. Everything here is expressed as single-key reads
// and writes against sandbox.Deps.StorageDeps, and assumes a single writer.
// The record operations built on top of these helpers live in the
// record package.
//
// This package names no object type of the api beyond the plain-data ones,
// with the single exception of api.Record, which EncodeValue accepts as
// the value of a Link field and reads nothing but the ID off. It names the
// type, never the record package, which is what keeps the package graph
// acyclic: record may import dense, dense may never import record.

// stringer is what a caller-supplied value may implement to be stored in a
// Key field without being a string. The sandbox may not import `fmt`, so
// the one method of fmt.Stringer is restated here.
type stringer interface {
	String() string
}

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
	for _, segment := range sandbox.Deps.StringsDeps.Split(path, "/") {
		if segment == "" {
			continue
		}
		prefix = append(prefix, segment)
	}
	return append(prefix, name)
}

// SizeKey holds the number of live records of a collection — the highest
// occupied position of its dense list.
func SizeKey(sandbox *api.Sandbox, prefix []string) []string {
	return Key(sandbox, prefix, "size")
}

// LastIDKey holds the highest id ever allocated in a collection. It only
// grows, which is what makes an id never reused.
func LastIDKey(sandbox *api.Sandbox, prefix []string) []string {
	return Key(sandbox, prefix, "last-id")
}

// ListKey holds the id living at one position of a collection's dense list.
// Positions run from 1 to the value of SizeKey with no gap, which is what
// makes iteration possible without listing keys.
func ListKey(sandbox *api.Sandbox, prefix []string, position int64) []string {
	return Key(sandbox, prefix, "list", formatID(sandbox, position))
}

// PositionKey holds the position a record currently occupies in the dense
// list. It is the back-pointer that makes a removal cost the same whatever
// the size of the collection, and its presence is what marks a record live.
func PositionKey(sandbox *api.Sandbox, prefix []string, id int64) []string {
	return Key(sandbox, prefix, formatID(sandbox, id), "position")
}

// ValueKey holds one field value of one record.
func ValueKey(sandbox *api.Sandbox, prefix []string, id int64, field string) []string {
	return Key(sandbox, prefix, formatID(sandbox, id), "values", field)
}

// IndexKey holds the id owning one value of one Key field — the unique
// index, addressed by the hash of the value so a lookup is a single read.
func IndexKey(sandbox *api.Sandbox, prefix []string, field string, hash string) []string {
	return Key(sandbox, prefix, "keys", field, hash)
}

// SubPrefix is the collection prefix of a Nested field of one record. A
// nested collection is a collection like any other, which is why every
// helper here works on it unchanged.
func SubPrefix(sandbox *api.Sandbox, prefix []string, id int64, field string) []string {
	return Key(sandbox, prefix, formatID(sandbox, id), field)
}

// formatID renders an id or a position as the decimal segment every key
// above carries it as.
func formatID(sandbox *api.Sandbox, value int64) string {
	return sandbox.Deps.StringsDeps.FormatInt(value, 10)
}

// HashIndexValue normalizes and hashes an encoded value, so index lookups
// are case-insensitive and a key never grows with the value it indexes.
func HashIndexValue(sandbox *api.Sandbox, encoded string) string {
	lowered := sandbox.Deps.StringsDeps.ToLower(encoded)
	return sandbox.Deps.HashDeps.Sha256Hex([]byte(lowered))
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

// InternalError wraps a storage failure as a typed *api.Error, which is the
// only way a backend error ever reaches a caller.
func InternalError(sandbox *api.Sandbox, err error) *api.Error {
	return liberror.New(sandbox, api.Internal, "", err.Error())
}

// ParseID reads an id back from the decimal form WriteInt stores it in.
func ParseID(sandbox *api.Sandbox, raw []byte) (int64, error) {
	id, err := sandbox.Deps.StringsDeps.ParseInt(string(raw), 10, 64)
	if err != nil {
		return 0, sandbox.Deps.StdDeps.Errorf("keep: invalid id: %s", string(raw))
	}
	return id, nil
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
			return typed.String(), nil
		default:
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects a string value, got %T", field.Name, value))
		}
	case api.Int:
		switch typed := value.(type) {
		case int:
			return sandbox.Deps.StringsDeps.FormatInt(int64(typed), 10), nil
		case int32:
			return sandbox.Deps.StringsDeps.FormatInt(int64(typed), 10), nil
		case int64:
			return sandbox.Deps.StringsDeps.FormatInt(typed, 10), nil
		default:
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects an integer value, got %T", field.Name, value))
		}
	case api.Float:
		switch typed := value.(type) {
		case float64:
			return formatFloat(sandbox, typed), nil
		case float32:
			return formatFloat(sandbox, float64(typed)), nil
		case int:
			return formatFloat(sandbox, float64(typed)), nil
		case int32:
			return formatFloat(sandbox, float64(typed)), nil
		case int64:
			return formatFloat(sandbox, float64(typed)), nil
		default:
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects a floating-point value, got %T", field.Name, value))
		}
	case api.Link:
		// A Link naming no collection is a mistake in the schema, and the
		// one thing about a link this function can catch: whether the id
		// still names a live record is a read, and is left to GetLink.
		if field.Target == "" {
			return "", liberror.New(sandbox, api.InvalidField, field.Name,
				sandbox.Deps.StdDeps.Sprintf("link field %q declares no Target schema", field.Name))
		}
		switch typed := value.(type) {
		case api.Record:
			return sandbox.Deps.StringsDeps.FormatInt(typed.ID, 10), nil
		case int:
			return sandbox.Deps.StringsDeps.FormatInt(int64(typed), 10), nil
		case int32:
			return sandbox.Deps.StringsDeps.FormatInt(int64(typed), 10), nil
		case int64:
			return sandbox.Deps.StringsDeps.FormatInt(typed, 10), nil
		default:
			return "", liberror.NewWithValue(sandbox, api.InvalidField, field.Name, value,
				sandbox.Deps.StdDeps.Sprintf("field %q expects a record or a record id, got %T", field.Name, value))
		}
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
