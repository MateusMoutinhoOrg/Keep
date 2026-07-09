package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

const (
	OldEmail = "mateus@gmail.com"
	NewEmail = "newmateus@gmail.com"
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
	users := db.GetSchema("users")

	// Create user before updating key
	_, err := users.NewItem(map[string]any{
		"email":    OldEmail,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		fmt.Println("Error creating user before key update:", err)
		return
	}

	// Find the user by the current key value
	foundUser := users.FindByKey("email", OldEmail)
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Update an indexed field (requires re-indexing: new index entry, update value, delete old index)
	// Uses the same Update method, but internally detects that email is a Key
	// and performs the safe re-index sequence described in the documentation
	errUpdate := foundUser.Update("email", NewEmail)
	if errUpdate != nil {
		fmt.Println("Error updating user key", errUpdate)
		return
	}
	fmt.Println("User email updated successfully")
}
