package record

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
)

// One record of a collection: the factories filling the function fields of
// api.Record, and the record-level operations the collection package
// builds on.
//
// Every factory takes the sandbox and a pointer to the api.Record being
// built, and returns the closure for one field. Reading the record's Fields,
// Prefix and ID back off that pointer is what lets Build assign the fields
// in any order, and what keeps the record itself free of any dependency
// field a caller could reach.
//
// The one thing a record needs that its own Fields, Prefix and ID do not
// describe is the collection a Link field points at, so a dense.LinkResolver
// travels beside them, from the database that built it down through every
// nested collection. It is a parameter and never a field, for the same
// reason storage is: an api type carries no wiring a caller could read or
// replace.

// GetFactory fills api.Record.Get.
func GetFactory(sandbox *api.Sandbox, record *api.Record) func(fieldName string) (any, *api.Error) {
	return func(fieldName string) (any, *api.Error) {
		field, ok := dense.FindField(sandbox, record.Fields, fieldName)
		if !ok {
			return nil, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.StdDeps.Sprintf("field %q is not part of the schema", fieldName))
		}
		if field.Type == api.Nested {
			return nil, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.StdDeps.Sprintf("field %q is a nested collection, use ListNested(%q)", fieldName, fieldName))
		}
		raw, found, err := sandbox.Deps.StorageDeps.Read(dense.ValueKey(sandbox, record.Prefix, record.ID, fieldName))
		if err != nil {
			return nil, dense.InternalError(sandbox, err)
		}
		if !found {
			return nil, liberror.New(sandbox, api.NoValue, fieldName,
				sandbox.Deps.StdDeps.Sprintf("field %q has no value for this record", fieldName))
		}
		return dense.DecodeValue(sandbox, field, raw)
	}
}

// GetLinkFactory fills api.Record.GetLink, following a Link field to the
// record it names in the collection its Target declares. Resolution goes
// through ResolveByID, so a link to a record that has been removed reports
// ok == false — and since ids are never reused, it never reports a
// different record instead.
func GetLinkFactory(sandbox *api.Sandbox, record *api.Record, resolve dense.LinkResolver) func(fieldName string) (api.Record, bool) {
	return func(fieldName string) (api.Record, bool) {
		field, ok := dense.FindField(sandbox, record.Fields, fieldName)
		if !ok || field.Type != api.Link {
			return api.Record{}, false
		}
		targetFields, targetPrefix, ok := resolve(field.Target)
		if !ok {
			return api.Record{}, false
		}
		raw, found, err := sandbox.Deps.StorageDeps.Read(dense.ValueKey(sandbox, record.Prefix, record.ID, fieldName))
		if err != nil || !found {
			return api.Record{}, false
		}
		return ResolveRawID(sandbox, targetFields, targetPrefix, resolve, raw)
	}
}

// UpdateFactory fills api.Record.Update. For an indexed (Key) field it
// runs the safe re-index sequence — write the new index entry, write the
// value, then delete the old index entry — so a crash part-way through
// never leaves the record unreachable: the worst outcome is an index entry
// pointing at a record that no longer holds the value, which resolves to
// nothing and is overwritten by the next write.
func UpdateFactory(sandbox *api.Sandbox, record *api.Record) func(fieldName string, value any) *api.Error {
	return func(fieldName string, value any) *api.Error {
		field, ok := dense.FindField(sandbox, record.Fields, fieldName)
		if !ok {
			return liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.StdDeps.Sprintf("field %q is not part of the schema", fieldName))
		}
		if field.Type == api.Nested {
			return liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.StdDeps.Sprintf("field %q is a nested collection and cannot be updated directly", fieldName))
		}
		encoded, failure := dense.EncodeValue(sandbox, field, value)
		if failure != nil {
			return failure
		}
		storage := sandbox.Deps.StorageDeps
		valueKey := dense.ValueKey(sandbox, record.Prefix, record.ID, fieldName)

		if field.Type != api.Key {
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
			otherID, parseErr := dense.ParseID(sandbox, existing)
			if parseErr != nil {
				return dense.InternalError(sandbox, parseErr)
			}
			if otherID != record.ID {
				return liberror.NewWithValue(sandbox, api.KeyConflict, fieldName, value,
					sandbox.Deps.StdDeps.Sprintf("value for key %q already exists", fieldName))
			}
		}

		// Step 3: write the new index entry before touching anything else.
		if err := dense.WriteInt(sandbox, dense.IndexKey(sandbox, record.Prefix, fieldName, newHash), record.ID); err != nil {
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

// RemoveFactory fills api.Record.Remove, deleting the record with the
// swap-with-last procedure that keeps the position list dense at a cost
// independent of the size of the collection. The last record moves into the
// freed position, so list order is not stable across removals.
func RemoveFactory(sandbox *api.Sandbox, record *api.Record, resolve dense.LinkResolver) func() *api.Error {
	return func() *api.Error {
		storage := sandbox.Deps.StorageDeps

		// Step 1: read the victim's position. Absent means already gone.
		position, live, err := readPosition(sandbox, record.Prefix, record.ID)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if !live {
			return nil
		}

		// Step 2: locate the record living at the last position.
		size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, record.Prefix))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		lastRaw, found, err := storage.Read(dense.ListKey(sandbox, record.Prefix, size))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if !found {
			return liberror.New(sandbox, api.Internal, "",
				sandbox.Deps.StdDeps.Sprintf("position %d of the list is missing", size))
		}
		lastID, err := dense.ParseID(sandbox, lastRaw)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}

		// Step 3: move the last record into the hole.
		if position != size {
			if err := storage.Write(dense.ListKey(sandbox, record.Prefix, position), lastRaw); err != nil {
				return dense.InternalError(sandbox, err)
			}
			if err := dense.WriteInt(sandbox, dense.PositionKey(sandbox, record.Prefix, lastID), position); err != nil {
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
		for _, field := range record.Fields {
			if field.Type != api.Key {
				continue
			}
			raw, found, err := storage.Read(dense.ValueKey(sandbox, record.Prefix, record.ID, field.Name))
			if err != nil {
				return dense.InternalError(sandbox, err)
			}
			if !found {
				continue
			}
			indexKey := dense.IndexKey(sandbox, record.Prefix, field.Name, dense.HashIndexValue(sandbox, string(raw)))
			if err := storage.Delete(indexKey); err != nil {
				return dense.InternalError(sandbox, err)
			}
		}

		// Step 6: remove the record's own data, nested collections included.
		for _, field := range record.Fields {
			if field.Type == api.Nested {
				nested := dense.SubPrefix(sandbox, record.Prefix, record.ID, field.Name)
				if failure := ClearCollection(sandbox, field.Fields, nested, resolve); failure != nil {
					return failure
				}
				continue
			}
			if err := storage.Delete(dense.ValueKey(sandbox, record.Prefix, record.ID, field.Name)); err != nil {
				return dense.InternalError(sandbox, err)
			}
		}
		if err := storage.Delete(dense.PositionKey(sandbox, record.Prefix, record.ID)); err != nil {
			return dense.InternalError(sandbox, err)
		}
		return nil
	}
}

// HasValuesFactory fills api.Record.HasValues.
func HasValuesFactory(sandbox *api.Sandbox, record *api.Record) func(fields []string) bool {
	return func(fields []string) bool {
		for _, field := range fields {
			exists, err := sandbox.Deps.StorageDeps.Exists(dense.ValueKey(sandbox, record.Prefix, record.ID, field))
			if err != nil || !exists {
				return false
			}
		}
		return true
	}
}

// ListNestedFactory fills api.Record.ListNested, returning every record of
// a Nested field.
func ListNestedFactory(sandbox *api.Sandbox, record *api.Record, resolve dense.LinkResolver) func(fieldName string) []api.Record {
	return func(fieldName string) []api.Record {
		field, ok := dense.FindField(sandbox, record.Fields, fieldName)
		if !ok || field.Type != api.Nested {
			return nil
		}
		nested := dense.SubPrefix(sandbox, record.Prefix, record.ID, fieldName)
		result, failure := ListRange(sandbox, field.Fields, nested, resolve, 1, 0)
		if failure != nil {
			return nil
		}
		return result
	}
}

// InsertNestedFactory fills api.Record.InsertNested, inserting a record
// into a Nested field of this record.
func InsertNestedFactory(sandbox *api.Sandbox, record *api.Record, resolve dense.LinkResolver) func(fieldName string, fields map[string]any) (api.Record, *api.Error) {
	return func(fieldName string, fields map[string]any) (api.Record, *api.Error) {
		field, ok := dense.FindField(sandbox, record.Fields, fieldName)
		if !ok || field.Type != api.Nested {
			return api.Record{}, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.StdDeps.Sprintf("field %q is not a nested collection of the schema", fieldName))
		}
		nested := dense.SubPrefix(sandbox, record.Prefix, record.ID, fieldName)
		return Insert(sandbox, field.Fields, nested, resolve, fields)
	}
}

// StringFactory fills api.Record.String, rendering the record's id and
// its plain fields so a caller printing a record sees its data.
func StringFactory(sandbox *api.Sandbox, record *api.Record) func() string {
	return func() string {
		parts := make([]string, 0, len(record.Fields))
		for _, field := range record.Fields {
			if field.Type == api.Nested {
				continue
			}
			value, failure := record.Get(field.Name)
			if failure != nil {
				continue
			}
			// A Bytes value prints as its length: its contents are binary,
			// and may be as large as a file.
			if raw, ok := value.([]byte); ok {
				parts = append(parts, sandbox.Deps.StdDeps.Sprintf("%s: %d bytes", field.Name, len(raw)))
				continue
			}
			parts = append(parts, sandbox.Deps.StdDeps.Sprintf("%s: %v", field.Name, value))
		}
		joined := sandbox.Deps.StringsDeps.Join(parts, ", ")
		return sandbox.Deps.StdDeps.Sprintf("{id: %d, %s}", record.ID, joined)
	}
}

// readPosition reads the back-pointer that marks a record live. live is
// false when the record was never written or has been removed.
func readPosition(sandbox *api.Sandbox, prefix []string, id int64) (position int64, live bool, err error) {
	raw, found, err := sandbox.Deps.StorageDeps.Read(dense.PositionKey(sandbox, prefix, id))
	if err != nil || !found {
		return 0, false, err
	}
	parsed, err := dense.ParseID(sandbox, raw)
	if err != nil {
		return 0, false, err
	}
	return parsed, true, nil
}

// Build assembles an api.Record for an existing record id, running
// every field factory over it. It is the shared aggregate behind Insert,
// ResolveByID, ResolveRawID, ListRange and ClearCollection: adding a
// function field to api.Record means adding its factory call here.
func Build(sandbox *api.Sandbox, fields []api.Field, prefix []string, resolve dense.LinkResolver, id int64) api.Record {
	record := api.Record{Fields: fields, Prefix: prefix, ID: id}
	record.Get = GetFactory(sandbox, &record)
	record.GetLink = GetLinkFactory(sandbox, &record, resolve)
	record.Update = UpdateFactory(sandbox, &record)
	record.Remove = RemoveFactory(sandbox, &record, resolve)
	record.HasValues = HasValuesFactory(sandbox, &record)
	record.ListNested = ListNestedFactory(sandbox, &record, resolve)
	record.InsertNested = InsertNestedFactory(sandbox, &record, resolve)
	record.String = StringFactory(sandbox, &record)
	return record
}

// Insert writes a record into the collection identified by prefix,
// following the insertion procedure of the dense record pattern: validate,
// reserve an id, write the data, then publish by growing the list — the
// size key is the commit point, so a crash before it leaves an orphan
// nothing reads.
func Insert(sandbox *api.Sandbox, fields []api.Field, prefix []string, resolve dense.LinkResolver, values map[string]any) (api.Record, *api.Error) {
	storage := sandbox.Deps.StorageDeps

	// Every provided field has to be a plain field of the schema.
	for name := range values {
		field, ok := dense.FindField(sandbox, fields, name)
		if !ok {
			return api.Record{}, liberror.New(sandbox, api.InvalidField, name,
				sandbox.Deps.StdDeps.Sprintf("field %q is not part of the schema", name))
		}
		if field.Type == api.Nested {
			return api.Record{}, liberror.New(sandbox, api.InvalidField, name,
				sandbox.Deps.StdDeps.Sprintf("field %q is a nested collection and cannot be set directly", name))
		}
	}
	// Every required field has to be provided.
	for _, field := range fields {
		if field.Required && field.Type != api.Nested {
			if _, ok := values[field.Name]; !ok {
				return api.Record{}, liberror.New(sandbox, api.MissingField, field.Name,
					sandbox.Deps.StdDeps.Sprintf("required field %q is missing", field.Name))
			}
		}
	}

	encoded := make(map[string]string, len(values))
	for name, value := range values {
		field, _ := dense.FindField(sandbox, fields, name)
		encodedValue, failure := dense.EncodeValue(sandbox, field, value)
		if failure != nil {
			return api.Record{}, failure
		}
		encoded[name] = encodedValue
	}

	// Step 1: uniqueness check on every indexed (Key) field.
	for _, field := range fields {
		if field.Type != api.Key {
			continue
		}
		encodedValue, ok := encoded[field.Name]
		if !ok {
			continue
		}
		indexKey := dense.IndexKey(sandbox, prefix, field.Name, dense.HashIndexValue(sandbox, encodedValue))
		exists, err := storage.Exists(indexKey)
		if err != nil {
			return api.Record{}, dense.InternalError(sandbox, err)
		}
		if exists {
			return api.Record{}, liberror.NewWithValue(sandbox, api.KeyConflict, field.Name, values[field.Name],
				sandbox.Deps.StdDeps.Sprintf("value for key %q already exists", field.Name))
		}
	}

	// Step 2: allocate the id. The counter only grows, so ids are never reused.
	lastID, err := dense.ReadInt(sandbox, dense.LastIDKey(sandbox, prefix))
	if err != nil {
		return api.Record{}, dense.InternalError(sandbox, err)
	}
	id := lastID + 1
	if err := dense.WriteInt(sandbox, dense.LastIDKey(sandbox, prefix), id); err != nil {
		return api.Record{}, dense.InternalError(sandbox, err)
	}

	// Step 3: the new position is size+1.
	size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
	if err != nil {
		return api.Record{}, dense.InternalError(sandbox, err)
	}
	position := size + 1

	// Step 4: write the record's data and its back-pointer.
	for name, encodedValue := range encoded {
		if err := storage.Write(dense.ValueKey(sandbox, prefix, id, name), []byte(encodedValue)); err != nil {
			return api.Record{}, dense.InternalError(sandbox, err)
		}
	}
	if err := dense.WriteInt(sandbox, dense.PositionKey(sandbox, prefix, id), position); err != nil {
		return api.Record{}, dense.InternalError(sandbox, err)
	}

	// Step 5: write the unique index entries.
	for _, field := range fields {
		if field.Type != api.Key {
			continue
		}
		encodedValue, ok := encoded[field.Name]
		if !ok {
			continue
		}
		indexKey := dense.IndexKey(sandbox, prefix, field.Name, dense.HashIndexValue(sandbox, encodedValue))
		if err := dense.WriteInt(sandbox, indexKey, id); err != nil {
			return api.Record{}, dense.InternalError(sandbox, err)
		}
	}

	// Step 6: publish — list slot first, size last. Size is the commit point.
	if err := dense.WriteInt(sandbox, dense.ListKey(sandbox, prefix, position), id); err != nil {
		return api.Record{}, dense.InternalError(sandbox, err)
	}
	if err := dense.WriteInt(sandbox, dense.SizeKey(sandbox, prefix), position); err != nil {
		return api.Record{}, dense.InternalError(sandbox, err)
	}

	return Build(sandbox, fields, prefix, resolve, id), nil
}

// ResolveByID returns the record carrying the given id, but only while it
// is still live — that is, while its position back-pointer exists. ok is
// false for an id that was never allocated and for one whose record was
// removed; ids are never reused, so a stale id never resolves to a
// different record.
func ResolveByID(sandbox *api.Sandbox, fields []api.Field, prefix []string, resolve dense.LinkResolver, id int64) (api.Record, bool) {
	exists, err := sandbox.Deps.StorageDeps.Exists(dense.PositionKey(sandbox, prefix, id))
	if err != nil || !exists {
		return api.Record{}, false
	}
	return Build(sandbox, fields, prefix, resolve, id), true
}

// ResolveRawID parses an id read out of an index entry or a Link field and
// returns the record only while it is still live. ok is false when it is
// not.
func ResolveRawID(sandbox *api.Sandbox, fields []api.Field, prefix []string, resolve dense.LinkResolver, rawID []byte) (api.Record, bool) {
	id, err := dense.ParseID(sandbox, rawID)
	if err != nil {
		return api.Record{}, false
	}
	return ResolveByID(sandbox, fields, prefix, resolve, id)
}

// ListRange reads records out of the dense position list, starting at from
// (counted from 1). A chunk of 0 means "to the end of the collection".
func ListRange(sandbox *api.Sandbox, fields []api.Field, prefix []string, resolve dense.LinkResolver, from int64, chunk int64) ([]api.Record, *api.Error) {
	size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
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
	result := make([]api.Record, 0)
	for position := from; position <= to; position++ {
		id, err := dense.ReadInt(sandbox, dense.ListKey(sandbox, prefix, position))
		if err != nil {
			return nil, dense.InternalError(sandbox, err)
		}
		result = append(result, Build(sandbox, fields, prefix, resolve, id))
	}
	return result, nil
}

// ClearCollection removes every record of a collection, and is what a
// removal runs over each nested collection of the record it deletes.
// Records go from the last position backwards, so no swap is ever needed.
func ClearCollection(sandbox *api.Sandbox, fields []api.Field, prefix []string, resolve dense.LinkResolver) *api.Error {
	for {
		size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if size == 0 {
			break
		}
		id, err := dense.ReadInt(sandbox, dense.ListKey(sandbox, prefix, size))
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		record := Build(sandbox, fields, prefix, resolve, id)
		if failure := record.Remove(); failure != nil {
			return failure
		}
	}
	if err := sandbox.Deps.StorageDeps.Delete(dense.SizeKey(sandbox, prefix)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	if err := sandbox.Deps.StorageDeps.Delete(dense.LastIDKey(sandbox, prefix)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	return nil
}
