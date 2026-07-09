package database

import (
	"errors"
	"strconv"
)

type SchemaInstance struct {
	db     *KeepDatabase
	schema *Schema
	items  []Item
	prefix string
}

// FindByKey looks a record up through its unique index: normalize and
// hash the value, resolve the id, and check the record is live.
// Returns nil when the key field is unknown or no record matches.
func (si *SchemaInstance) FindByKey(key string, keyValue any) *SchemaItem {
	item := findItem(si.items, key)
	if item == nil || item.Type != Key {
		return nil
	}
	encoded, e := encodeValue(item, keyValue)
	if e != nil {
		return nil
	}
	raw, err := si.db.Deps.Read(indexKey(si.prefix, key, hashIndexValue(encoded)))
	if err != nil {
		return nil
	}
	return resolveLive(si.db, si.items, si.prefix, raw)
}

// resolveLive parses an id read from an index entry and returns the
// record only if it is still live (its position back-pointer exists).
func resolveLive(db *KeepDatabase, items []Item, prefix string, rawID []byte) *SchemaItem {
	id, err := parseID(rawID)
	if err != nil {
		return nil
	}
	exists, err := db.Deps.Exists(positionKey(prefix, id))
	if err != nil || !exists {
		return nil
	}
	return &SchemaItem{db: db, items: items, prefix: prefix, id: id}
}

func (si *SchemaInstance) NewItem(fields map[string]any) (*SchemaItem, *Error) {
	return newItem(si.db, si.items, si.prefix, fields)
}

// ListAll iterates the dense list from position 1 through size.
func (si *SchemaInstance) ListAll() ([]*SchemaItem, *Error) {
	return listRange(si.db, si.items, si.prefix, 1, 0)
}

// List returns up to `chunk` records starting at `position` (1-based).
func (si *SchemaInstance) List(position int, chunk int) ([]*SchemaItem, *Error) {
	return listRange(si.db, si.items, si.prefix, int64(position), int64(chunk))
}

func parseID(raw []byte) (int64, error) {
	id, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return 0, errors.New("keep: invalid id: " + string(raw))
	}
	return id, nil
}
