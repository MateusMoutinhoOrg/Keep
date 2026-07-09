package main

import (
	"fmt"
	"os"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
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
	// Start from a clean database directory so the sample is deterministic
	os.RemoveAll("testDatabase")

	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("user")

	// Create 3 users before listing
	usersToCreate := []map[string]any{
		{"email": "mateus1@gmail.com", "username": "mateus1", "age": 20},
		{"email": "mateus2@gmail.com", "username": "mateus2", "age": 25},
		{"email": "mateus3@gmail.com", "username": "mateus3", "age": 30},
	}

	for _, u := range usersToCreate {
		_, err := users.NewItem(u)
		if err != nil {
			fmt.Println("Error creating user before listing all:", err)
			return
		}
	}

	// Iterate all records in the dense list (positions 1..size)
	allUsers, err := users.ListAll()
	if err != nil {
		fmt.Println("Error listing users", err)
		return
	}
	for _, user := range allUsers {
		fmt.Println("User:", user)
	}
}
