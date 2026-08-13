package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	// StartPosition is the starting index in the dense list (1-based)
	StartPosition = 1
	// ChunkSize is the number of records to fetch per page
	ChunkSize = 10
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
			if err.Type == keeptypes.KeyConflict {
				// Already created by a previous run, keep going
				fmt.Printf("User %v already exists, skipping\n", u["email"])
				continue
			}
			fmt.Println("Error creating user before paginated listing:", err.Message)
			return
		}
	}

	// Fetch a chunk of users starting at a given position
	// List(position, chunk) returns up to `chunk` records starting from `position`
	userPage, err := users.List(StartPosition, ChunkSize)
	if err != nil {
		fmt.Println("Error listing users", err.Message)
		return
	}
	for _, user := range userPage {
		fmt.Println("User:", user.String())
	}
}
