package keep

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type KeepLib struct {
	Deps deps.Deps
}

func (l KeepLib) NewDatabase(props database.Props) *database.Database {
	return &database.Database{
		Deps:  l.Deps,
		Props: props,
	}
}

func New(d deps.Deps) *KeepLib {
	return &KeepLib{
		Deps: d,
	}
}
