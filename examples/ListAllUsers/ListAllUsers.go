package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
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

	// Create 3 users before listing
	usersToCreate := []map[string]any{
		{"email": "mateus1@gmail.com", "username": "mateus1", "age": 20},
		{"email": "mateus2@gmail.com", "username": "mateus2", "age": 25},
		{"email": "mateus3@gmail.com", "username": "mateus3", "age": 30},
	}

	for _, u := range usersToCreate {
		_, err := users.NewItem(u)
		if err != nil {
			if err.Type() == api.KeyConflict {
				// Already created by a previous run, keep going
				fmt.Printf("User %v already exists, skipping\n", u["email"])
				continue
			}
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
