package keep

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type KeepLib struct {
	deps deps.Deps
}

func (l KeepLib) NewDatabase(schema database.Schema) database.Database {
	return database.Database{
		deps:   l.deps,
		schema: schema,
	}
}

func New(d deps.Deps) *KeepLib {
	return &KeepLib{
		deps: d,
	}
}
