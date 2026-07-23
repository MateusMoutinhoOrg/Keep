package database

// Implementation of the Dense Record Pattern described in
// docs/Explanation/DenseRecordPattern.md. Every operation is expressed as single-key
// reads/writes/deletes against the deps.Deps backend, and assumes a
// single writer.

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

func sizeKey(prefix string) string   { return prefix + "-size" }
func lastIDKey(prefix string) string { return prefix + "-last-id" }

func listKey(prefix string, position int64) string {
	return fmt.Sprintf("%s-list-%d", prefix, position)
}

func positionKey(prefix string, id int64) string {
	return fmt.Sprintf("%s-%d-position", prefix, id)
}

func valueKey(prefix string, id int64, field string) string {
	return fmt.Sprintf("%s-%d-values-%s", prefix, id, field)
}

func indexKey(prefix string, field string, hash string) string {
	return fmt.Sprintf("%s-keys-%s-%s", prefix, field, hash)
}

// subPrefix is the collection prefix of a nested (Database type) field
// of a given record.
func subPrefix(prefix string, id int64, field string) string {
	return fmt.Sprintf("%s-%d-%s", prefix, id, field)
}

// hashIndexValue normalizes (lowercases) and hashes an encoded value so
// index lookups are case-insensitive and key length stays bounded.
func hashIndexValue(encoded string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(encoded)))
	return hex.EncodeToString(sum[:])
}

func findItem(items []Item, name string) *Item {
	for i := range items {
		if items[i].Name == name {
			return &items[i]
		}
	}
	return nil
}

func internalError(err error) *Error {
	return &Error{Type: Internal, Msg: err.Error()}
}

// encodeValue converts a caller-provided value to its canonical stored
// string form, validating it against the item's type.
func encodeValue(item *Item, value any) (string, *Error) {
	switch item.Type {
	case Key:
		switch v := value.(type) {
		case string:
			return v, nil
		case fmt.Stringer:
			return v.String(), nil
		default:
			return "", &Error{Type: InvalidField, Key: item.Name, KeyValue: value,
				Msg: fmt.Sprintf("field %q expects a string value, got %T", item.Name, value)}
		}
	case Int:
		switch v := value.(type) {
		case int:
			return strconv.Itoa(v), nil
		case int32:
			return strconv.FormatInt(int64(v), 10), nil
		case int64:
			return strconv.FormatInt(v, 10), nil
		default:
			return "", &Error{Type: InvalidField, Key: item.Name, KeyValue: value,
				Msg: fmt.Sprintf("field %q expects an integer value, got %T", item.Name, value)}
		}
	default:
		return "", &Error{Type: InvalidField, Key: item.Name,
			Msg: fmt.Sprintf("field %q cannot be encoded as a plain value", item.Name)}
	}
}

// decodeValue converts a stored value back to its typed form.
func decodeValue(item *Item, raw []byte) (any, *Error) {
	switch item.Type {
	case Int:
		n, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return nil, internalError(err)
		}
		return n, nil
	default:
		return string(raw), nil
	}
}

// readCount reads an integer key, treating a missing key as zero.
func readCount(d deps.Deps, key string) (int64, error) {
	raw, err := d.Read(key)
	if errors.Is(err, deps.ErrKeyNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(raw), 10, 64)
}

func writeInt(d deps.Deps, key string, value int64) error {
	return d.Write(key, []byte(strconv.FormatInt(value, 10)))
}

// newItem inserts a record into the collection identified by prefix,
// following the insertion procedure of the dense record pattern.
func newItem(db *KeepDatabase, items []Item, prefix string, fields map[string]any) (*SchemaItem, *Error) {
	d := db.Deps

	// Validate provided fields against the schema.
	for name := range fields {
		item := findItem(items, name)
		if item == nil {
			return nil, &Error{Type: InvalidField, Key: name,
				Msg: fmt.Sprintf("field %q is not part of the schema", name)}
		}
		if item.Type == Database {
			return nil, &Error{Type: InvalidField, Key: name,
				Msg: fmt.Sprintf("field %q is a sub-database and cannot be set directly", name)}
		}
	}
	for i := range items {
		item := &items[i]
		if item.Required && item.Type != Database {
			if _, ok := fields[item.Name]; !ok {
				return nil, &Error{Type: MissingField, Key: item.Name,
					Msg: fmt.Sprintf("required field %q is missing", item.Name)}
			}
		}
	}

	encoded := make(map[string]string, len(fields))
	for name, value := range fields {
		enc, e := encodeValue(findItem(items, name), value)
		if e != nil {
			return nil, e
		}
		encoded[name] = enc
	}

	// Step 1: uniqueness check on every indexed (Key) field.
	for i := range items {
		item := &items[i]
		if item.Type != Key {
			continue
		}
		enc, ok := encoded[item.Name]
		if !ok {
			continue
		}
		exists, err := d.Exists(indexKey(prefix, item.Name, hashIndexValue(enc)))
		if err != nil {
			return nil, internalError(err)
		}
		if exists {
			return nil, &Error{Type: KeyConflict, Key: item.Name, KeyValue: fields[item.Name],
				Msg: fmt.Sprintf("value for key %q already exists", item.Name)}
		}
	}

	// Step 2: allocate the id (never reused, only grows).
	lastID, err := readCount(d, lastIDKey(prefix))
	if err != nil {
		return nil, internalError(err)
	}
	id := lastID + 1
	if err := writeInt(d, lastIDKey(prefix), id); err != nil {
		return nil, internalError(err)
	}

	// Step 3: the new position is size+1.
	size, err := readCount(d, sizeKey(prefix))
	if err != nil {
		return nil, internalError(err)
	}
	position := size + 1

	// Step 4: write the record's data and back-pointer.
	for name, enc := range encoded {
		if err := d.Write(valueKey(prefix, id, name), []byte(enc)); err != nil {
			return nil, internalError(err)
		}
	}
	if err := writeInt(d, positionKey(prefix, id), position); err != nil {
		return nil, internalError(err)
	}

	// Step 5: write the unique index entries.
	for i := range items {
		item := &items[i]
		if item.Type != Key {
			continue
		}
		enc, ok := encoded[item.Name]
		if !ok {
			continue
		}
		if err := writeInt(d, indexKey(prefix, item.Name, hashIndexValue(enc)), id); err != nil {
			return nil, internalError(err)
		}
	}

	// Step 6: publish — list slot first, size last (commit point).
	if err := writeInt(d, listKey(prefix, position), id); err != nil {
		return nil, internalError(err)
	}
	if err := writeInt(d, sizeKey(prefix), position); err != nil {
		return nil, internalError(err)
	}

	return &SchemaItem{db: db, items: items, prefix: prefix, id: id}, nil
}

// listRange reads records from the dense list starting at `from`
// (1-based). A chunk of 0 means "until the end of the list".
func listRange(db *KeepDatabase, items []Item, prefix string, from int64, chunk int64) ([]*SchemaItem, *Error) {
	d := db.Deps
	size, err := readCount(d, sizeKey(prefix))
	if err != nil {
		return nil, internalError(err)
	}
	if from < 1 {
		from = 1
	}
	to := size
	if chunk > 0 && from+chunk-1 < size {
		to = from + chunk - 1
	}
	result := make([]*SchemaItem, 0)
	for p := from; p <= to; p++ {
		id, err := readCount(d, listKey(prefix, p))
		if err != nil {
			return nil, internalError(err)
		}
		result = append(result, &SchemaItem{db: db, items: items, prefix: prefix, id: id})
	}
	return result, nil
}

// clearCollection removes every record of a collection (used when a
// parent record owning a sub-database is deleted). Records are removed
// from the last position backwards so no swap is ever needed.
func clearCollection(db *KeepDatabase, items []Item, prefix string) *Error {
	d := db.Deps
	for {
		size, err := readCount(d, sizeKey(prefix))
		if err != nil {
			return internalError(err)
		}
		if size == 0 {
			break
		}
		id, err := readCount(d, listKey(prefix, size))
		if err != nil {
			return internalError(err)
		}
		item := &SchemaItem{db: db, items: items, prefix: prefix, id: id}
		if e := item.Remove(); e.Msg != "" {
			return &e
		}
	}
	if err := d.Delete(sizeKey(prefix)); err != nil {
		return internalError(err)
	}
	if err := d.Delete(lastIDKey(prefix)); err != nil {
		return internalError(err)
	}
	return nil
}
