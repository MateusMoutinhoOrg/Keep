//go:build ignore

// This file is an illustrative sample, not part of the build.
package lib

// KeepDatabase is an object created by the lib: it carries the injected
// deps plus the schema description it was built from. Declared in
// pkg/lib/types.go.
type KeepDatabase struct {
	deps  deps.Deps
	Props Props
}

// GetSchema returns the collection with the given name, copying the
// database's deps into it so its methods reach storage through the same
// injected functions. Returns nil when no schema matches.
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
