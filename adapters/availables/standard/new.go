package standard

import (
	filestorage "github.com/MateusMoutinhoOrg/Keep/adapters/libs/filestorage"
	hashdeps "github.com/MateusMoutinhoOrg/Keep/adapters/libs/hashdeps"
	std "github.com/MateusMoutinhoOrg/Keep/adapters/libs/std"
	stringsdeps "github.com/MateusMoutinhoOrg/Keep/adapters/libs/stringsdeps"
	deps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	filestorage.Bind(&deps)
	hashdeps.Bind(&deps)
	std.Bind(&deps)
	stringsdeps.Bind(&deps)
	return deps
}
