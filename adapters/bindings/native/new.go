package native

import (
	memstorage "github.com/MateusMoutinhoOrg/Keep/adapters/impls/memstorage"
	osstd "github.com/MateusMoutinhoOrg/Keep/adapters/impls/osstd"
	sha256hash "github.com/MateusMoutinhoOrg/Keep/adapters/impls/sha256hash"
	stdstrings "github.com/MateusMoutinhoOrg/Keep/adapters/impls/stdstrings"
	deps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	memstorage.Bind(&deps)
	osstd.Bind(&deps)
	sha256hash.Bind(&deps)
	stdstrings.Bind(&deps)
	return deps
}
