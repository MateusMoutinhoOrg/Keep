package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

const (
	EmailToSearch = "mateus@gmail.com"
	NewAge        = 28
)

var Schemas = []database.Schema{
	{
		Name: "user",
		Itens: []database.Item{
			{
				Type:     database.Key,
				Required: true,
				Name:     "email",
			},
			{
				Type:     database.Key,
				Required: true,
				Name:     "username",
			},
			{
				Name:     "age",
				Required: true,
				Type:     database.Int,
			},

			{
				Name: "sessions",
				Type: database.Database,
				Itens: []database.Item{
					{
						Name:     "token",
						Type:     database.Key,
						Required: true,
					},
					{
						Name:     "creation",
						Type:     database.Int,
						Required: true,
					},
					{
						Name:     "expiration",
						Type:     database.Int,
						Required: true,
					},
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
	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users, err := db.GetSchema("users")
	if err != nil {
		fmt.Println("Error getting schema:", err)
		return
	}

	// Find the user to update
	foundUser, err := users.FindByKey("email", EmailToSearch)
	if err != nil {
		fmt.Println("Error finding user:", err)
		return
	}
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Update a non-indexed field (simple single key write)
	err = foundUser.Update("age", NewAge)
	if err != nil {
		fmt.Println("Error updating user:", err)
		return
	}
	fmt.Println("User updated successfully")
}
