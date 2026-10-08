package databases

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	writelock "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/writelock"
)

// NewFactory fills api.Databases.New with the closure that checks a Props
// description, then binds a copy of it to the sandbox, handing back the
// database every collection of it is reached through. Every database it
// builds shares locks, so two databases over the same Path still take one
// write lock per collection.
func NewFactory(sandbox *api.Sandbox, databases *api.Databases, locks *writelock.Registry) func(props api.Props) (api.Database, *api.Error) {
	return func(props api.Props) (api.Database, *api.Error) {
		if failure := Validate(sandbox, props); failure != nil {
			return api.Database{}, failure
		}
		return NewDatabase(sandbox, CopyProps(sandbox, props), locks), nil
	}
}

// NewDatabases builds the api.Databases contract, running every factory
// over it to fill its function fields. Adding a function field to
// api.Databases means adding its factory call here. It is also where the
// sandbox's registry of in-process write locks is built: once per sandbox.
func NewDatabases(sandbox *api.Sandbox) api.Databases {
	locks := writelock.NewRegistry(sandbox)
	databases := api.Databases{}
	databases.New = NewFactory(sandbox, &databases, locks)
	return databases
}
