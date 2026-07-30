package lib

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	internallib "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/lib"
)

// New injects a Deps implementation into the library and returns
// the api.Lib entry point.
func New(d deps.Deps) api.Lib {
	return &internallib.Lib{Deps: d}
}
