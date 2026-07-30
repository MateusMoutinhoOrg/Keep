package description

// The structs backing the caller-supplied description. They implement
// api.Props, api.Schema and api.Item so the contracts package can stay
// interfaces-only; callers build them through the constructors exported
// by the sandbox package.

import "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"

// Item implements api.Item.
type Item struct {
	ItemName string
	ItemType int
	IsNeeded bool
	SubItens []api.Item
}

func (i *Item) Name() string      { return i.ItemName }
func (i *Item) Type() int         { return i.ItemType }
func (i *Item) Required() bool    { return i.IsNeeded }
func (i *Item) Itens() []api.Item { return i.SubItens }

// Schema implements api.Schema.
type Schema struct {
	SchemaName string
	Fields     []api.Item
}

func (s *Schema) Name() string      { return s.SchemaName }
func (s *Schema) Itens() []api.Item { return s.Fields }

// Props implements api.Props.
type Props struct {
	KeyPrefix   string
	Collections []api.Schema
}

func (p *Props) Path() string          { return p.KeyPrefix }
func (p *Props) Schemas() []api.Schema { return p.Collections }
