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
