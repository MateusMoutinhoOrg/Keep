package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	EmailToSearch = "mateus@gmail.com"
	NewAge        = 28
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

	// Create user before updating (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type() != api.KeyConflict {
			fmt.Println("Error creating user before update:", err)
			return
		}
		fmt.Println("User already exists, updating the existing one")
	}

	// Find the user to update
	foundUser := users.FindByKey("email", EmailToSearch)
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Update a non-indexed field (simple single key write)
	errUpdate := foundUser.Update("age", NewAge)
	if errUpdate != nil {
		fmt.Println("Error updating user", errUpdate)
		return
	}
	fmt.Println("User updated successfully")
}
