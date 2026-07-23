package lib

// NewDatabase creates a KeepDatabase from a Props description, with the
// lib's injected dependencies wired in.
func (l *Lib) NewDatabase(props Props) *KeepDatabase {
	return &KeepDatabase{
		deps:  l.deps,
		Props: props,
	}
}

// GetSchema returns the collection with the given name, or nil when no
// schema of the database has that name.
func (d *KeepDatabase) GetSchema(name string) *SchemaInstance {
	for i := range d.Props.Schemas {
		if d.Props.Schemas[i].Name == name {
			schema := &d.Props.Schemas[i]
			return &SchemaInstance{
				deps:   d.deps,
				items:  schema.Itens,
				prefix: d.Props.Path + schema.Name,
			}
		}
	}
	return nil
}
