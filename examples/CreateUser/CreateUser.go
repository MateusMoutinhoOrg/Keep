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

func createProps() api.Props {

	//========================Sessions==========================
	token := lib.NewKeyItem("token", true)
	creation := lib.NewIntItem("creation", true)
	expiration := lib.NewIntItem("expiration", true)
	sessions := lib.NewDatabaseItem("sessions", token, creation, expiration)

	//========================User==========================
	email := lib.NewKeyItem("email", true)
	username := lib.NewKeyItem("username", true)
	age := lib.NewIntItem("age", true)
	user := lib.NewSchema("user", email, username, age, sessions)

	//========================Props==========================
	return lib.NewProps("testDatabase/", user)
}
func main() {
	deps := standard.New()
	keep := lib.New(deps)
	props := createProps()
	db := keep.NewDatabase(props)
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
