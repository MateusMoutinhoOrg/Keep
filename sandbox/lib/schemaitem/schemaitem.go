package schemaitem

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/lib/dense"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/lib/liberror"
)

// Factories filling the function fields of api.SchemaItem. Each takes a
// pointer to the api.SchemaItem being built (the carrier) and returns
// the closure for one field, reading s.Deps, s.Items, s.Prefix and s.Id
// through the pointer at call time.

// GetFactory fills api.SchemaItem.Get.
func GetFactory(s *api.SchemaItem) func(fieldName string) (any, *api.Error) {
	return func(fieldName string) (any, *api.Error) {
		item, ok := dense.FindItem(s.Items, fieldName)
		if !ok {
			return nil, liberror.New(api.InvalidField, fieldName,
				fmt.Sprintf("field %q is not part of the schema", fieldName))
		}
		if item.Type == api.Database {
			return nil, liberror.New(api.InvalidField, fieldName,
				fmt.Sprintf("field %q is a sub-database, use ListAll(%q)", fieldName, fieldName))
		}
		raw, err := s.Deps.Read(dense.ValueKey(s.Prefix, s.Id, fieldName))
		if errors.Is(err, deps.ErrKeyNotFound) {
			return nil, liberror.New(api.NotFound, fieldName,
				fmt.Sprintf("field %q has no value for this record", fieldName))
		}
		if err != nil {
			return nil, dense.InternalError(err)
		}
		return dense.DecodeValue(item, raw)
	}
}

// UpdateFactory fills api.SchemaItem.Update. For indexed (Key) fields it
// performs the safe re-index sequence: write the new index entry, write
// the value, then delete the old index entry — so a crash mid-update
// never leaves the record unreachable.
func UpdateFactory(s *api.SchemaItem) func(fieldName string, value any) *api.Error {
	return func(fieldName string, value any) *api.Error {
		item, ok := dense.FindItem(s.Items, fieldName)
		if !ok {
			return liberror.New(api.InvalidField, fieldName,
				fmt.Sprintf("field %q is not part of the schema", fieldName))
		}
		if item.Type == api.Database {
			return liberror.New(api.InvalidField, fieldName,
				fmt.Sprintf("field %q is a sub-database and cannot be updated directly", fieldName))
		}
		encoded, e := dense.EncodeValue(item, value)
		if e != nil {
			return e
		}
		d := s.Deps
		vk := dense.ValueKey(s.Prefix, s.Id, fieldName)

		if item.Type != api.Key {
			if err := d.Write(vk, []byte(encoded)); err != nil {
				return dense.InternalError(err)
			}
			return nil
		}

		// Step 1: read the old value to locate the old index entry.
		oldRaw, err := d.Read(vk)
		if err != nil && !errors.Is(err, deps.ErrKeyNotFound) {
			return dense.InternalError(err)
		}

		// Step 2: reject the new value if another record already owns it.
		newHash := dense.HashIndexValue(encoded)
		existing, err := d.Read(dense.IndexKey(s.Prefix, fieldName, newHash))
		if err == nil {
			otherID, parseErr := dense.ParseID(existing)
			if parseErr != nil {
				return dense.InternalError(parseErr)
			}
			if otherID != s.Id {
				return liberror.NewWithValue(api.KeyConflict, fieldName, value,
					fmt.Sprintf("value for key %q already exists", fieldName))
			}
		} else if !errors.Is(err, deps.ErrKeyNotFound) {
			return dense.InternalError(err)
		}

		// Step 3: write the new index entry before touching anything else.
		if err := dense.WriteInt(d, dense.IndexKey(s.Prefix, fieldName, newHash), s.Id); err != nil {
			return dense.InternalError(err)
		}
		// Step 4: write the new value.
		if err := d.Write(vk, []byte(encoded)); err != nil {
			return dense.InternalError(err)
		}
		// Step 5: delete the old index entry.
		if oldRaw != nil {
			oldHash := dense.HashIndexValue(string(oldRaw))
			if oldHash != newHash {
				if err := d.Delete(dense.IndexKey(s.Prefix, fieldName, oldHash)); err != nil {
					return dense.InternalError(err)
				}
			}
		}
		return nil
	}
}

// RemoveFactory fills api.SchemaItem.Remove, deleting the record using
// the swap-with-last procedure, keeping the list dense at constant cost.
// Note that this moves the last record into the freed position, so list
// order is not stable.
func RemoveFactory(s *api.SchemaItem) func() *api.Error {
	return func() *api.Error {
		d := s.Deps

		// Step 1: read the victim's position; missing means already gone.
		position, err := readPosition(d, s.Prefix, s.Id)
		if errors.Is(err, deps.ErrKeyNotFound) {
			return nil
		}
		if err != nil {
			return dense.InternalError(err)
		}

		// Step 2: locate the last record.
		size, err := dense.ReadCount(d, dense.SizeKey(s.Prefix))
		if err != nil {
			return dense.InternalError(err)
		}
		lastRaw, err := d.Read(dense.ListKey(s.Prefix, size))
		if err != nil {
			return dense.InternalError(err)
		}
		lastID, err := dense.ParseID(lastRaw)
		if err != nil {
			return dense.InternalError(err)
		}

		// Step 3: move the last record into the hole.
		if position != size {
			if err := d.Write(dense.ListKey(s.Prefix, position), lastRaw); err != nil {
				return dense.InternalError(err)
			}
			if err := dense.WriteInt(d, dense.PositionKey(s.Prefix, lastID), position); err != nil {
				return dense.InternalError(err)
			}
		}

		// Step 4: shrink the list.
		if err := d.Delete(dense.ListKey(s.Prefix, size)); err != nil {
			return dense.InternalError(err)
		}
		if err := dense.WriteInt(d, dense.SizeKey(s.Prefix), size-1); err != nil {
			return dense.InternalError(err)
		}

		// Step 5: remove the unique index entries.
		for _, item := range s.Items {
			if item.Type != api.Key {
				continue
			}
			raw, err := d.Read(dense.ValueKey(s.Prefix, s.Id, item.Name))
			if errors.Is(err, deps.ErrKeyNotFound) {
				continue
			}
			if err != nil {
				return dense.InternalError(err)
			}
			if err := d.Delete(dense.IndexKey(s.Prefix, item.Name, dense.HashIndexValue(string(raw)))); err != nil {
				return dense.InternalError(err)
			}
		}

		// Step 6: remove the record's data, including nested sub-databases.
		for _, item := range s.Items {
			if item.Type == api.Database {
				if e := ClearCollection(s.Deps, item.Itens, dense.SubPrefix(s.Prefix, s.Id, item.Name)); e != nil {
					return e
				}
				continue
			}
			if err := d.Delete(dense.ValueKey(s.Prefix, s.Id, item.Name)); err != nil {
				return dense.InternalError(err)
			}
		}
		if err := d.Delete(dense.PositionKey(s.Prefix, s.Id)); err != nil {
			return dense.InternalError(err)
		}
		return nil
	}
}

// CheckKeysPresenceFactory fills api.SchemaItem.CheckKeysPresence.
func CheckKeysPresenceFactory(s *api.SchemaItem) func(keys []string) bool {
	return func(keys []string) bool {
		for _, key := range keys {
			exists, err := s.Deps.Exists(dense.ValueKey(s.Prefix, s.Id, key))
			if err != nil || !exists {
				return false
			}
		}
		return true
	}
}

// ListAllFactory fills api.SchemaItem.ListAll, returning every record of
// a sub-database (Database) field.
func ListAllFactory(s *api.SchemaItem) func(fieldName string) []api.SchemaItem {
	return func(fieldName string) []api.SchemaItem {
		item, ok := dense.FindItem(s.Items, fieldName)
		if !ok || item.Type != api.Database {
			return nil
		}
		result, e := ListRange(s.Deps, item.Itens, dense.SubPrefix(s.Prefix, s.Id, fieldName), 1, 0)
		if e != nil {
			return nil
		}
		return result
	}
}

// NewSubItemFactory fills api.SchemaItem.NewSubItem, inserting a record
// into a sub-database (Database) field of this record.
func NewSubItemFactory(s *api.SchemaItem) func(fieldName string, fields map[string]any) (api.SchemaItem, *api.Error) {
	return func(fieldName string, fields map[string]any) (api.SchemaItem, *api.Error) {
		item, ok := dense.FindItem(s.Items, fieldName)
		if !ok || item.Type != api.Database {
			return api.SchemaItem{}, liberror.New(api.InvalidField, fieldName,
				fmt.Sprintf("field %q is not a sub-database of the schema", fieldName))
		}
		return New(s.Deps, item.Itens, dense.SubPrefix(s.Prefix, s.Id, fieldName), fields)
	}
}

// StringFactory fills api.SchemaItem.String, rendering the record's
// plain fields so samples printing a record show its data.
func StringFactory(s *api.SchemaItem) func() string {
	return func() string {
		parts := make([]string, 0, len(s.Items))
		for _, item := range s.Items {
			if item.Type == api.Database {
				continue
			}
			value, e := s.Get(item.Name)
			if e != nil {
				continue
			}
			parts = append(parts, fmt.Sprintf("%s: %v", item.Name, value))
		}
		return fmt.Sprintf("{id: %d, %s}", s.Id, strings.Join(parts, ", "))
	}
}

func readPosition(d deps.Deps, prefix string, id int64) (int64, error) {
	raw, err := d.Read(dense.PositionKey(prefix, id))
	if err != nil {
		return 0, err
	}
	return dense.ParseID(raw)
}

// build assembles an api.SchemaItem for an existing record id, running
// every field factory over it. It is the shared aggregate behind New,
// ResolveLive, ListRange and ClearCollection.
func build(d deps.Deps, items []api.Item, prefix string, id int64) api.SchemaItem {
	s := api.SchemaItem{Deps: d, Items: items, Prefix: prefix, Id: id}
	s.Get = GetFactory(&s)
	s.Update = UpdateFactory(&s)
	s.Remove = RemoveFactory(&s)
	s.CheckKeysPresence = CheckKeysPresenceFactory(&s)
	s.ListAll = ListAllFactory(&s)
	s.NewSubItem = NewSubItemFactory(&s)
	s.String = StringFactory(&s)
	return s
}

// New inserts a record into the collection identified by prefix,
// following the insertion procedure of the dense record pattern.
func New(d deps.Deps, items []api.Item, prefix string, fields map[string]any) (api.SchemaItem, *api.Error) {

	// Validate provided fields against the schema.
	for name := range fields {
		item, ok := dense.FindItem(items, name)
		if !ok {
			return api.SchemaItem{}, liberror.New(api.InvalidField, name,
				fmt.Sprintf("field %q is not part of the schema", name))
		}
		if item.Type == api.Database {
			return api.SchemaItem{}, liberror.New(api.InvalidField, name,
				fmt.Sprintf("field %q is a sub-database and cannot be set directly", name))
		}
	}
	for _, item := range items {
		if item.Required && item.Type != api.Database {
			if _, ok := fields[item.Name]; !ok {
				return api.SchemaItem{}, liberror.New(api.MissingField, item.Name,
					fmt.Sprintf("required field %q is missing", item.Name))
			}
		}
	}

	encoded := make(map[string]string, len(fields))
	for name, value := range fields {
		item, _ := dense.FindItem(items, name)
		enc, e := dense.EncodeValue(item, value)
		if e != nil {
			return api.SchemaItem{}, e
		}
		encoded[name] = enc
	}

	// Step 1: uniqueness check on every indexed (Key) field.
	for _, item := range items {
		if item.Type != api.Key {
			continue
		}
		enc, ok := encoded[item.Name]
		if !ok {
			continue
		}
		exists, err := d.Exists(dense.IndexKey(prefix, item.Name, dense.HashIndexValue(enc)))
		if err != nil {
			return api.SchemaItem{}, dense.InternalError(err)
		}
		if exists {
			return api.SchemaItem{}, liberror.NewWithValue(api.KeyConflict, item.Name, fields[item.Name],
				fmt.Sprintf("value for key %q already exists", item.Name))
		}
	}

	// Step 2: allocate the id (never reused, only grows).
	lastID, err := dense.ReadCount(d, dense.LastIDKey(prefix))
	if err != nil {
		return api.SchemaItem{}, dense.InternalError(err)
	}
	id := lastID + 1
	if err := dense.WriteInt(d, dense.LastIDKey(prefix), id); err != nil {
		return api.SchemaItem{}, dense.InternalError(err)
	}

	// Step 3: the new position is size+1.
	size, err := dense.ReadCount(d, dense.SizeKey(prefix))
	if err != nil {
		return api.SchemaItem{}, dense.InternalError(err)
	}
	position := size + 1

	// Step 4: write the record's data and back-pointer.
	for name, enc := range encoded {
		if err := d.Write(dense.ValueKey(prefix, id, name), []byte(enc)); err != nil {
			return api.SchemaItem{}, dense.InternalError(err)
		}
	}
	if err := dense.WriteInt(d, dense.PositionKey(prefix, id), position); err != nil {
		return api.SchemaItem{}, dense.InternalError(err)
	}

	// Step 5: write the unique index entries.
	for _, item := range items {
		if item.Type != api.Key {
			continue
		}
		enc, ok := encoded[item.Name]
		if !ok {
			continue
		}
		if err := dense.WriteInt(d, dense.IndexKey(prefix, item.Name, dense.HashIndexValue(enc)), id); err != nil {
			return api.SchemaItem{}, dense.InternalError(err)
		}
	}

	// Step 6: publish — list slot first, size last (commit point).
	if err := dense.WriteInt(d, dense.ListKey(prefix, position), id); err != nil {
		return api.SchemaItem{}, dense.InternalError(err)
	}
	if err := dense.WriteInt(d, dense.SizeKey(prefix), position); err != nil {
		return api.SchemaItem{}, dense.InternalError(err)
	}

	return build(d, items, prefix, id), nil
}

// ResolveLive parses an id read from an index entry and returns the
// record only if it is still live (its position back-pointer exists).
// ok is false when it is not.
func ResolveLive(d deps.Deps, items []api.Item, prefix string, rawID []byte) (api.SchemaItem, bool) {
	id, err := dense.ParseID(rawID)
	if err != nil {
		return api.SchemaItem{}, false
	}
	exists, err := d.Exists(dense.PositionKey(prefix, id))
	if err != nil || !exists {
		return api.SchemaItem{}, false
	}
	return build(d, items, prefix, id), true
}

// ListRange reads records from the dense list starting at `from`
// (1-based). A chunk of 0 means "until the end of the list".
func ListRange(d deps.Deps, items []api.Item, prefix string, from int64, chunk int64) ([]api.SchemaItem, *api.Error) {

	size, err := dense.ReadCount(d, dense.SizeKey(prefix))
	if err != nil {
		return nil, dense.InternalError(err)
	}
	if from < 1 {
		from = 1
	}
	to := size
	if chunk > 0 && from+chunk-1 < size {
		to = from + chunk - 1
	}
	result := make([]api.SchemaItem, 0)
	for p := from; p <= to; p++ {
		id, err := dense.ReadCount(d, dense.ListKey(prefix, p))
		if err != nil {
			return nil, dense.InternalError(err)
		}
		result = append(result, build(d, items, prefix, id))
	}
	return result, nil
}

// ClearCollection removes every record of a collection (used when a
// parent record owning a sub-database is deleted). Records are removed
// from the last position backwards so no swap is ever needed.
func ClearCollection(d deps.Deps, items []api.Item, prefix string) *api.Error {

	for {
		size, err := dense.ReadCount(d, dense.SizeKey(prefix))
		if err != nil {
			return dense.InternalError(err)
		}
		if size == 0 {
			break
		}
		id, err := dense.ReadCount(d, dense.ListKey(prefix, size))
		if err != nil {
			return dense.InternalError(err)
		}
		item := build(d, items, prefix, id)
		if e := item.Remove(); e != nil {
			return e
		}
	}
	if err := d.Delete(dense.SizeKey(prefix)); err != nil {
		return dense.InternalError(err)
	}
	if err := d.Delete(dense.LastIDKey(prefix)); err != nil {
		return dense.InternalError(err)
	}
	return nil
}
