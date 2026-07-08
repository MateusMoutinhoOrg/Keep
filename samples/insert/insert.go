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

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)

	schema := []database.Schema{
           database.Schema{
			 Name:"User",
			 Itens: []database.Item{
				   database.Item{
					  Type:database.String
                       Key:true
                       Name:"Email"
				   },
                   database.Item{
					Type: database.String
					Key: true
					Name: "UserName"
				   },
				   database.Item{
					Name:"Age",
					Type:database.Int
				   }

				 
			 }
		   }

	}

	database := keep.NewDatabase(schema)
	users  := database.GetSchema("Users")

    possibleUser := users.FindByKey("email",EmailToInsert)
	if possibleUser != nil {
		fmt.Println("User already exists")
		return
	}

	possibleUser = users.FindByKey("username",UserNameToInsert)
	if possibleUser != nil {
		fmt.Println("User already exists")
		return
	}

	createdUser := users.newItem()

    err := createdUser.setValue("Email",EmailToInsert)
	if err != nil {
		//Revert changes 
		createdUser.Remove()
		fmt.Println(err)
		return
	}
	err = createdUser.setValue("Username",UserNameToInsert)
	if err != nil {
		//Revert changes 
		createdUser.Remove()	
		fmt.Println(err)
		return
	}
	err = createdUser.setValue("Age",AgeToInsert)
	if err != nil {
		//Revert changes 
		createdUser.Remove()	
		fmt.Println(err)
		return
	}

	
	
	fmt.Println("User created successfully")

}
