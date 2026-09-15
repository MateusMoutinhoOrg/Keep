package databases

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	dense "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/dense"
	schemainstance "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/schemainstance"
)

// One database: the factories filling the function fields of
// api.DatabaseHandle. A handle owns no state beyond the Props it was built
// from — every collection it hands out is derived from that description and
// from the prefix Props.Path names, so building one writes nothing.

// GetSchemaFactory fills api.DatabaseHandle.GetSchema. ok is false when the
// Props declares no schema under the given name.
func GetSchemaFactory(sandbox *api.Sandbox, handle *api.DatabaseHandle) func(name string) (api.SchemaInstance, bool) {
	return func(name string) (api.SchemaInstance, bool) {
		for _, schema := range handle.Props.Schemas {
			if schema.Name == name {
				prefix := dense.RootPrefix(sandbox, handle.Props.Path, schema.Name)
				return schemainstance.New(sandbox, schema.Itens, prefix), true
			}
		}
		return api.SchemaInstance{}, false
	}
}

// NewHandle builds an api.DatabaseHandle over a Props description, running
// every factory over it to fill its function fields. Adding a function
// field to api.DatabaseHandle means adding its factory call here.
func NewHandle(sandbox *api.Sandbox, props api.Props) api.DatabaseHandle {
	handle := api.DatabaseHandle{Props: props}
	handle.GetSchema = GetSchemaFactory(sandbox, &handle)
	return handle
}
