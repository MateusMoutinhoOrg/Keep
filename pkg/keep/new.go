package keep

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type KeepLib struct {
	Deps deps.Deps
}

func (l KeepLib) NewDatabase(schema []database.Schema) *database.Database {
	return &database.Database{
		Deps:    l.Deps,
		Schemas: schema,
	}
}

func New(d deps.Deps) *KeepLib {
	return &KeepLib{
		Deps: d,
	}
}
