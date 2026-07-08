package database

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type Props struct {
	Path         string
	FixIntegrity bool
	Schemas      []Schema
}

type Database struct {
	Deps  deps.Deps
	Props Props
}

func (d *Database) GetSchema(name string) *SchemaInstance { return nil }
