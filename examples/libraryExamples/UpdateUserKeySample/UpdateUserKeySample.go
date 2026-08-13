package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	OldEmail = "mateus@gmail.com"
	NewEmail = "newmateus@gmail.com"
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

	// Create user before updating key (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    OldEmail,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != keeptypes.KeyConflict {
			fmt.Println("Error creating user before key update:", err.Message)
			return
		}
		fmt.Println("User already exists, updating the existing one")
	}

	// Find the user by the current key value. On a re-run the email was
	// already changed to NewEmail, so fall back to it and swap back.
	targetEmail := NewEmail
	foundUser, ok := users.FindByKey("email", OldEmail)
	if !ok {
		foundUser, ok = users.FindByKey("email", NewEmail)
		targetEmail = OldEmail
	}
	if !ok {
		fmt.Println("User not found")
		return
	}

	// Update an indexed field (requires re-indexing: new index entry, update value, delete old index)
	// Uses the same Update method, but internally detects that email is a Key
	// and performs the safe re-index sequence described in the documentation
	if errUpdate := foundUser.Update("email", targetEmail); errUpdate != nil {
		fmt.Println("Error updating user key", errUpdate.Message)
		return
	}
	fmt.Println("User email updated successfully to", targetEmail)
}
