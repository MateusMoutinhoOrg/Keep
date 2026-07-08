package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

const (
	UserNameToInsert = "mateus"
	PasswordToInsert = "12345"
	EmailToInsert    = "mateus@gmail.com"
	AgeToInsert      = 27
)

var Props = database.Props{
	Path: "testDatabase/",
}

func remove_user(user database.Instance) {
	//revert the transactions
	user.RemoveKey("Password")
	user.RemoveKey("Username")
	user.RemoveKey("Email")
	user.Remove()
}

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.Table("Users")

	// check for Already Existent user
	possibleUser := users.FindByKey("Username", UserNameToInsert)
	if possibleUser != nil {
		fmt.Print("User already exist")
		return
	}
	possibleUser = users.FindByKey("Email", EmailToInsert)
	if possibleUser != nil {
		fmt.Print("User already exist")
		return
	}

	createdUser := users.newItem()

	err := createdUser.SetUniqueItem("Email", EmailToInsert)
	if err != nil {
		remove_user(createdUser)
	}

	err = createdUser.SetUniqueItem("Username", UserNameToInsert)
	if err != nil {
		remove_user(createdUser)
	}

	err = createdUser.SetUniqueItem("Password", PasswordToInsert)
	if err != nil {
		remove_user(createdUser)
	}

	err = createdUser.SetUniqueItem("Age", AgeToInsert)
	if err != nil {
		remove_user(createdUser)
	}

	fmt.Println("User created successfully", createdUser)
}
