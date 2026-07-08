package database

import (
	"github.com/MateusMoutinhoOrg/Keep/pkg/deps"
)

type Schema struct {
}
type Props struct {
	Path    string
	Schemas []Schema
}
type Database struct {
	Deps  deps.Deps
	Props Props
}
