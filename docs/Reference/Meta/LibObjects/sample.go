//go:build ignore

// This file is an illustrative sample, not part of the build.
package database

// KeepDatabase is an object created by the lib: it carries the injected
// Deps plus the schema description it was built from.
type KeepDatabase struct {
	Deps  deps.Deps
	Props Props
}

// GetSchema returns the collection with the given name, wiring the
// owning database into it so its methods reach storage through the same
// injected Deps. Returns nil when no schema matches.
func (d *KeepDatabase) GetSchema(name string) *SchemaInstance {
	for i := range d.Props.Schemas {
		if d.Props.Schemas[i].Name == name {
			schema := &d.Props.Schemas[i]
			return &SchemaInstance{
				db:     d,
				schema: schema,
				items:  schema.Itens,
				prefix: d.Props.Path + schema.Name,
			}
		}
	}
	return nil
}
