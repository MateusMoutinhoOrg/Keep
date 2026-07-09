package database

type SchemaInstance struct {
	db     *KeepDatabase
	schema *Schema
}

func (si *SchemaInstance) FindByKey(key string, keyValue any) *SchemaItem { return nil }

func (si *SchemaInstance) NewItem(fields map[string]any) (*SchemaItem, *Error) { return nil, nil }

func (si *SchemaInstance) ListAll() ([]*SchemaItem, *Error) { return nil, nil }

func (si *SchemaInstance) List(position int, chunk int) ([]*SchemaItem, *Error) { return nil, nil }
