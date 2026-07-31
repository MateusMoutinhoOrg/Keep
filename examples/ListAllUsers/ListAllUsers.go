package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	database "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

var Schemas = []database.Schema{
	{
		Name: "user",
		Itens: []database.Item{
			{Name: "email", Type: database.Key, Required: true},
			{Name: "username", Type: database.Key, Required: true},
			{Name: "age", Type: database.Int, Required: true},
			{
				Name: "sessions",
				Type: database.Database,
				Itens: []database.Item{
					{Name: "token", Type: database.Key, Required: true},
					{Name: "creation", Type: database.Int, Required: true},
					{Name: "expiration", Type: database.Int, Required: true},
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
	deps := keepadapter.New()
	keep := keeplib.New(deps)
	db := keep.NewDatabase(Props)
	users, _ := db.GetSchema("user")

	// Create 3 users before listing
	usersToCreate := []map[string]any{
		{"email": "mateus1@gmail.com", "username": "mateus1", "age": 20},
		{"email": "mateus2@gmail.com", "username": "mateus2", "age": 25},
		{"email": "mateus3@gmail.com", "username": "mateus3", "age": 30},
	}

	for _, u := range usersToCreate {
		_, err := users.NewItem(u)
		if err != nil {
			if err.Type == database.KeyConflict {
				// Already created by a previous run, keep going
				fmt.Printf("User %v already exists, skipping\n", u["email"])
				continue
			}
			fmt.Println("Error creating user before listing all:", err.Message)
			return
		}
	}

	// Iterate all records in the dense list (positions 1..size)
	allUsers, err := users.ListAll()
	if err != nil {
		fmt.Println("Error listing users", err.Message)
		return
	}
	for _, user := range allUsers {
		fmt.Println("User:", user.String())
	}
}
