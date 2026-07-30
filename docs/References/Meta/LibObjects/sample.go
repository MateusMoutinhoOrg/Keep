//go:build ignore

// This file is an illustrative sample, not part of the build.
package database

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/schemainstance"
)

// KeepDatabase implements api.KeepDatabase. It carries the injected
// Deps plus the description it was built from, and propagates both to
// every collection it hands back.
type KeepDatabase struct {
	Deps        deps.Deps
	Description api.Props
}

// Props returns the description the database was created from.
func (d *KeepDatabase) Props() api.Props { return d.Description }

// GetSchema returns the collection with the given name, copying the
// database's Deps into it so its methods reach storage through the same
// injected backend. It returns a literal nil when no schema matches, so
// the caller's `== nil` check works.
func (d *KeepDatabase) GetSchema(name string) api.SchemaInstance {
	// Props, Schema and Item are interfaces, so their data is read
	// through methods rather than fields.
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
