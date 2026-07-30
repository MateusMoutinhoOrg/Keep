package database

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/schemainstance"
)

// KeepDatabase implements api.KeepDatabase. It holds the Props it was
// described with and the injected Deps, propagating both to every
// collection it hands back.
type KeepDatabase struct {
	Deps        deps.Deps
	Description api.Props
}

// Props returns the description the database was created from.
func (d *KeepDatabase) Props() api.Props { return d.Description }

// GetSchema returns the collection with the given name, or nil when no
// schema of the database has that name.
func (d *KeepDatabase) GetSchema(name string) api.SchemaInstance {
	for _, schema := range d.Description.Schemas() {
		if schema.Name() == name {
			return &schemainstance.SchemaInstance{
				Deps:   d.Deps,
				Items:  schema.Itens(),
				Prefix: d.Description.Path() + schema.Name(),
			}
		}
	}
	return nil
}
