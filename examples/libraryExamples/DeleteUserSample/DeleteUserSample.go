package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	EmailToDelete = "mateus@gmail.com"
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

	// Create the user first before deleting (skip if it survived a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToDelete,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != keeptypes.KeyConflict {
			fmt.Println("Error creating user before delete:", err.Message)
			return
		}
		fmt.Println("User already exists, deleting the existing one")
	}

	// First, find the user by key
	foundUser, ok := users.FindByKey("email", EmailToDelete)
	if !ok {
		fmt.Println("User not found")
		return
	}

	// Then, remove the user (swap-with-last deletion)
	if errRemove := foundUser.Remove(); errRemove != nil {
		fmt.Println("Error deleting user:", errRemove.Message)
		return
	}
	fmt.Println("User deleted successfully")
}
