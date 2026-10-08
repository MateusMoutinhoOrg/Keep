package databases

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	collection "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/collection"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	writelock "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/writelock"
)

// One database: the factories filling the function fields of
// api.Database. A database owns no state beyond the Props it was built
// from — every collection it hands out is derived from that description and
// from the prefix Props.Path names, so building one writes nothing.

// CollectionFactory fills api.Database.Collection. ok is false when the
// Props declares no schema under the given name. It is also the one place a
// dense.Scope is built: the database is the only thing that holds the
// whole Props, so it is the only thing that can say what collection a Link
// field points at, and it hands every collection the sandbox's write locks.
func CollectionFactory(sandbox *api.Sandbox, database *api.Database, locks *writelock.Registry) func(name string) (api.Collection, bool) {
	return func(name string) (api.Collection, bool) {
		resolve := dense.NewLinkResolver(sandbox, database.Props)
		for _, schema := range database.Props.Schemas {
			if schema.Name == name {
				prefix := dense.RootPrefix(sandbox, database.Props.Path, schema.Name)
				scope := dense.Scope{
					Resolve:       resolve,
					Locks:         locks,
					Root:          prefix,
					NewCollection: collection.New,
				}
				return collection.New(sandbox, schema.Fields, prefix, scope), true
			}
		}
		return api.Collection{}, false
	}
}

// NewDatabase builds an api.Database over a Props description, running
// every factory over it to fill its function fields. Adding a function
// field to api.Database means adding its factory call here. The closures
// read the database NewDatabase keeps for itself; the one handed back
// carries a copy of the Props.
func NewDatabase(sandbox *api.Sandbox, props api.Props, locks *writelock.Registry) api.Database {
	database := api.Database{Props: props}
	database.Collection = CollectionFactory(sandbox, &database, locks)

	public := database
	public.Props = CopyProps(sandbox, props)
	return public
}

// CopyProps returns a deep copy of a Props, so a caller editing the one it
// passed in, or the one a Database hands back, changes nothing a database
// does.
func CopyProps(sandbox *api.Sandbox, props api.Props) api.Props {
	copied := api.Props{Path: props.Path}
	if props.Schemas != nil {
		copied.Schemas = make([]api.Schema, len(props.Schemas))
		for index, schema := range props.Schemas {
			copied.Schemas[index] = api.Schema{Name: schema.Name, Fields: dense.CopyFields(sandbox, schema.Fields)}
		}
	}
	return copied
}
