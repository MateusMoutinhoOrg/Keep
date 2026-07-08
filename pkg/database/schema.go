package database

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type Schema struct {
}

type Database struct {
	Deps    deps.Deps
	Schemas []Schema
}
