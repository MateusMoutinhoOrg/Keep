package database

type SchemaInstance struct{}

func (si *SchemaInstance) FindByKey(key string, keyValue any) *SchemaItem { return nil }

func (si *SchemaInstance) NewItem(fields map[string]any) (*SchemaItem, Error) { return nil, Error{} }
