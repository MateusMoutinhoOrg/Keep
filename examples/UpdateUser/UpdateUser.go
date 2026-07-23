package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

const (
	EmailToSearch = "mateus@gmail.com"
	NewAge        = 28
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

	// Create user before updating (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != lib.KeyConflict {
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
