package standard

import (
	filestorage "github.com/MateusMoutinhoOrg/Keep/adapters/impls/filestorage"
	osstd "github.com/MateusMoutinhoOrg/Keep/adapters/impls/osstd"
	sha256hash "github.com/MateusMoutinhoOrg/Keep/adapters/impls/sha256hash"
	stdstrings "github.com/MateusMoutinhoOrg/Keep/adapters/impls/stdstrings"
	timesleep "github.com/MateusMoutinhoOrg/Keep/adapters/impls/timesleep"
	xtextfold "github.com/MateusMoutinhoOrg/Keep/adapters/impls/xtextfold"
	deps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	filestorage.Bind(&deps)
	osstd.Bind(&deps)
	sha256hash.Bind(&deps)
	stdstrings.Bind(&deps)
	timesleep.Bind(&deps)
	xtextfold.Bind(&deps)
	return deps
}
