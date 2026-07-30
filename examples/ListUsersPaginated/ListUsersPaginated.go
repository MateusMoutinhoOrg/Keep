package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	// Position is the starting index in the dense list (1-based)
	StartPosition = 1
	// ChunkSize is the number of records to fetch per page
	ChunkSize = 10
)

var Props = lib.NewProps("testDatabase/",
	lib.NewSchema("user",
		lib.NewKeyItem("email", true),
		lib.NewKeyItem("username", true),
		lib.NewIntItem("age", true),
		lib.NewDatabaseItem("sessions",
			lib.NewKeyItem("token", true),
			lib.NewIntItem("creation", true),
			lib.NewIntItem("expiration", true),
		),
	),
)

func main() {
	deps := standard.New()
	keep := lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("user")

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
			if err.Type() == api.KeyConflict {
				// Already created by a previous run, keep going
				fmt.Printf("User %v already exists, skipping\n", u["email"])
				continue
			}
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
