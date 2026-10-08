package record

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
	writelock "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/writelock"
)

// One record of a collection: the factories filling the function fields of
// api.Record, and the record-level operations the collection package
// builds on. The removal procedure is in remove.go, the repair of a whole
// collection in repair.go.
//
// Every factory takes the sandbox and a pointer to the api.Record being
// built, and returns the closure for one field. Reading the record's Fields,
// Prefix and ID back off that pointer is what lets Build assign the fields
// in any order, and what keeps the record itself free of any dependency
// field a caller could reach. The pointer is to the record Build keeps for
// itself: the one it hands back carries copies of Fields and Prefix, so a
// caller editing them changes nothing the closures do.
//
// The one thing a record needs that its own Fields, Prefix and ID do not
// describe is the database around it — the collection a Link field points
// at, the write lock to take, the record owning it when nested — so a
// dense.Scope travels beside them, from the database that built it down
// through every nested collection. It is a parameter and never a field, for
// the same reason storage is: an api type carries no wiring a caller could
// read or replace.
//
// Every closure that writes takes the write lock of the record's top-level
// collection, and calls only functions of this package that take none — a
// public closure is never called while the lock is held, which is what keeps
// a non-reentrant lock from deadlocking on a nested removal.

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
				sandbox.Deps.StdDeps.Sprintf("field %q is a nested collection, use Nested(%q)", fieldName, fieldName))
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
func GetLinkFactory(sandbox *api.Sandbox, record *api.Record, scope dense.Scope) func(fieldName string) (api.Record, bool, *api.Error) {
	return func(fieldName string) (api.Record, bool, *api.Error) {
		field, ok := dense.FindField(sandbox, record.Fields, fieldName)
		if !ok || field.Type != api.Link {
			return api.Record{}, false, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.StdDeps.Sprintf("field %q is not a link field of the schema", fieldName))
		}
		targetFields, targetPrefix, ok := scope.Resolve(field.Target)
		if !ok {
			return api.Record{}, false, liberror.New(sandbox, api.InvalidField, fieldName,
				sandbox.Deps.StdDeps.Sprintf("link field %q targets %q, which is not a schema of the database", fieldName, field.Target))
		}
		raw, found, err := sandbox.Deps.StorageDeps.Read(dense.ValueKey(sandbox, record.Prefix, record.ID, fieldName))
		if err != nil {
			return api.Record{}, false, dense.InternalError(sandbox, err)
		}
		if !found {
			return api.Record{}, false, nil
		}
		id, err := dense.ParseID(sandbox, raw)
		if err != nil {
			return api.Record{}, false, dense.InternalError(sandbox, err)
		}
		targetScope := dense.TargetScope(sandbox, scope, targetPrefix)
		return ResolveByID(sandbox, targetFields, targetPrefix, targetScope, id)
	}
}

// UpdateFactory fills api.Record.Update. For an indexed (Key) field it
// runs the safe re-index sequence — write the new index entry, write the
// value, then delete the old index entry — so a crash part-way through
// never leaves the record unreachable. The worst it leaves is an index entry
// naming a record that no longer holds its value: such an entry is not valid
// (dense.IndexOwner), so a lookup through it finds nothing and the next
// writer claims it.
func UpdateFactory(sandbox *api.Sandbox, record *api.Record, scope dense.Scope) func(fieldName string, value any) *api.Error {
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
		encoded := ""
		if value == nil {
			if field.Required {
				return liberror.New(sandbox, api.MissingField, fieldName,
					sandbox.Deps.StdDeps.Sprintf("required field %q cannot be cleared", fieldName))
			}
		} else {
			var failure *api.Error
			encoded, failure = dense.EncodeFieldValue(sandbox, scope, field, value)
			if failure != nil {
				return failure
			}
		}

		release, failure := writelock.Acquire(sandbox, scope.Locks, scope.Root)
		if failure != nil {
			return failure
		}
		defer release()

		live, err := dense.Live(sandbox, record.Prefix, record.ID)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if !live {
			return removedError(sandbox, record.ID)
		}
		if value == nil {
			return clearValue(sandbox, record.Prefix, record.ID, field)
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

		// Step 2: refuse the new value when another live record owns it.
		newHash := dense.HashIndexValue(sandbox, encoded)
		owner, valid, err := dense.IndexOwner(sandbox, record.Prefix, fieldName, newHash)
		if err != nil {
			return dense.InternalError(sandbox, err)
		}
		if valid && owner != record.ID {
			return liberror.NewWithValue(sandbox, api.KeyConflict, fieldName, value,
				sandbox.Deps.StdDeps.Sprintf("value for key %q already exists", fieldName))
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
				if failure := deleteOwnEntry(sandbox, record.Prefix, fieldName, oldHash, record.ID); failure != nil {
					return failure
				}
			}
		}
		return nil
	}
}

// RemoveFactory fills api.Record.Remove, deleting the record with the
// swap-with-last procedure of remove.go under the collection's write lock.
func RemoveFactory(sandbox *api.Sandbox, record *api.Record, scope dense.Scope) func() *api.Error {
	return func() *api.Error {
		release, failure := writelock.Acquire(sandbox, scope.Locks, scope.Root)
		if failure != nil {
			return failure
		}
		defer release()
		return remove(sandbox, record.Fields, record.Prefix, scope, record.ID)
	}
}

// HasValuesFactory fills api.Record.HasValues, checking each value key
// exists without reading it.
func HasValuesFactory(sandbox *api.Sandbox, record *api.Record) func(fields []string) (bool, *api.Error) {
	return func(fields []string) (bool, *api.Error) {
		for _, name := range fields {
			field, ok := dense.FindField(sandbox, record.Fields, name)
			if !ok || field.Type == api.Nested {
				return false, liberror.New(sandbox, api.InvalidField, name,
					sandbox.Deps.StdDeps.Sprintf("field %q is not a plain field of the schema", name))
			}
		}
		for _, name := range fields {
			exists, err := sandbox.Deps.StorageDeps.Exists(dense.ValueKey(sandbox, record.Prefix, record.ID, name))
			if err != nil {
				return false, dense.InternalError(sandbox, err)
			}
			if !exists {
				return false, nil
			}
		}
		return true, nil
	}
}

// NestedFactory fills api.Record.Nested, handing out a Nested field of this
// record as a Collection built exactly like a top-level one, with this
// record as its owner.
func NestedFactory(sandbox *api.Sandbox, record *api.Record, scope dense.Scope) func(fieldName string) (api.Collection, *api.Error) {
	return func(fieldName string) (api.Collection, *api.Error) {
		return nestedCollection(sandbox, record, scope, fieldName)
	}
}

// ListNestedFactory fills api.Record.ListNested, returning every record of
// a Nested field.
func ListNestedFactory(sandbox *api.Sandbox, record *api.Record, scope dense.Scope) func(fieldName string) ([]api.Record, *api.Error) {
	return func(fieldName string) ([]api.Record, *api.Error) {
		nested, failure := nestedCollection(sandbox, record, scope, fieldName)
		if failure != nil {
			return nil, failure
		}
		return nested.ListAll()
	}
}

// InsertNestedFactory fills api.Record.InsertNested, inserting a record
// into a Nested field of this record.
func InsertNestedFactory(sandbox *api.Sandbox, record *api.Record, scope dense.Scope) func(fieldName string, fields map[string]any) (api.Record, *api.Error) {
	return func(fieldName string, fields map[string]any) (api.Record, *api.Error) {
		nested, failure := nestedCollection(sandbox, record, scope, fieldName)
		if failure != nil {
			return api.Record{}, failure
		}
		return nested.Insert(fields)
	}
}

// StringFactory fills api.Record.String, rendering the record's id and
// its plain fields so a caller printing a record sees its data. Text is
// quoted, so no value can pass for a separator.
func StringFactory(sandbox *api.Sandbox, record *api.Record) func() string {
	return func() string {
		parts := make([]string, 0, len(record.Fields)+1)
		parts = append(parts, sandbox.Deps.StdDeps.Sprintf("id: %d", record.ID))
		for _, field := range record.Fields {
			if field.Type == api.Nested {
				continue
			}
			value, failure := record.Get(field.Name)
			if failure != nil {
				continue
			}
			switch typed := value.(type) {
			case []byte:
				// A Bytes value prints as its length: its contents are
				// binary, and may be as large as a file.
				parts = append(parts, sandbox.Deps.StdDeps.Sprintf("%s: %d bytes", field.Name, len(typed)))
			case string:
				parts = append(parts, sandbox.Deps.StdDeps.Sprintf("%s: %q", field.Name, typed))
			default:
				parts = append(parts, sandbox.Deps.StdDeps.Sprintf("%s: %v", field.Name, typed))
			}
		}
		return "{" + sandbox.Deps.StringsDeps.Join(parts, ", ") + "}"
	}
}

// nestedCollection builds the collection a Nested field of record holds.
func nestedCollection(sandbox *api.Sandbox, record *api.Record, scope dense.Scope, fieldName string) (api.Collection, *api.Error) {
	field, ok := dense.FindField(sandbox, record.Fields, fieldName)
	if !ok || field.Type != api.Nested {
		return api.Collection{}, liberror.New(sandbox, api.InvalidField, fieldName,
			sandbox.Deps.StdDeps.Sprintf("field %q is not a nested collection of the schema", fieldName))
	}
	prefix := dense.SubPrefix(sandbox, record.Prefix, record.ID, fieldName)
	nestedScope := dense.NestedScope(sandbox, scope, record.Prefix, record.ID)
	return scope.NewCollection(sandbox, field.Fields, prefix, nestedScope), nil
}

// clearValue deletes one value of a live record — the value first, then the
// index entry it owned, so a crash between the two leaves an entry naming a
// record with no such value: one that is not valid, and so free.
func clearValue(sandbox *api.Sandbox, prefix []string, id int64, field api.Field) *api.Error {
	storage := sandbox.Deps.StorageDeps
	valueKey := dense.ValueKey(sandbox, prefix, id, field.Name)
	raw, hadValue, err := storage.Read(valueKey)
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	if !hadValue {
		return nil
	}
	if err := storage.Delete(valueKey); err != nil {
		return dense.InternalError(sandbox, err)
	}
	if field.Type == api.Key {
		return deleteOwnEntry(sandbox, prefix, field.Name, dense.HashIndexValue(sandbox, string(raw)), id)
	}
	return nil
}

// deleteOwnEntry deletes an index entry only while it still names id. An
// entry a record no longer holds may since have been claimed by another
// record, and deleting it would make that record unreachable by key.
func deleteOwnEntry(sandbox *api.Sandbox, prefix []string, field string, hash string, id int64) *api.Error {
	entryID, found, err := dense.ReadIndex(sandbox, prefix, field, hash)
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	if !found || entryID != id {
		return nil
	}
	if err := sandbox.Deps.StorageDeps.Delete(dense.IndexKey(sandbox, prefix, field, hash)); err != nil {
		return dense.InternalError(sandbox, err)
	}
	return nil
}

// removedError is the failure of a write to a record that is not live.
func removedError(sandbox *api.Sandbox, id int64) *api.Error {
	return liberror.New(sandbox, api.Removed, "",
		sandbox.Deps.StdDeps.Sprintf("record %d is no longer live", id))
}

// CheckOwner refuses a write into a nested collection whose owning record is
// no longer live, so nothing is written under a record that has been
// removed. A top-level collection has no owner and always passes. The caller
// holds the write lock.
func CheckOwner(sandbox *api.Sandbox, scope dense.Scope) *api.Error {
	if scope.Owner == nil {
		return nil
	}
	live, err := dense.Live(sandbox, scope.Owner, scope.OwnerID)
	if err != nil {
		return dense.InternalError(sandbox, err)
	}
	if !live {
		return removedError(sandbox, scope.OwnerID)
	}
	return nil
}

// Build assembles an api.Record for an existing record id, running
// every field factory over it. It is the shared aggregate behind Insert,
// ResolveByID, ListRange and the collection lookups: adding a function field
// to api.Record means adding its factory call here.
func Build(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope, id int64) api.Record {
	record := api.Record{Fields: fields, Prefix: prefix, ID: id}
	record.Get = GetFactory(sandbox, &record)
	record.GetLink = GetLinkFactory(sandbox, &record, scope)
	record.Update = UpdateFactory(sandbox, &record, scope)
	record.Remove = RemoveFactory(sandbox, &record, scope)
	record.HasValues = HasValuesFactory(sandbox, &record)
	record.Nested = NestedFactory(sandbox, &record, scope)
	record.ListNested = ListNestedFactory(sandbox, &record, scope)
	record.InsertNested = InsertNestedFactory(sandbox, &record, scope)
	record.String = StringFactory(sandbox, &record)

	public := record
	public.Fields = dense.CopyFields(sandbox, fields)
	public.Prefix = dense.Key(sandbox, prefix)
	return public
}

// Insert writes a record into the collection identified by prefix,
// following the insertion procedure of the dense record pattern: validate,
// lock, reserve an id, write the data, then publish by growing the list — the
// size key is the commit point, so a crash before it leaves an orphan
// nothing reads.
func Insert(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope, values map[string]any) (api.Record, *api.Error) {
	// Step 0: validate and encode, with no I/O at all.
	encoded, failure := encodeInsert(sandbox, fields, scope, values)
	if failure != nil {
		return api.Record{}, failure
	}

	release, failure := writelock.Acquire(sandbox, scope.Locks, scope.Root)
	if failure != nil {
		return api.Record{}, failure
	}
	defer release()
	if failure := CheckOwner(sandbox, scope); failure != nil {
		return api.Record{}, failure
	}

	// Step 1: uniqueness check on every indexed (Key) field. An entry that
	// is not valid is free.
	for _, field := range fields {
		encodedValue, ok := encoded[field.Name]
		if field.Type != api.Key || !ok {
			continue
		}
		hash := dense.HashIndexValue(sandbox, encodedValue)
		_, valid, err := dense.IndexOwner(sandbox, prefix, field.Name, hash)
		if err != nil {
			return api.Record{}, dense.InternalError(sandbox, err)
		}
		if valid {
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

	// From here on the id is spent. A step that fails abandons the record:
	// what was written of it is deleted again, as far as the backend lets
	// it, so a backend failing for a moment leaves no debris for Repair.
	if err := publish(sandbox, fields, prefix, encoded, id); err != nil {
		_ = cleanup(sandbox, fields, prefix, scope, id)
		return api.Record{}, dense.InternalError(sandbox, err)
	}

	return Build(sandbox, fields, prefix, scope, id), nil
}

// publish runs steps 3 to 6 of an insert for the id step 2 allocated: the
// record's data and back-pointer, its index entries, then the list slot and
// the size — the commit point, written last.
func publish(sandbox *api.Sandbox, fields []api.Field, prefix []string, encoded map[string]string, id int64) error {
	storage := sandbox.Deps.StorageDeps

	// Step 3: the new position is size+1.
	size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
	if err != nil {
		return err
	}
	position := size + 1

	// Step 4: write the record's data and its back-pointer.
	for _, field := range fields {
		encodedValue, ok := encoded[field.Name]
		if !ok {
			continue
		}
		if err := storage.Write(dense.ValueKey(sandbox, prefix, id, field.Name), []byte(encodedValue)); err != nil {
			return err
		}
	}
	if err := dense.WriteInt(sandbox, dense.PositionKey(sandbox, prefix, id), position); err != nil {
		return err
	}

	// Step 5: write the unique index entries.
	for _, field := range fields {
		encodedValue, ok := encoded[field.Name]
		if field.Type != api.Key || !ok {
			continue
		}
		indexKey := dense.IndexKey(sandbox, prefix, field.Name, dense.HashIndexValue(sandbox, encodedValue))
		if err := dense.WriteInt(sandbox, indexKey, id); err != nil {
			return err
		}
	}

	// Step 6: publish — list slot first, size last. Size is the commit point.
	if err := dense.WriteInt(sandbox, dense.ListKey(sandbox, prefix, position), id); err != nil {
		return err
	}
	return dense.WriteInt(sandbox, dense.SizeKey(sandbox, prefix), position)
}

// encodeInsert validates the fields map of an insert against the schema and
// encodes every value, reporting the first failure in schema order — an
// unknown name before anything else, the smallest one when there are
// several — so the same bad insert always reports the same field. A nil
// value is the same as a value left out.
func encodeInsert(sandbox *api.Sandbox, fields []api.Field, scope dense.Scope, values map[string]any) (map[string]string, *api.Error) {
	unknown, hasUnknown := "", false
	for name := range values {
		if _, ok := dense.FindField(sandbox, fields, name); ok {
			continue
		}
		if !hasUnknown || name < unknown {
			unknown, hasUnknown = name, true
		}
	}
	if hasUnknown {
		return nil, liberror.New(sandbox, api.InvalidField, unknown,
			sandbox.Deps.StdDeps.Sprintf("field %q is not part of the schema", unknown))
	}

	encoded := make(map[string]string, len(values))
	for _, field := range fields {
		value, provided := values[field.Name]
		if field.Type == api.Nested {
			if provided {
				return nil, liberror.New(sandbox, api.InvalidField, field.Name,
					sandbox.Deps.StdDeps.Sprintf("field %q is a nested collection and cannot be set directly", field.Name))
			}
			continue
		}
		if !provided || value == nil {
			if field.Required {
				return nil, liberror.New(sandbox, api.MissingField, field.Name,
					sandbox.Deps.StdDeps.Sprintf("required field %q is missing", field.Name))
			}
			continue
		}
		encodedValue, failure := dense.EncodeFieldValue(sandbox, scope, field, value)
		if failure != nil {
			return nil, failure
		}
		encoded[field.Name] = encodedValue
	}
	return encoded, nil
}

// ResolveByID returns the record carrying the given id, but only while it
// is live (dense.Live). ok is false for an id that was never allocated, for
// one whose insert never committed and for one whose record was removed; ids
// are never reused, so a stale id never resolves to a different record.
func ResolveByID(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope, id int64) (api.Record, bool, *api.Error) {
	live, err := dense.Live(sandbox, prefix, id)
	if err != nil {
		return api.Record{}, false, dense.InternalError(sandbox, err)
	}
	if !live {
		return api.Record{}, false, nil
	}
	return Build(sandbox, fields, prefix, scope, id), true, nil
}

// ListRange reads records out of the dense position list, starting at from
// (counted from 1). A chunk of 0 means "to the end of the collection". Only
// live slots are read back (dense.SlotLive), and no id twice: a slot a
// failed removal left behind is skipped rather than handed out as a record,
// and so is one a removal running beside this read has just emptied.
func ListRange(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope, from int64, chunk int64) ([]api.Record, *api.Error) {
	if from < 1 {
		return nil, liberror.NewWithValue(sandbox, api.InvalidArgument, "", from,
			sandbox.Deps.StdDeps.Sprintf("position %d is below 1: positions are counted from 1", from))
	}
	if chunk < 0 {
		return nil, liberror.NewWithValue(sandbox, api.InvalidArgument, "", chunk,
			sandbox.Deps.StdDeps.Sprintf("chunk %d is negative: 0 means to the end", chunk))
	}
	size, err := dense.ReadInt(sandbox, dense.SizeKey(sandbox, prefix))
	if err != nil {
		return nil, dense.InternalError(sandbox, err)
	}
	to := size
	// chunk <= size-from is from+chunk-1 < size written so it cannot
	// overflow, whatever chunk a caller passes.
	if chunk > 0 && chunk <= size-from {
		to = from + chunk - 1
	}
	result := make([]api.Record, 0)
	seen := map[int64]bool{}
	for position := from; position <= to; position++ {
		id, live, err := dense.SlotLive(sandbox, prefix, position)
		if err != nil {
			return nil, dense.InternalError(sandbox, err)
		}
		if !live || seen[id] {
			continue
		}
		seen[id] = true
		result = append(result, Build(sandbox, fields, prefix, scope, id))
	}
	return result, nil
}
