package database

type SchemaItem struct {
	db     *KeepDatabase
	schema *Schema
	fields map[string]any
}

func (s *SchemaItem) CheckKeysPresence(keys []string) bool { return false }

func (s *SchemaItem) Get(fieldName string) (any, *Error) { return nil, nil }

func (s *SchemaItem) Update(fieldName string, value any) *Error { return nil }

func (s *SchemaItem) Remove() Error { return Error{} }

func (s *SchemaItem) ListAll(fieldName string) []*SchemaItem { return nil }
