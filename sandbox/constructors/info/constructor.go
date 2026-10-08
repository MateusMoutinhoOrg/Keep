package info

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	info "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/info"
)

// Constructor fills Sandbox.Info, building it with the
// NewInfo of sandbox/internal/info. sandbox/new.go calls it
// once, along with the Constructor of every other package under
// sandbox/constructors/.
//
// Written once by `agnos build` and then yours: wrap the
// implementation, decorate the contract, or build a different one entirely.
// No build rewrites this file once it is there.
func Constructor(sandbox *api.Sandbox) {
	sandbox.Info = info.NewInfo(sandbox)
}
