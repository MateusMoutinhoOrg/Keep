package database

type SchemaItem struct{}

func (s *SchemaItem) CheckKeysPresence(keys []string) bool { return false }

func (s *SchemaItem) Remove() Error { return Error{} }
