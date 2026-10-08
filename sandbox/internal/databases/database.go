package databases

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	collection "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/collection"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
)

// One database: the factories filling the function fields of
// api.Database. A database owns no state beyond the Props it was built
// from — every collection it hands out is derived from that description and
// from the prefix Props.Path names, so building one writes nothing.

// CollectionFactory fills api.Database.Collection. ok is false when the
// Props declares no schema under the given name. It is also the one place a
// dense.LinkResolver is built: the database is the only thing that holds the
// whole Props, so it is the only thing that can say what collection a Link
// field points at.
func CollectionFactory(sandbox *api.Sandbox, database *api.Database) func(name string) (api.Collection, bool) {
	return func(name string) (api.Collection, bool) {
		resolve := dense.NewLinkResolver(sandbox, database.Props)
		for _, schema := range database.Props.Schemas {
			if schema.Name == name {
				prefix := dense.RootPrefix(sandbox, database.Props.Path, schema.Name)
				return collection.New(sandbox, schema.Fields, prefix, resolve), true
			}
		}
		return api.Collection{}, false
	}
}

// NewDatabase builds an api.Database over a Props description, running
// every factory over it to fill its function fields. Adding a function
// field to api.Database means adding its factory call here.
func NewDatabase(sandbox *api.Sandbox, props api.Props) api.Database {
	database := api.Database{Props: props}
	database.Collection = CollectionFactory(sandbox, &database)
	return database
}
