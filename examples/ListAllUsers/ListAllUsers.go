package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

func createProps() api.Props {

	//========================Sessions==========================
	token := lib.NewKeyItem("token", true)
	creation := lib.NewIntItem("creation", true)
	expiration := lib.NewIntItem("expiration", true)
	sessions := lib.NewDatabaseItem("sessions", token, creation, expiration)

	//========================User==========================
	email := lib.NewKeyItem("email", true)
	username := lib.NewKeyItem("username", true)
	age := lib.NewIntItem("age", true)
	user := lib.NewSchema("user", email, username, age, sessions)

	//========================Props==========================
	return lib.NewProps("testDatabase/", user)
}
func main() {
	deps := standard.New()
	keep := lib.New(deps)
	props := createProps()
	db := keep.NewDatabase(props)
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
