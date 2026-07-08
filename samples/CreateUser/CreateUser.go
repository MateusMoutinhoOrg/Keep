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
				Name:     "Age",
				Required: true,
				Type:     database.Int,
			},
		},
	},
}

var Props = database.Props{
	Path:    "testDatabase/",
	Schemas: Schemas,
}

func FixIntegrity(db *database.Database, users *database.SchemaInstance, err database.Error) database.Error {
	if err.Type == database.KeyConflict {
		AlreadyExistUser := users.FindByKey(err.Key, err.KeyValue)
		UserOk := AlreadyExistUser.RequiredKeysExist()
		// if Email or Username not exist, means is a integrity Error, and these user can be removed
		if !UserOk {
			//remove these User since it was not fully created
			err := AlreadyExistUser.Remove()
			return err
		}
	}
	return err
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

	if err.Type != 0 {
		err = FixIntegrity(db, users, err)
		if err.Type != 0 {
			fmt.Println(err)
			return
		}
		// try Again
		createdUser, err = users.NewItem(map[string]any{
			"Email":    EmailToInsert,
			"UserName": UserNameToInsert,
			"Age":      AgeToInsert,
		})
		if err.Type != 0 {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("User created successfully", createdUser)
}
