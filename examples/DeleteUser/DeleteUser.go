package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

const (
	EmailToDelete = "mateus@gmail.com"
)

var Schemas = []lib.Schema{
	{
		Name: "user",
		Itens: []lib.Item{
			{
				Type:     lib.Key,
				Required: true,
				Name:     "email",
			},
			{
				Type:     lib.Key,
				Required: true,
				Name:     "username",
			},
			{
				Name:     "age",
				Required: true,
				Type:     lib.Int,
			},

			{
				Name: "sessions",
				Type: lib.Database,
				Itens: []lib.Item{
					{
						Name:     "token",
						Type:     lib.Key,
						Required: true,
					},
					{
						Name:     "creation",
						Type:     lib.Int,
						Required: true,
					},
					{
						Name:     "expiration",
						Type:     lib.Int,
						Required: true,
					},
				},
			},
		},
	},
}

var Props = lib.Props{
	Path:    "testDatabase/",
	Schemas: Schemas,
}

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
		if err.Type != lib.KeyConflict {
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
	if errRemove.Msg != "" {
		fmt.Println("Error deleting user:", errRemove.Error())
		return
	}
	fmt.Println("User deleted successfully")
}
