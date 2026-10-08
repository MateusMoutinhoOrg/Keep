package collection

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	liberror "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/liberror"
	record "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/record"
	writelock "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/writelock"
)

// One collection of records: the factories filling the function fields of
// api.Collection. A collection is nothing but a key prefix and the list
// of fields its records hold, so the same New builds a top-level collection
// and a nested one. The dense.Scope it is built with is the database the
// collection belongs to, the write lock it takes and, when nested, the record
// owning it; every factory passes it on unchanged.

// InsertFactory fills api.Collection.Insert.
func InsertFactory(sandbox *api.Sandbox, collection *api.Collection, scope dense.Scope) func(fields map[string]any) (api.Record, *api.Error) {
	return func(fields map[string]any) (api.Record, *api.Error) {
		return record.Insert(sandbox, collection.Fields, collection.Prefix, scope, fields)
	}
}

// FindByKeyFactory fills api.Collection.FindByKey: encode and hash the
// value, then read the index entry and check it is valid — it names a live
// record whose current value hashes back to it. An entry a failed write left
// behind is not valid, so it is reported as no match rather than as the
// record it names.
func FindByKeyFactory(sandbox *api.Sandbox, collection *api.Collection, scope dense.Scope) func(field string, value any) (api.Record, bool, *api.Error) {
	return func(field string, value any) (api.Record, bool, *api.Error) {
		declared, ok := dense.FindField(sandbox, collection.Fields, field)
		if !ok || declared.Type != api.Key {
			return api.Record{}, false, liberror.New(sandbox, api.InvalidField, field,
				sandbox.Deps.StdDeps.Sprintf("field %q is not a Key field of the schema", field))
		}
		encoded, failure := dense.EncodeValue(sandbox, declared, value)
		if failure != nil {
			return api.Record{}, false, failure
		}
		hash := dense.HashIndexValue(sandbox, encoded)
		id, valid, err := dense.IndexOwner(sandbox, collection.Prefix, field, hash)
		if err != nil {
			return api.Record{}, false, dense.InternalError(sandbox, err)
		}
		if !valid {
			return api.Record{}, false, nil
		}
		return record.Build(sandbox, collection.Fields, collection.Prefix, scope, id), true, nil
	}
}

// FindByIDFactory fills api.Collection.FindByID: resolve the record
// straight from its permanent id, with no index lookup at all. ok is false
// when the collection holds no live record under that id.
func FindByIDFactory(sandbox *api.Sandbox, collection *api.Collection, scope dense.Scope) func(id int64) (api.Record, bool, *api.Error) {
	return func(id int64) (api.Record, bool, *api.Error) {
		return record.ResolveByID(sandbox, collection.Fields, collection.Prefix, scope, id)
	}
}

// ListAllFactory fills api.Collection.ListAll, walking the dense list
// from position 1 to the end.
func ListAllFactory(sandbox *api.Sandbox, collection *api.Collection, scope dense.Scope) func() ([]api.Record, *api.Error) {
	return func() ([]api.Record, *api.Error) {
		return record.ListRange(sandbox, collection.Fields, collection.Prefix, scope, 1, 0)
	}
}

// ListFactory fills api.Collection.List.
func ListFactory(sandbox *api.Sandbox, collection *api.Collection, scope dense.Scope) func(position int, chunk int) ([]api.Record, *api.Error) {
	return func(position int, chunk int) ([]api.Record, *api.Error) {
		return record.ListRange(sandbox, collection.Fields, collection.Prefix, scope, int64(position), int64(chunk))
	}
}

// RepairFactory fills api.Collection.Repair, running the repair of
// record.Repair under the collection's write lock.
func RepairFactory(sandbox *api.Sandbox, collection *api.Collection, scope dense.Scope) func() *api.Error {
	return func() *api.Error {
		release, failure := writelock.Acquire(sandbox, scope.Locks, scope.Root)
		if failure != nil {
			return failure
		}
		defer release()
		if failure := record.CheckOwner(sandbox, scope); failure != nil {
			return failure
		}
		return record.Repair(sandbox, collection.Fields, collection.Prefix, scope)
	}
}

// New builds an api.Collection over a prefix and the fields its records
// hold, running every factory over it to fill its function fields. Adding a
// function field to api.Collection means adding its factory call here. The
// closures read the collection New keeps for itself; the one handed back
// carries copies of Fields and Prefix.
func New(sandbox *api.Sandbox, fields []api.Field, prefix []string, scope dense.Scope) api.Collection {
	collection := api.Collection{Fields: fields, Prefix: prefix}
	collection.Insert = InsertFactory(sandbox, &collection, scope)
	collection.FindByKey = FindByKeyFactory(sandbox, &collection, scope)
	collection.FindByID = FindByIDFactory(sandbox, &collection, scope)
	collection.ListAll = ListAllFactory(sandbox, &collection, scope)
	collection.List = ListFactory(sandbox, &collection, scope)
	collection.Repair = RepairFactory(sandbox, &collection, scope)

	public := collection
	public.Fields = dense.CopyFields(sandbox, fields)
	public.Prefix = dense.Key(sandbox, prefix)
	return public
}
