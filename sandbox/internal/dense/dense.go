package dense

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
)

// Key layout and value encoding of the Dense Record Pattern, documented in
// docs/DenseRecordPattern. Everything here is expressed as single-key reads
// and writes against sandbox.Deps.Storagedeps, and assumes a single writer.
// The record operations built on top of these helpers live in the
// schemaitem package.
//
// This package names no object type of the api beyond the plain-data ones,
// which is what keeps the package graph acyclic: schemaitem may import
// dense, dense may never import schemaitem.

// stringer is what a caller-supplied value may implement to be stored in a
// Key field without being a string. The sandbox may not import `fmt`, so
// the one method of fmt.Stringer is restated here.
type stringer interface {
	String() string
}

// SizeKey holds the number of live records of a collection — the highest
// occupied position of its dense list.
func SizeKey(sandbox *api.Sandbox, prefix string) string {
	return prefix + "-size"
}

// LastIdKey holds the highest id ever allocated in a collection. It only
// grows, which is what makes an id never reused.
func LastIdKey(sandbox *api.Sandbox, prefix string) string {
	return prefix + "-last-id"
}

// ListKey holds the id living at one position of a collection's dense list.
// Positions run from 1 to the value of SizeKey with no gap, which is what
// makes iteration possible without listing keys.
func ListKey(sandbox *api.Sandbox, prefix string, position int64) string {
	return sandbox.Deps.Std.Sprintf("%s-list-%d", prefix, position)
}

// PositionKey holds the position a record currently occupies in the dense
// list. It is the back-pointer that makes a removal cost the same whatever
// the size of the collection, and its presence is what marks a record live.
func PositionKey(sandbox *api.Sandbox, prefix string, id int64) string {
	return sandbox.Deps.Std.Sprintf("%s-%d-position", prefix, id)
}

// ValueKey holds one field value of one record.
func ValueKey(sandbox *api.Sandbox, prefix string, id int64, field string) string {
	return sandbox.Deps.Std.Sprintf("%s-%d-values-%s", prefix, id, field)
}

// IndexKey holds the id owning one value of one Key field — the unique
// index, addressed by the hash of the value so a lookup is a single read.
func IndexKey(sandbox *api.Sandbox, prefix string, field string, hash string) string {
	return sandbox.Deps.Std.Sprintf("%s-keys-%s-%s", prefix, field, hash)
}

// SubPrefix is the collection prefix of a nested (Database) field of one
// record. A nested collection is a collection like any other, which is why
// every helper here works on it unchanged.
func SubPrefix(sandbox *api.Sandbox, prefix string, id int64, field string) string {
	return sandbox.Deps.Std.Sprintf("%s-%d-%s", prefix, id, field)
}

// HashIndexValue normalizes and hashes an encoded value, so index lookups
// are case-insensitive and a key never grows with the value it indexes.
func HashIndexValue(sandbox *api.Sandbox, encoded string) string {
	lowered := sandbox.Deps.Stringsdeps.ToLower(encoded)
	return sandbox.Deps.Hashdeps.Sha256Hex([]byte(lowered))
}

// FindItem returns the schema field with the given name. ok is false when
// the schema declares no such field.
func FindItem(sandbox *api.Sandbox, items []api.Item, name string) (item api.Item, ok bool) {
	for _, candidate := range items {
		if candidate.Name == name {
			return candidate, true
		}
	}
	return api.Item{}, false
}

// InternalError wraps a storage failure as a typed *api.Error, which is the
// only way a backend error ever reaches a caller.
func InternalError(sandbox *api.Sandbox, err error) *api.Error {
	return liberror.New(sandbox, api.Internal, "", err.Error())
}

// ParseId reads an id back from the decimal form WriteInt stores it in.
func ParseId(sandbox *api.Sandbox, raw []byte) (int64, error) {
	id, err := sandbox.Deps.Stringsdeps.ParseInt(string(raw), 10, 64)
	if err != nil {
		return 0, sandbox.Deps.Std.Errorf("keep: invalid id: %s", string(raw))
	}
	return id, nil
}

// EncodeValue converts a caller-provided value to the canonical string form
// it is stored in, validating it against the field's type on the way.
func EncodeValue(sandbox *api.Sandbox, item api.Item, value any) (string, *api.Error) {
	switch item.Type {
	case api.Key:
		switch typed := value.(type) {
		case string:
			return typed, nil
		case stringer:
			return typed.String(), nil
		default:
			return "", liberror.NewWithValue(sandbox, api.InvalidField, item.Name, value,
				sandbox.Deps.Std.Sprintf("field %q expects a string value, got %T", item.Name, value))
		}
	case api.Int:
		switch typed := value.(type) {
		case int:
			return sandbox.Deps.Stringsdeps.FormatInt(int64(typed), 10), nil
		case int32:
			return sandbox.Deps.Stringsdeps.FormatInt(int64(typed), 10), nil
		case int64:
			return sandbox.Deps.Stringsdeps.FormatInt(typed, 10), nil
		default:
			return "", liberror.NewWithValue(sandbox, api.InvalidField, item.Name, value,
				sandbox.Deps.Std.Sprintf("field %q expects an integer value, got %T", item.Name, value))
		}
	default:
		return "", liberror.New(sandbox, api.InvalidField, item.Name,
			sandbox.Deps.Std.Sprintf("field %q cannot be encoded as a plain value", item.Name))
	}
}

// DecodeValue converts a stored value back to the typed form a caller of
// SchemaItem.Get receives: an int64 for an Int field, a string otherwise.
func DecodeValue(sandbox *api.Sandbox, item api.Item, raw []byte) (any, *api.Error) {
	switch item.Type {
	case api.Int:
		number, err := sandbox.Deps.Stringsdeps.ParseInt(string(raw), 10, 64)
		if err != nil {
			return nil, InternalError(sandbox, err)
		}
		return number, nil
	default:
		return string(raw), nil
	}
}

// ReadCount reads an integer key, treating a key that holds nothing as
// zero: a collection nothing was ever written to has no size key, and its
// size is zero.
func ReadCount(sandbox *api.Sandbox, key string) (int64, error) {
	raw, found, err := sandbox.Deps.Storagedeps.Read(key)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, nil
	}
	return sandbox.Deps.Stringsdeps.ParseInt(string(raw), 10, 64)
}

// WriteInt stores an integer under key in the canonical decimal form every
// reader here expects.
func WriteInt(sandbox *api.Sandbox, key string, value int64) error {
	encoded := sandbox.Deps.Stringsdeps.FormatInt(value, 10)
	return sandbox.Deps.Storagedeps.Write(key, []byte(encoded))
}
