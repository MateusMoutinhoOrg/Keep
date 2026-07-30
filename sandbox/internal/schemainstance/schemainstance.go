package schemainstance

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/schemaitem"
)

// SchemaInstance implements api.SchemaInstance. It carries the Deps it
// was created with and propagates them to every record it hands back.
type SchemaInstance struct {
	Deps   deps.Deps
	Items  []api.Item
	Prefix string
}

// FindByKey looks a record up through its unique index: normalize and
// hash the value, resolve the id, and check the record is live.
// Returns nil when the key field is unknown or no record matches.
func (si *SchemaInstance) FindByKey(key string, keyValue any) api.SchemaItem {
	item := dense.FindItem(si.Items, key)
	if item == nil || item.Type() != api.KeyItem {
		return nil
	}
	encoded, e := dense.EncodeValue(item, keyValue)
	if e != nil {
		return nil
	}
	raw, err := si.Deps.Read(dense.IndexKey(si.Prefix, key, dense.HashIndexValue(encoded)))
	if err != nil {
		return nil
	}
	return schemaitem.ResolveLive(si.Deps, si.Items, si.Prefix, raw)
}

// NewItem inserts a record into the collection, validating the provided
// fields against the schema.
func (si *SchemaInstance) NewItem(fields map[string]any) (api.SchemaItem, api.Error) {
	return schemaitem.New(si.Deps, si.Items, si.Prefix, fields)
}

// ListAll iterates the dense list from position 1 through size.
func (si *SchemaInstance) ListAll() ([]api.SchemaItem, api.Error) {
	return schemaitem.ListRange(si.Deps, si.Items, si.Prefix, 1, 0)
}

// List returns up to `chunk` records starting at `position` (1-based).
func (si *SchemaInstance) List(position int, chunk int) ([]api.SchemaItem, api.Error) {
	return schemaitem.ListRange(si.Deps, si.Items, si.Prefix, int64(position), int64(chunk))
}
