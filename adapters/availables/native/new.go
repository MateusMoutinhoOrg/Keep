package native

import (
	hashdeps "github.com/MateusMoutinhoOrg/Keep/adapters/libs/hashdeps"
	memstorage "github.com/MateusMoutinhoOrg/Keep/adapters/libs/memstorage"
	std "github.com/MateusMoutinhoOrg/Keep/adapters/libs/std"
	stringsdeps "github.com/MateusMoutinhoOrg/Keep/adapters/libs/stringsdeps"
	deps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	hashdeps.Bind(&deps)
	memstorage.Bind(&deps)
	std.Bind(&deps)
	stringsdeps.Bind(&deps)
	return deps
}
