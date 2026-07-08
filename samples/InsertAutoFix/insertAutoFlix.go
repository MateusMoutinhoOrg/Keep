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

var Schemas = []database.Schema{
	{
		Name: "User",
		Itens: []database.Item{
			{
				Type:     database.Key,
				Required: true,
				Name:     "Email",
			},
			{
				Type:     database.Key,
				Required: true,
				Name:     "UserName",
			},
			{
				Name: "Age",
				Type: database.Int,
			},
		},
	},
}

var Props = database.Props{
	FixIntegrity: true,
	Path:         "testDatabase/",
	Schemas:      Schemas,
}

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("Users")

	createdUser, err := users.NewItem(map[string]any{
		"Email":    EmailToInsert,
		"UserName": UserNameToInsert,
		"Age":      AgeToInsert,
	})

	//the lib iteself fix Integrity Brokes
	if err.Type != 0 {
		fmt.Println(err)
		return
	}

	fmt.Println("User created successfully", createdUser)
}
