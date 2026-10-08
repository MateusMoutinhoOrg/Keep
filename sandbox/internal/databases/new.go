package databases

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// NewFactory fills api.Databases.New with the closure that binds a Props
// description to the sandbox, handing back the database every collection of
// it is reached through.
func NewFactory(sandbox *api.Sandbox, databases *api.Databases) func(props api.Props) api.Database {
	return func(props api.Props) api.Database {
		return NewDatabase(sandbox, props)
	}
}

// NewDatabases builds the api.Databases contract, running every factory
// over it to fill its function fields. Adding a function field to
// api.Databases means adding its factory call here.
func NewDatabases(sandbox *api.Sandbox) api.Databases {
	databases := api.Databases{}
	databases.New = NewFactory(sandbox, &databases)
	return databases
}
