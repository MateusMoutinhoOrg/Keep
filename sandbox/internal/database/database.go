package database

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/schemainstance"
)

// GetSchemaFactory fills api.KeepDatabase.GetSchema. ok is false when no
// schema of the database has the given name.
func GetSchemaFactory(kd *api.KeepDatabase) func(name string) (api.SchemaInstance, bool) {
	return func(name string) (api.SchemaInstance, bool) {
		for _, schema := range kd.Props.Schemas {
			if schema.Name == name {
				return schemainstance.New(kd.Deps, schema.Itens, kd.Props.Path+schema.Name), true
			}
		}
		return api.SchemaInstance{}, false
	}
}

// New builds an api.KeepDatabase, storing the injected Deps and the
// Props it was described with, and runs every factory over it to fill
// its function fields. Adding a function field to api.KeepDatabase
// means adding its factory call here.
func New(d deps.Deps, props api.Props) api.KeepDatabase {
	kd := api.KeepDatabase{Deps: d, Props: props}
	kd.GetSchema = GetSchemaFactory(&kd)
	return kd
}
