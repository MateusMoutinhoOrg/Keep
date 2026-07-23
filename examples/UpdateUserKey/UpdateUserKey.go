package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

const (
	OldEmail = "mateus@gmail.com"
	NewEmail = "newmateus@gmail.com"
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

	// Create user before updating key (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    OldEmail,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != lib.KeyConflict {
			fmt.Println("Error creating user before key update:", err)
			return
		}
		fmt.Println("User already exists, updating the existing one")
	}

	// Find the user by the current key value. On a re-run the email was
	// already changed to NewEmail, so fall back to it and swap back.
	targetEmail := NewEmail
	foundUser := users.FindByKey("email", OldEmail)
	if foundUser == nil {
		foundUser = users.FindByKey("email", NewEmail)
		targetEmail = OldEmail
	}
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Update an indexed field (requires re-indexing: new index entry, update value, delete old index)
	// Uses the same Update method, but internally detects that email is a Key
	// and performs the safe re-index sequence described in the documentation
	errUpdate := foundUser.Update("email", targetEmail)
	if errUpdate != nil {
		fmt.Println("Error updating user key", errUpdate)
		return
	}
	fmt.Println("User email updated successfully to", targetEmail)
}
