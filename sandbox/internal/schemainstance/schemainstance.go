package schemainstance

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	schemaitem "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/schemaitem"
)

// One collection of records: the factories filling the function fields of
// api.SchemaInstance. A collection is nothing but a key prefix and the list
// of fields its records hold, so the same New builds a top-level collection
// and a nested one. The dense.LinkResolver it is built with is the database
// the collection belongs to, as far as a Link field is concerned; every
// factory passes it on unchanged.

// NewItemFactory fills api.SchemaInstance.NewItem.
func NewItemFactory(sandbox *api.Sandbox, instance *api.SchemaInstance, resolve dense.LinkResolver) func(fields map[string]any) (api.SchemaItem, *api.Error) {
	return func(fields map[string]any) (api.SchemaItem, *api.Error) {
		return schemaitem.New(sandbox, instance.Items, instance.Prefix, resolve, fields)
	}
}

// FindByKeyFactory fills api.SchemaInstance.FindByKey: encode and hash the
// value, resolve the id through the unique index in one read, then check
// the record it names is still live. ok is false when the field is not an
// indexed Key or no live record matches.
func FindByKeyFactory(sandbox *api.Sandbox, instance *api.SchemaInstance, resolve dense.LinkResolver) func(key string, keyValue any) (api.SchemaItem, bool) {
	return func(key string, keyValue any) (api.SchemaItem, bool) {
		item, ok := dense.FindItem(sandbox, instance.Items, key)
		if !ok || item.Type != api.Key {
			return api.SchemaItem{}, false
		}
		encoded, failure := dense.EncodeValue(sandbox, item, keyValue)
		if failure != nil {
			return api.SchemaItem{}, false
		}
		indexKey := dense.IndexKey(sandbox, instance.Prefix, key, dense.HashIndexValue(sandbox, encoded))
		raw, found, err := sandbox.Deps.Storagedeps.Read(indexKey)
		if err != nil || !found {
			return api.SchemaItem{}, false
		}
		return schemaitem.ResolveLive(sandbox, instance.Items, instance.Prefix, resolve, raw)
	}
}

// FindByIdFactory fills api.SchemaInstance.FindById: resolve the record
// straight from its permanent id, with no index lookup at all. ok is false
// when the collection holds no live record under that id.
func FindByIdFactory(sandbox *api.Sandbox, instance *api.SchemaInstance, resolve dense.LinkResolver) func(id int64) (api.SchemaItem, bool) {
	return func(id int64) (api.SchemaItem, bool) {
		return schemaitem.ResolveById(sandbox, instance.Items, instance.Prefix, resolve, id)
	}
}

// ListAllFactory fills api.SchemaInstance.ListAll, walking the dense list
// from position 1 to the end.
func ListAllFactory(sandbox *api.Sandbox, instance *api.SchemaInstance, resolve dense.LinkResolver) func() ([]api.SchemaItem, *api.Error) {
	return func() ([]api.SchemaItem, *api.Error) {
		return schemaitem.ListRange(sandbox, instance.Items, instance.Prefix, resolve, 1, 0)
	}
}

// ListFactory fills api.SchemaInstance.List.
func ListFactory(sandbox *api.Sandbox, instance *api.SchemaInstance, resolve dense.LinkResolver) func(position int, chunk int) ([]api.SchemaItem, *api.Error) {
	return func(position int, chunk int) ([]api.SchemaItem, *api.Error) {
		return schemaitem.ListRange(sandbox, instance.Items, instance.Prefix, resolve, int64(position), int64(chunk))
	}
}

// New builds an api.SchemaInstance over a prefix and the fields its records
// hold, running every factory over it to fill its function fields. Adding a
// function field to api.SchemaInstance means adding its factory call here.
func New(sandbox *api.Sandbox, items []api.Item, prefix []string, resolve dense.LinkResolver) api.SchemaInstance {
	instance := api.SchemaInstance{Items: items, Prefix: prefix}
	instance.NewItem = NewItemFactory(sandbox, &instance, resolve)
	instance.FindByKey = FindByKeyFactory(sandbox, &instance, resolve)
	instance.FindById = FindByIdFactory(sandbox, &instance, resolve)
	instance.ListAll = ListAllFactory(sandbox, &instance, resolve)
	instance.List = ListFactory(sandbox, &instance, resolve)
	return instance
}
