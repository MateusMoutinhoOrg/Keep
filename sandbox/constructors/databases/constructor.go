package databases

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	databases "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/databases"
)

// Constructor fills Sandbox.Databases, building it with the
// NewDatabases of sandbox/internal/databases. sandbox/new.go calls it
// once, along with the Constructor of every other package under
// sandbox/constructors/.
//
// Written once by `agnos build` and then yours: wrap the
// implementation, decorate the contract, or build a different one entirely.
// No build rewrites this file once it is there.
func Constructor(sandbox *api.Sandbox) {
	sandbox.Databases = databases.NewDatabases(sandbox)
}
