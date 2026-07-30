package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	UserNameToInsert = "mateus"
	EmailToInsert    = "mateus@gmail.com"
	AgeToInsert      = 27
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

	createdUser, err := users.NewItem(map[string]any{
		"email":    EmailToInsert,
		"username": UserNameToInsert,
		"age":      AgeToInsert,
	})
	if err != nil {
		if err.Type() == api.KeyConflict {
			// Second run: the unique index already holds this key
			fmt.Println("User already exists:", err)
			fmt.Println("Existing user:", users.FindByKey("email", EmailToInsert))
			return
		}
		fmt.Println("Error creating user", err)
		return
	}
	fmt.Println("User created successfully", createdUser)
}
