package lib

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type Lib struct {
	deps deps.Deps
}

func New(d deps.Deps) Lib {
	return Lib{deps: d}
}
