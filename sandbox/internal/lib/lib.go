package lib

import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/internal/database"
)

// Lib implements api.Lib. It holds the injected Deps and propagates
// them to every database it creates.
type Lib struct {
	Deps deps.Deps
}

// NewDatabase creates a database from a Props description, with the
// lib's injected dependencies wired in.
func (l *Lib) NewDatabase(props api.Props) api.KeepDatabase {
	return &database.KeepDatabase{
		Deps:        l.Deps,
		Description: props,
	}
}
