package lib

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/database"
)

// NewDatabaseFactory fills api.Lib.NewDatabase.
func NewDatabaseFactory(l *api.Lib) func(props api.Props) api.KeepDatabase {
	return func(props api.Props) api.KeepDatabase {
		return database.New(l.Deps, props)
	}
}

// New builds the api.Lib entry point, storing the injected deps on it
// and running every lib factory over it to fill its function fields.
// Adding a function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.NewDatabase = NewDatabaseFactory(&l)
	return l
}
