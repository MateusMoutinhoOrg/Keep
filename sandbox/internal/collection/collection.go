package collection

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	record "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/record"
)

// One collection of records: the factories filling the function fields of
// api.Collection. A collection is nothing but a key prefix and the list
// of fields its records hold, so the same New builds a top-level collection
// and a nested one. The dense.LinkResolver it is built with is the database
// the collection belongs to, as far as a Link field is concerned; every
// factory passes it on unchanged.

// InsertFactory fills api.Collection.Insert.
func InsertFactory(sandbox *api.Sandbox, collection *api.Collection, resolve dense.LinkResolver) func(fields map[string]any) (api.Record, *api.Error) {
	return func(fields map[string]any) (api.Record, *api.Error) {
		return record.Insert(sandbox, collection.Fields, collection.Prefix, resolve, fields)
	}
}

// FindByKeyFactory fills api.Collection.FindByKey: encode and hash the
// value, resolve the id through the unique index in one read, then check
// the record it names is still live. ok is false when the field is not an
// indexed Key or no live record matches.
func FindByKeyFactory(sandbox *api.Sandbox, collection *api.Collection, resolve dense.LinkResolver) func(field string, value any) (api.Record, bool) {
	return func(field string, value any) (api.Record, bool) {
		declared, ok := dense.FindField(sandbox, collection.Fields, field)
		if !ok || declared.Type != api.Key {
			return api.Record{}, false
		}
		encoded, failure := dense.EncodeValue(sandbox, declared, value)
		if failure != nil {
			return api.Record{}, false
		}
		indexKey := dense.IndexKey(sandbox, collection.Prefix, field, dense.HashIndexValue(sandbox, encoded))
		raw, found, err := sandbox.Deps.StorageDeps.Read(indexKey)
		if err != nil || !found {
			return api.Record{}, false
		}
		return record.ResolveRawID(sandbox, collection.Fields, collection.Prefix, resolve, raw)
	}
}

// FindByIDFactory fills api.Collection.FindByID: resolve the record
// straight from its permanent id, with no index lookup at all. ok is false
// when the collection holds no live record under that id.
func FindByIDFactory(sandbox *api.Sandbox, collection *api.Collection, resolve dense.LinkResolver) func(id int64) (api.Record, bool) {
	return func(id int64) (api.Record, bool) {
		return record.ResolveByID(sandbox, collection.Fields, collection.Prefix, resolve, id)
	}
}

// ListAllFactory fills api.Collection.ListAll, walking the dense list
// from position 1 to the end.
func ListAllFactory(sandbox *api.Sandbox, collection *api.Collection, resolve dense.LinkResolver) func() ([]api.Record, *api.Error) {
	return func() ([]api.Record, *api.Error) {
		return record.ListRange(sandbox, collection.Fields, collection.Prefix, resolve, 1, 0)
	}
}

// ListFactory fills api.Collection.List.
func ListFactory(sandbox *api.Sandbox, collection *api.Collection, resolve dense.LinkResolver) func(position int, chunk int) ([]api.Record, *api.Error) {
	return func(position int, chunk int) ([]api.Record, *api.Error) {
		return record.ListRange(sandbox, collection.Fields, collection.Prefix, resolve, int64(position), int64(chunk))
	}
}

// New builds an api.Collection over a prefix and the fields its records
// hold, running every factory over it to fill its function fields. Adding a
// function field to api.Collection means adding its factory call here.
func New(sandbox *api.Sandbox, fields []api.Field, prefix []string, resolve dense.LinkResolver) api.Collection {
	collection := api.Collection{Fields: fields, Prefix: prefix}
	collection.Insert = InsertFactory(sandbox, &collection, resolve)
	collection.FindByKey = FindByKeyFactory(sandbox, &collection, resolve)
	collection.FindByID = FindByIDFactory(sandbox, &collection, resolve)
	collection.ListAll = ListAllFactory(sandbox, &collection, resolve)
	collection.List = ListFactory(sandbox, &collection, resolve)
	return collection
}
