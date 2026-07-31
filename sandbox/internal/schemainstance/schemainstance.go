package schemainstance

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/schemaitem"
)

// FindByKeyFactory fills api.SchemaInstance.FindByKey: normalize and
// hash the value, resolve the id through the unique index, and check
// the record is live. ok is false when the key field is unknown or no
// record matches.
func FindByKeyFactory(si *api.SchemaInstance) func(key string, keyValue any) (api.SchemaItem, bool) {
	return func(key string, keyValue any) (api.SchemaItem, bool) {
		item, ok := dense.FindItem(si.Items, key)
		if !ok || item.Type != api.Key {
			return api.SchemaItem{}, false
		}
		encoded, e := dense.EncodeValue(item, keyValue)
		if e != nil {
			return api.SchemaItem{}, false
		}
		raw, err := si.Deps.Read(dense.IndexKey(si.Prefix, key, dense.HashIndexValue(encoded)))
		if err != nil {
			return api.SchemaItem{}, false
		}
		return schemaitem.ResolveLive(si.Deps, si.Items, si.Prefix, raw)
	}
}

// NewItemFactory fills api.SchemaInstance.NewItem.
func NewItemFactory(si *api.SchemaInstance) func(fields map[string]any) (api.SchemaItem, *api.Error) {
	return func(fields map[string]any) (api.SchemaItem, *api.Error) {
		return schemaitem.New(si.Deps, si.Items, si.Prefix, fields)
	}
}

// ListAllFactory fills api.SchemaInstance.ListAll, iterating the dense
// list from position 1 through size.
func ListAllFactory(si *api.SchemaInstance) func() ([]api.SchemaItem, *api.Error) {
	return func() ([]api.SchemaItem, *api.Error) {
		return schemaitem.ListRange(si.Deps, si.Items, si.Prefix, 1, 0)
	}
}

// ListFactory fills api.SchemaInstance.List.
func ListFactory(si *api.SchemaInstance) func(position int, chunk int) ([]api.SchemaItem, *api.Error) {
	return func(position int, chunk int) ([]api.SchemaItem, *api.Error) {
		return schemaitem.ListRange(si.Deps, si.Items, si.Prefix, int64(position), int64(chunk))
	}
}

// New builds an api.SchemaInstance, storing the injected Deps and the
// schema's fields and prefix, and runs every factory over it to fill
// its function fields. Adding a function field to api.SchemaInstance
// means adding its factory call here.
func New(d deps.Deps, items []api.Item, prefix string) api.SchemaInstance {
	si := api.SchemaInstance{Deps: d, Items: items, Prefix: prefix}
	si.FindByKey = FindByKeyFactory(&si)
	si.NewItem = NewItemFactory(&si)
	si.ListAll = ListAllFactory(&si)
	si.List = ListFactory(&si)
	return si
}
