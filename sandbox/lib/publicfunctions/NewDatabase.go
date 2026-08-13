package publicfunctions

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/lib/database"
)

// NewDatabaseFactory fills api.Lib.NewDatabase with a closure that binds a
// Props description to the injected deps, handing back the api.KeepDatabase
// every schema is reached through.
func NewDatabaseFactory(l *api.Lib) func(props api.Props) api.KeepDatabase {
	return func(props api.Props) api.KeepDatabase {
		return database.New(l.Deps, props)
	}
}
