package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
	databases "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/databases"
	info "github.com/MateusMoutinhoOrg/Keep/sandbox/internal/info"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}

	self.Databases = databases.NewDatabases(&self)
	self.Info = info.NewInfo(&self)

	return &self
}
