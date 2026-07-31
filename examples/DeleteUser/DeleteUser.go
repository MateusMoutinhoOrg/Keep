package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	database "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	EmailToDelete = "mateus@gmail.com"
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

	// Create the user first before deleting (skip if it survived a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToDelete,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != database.KeyConflict {
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
