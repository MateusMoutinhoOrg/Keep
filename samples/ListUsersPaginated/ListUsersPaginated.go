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
		Name: "user",
		Itens: []database.Item{
			{
				Type:     database.Key,
				Required: true,
				Name:     "email",
			},
			{
				Type:     database.Key,
				Required: true,
				Name:     "username",
			},
			{
				Name:     "age",
				Required: true,
				Type:     database.Int,
			},

			{
				Name: "sessions",
				Type: database.Database,
				Itens: []database.Item{
					{
						Name:     "token",
						Type:     database.Key,
						Required: true,
					},
					{
						Name:     "creation",
						Type:     database.Int,
						Required: true,
					},
					{
						Name:     "expiration",
						Type:     database.Int,
						Required: true,
					},
				},
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
	users := db.GetSchema("users")

	// Create 5 users before listing paginated
	usersToCreate := []map[string]any{
		{"email": "mateus1@gmail.com", "username": "mateus1", "age": 21},
		{"email": "mateus2@gmail.com", "username": "mateus2", "age": 22},
		{"email": "mateus3@gmail.com", "username": "mateus3", "age": 23},
		{"email": "mateus4@gmail.com", "username": "mateus4", "age": 24},
		{"email": "mateus5@gmail.com", "username": "mateus5", "age": 25},
	}

	for _, u := range usersToCreate {
		_, err := users.NewItem(u)
		if err != nil {
			fmt.Println("Error creating user before paginated listing:", err)
			return
		}
	}

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
