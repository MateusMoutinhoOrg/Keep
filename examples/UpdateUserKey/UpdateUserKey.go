package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	OldEmail = "mateus@gmail.com"
	NewEmail = "newmateus@gmail.com"
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

	// Create user before updating key (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    OldEmail,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type() != api.KeyConflict {
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
