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
			return &SchemaInstance{
				db:     d,
				schema: &d.Props.Schemas[i],
			}
		}
	}
	return nil
}
