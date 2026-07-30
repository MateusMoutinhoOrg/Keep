package lib

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/description"
)

// The constructors callers use to describe a database. The contracts
// package declares Props, Schema and Item as interfaces, so they cannot
// be built with a composite literal — these functions are the way in.
// Each takes primitives and interfaces, and hands back an interface.

// NewProps describes a database: the prefix every key is written under,
// and the collections it holds.
func NewProps(path string, schemas ...api.Schema) api.Props {
	return &description.Props{KeyPrefix: path, Collections: schemas}
}

// NewSchema describes one collection of records and its fields.
func NewSchema(name string, itens ...api.Item) api.Schema {
	return &description.Schema{SchemaName: name, Fields: itens}
}

// NewKeyItem describes a unique, indexed string field. Values are
// matched case-insensitively and can be looked up with FindByKey.
func NewKeyItem(name string, required bool) api.Item {
	return &description.Item{ItemName: name, ItemType: api.KeyItem, IsNeeded: required}
}

// NewIntItem describes a plain integer field. Values are accepted as
// int, int32, or int64, and always read back as int64.
func NewIntItem(name string, required bool) api.Item {
	return &description.Item{ItemName: name, ItemType: api.IntItem, IsNeeded: required}
}

// NewDatabaseItem describes a nested collection owned by each record of
// the schema, with its own fields.
func NewDatabaseItem(name string, itens ...api.Item) api.Item {
	return &description.Item{ItemName: name, ItemType: api.DatabaseItem, SubItens: itens}
}
