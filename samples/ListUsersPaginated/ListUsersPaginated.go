package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

const (
	// Position is the starting index in the dense list (1-based)
	StartPosition = 1
	// ChunkSize is the number of records to fetch per page
	ChunkSize = 10
)

var Schemas = []database.Schema{
	{
		Name: "User",
		Itens: []database.Item{
			{
				Type:     database.Key,
				Required: true,
				Name:     "Email",
			},
			{
				Type:     database.Key,
				Required: true,
				Name:     "UserName",
			},
			{
				Name:     "Age",
				Required: true,
				Type:     database.Int,
			},
		},
	},
}

var Props = database.Props{
	Path:    "testDatabase/",
	Schemas: Schemas,
}

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("Users")

	// Fetch a chunk of users starting at a given position
	// List(position, chunk) returns up to `chunk` records starting from `position`
	userPage, err := users.List(StartPosition, ChunkSize)
	if err != nil {
		fmt.Println("Error listing users", err)
		return
	}
	for _, user := range userPage {
		fmt.Println("User:", user)
	}
}
