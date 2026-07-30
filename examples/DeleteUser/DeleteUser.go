package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	EmailToDelete = "mateus@gmail.com"
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

	// Create the user first before deleting (skip if it survived a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToDelete,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type() != api.KeyConflict {
			fmt.Println("Error creating user before delete:", err)
			return
		}
		fmt.Println("User already exists, deleting the existing one")
	}

	// First, find the user by key
	foundUser := users.FindByKey("email", EmailToDelete)
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Then, remove the user (swap-with-last deletion)
	errRemove := foundUser.Remove()
	if errRemove != nil {
		fmt.Println("Error deleting user:", errRemove)
		return
	}
	fmt.Println("User deleted successfully")
}
