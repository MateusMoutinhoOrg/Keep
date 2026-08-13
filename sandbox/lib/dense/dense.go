package dense

// Key layout and value encoding of the Dense Record Pattern described in
// docs/References/DenseRecordPattern.md. Everything here is expressed as
// single-key reads/writes against the deps.Deps backend, and assumes a
// single writer. The record operations built on top of these helpers live
// in the schemaitem package.

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/lib/liberror"
)

func SizeKey(prefix string) string   { return prefix + "-size" }
func LastIDKey(prefix string) string { return prefix + "-last-id" }

func ListKey(prefix string, position int64) string {
	return fmt.Sprintf("%s-list-%d", prefix, position)
}

func PositionKey(prefix string, id int64) string {
	return fmt.Sprintf("%s-%d-position", prefix, id)
}

func ValueKey(prefix string, id int64, field string) string {
	return fmt.Sprintf("%s-%d-values-%s", prefix, id, field)
}

func IndexKey(prefix string, field string, hash string) string {
	return fmt.Sprintf("%s-keys-%s-%s", prefix, field, hash)
}

// SubPrefix is the collection prefix of a nested (Database) field of a
// given record.
func SubPrefix(prefix string, id int64, field string) string {
	return fmt.Sprintf("%s-%d-%s", prefix, id, field)
}

// HashIndexValue normalizes (lowercases) and hashes an encoded value so
// index lookups are case-insensitive and key length stays bounded.
func HashIndexValue(encoded string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(encoded)))
	return hex.EncodeToString(sum[:])
}

// FindItem returns the schema field with the given name. ok is false
// when the schema declares no such field.
func FindItem(items []api.Item, name string) (item api.Item, ok bool) {
	for _, it := range items {
		if it.Name == name {
			return it, true
		}
	}
	return api.Item{}, false
}

// InternalError wraps a backend failure as a typed *api.Error.
func InternalError(err error) *api.Error {
	return liberror.New(api.Internal, "", err.Error())
}

// ParseID reads an id written by WriteInt.
func ParseID(raw []byte) (int64, error) {
	id, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return 0, errors.New("keep: invalid id: " + string(raw))
	}
	return id, nil
}

// EncodeValue converts a caller-provided value to its canonical stored
// string form, validating it against the item's type.
func EncodeValue(item api.Item, value any) (string, *api.Error) {
	switch item.Type {
	case api.Key:
		switch v := value.(type) {
		case string:
			return v, nil
		case fmt.Stringer:
			return v.String(), nil
		default:
			return "", liberror.NewWithValue(api.InvalidField, item.Name, value,
				fmt.Sprintf("field %q expects a string value, got %T", item.Name, value))
		}
	case api.Int:
		switch v := value.(type) {
		case int:
			return strconv.Itoa(v), nil
		case int32:
			return strconv.FormatInt(int64(v), 10), nil
		case int64:
			return strconv.FormatInt(v, 10), nil
		default:
			return "", liberror.NewWithValue(api.InvalidField, item.Name, value,
				fmt.Sprintf("field %q expects an integer value, got %T", item.Name, value))
		}
	default:
		return "", liberror.New(api.InvalidField, item.Name,
			fmt.Sprintf("field %q cannot be encoded as a plain value", item.Name))
	}
}

// DecodeValue converts a stored value back to its typed form.
func DecodeValue(item api.Item, raw []byte) (any, *api.Error) {
	switch item.Type {
	case api.Int:
		n, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return nil, InternalError(err)
		}
		return n, nil
	default:
		return string(raw), nil
	}
}

// ReadCount reads an integer key, treating a missing key as zero.
func ReadCount(d deps.Deps, key string) (int64, error) {
	raw, err := d.Read(key)
	if errors.Is(err, deps.ErrKeyNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(raw), 10, 64)
}

// WriteInt stores an integer under key in the canonical decimal form.
func WriteInt(d deps.Deps, key string, value int64) error {
	return d.Write(key, []byte(strconv.FormatInt(value, 10)))
}
