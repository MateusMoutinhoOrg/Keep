package schemaitem

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
)

// One record of a collection: the factories filling the function fields of
// api.SchemaItem, and the record-level operations the schemainstance
// package builds on.
//
// Every factory takes the sandbox and a pointer to the api.SchemaItem being
// built, and returns the closure for one field. Reading the record's Items,
// Prefix and Id back off that pointer is what lets Build assign the fields
// in any order, and what keeps the record itself free of any dependency
// field a caller could reach.
//
// The one thing a record needs that its own Items, Prefix and Id do not
// describe is the collection a Link field points at, so a dense.LinkResolver
// travels beside them, from the handle that built it down through every
// nested collection. It is a parameter and never a field, for the same
// reason storage is: an api type carries no wiring a caller could read or
// replace.

// GetFactory fills api.SchemaItem.Get.
func GetFactory(sandbox *api.Sandbox, record *api.SchemaItem) func(fieldName string) (any, *api.Error) {
	return func(fieldName string) (any, *api.Error) {
		item, ok := dense.FindItem(sandbox, record.Items, fieldName)
		if !ok {
			return nil, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.Std.Sprintf("field %q is not part of the schema", fieldName))
		}
		if item.Type == api.Database {
			return nil, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.Std.Sprintf("field %q is a nested collection, use ListAll(%q)", fieldName, fieldName))
		}
		raw, found, err := sandbox.Deps.Storagedeps.Read(dense.ValueKey(sandbox, record.Prefix, record.Id, fieldName))
		if err != nil {
			return nil, dense.InternalError(sandbox, err)
		}
		if !found {
			return nil, liberror.New(sandbox, api.NotFound, fieldName,
				sandbox.Deps.Std.Sprintf("field %q has no value for this record", fieldName))
		}
		return dense.DecodeValue(sandbox, item, raw)
	}
}

// GetLinkFactory fills api.SchemaItem.GetLink, following a Link field to the
// record it names in the collection its Target declares. Resolution goes
// through ResolveById, so a link to a record that has been removed reports
// ok == false — and since ids are never reused, it never reports a
// different record instead.
func GetLinkFactory(sandbox *api.Sandbox, record *api.SchemaItem, resolve dense.LinkResolver) func(fieldName string) (api.SchemaItem, bool) {
	return func(fieldName string) (api.SchemaItem, bool) {
		item, ok := dense.FindItem(sandbox, record.Items, fieldName)
		if !ok || item.Type != api.Link {
			return api.SchemaItem{}, false
		}
		targetItems, targetPrefix, ok := resolve(item.Target)
		if !ok {
			return api.SchemaItem{}, false
		}
		raw, found, err := sandbox.Deps.Storagedeps.Read(dense.ValueKey(sandbox, record.Prefix, record.Id, fieldName))
		if err != nil || !found {
			return api.SchemaItem{}, false
		}
		return ResolveLive(sandbox, targetItems, targetPrefix, resolve, raw)
	}
}

// UpdateFactory fills api.SchemaItem.Update. For an indexed (Key) field it
// runs the safe re-index sequence — write the new index entry, write the
// value, then delete the old index entry — so a crash part-way through
// never leaves the record unreachable: the worst outcome is an index entry
// pointing at a record that no longer holds the value, which resolves to
// nothing and is overwritten by the next write.
func UpdateFactory(sandbox *api.Sandbox, record *api.SchemaItem) func(fieldName string, value any) *api.Error {
	return func(fieldName string, value any) *api.Error {
		item, ok := dense.FindItem(sandbox, record.Items, fieldName)
		if !ok {
			return liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.Std.Sprintf("field %q is not part of the schema", fieldName))
		}
		if item.Type == api.Database {
			return liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.Std.Sprintf("field %q is a nested collection and cannot be updated directly", fieldName))
		}
		encoded, failure := dense.EncodeValue(sandbox, item, value)
		if failure != nil {
			return failure
		}
		storage := sandbox.Deps.Storagedeps
		valueKey := dense.ValueKey(sandbox, record.Prefix, record.Id, fieldName)

		if item.Type != api.Key {
			if err := storage.Write(valueKey, []byte(encoded)); err != nil {
				return dense.InternalError(sandbox, err)
			}
			return nil
		}

		// Step 1: read the old value, to locate the index entry it owns.
		oldRaw, hadValue, err := storage.Read(valueKey)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}

		// Step 2: refuse the new value when another record already owns it.
		newHash := dense.HashIndexValue(sandbox, encoded)
		existing, indexed, err := storage.Read(dense.IndexKey(sandbox, record.Prefix, fieldName, newHash))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if indexed {
			otherId, parseErr := dense.ParseId(sandbox, existing)
			if parseErr != nil {
				return dense.InternalError(sandbox, parseErr)
			}
			if otherId != record.Id {
				return liberror.NewWithValue(sandbox, api.KeyConflict, fieldName, value,
					sandbox.Deps.Std.Sprintf("value for key %q already exists", fieldName))
			}
		}

		// Step 3: write the new index entry before touching anything else.
		if err := dense.WriteInt(sandbox, dense.IndexKey(sandbox, record.Prefix, fieldName, newHash), record.Id); err != nil {
			return dense.InternalError(sandbox, err)
		}
		// Step 4: write the new value.
		if err := storage.Write(valueKey, []byte(encoded)); err != nil {
			return dense.InternalError(sandbox, err)
		}
		// Step 5: drop the index entry the old value owned.
		if hadValue {
			oldHash := dense.HashIndexValue(sandbox, string(oldRaw))
			if oldHash != newHash {
				if err := storage.Delete(dense.IndexKey(sandbox, record.Prefix, fieldName, oldHash)); err != nil {
					return dense.InternalError(sandbox, err)
				}
			}
		}
		return nil
	}
}

// RemoveFactory fills api.SchemaItem.Remove, deleting the record with the
// swap-with-last procedure that keeps the position list dense at a cost
// independent of the size of the collection. The last record moves into the
// freed position, so list order is not stable across removals.
func RemoveFactory(sandbox *api.Sandbox, record *api.SchemaItem, resolve dense.LinkResolver) func() *api.Error {
	return func() *api.Error {
		storage := sandbox.Deps.Storagedeps

		// Step 1: read the victim's position. Absent means already gone.
		position, live, err := readPosition(sandbox, record.Prefix, record.Id)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if !live {
			return nil
		}

		// Step 2: locate the record living at the last position.
		size, err := dense.ReadCount(sandbox, dense.SizeKey(sandbox, record.Prefix))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		lastRaw, found, err := storage.Read(dense.ListKey(sandbox, record.Prefix, size))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if !found {
			return liberror.New(sandbox, api.Internal, "",
				sandbox.Deps.Std.Sprintf("position %d of the list is missing", size))
		}
		lastId, err := dense.ParseId(sandbox, lastRaw)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}

		// Step 3: move the last record into the hole.
		if position != size {
			if err := storage.Write(dense.ListKey(sandbox, record.Prefix, position), lastRaw); err != nil {
				return dense.InternalError(sandbox, err)
			}
			if err := dense.WriteInt(sandbox, dense.PositionKey(sandbox, record.Prefix, lastId), position); err != nil {
				return dense.InternalError(sandbox, err)
			}
		}

		// Step 4: shrink the list.
		if err := storage.Delete(dense.ListKey(sandbox, record.Prefix, size)); err != nil {
			return dense.InternalError(sandbox, err)
		}
		if err := dense.WriteInt(sandbox, dense.SizeKey(sandbox, record.Prefix), size-1); err != nil {
			return dense.InternalError(sandbox, err)
		}

		// Step 5: drop the unique index entries the record owned.
		for _, item := range record.Items {
			if item.Type != api.Key {
				continue
			}
			raw, found, err := storage.Read(dense.ValueKey(sandbox, record.Prefix, record.Id, item.Name))
			if err != nil {
				return dense.InternalError(sandbox, err)
			}
			if !found {
				continue
			}
			indexKey := dense.IndexKey(sandbox, record.Prefix, item.Name, dense.HashIndexValue(sandbox, string(raw)))
			if err := storage.Delete(indexKey); err != nil {
				return dense.InternalError(sandbox, err)
			}
		}

		// Step 6: remove the record's own data, nested collections included.
		for _, item := range record.Items {
			if item.Type == api.Database {
				nested := dense.SubPrefix(sandbox, record.Prefix, record.Id, item.Name)
				if failure := ClearCollection(sandbox, item.Itens, nested, resolve); failure != nil {
					return failure
				}
				continue
			}
			if err := storage.Delete(dense.ValueKey(sandbox, record.Prefix, record.Id, item.Name)); err != nil {
				return dense.InternalError(sandbox, err)
			}
		}
		if err := storage.Delete(dense.PositionKey(sandbox, record.Prefix, record.Id)); err != nil {
			return dense.InternalError(sandbox, err)
		}
		return nil
	}
}

// CheckKeysPresenceFactory fills api.SchemaItem.CheckKeysPresence.
func CheckKeysPresenceFactory(sandbox *api.Sandbox, record *api.SchemaItem) func(keys []string) bool {
	return func(keys []string) bool {
		for _, key := range keys {
			exists, err := sandbox.Deps.Storagedeps.Exists(dense.ValueKey(sandbox, record.Prefix, record.Id, key))
			if err != nil || !exists {
				return false
			}
		}
		return true
	}
}

// ListAllFactory fills api.SchemaItem.ListAll, returning every record of a
// nested (Database) field.
func ListAllFactory(sandbox *api.Sandbox, record *api.SchemaItem, resolve dense.LinkResolver) func(fieldName string) []api.SchemaItem {
	return func(fieldName string) []api.SchemaItem {
		item, ok := dense.FindItem(sandbox, record.Items, fieldName)
		if !ok || item.Type != api.Database {
			return nil
		}
		nested := dense.SubPrefix(sandbox, record.Prefix, record.Id, fieldName)
		result, failure := ListRange(sandbox, item.Itens, nested, resolve, 1, 0)
		if failure != nil {
			return nil
		}
		return result
	}
}

// NewSubItemFactory fills api.SchemaItem.NewSubItem, inserting a record
// into a nested (Database) field of this record.
func NewSubItemFactory(sandbox *api.Sandbox, record *api.SchemaItem, resolve dense.LinkResolver) func(fieldName string, fields map[string]any) (api.SchemaItem, *api.Error) {
	return func(fieldName string, fields map[string]any) (api.SchemaItem, *api.Error) {
		item, ok := dense.FindItem(sandbox, record.Items, fieldName)
		if !ok || item.Type != api.Database {
			return api.SchemaItem{}, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.Std.Sprintf("field %q is not a nested collection of the schema", fieldName))
		}
		nested := dense.SubPrefix(sandbox, record.Prefix, record.Id, fieldName)
		return New(sandbox, item.Itens, nested, resolve, fields)
	}
}

// StringFactory fills api.SchemaItem.String, rendering the record's id and
// its plain fields so a caller printing a record sees its data.
func StringFactory(sandbox *api.Sandbox, record *api.SchemaItem) func() string {
	return func() string {
		parts := make([]string, 0, len(record.Items))
		for _, item := range record.Items {
			if item.Type == api.Database {
				continue
			}
			value, failure := record.Get(item.Name)
			if failure != nil {
				continue
			}
			parts = append(parts, sandbox.Deps.Std.Sprintf("%s: %v", item.Name, value))
		}
		joined := sandbox.Deps.Stringsdeps.Join(parts, ", ")
		return sandbox.Deps.Std.Sprintf("{id: %d, %s}", record.Id, joined)
	}
}

// readPosition reads the back-pointer that marks a record live. live is
// false when the record was never written or has been removed.
func readPosition(sandbox *api.Sandbox, prefix []string, id int64) (position int64, live bool, err error) {
	raw, found, err := sandbox.Deps.Storagedeps.Read(dense.PositionKey(sandbox, prefix, id))
	if err != nil || !found {
		return 0, false, err
	}
	parsed, err := dense.ParseId(sandbox, raw)
	if err != nil {
		return 0, false, err
	}
	return parsed, true, nil
}

// Build assembles an api.SchemaItem for an existing record id, running
// every field factory over it. It is the shared aggregate behind New,
// ResolveById, ResolveLive, ListRange and ClearCollection: adding a
// function field to api.SchemaItem means adding its factory call here.
func Build(sandbox *api.Sandbox, items []api.Item, prefix []string, resolve dense.LinkResolver, id int64) api.SchemaItem {
	record := api.SchemaItem{Items: items, Prefix: prefix, Id: id}
	record.Get = GetFactory(sandbox, &record)
	record.GetLink = GetLinkFactory(sandbox, &record, resolve)
	record.Update = UpdateFactory(sandbox, &record)
	record.Remove = RemoveFactory(sandbox, &record, resolve)
	record.CheckKeysPresence = CheckKeysPresenceFactory(sandbox, &record)
	record.ListAll = ListAllFactory(sandbox, &record, resolve)
	record.NewSubItem = NewSubItemFactory(sandbox, &record, resolve)
	record.String = StringFactory(sandbox, &record)
	return record
}

// New inserts a record into the collection identified by prefix, following
// the insertion procedure of the dense record pattern: validate, reserve an
// id, write the data, then publish by growing the list — the size key is
// the commit point, so a crash before it leaves an orphan nothing reads.
func New(sandbox *api.Sandbox, items []api.Item, prefix []string, resolve dense.LinkResolver, fields map[string]any) (api.SchemaItem, *api.Error) {
	storage := sandbox.Deps.Storagedeps

	// Every provided field has to be a plain field of the schema.
	for name := range fields {
		item, ok := dense.FindItem(sandbox, items, name)
		if !ok {
			return api.SchemaItem{}, liberror.New(sandbox, api.InvalidField, name,
				sandbox.Deps.Std.Sprintf("field %q is not part of the schema", name))
		}
		if item.Type == api.Database {
			return api.SchemaItem{}, liberror.New(sandbox, api.InvalidField, name,
				sandbox.Deps.Std.Sprintf("field %q is a nested collection and cannot be set directly", name))
		}
	}
	// Every required field has to be provided.
	for _, item := range items {
		if item.Required && item.Type != api.Database {
			if _, ok := fields[item.Name]; !ok {
				return api.SchemaItem{}, liberror.New(sandbox, api.MissingField, item.Name,
					sandbox.Deps.Std.Sprintf("required field %q is missing", item.Name))
			}
		}
	}

	encoded := make(map[string]string, len(fields))
	for name, value := range fields {
		item, _ := dense.FindItem(sandbox, items, name)
		encodedValue, failure := dense.EncodeValue(sandbox, item, value)
		if failure != nil {
			return api.SchemaItem{}, failure
		}
		encoded[name] = encodedValue
	}

	// Step 1: uniqueness check on every indexed (Key) field.
	for _, item := range items {
		if item.Type != api.Key {
			continue
		}
		encodedValue, ok := encoded[item.Name]
		if !ok {
			continue
		}
		indexKey := dense.IndexKey(sandbox, prefix, item.Name, dense.HashIndexValue(sandbox, encodedValue))
		exists, err := storage.Exists(indexKey)
		if err != nil {
			return api.SchemaItem{}, dense.InternalError(sandbox, err)
		}
		if exists {
			return api.SchemaItem{}, liberror.NewWithValue(sandbox, api.KeyConflict, item.Name, fields[item.Name],
				sandbox.Deps.Std.Sprintf("value for key %q already exists", item.Name))
		}
	}

	// Step 2: allocate the id. The counter only grows, so ids are never reused.
	lastId, err := dense.ReadCount(sandbox, dense.LastIdKey(sandbox, prefix))
	if err != nil {
		return api.SchemaItem{}, dense.InternalError(sandbox, err)
	}
	id := lastId + 1
	if err := dense.WriteInt(sandbox, dense.LastIdKey(sandbox, prefix), id); err != nil {
		return api.SchemaItem{}, dense.InternalError(sandbox, err)
	}

	// Step 3: the new position is size+1.
	size, err := dense.ReadCount(sandbox, dense.SizeKey(sandbox, prefix))
	if err != nil {
		return api.SchemaItem{}, dense.InternalError(sandbox, err)
	}
	position := size + 1

	// Step 4: write the record's data and its back-pointer.
	for name, encodedValue := range encoded {
		if err := storage.Write(dense.ValueKey(sandbox, prefix, id, name), []byte(encodedValue)); err != nil {
			return api.SchemaItem{}, dense.InternalError(sandbox, err)
		}
	}
	if err := dense.WriteInt(sandbox, dense.PositionKey(sandbox, prefix, id), position); err != nil {
		return api.SchemaItem{}, dense.InternalError(sandbox, err)
	}

	// Step 5: write the unique index entries.
	for _, item := range items {
		if item.Type != api.Key {
			continue
		}
		encodedValue, ok := encoded[item.Name]
		if !ok {
			continue
		}
		indexKey := dense.IndexKey(sandbox, prefix, item.Name, dense.HashIndexValue(sandbox, encodedValue))
		if err := dense.WriteInt(sandbox, indexKey, id); err != nil {
			return api.SchemaItem{}, dense.InternalError(sandbox, err)
		}
	}

	// Step 6: publish — list slot first, size last. Size is the commit point.
	if err := dense.WriteInt(sandbox, dense.ListKey(sandbox, prefix, position), id); err != nil {
		return api.SchemaItem{}, dense.InternalError(sandbox, err)
	}
	if err := dense.WriteInt(sandbox, dense.SizeKey(sandbox, prefix), position); err != nil {
		return api.SchemaItem{}, dense.InternalError(sandbox, err)
	}

	return Build(sandbox, items, prefix, resolve, id), nil
}

// ResolveById returns the record carrying the given id, but only while it
// is still live — that is, while its position back-pointer exists. ok is
// false for an id that was never allocated and for one whose record was
// removed; ids are never reused, so a stale id never resolves to a
// different record.
func ResolveById(sandbox *api.Sandbox, items []api.Item, prefix []string, resolve dense.LinkResolver, id int64) (api.SchemaItem, bool) {
	exists, err := sandbox.Deps.Storagedeps.Exists(dense.PositionKey(sandbox, prefix, id))
	if err != nil || !exists {
		return api.SchemaItem{}, false
	}
	return Build(sandbox, items, prefix, resolve, id), true
}

// ResolveLive parses an id read out of an index entry and returns the
// record only while it is still live. ok is false when it is not.
func ResolveLive(sandbox *api.Sandbox, items []api.Item, prefix []string, resolve dense.LinkResolver, rawId []byte) (api.SchemaItem, bool) {
	id, err := dense.ParseId(sandbox, rawId)
	if err != nil {
		return api.SchemaItem{}, false
	}
	return ResolveById(sandbox, items, prefix, resolve, id)
}

// ListRange reads records out of the dense position list, starting at from
// (counted from 1). A chunk of 0 means "to the end of the collection".
func ListRange(sandbox *api.Sandbox, items []api.Item, prefix []string, resolve dense.LinkResolver, from int64, chunk int64) ([]api.SchemaItem, *api.Error) {
	size, err := dense.ReadCount(sandbox, dense.SizeKey(sandbox, prefix))
	if err != nil {
		return nil, dense.InternalError(sandbox, err)
	}
	if from < 1 {
		from = 1
	}
	to := size
	if chunk > 0 && from+chunk-1 < size {
		to = from + chunk - 1
	}
	result := make([]api.SchemaItem, 0)
	for position := from; position <= to; position++ {
		id, err := dense.ReadCount(sandbox, dense.ListKey(sandbox, prefix, position))
		if err != nil {
			return nil, dense.InternalError(sandbox, err)
		}
		result = append(result, Build(sandbox, items, prefix, resolve, id))
	}
	return result, nil
}

// ClearCollection removes every record of a collection, and is what a
// removal runs over each nested collection of the record it deletes.
// Records go from the last position backwards, so no swap is ever needed.
func ClearCollection(sandbox *api.Sandbox, items []api.Item, prefix []string, resolve dense.LinkResolver) *api.Error {
	for {
		size, err := dense.ReadCount(sandbox, dense.SizeKey(sandbox, prefix))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if size == 0 {
			break
		}
		id, err := dense.ReadCount(sandbox, dense.ListKey(sandbox, prefix, size))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		record := Build(sandbox, items, prefix, resolve, id)
		if failure := record.Remove(); failure != nil {
			return failure
		}
	}
	if err := sandbox.Deps.Storagedeps.Delete(dense.SizeKey(sandbox, prefix)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	if err := sandbox.Deps.Storagedeps.Delete(dense.LastIdKey(sandbox, prefix)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	return nil
}
