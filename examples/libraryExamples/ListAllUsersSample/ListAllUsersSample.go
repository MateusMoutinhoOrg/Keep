package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

var Schemas = []keeptypes.Schema{
	{
		Name: "user",
		Itens: []keeptypes.Item{
			{Name: "email", Type: keeptypes.Key, Required: true},
			{Name: "username", Type: keeptypes.Key, Required: true},
			{Name: "age", Type: keeptypes.Int, Required: true},
			{
				Name: "sessions",
				Type: keeptypes.Database,
				Itens: []keeptypes.Item{
					{Name: "token", Type: keeptypes.Key, Required: true},
					{Name: "creation", Type: keeptypes.Int, Required: true},
					{Name: "expiration", Type: keeptypes.Int, Required: true},
				},
			},
		},
	},
}

var Props = keeptypes.Props{
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
			if err.Type == keeptypes.KeyConflict {
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
