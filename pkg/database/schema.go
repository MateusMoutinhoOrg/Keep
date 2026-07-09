package database

type ItemType int

const (
	Key ItemType = iota
	Int
	Database
)

type Item struct {
	Name     string
	Type     ItemType
	Required bool
	Itens    []Item
}

type Schema struct {
	Name  string
	Itens []Item
}
