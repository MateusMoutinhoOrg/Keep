package database

type ItemType int

const (
	Key ItemType = iota
	Int
)

type Item struct {
	Name     string
	Type     ItemType
	Required bool
}

type Schema struct {
	Name  string
	Itens []Item
}
