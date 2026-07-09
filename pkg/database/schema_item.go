package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type SchemaItem struct {
	db     *KeepDatabase
	items  []Item
	prefix string
	id     int64
}

// Id returns the record's permanent, never-reused identifier.
func (s *SchemaItem) Id() int64 { return s.id }

// CheckKeysPresence reports whether every named field has a stored
// value for this record.
func (s *SchemaItem) CheckKeysPresence(keys []string) bool {
	for _, key := range keys {
		exists, err := s.db.Deps.Exists(valueKey(s.prefix, s.id, key))
		if err != nil || !exists {
			return false
		}
	}
	return true
}

func (s *SchemaItem) Get(fieldName string) (any, *Error) {
	item := findItem(s.items, fieldName)
	if item == nil {
		return nil, &Error{Type: InvalidField, Key: fieldName,
			Msg: fmt.Sprintf("field %q is not part of the schema", fieldName)}
	}
	if item.Type == Database {
		return nil, &Error{Type: InvalidField, Key: fieldName,
			Msg: fmt.Sprintf("field %q is a sub-database, use ListAll(%q)", fieldName, fieldName)}
	}
	raw, err := s.db.Deps.Read(valueKey(s.prefix, s.id, fieldName))
	if errors.Is(err, deps.ErrKeyNotFound) {
		return nil, &Error{Type: NotFound, Key: fieldName,
			Msg: fmt.Sprintf("field %q has no value for this record", fieldName)}
	}
	if err != nil {
		return nil, internalError(err)
	}
	return decodeValue(item, raw)
}

// Update writes a new value for a field. For indexed (Key) fields it
// performs the safe re-index sequence: write the new index entry, write
// the value, then delete the old index entry — so a crash mid-update
// never leaves the record unreachable.
func (s *SchemaItem) Update(fieldName string, value any) *Error {
	item := findItem(s.items, fieldName)
	if item == nil {
		return &Error{Type: InvalidField, Key: fieldName,
			Msg: fmt.Sprintf("field %q is not part of the schema", fieldName)}
	}
	if item.Type == Database {
		return &Error{Type: InvalidField, Key: fieldName,
			Msg: fmt.Sprintf("field %q is a sub-database and cannot be updated directly", fieldName)}
	}
	encoded, e := encodeValue(item, value)
	if e != nil {
		return e
	}
	d := s.db.Deps
	vk := valueKey(s.prefix, s.id, fieldName)

	if item.Type != Key {
		if err := d.Write(vk, []byte(encoded)); err != nil {
			return internalError(err)
		}
		return nil
	}

	// Step 1: read the old value to locate the old index entry.
	oldRaw, err := d.Read(vk)
	if err != nil && !errors.Is(err, deps.ErrKeyNotFound) {
		return internalError(err)
	}

	// Step 2: reject the new value if another record already owns it.
	newHash := hashIndexValue(encoded)
	existing, err := d.Read(indexKey(s.prefix, fieldName, newHash))
	if err == nil {
		otherID, parseErr := parseID(existing)
		if parseErr != nil {
			return internalError(parseErr)
		}
		if otherID != s.id {
			return &Error{Type: KeyConflict, Key: fieldName, KeyValue: value,
				Msg: fmt.Sprintf("value for key %q already exists", fieldName)}
		}
	} else if !errors.Is(err, deps.ErrKeyNotFound) {
		return internalError(err)
	}

	// Step 3: write the new index entry before touching anything else.
	if err := writeInt(d, indexKey(s.prefix, fieldName, newHash), s.id); err != nil {
		return internalError(err)
	}
	// Step 4: write the new value.
	if err := d.Write(vk, []byte(encoded)); err != nil {
		return internalError(err)
	}
	// Step 5: delete the old index entry.
	if oldRaw != nil {
		oldHash := hashIndexValue(string(oldRaw))
		if oldHash != newHash {
			if err := d.Delete(indexKey(s.prefix, fieldName, oldHash)); err != nil {
				return internalError(err)
			}
		}
	}
	return nil
}

// Remove deletes the record using the swap-with-last procedure, keeping
// the list dense at constant cost. Note that this moves the last record
// into the freed position, so list order is not stable.
func (s *SchemaItem) Remove() Error {
	d := s.db.Deps

	// Step 1: read the victim's position; missing means already gone.
	position, err := readPosition(d, s.prefix, s.id)
	if errors.Is(err, deps.ErrKeyNotFound) {
		return Error{}
	}
	if err != nil {
		return *internalError(err)
	}

	// Step 2: locate the last record.
	size, err := readCount(d, sizeKey(s.prefix))
	if err != nil {
		return *internalError(err)
	}
	lastRaw, err := d.Read(listKey(s.prefix, size))
	if err != nil {
		return *internalError(err)
	}
	lastID, err := parseID(lastRaw)
	if err != nil {
		return *internalError(err)
	}

	// Step 3: move the last record into the hole.
	if position != size {
		if err := d.Write(listKey(s.prefix, position), lastRaw); err != nil {
			return *internalError(err)
		}
		if err := writeInt(d, positionKey(s.prefix, lastID), position); err != nil {
			return *internalError(err)
		}
	}

	// Step 4: shrink the list.
	if err := d.Delete(listKey(s.prefix, size)); err != nil {
		return *internalError(err)
	}
	if err := writeInt(d, sizeKey(s.prefix), size-1); err != nil {
		return *internalError(err)
	}

	// Step 5: remove the unique index entries.
	for i := range s.items {
		item := &s.items[i]
		if item.Type != Key {
			continue
		}
		raw, err := d.Read(valueKey(s.prefix, s.id, item.Name))
		if errors.Is(err, deps.ErrKeyNotFound) {
			continue
		}
		if err != nil {
			return *internalError(err)
		}
		if err := d.Delete(indexKey(s.prefix, item.Name, hashIndexValue(string(raw)))); err != nil {
			return *internalError(err)
		}
	}

	// Step 6: remove the record's data, including nested sub-databases.
	for i := range s.items {
		item := &s.items[i]
		if item.Type == Database {
			if e := clearCollection(s.db, item.Itens, subPrefix(s.prefix, s.id, item.Name)); e != nil {
				return *e
			}
			continue
		}
		if err := d.Delete(valueKey(s.prefix, s.id, item.Name)); err != nil {
			return *internalError(err)
		}
	}
	if err := d.Delete(positionKey(s.prefix, s.id)); err != nil {
		return *internalError(err)
	}
	return Error{}
}

// ListAll returns every record of a sub-database (Database type) field.
func (s *SchemaItem) ListAll(fieldName string) []*SchemaItem {
	item := findItem(s.items, fieldName)
	if item == nil || item.Type != Database {
		return nil
	}
	result, e := listRange(s.db, item.Itens, subPrefix(s.prefix, s.id, fieldName), 1, 0)
	if e != nil {
		return nil
	}
	return result
}

// NewSubItem inserts a record into a sub-database (Database type) field
// of this record.
func (s *SchemaItem) NewSubItem(fieldName string, fields map[string]any) (*SchemaItem, *Error) {
	item := findItem(s.items, fieldName)
	if item == nil || item.Type != Database {
		return nil, &Error{Type: InvalidField, Key: fieldName,
			Msg: fmt.Sprintf("field %q is not a sub-database of the schema", fieldName)}
	}
	return newItem(s.db, item.Itens, subPrefix(s.prefix, s.id, fieldName), fields)
}

func readPosition(d deps.Deps, prefix string, id int64) (int64, error) {
	raw, err := d.Read(positionKey(prefix, id))
	if err != nil {
		return 0, err
	}
	return parseID(raw)
}

// String renders the record's plain fields, so samples printing a
// *SchemaItem show its data instead of internal pointers.
func (s *SchemaItem) String() string {
	parts := make([]string, 0, len(s.items))
	for i := range s.items {
		item := &s.items[i]
		if item.Type == Database {
			continue
		}
		value, e := s.Get(item.Name)
		if e != nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s: %v", item.Name, value))
	}
	return fmt.Sprintf("{id: %d, %s}", s.id, strings.Join(parts, ", "))
}
