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
					  Type:database.Key
                       Name:"Email"
				   },
                   database.Item{
					Type: database.Key 
					Name: "UserName"
				   },
				   database.Item{
					Name:"Age",
					Type:database.Int
				   }

				 
			 }
		   }

	}

	prop := database.Props{
		FixIntegrity: true,
		Path: "testDatabase/",
		Schemas: schema,
	}



	database := keep.NewDatabase(prop)
	users  := database.GetSchema("Users")

	createdUser,err := users.newItem(map[string]any{
        "Email": EmailToInsert,
        "UserName": UserNameToInsert,
		"Age":AgeToInsert,
	})

	if err != nil {
		fmt.Println(err)
		return
	}	
	
	fmt.Println("User created successfully")

}
