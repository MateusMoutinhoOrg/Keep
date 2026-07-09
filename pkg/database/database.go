package database

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type Props struct {
	Path    string
	Schemas []Schema
}

type KeepDatabase struct {
	Deps  deps.Deps
	Props Props
}

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
