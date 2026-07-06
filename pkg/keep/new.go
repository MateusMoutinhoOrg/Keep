package keep

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type KeepLib struct {
	deps deps.Deps
}

func New(d deps.Deps) *KeepLib {
	return &KeepLib{
		deps: d,
	}
}
