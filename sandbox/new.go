package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	config "github.com/MateusMoutinhoOrg/Keep/sandbox/constructors/config"
	databases "github.com/MateusMoutinhoOrg/Keep/sandbox/constructors/databases"
	info "github.com/MateusMoutinhoOrg/Keep/sandbox/constructors/info"
	deps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

// New builds the whole library: one call per package under
// sandbox/constructors/, each filling the field of the Sandbox it owns. The
// list is the directories themselves, so a constructor written by hand is
// called exactly like a generated one — this file is rendered around what is
// there, never the other way round.
func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}

	config.Constructor(&self)
	databases.Constructor(&self)
	info.Constructor(&self)

	return &self
}
