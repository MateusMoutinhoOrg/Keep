package lib

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/lib/publicfunctions"
)

// New builds the api.Lib entry point, storing the injected deps on it
// and running every lib factory over it to fill its function fields.
// Adding a function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.Version = publicfunctions.VersionFactory(&l)
	l.NewDatabase = publicfunctions.NewDatabaseFactory(&l)
	return l
}
