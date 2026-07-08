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


const Schemas = []database.Schema{
		database.Schema{
			Name:"User",
			Itens: []database.Item{
				database.Item{
					Type:database.Key,
					Required:true,
					Name:"Email"
				},
				database.Item{
				Type: database.Key,
				Required:true,
				Name: "UserName"
				},
				database.Item{
				Name:"Age",
				Type:database.Int
				}

				
			}
		}

}

const Props = database.Props{
	FixIntegrity: false,
	Path: "testDatabase/",
	Schemas: Schemas,
}

func FixIntegrity(database *database.Database,err database.Error) database.Error{
	if err.Type == database.KeyConflict  {
			AlreadyExistedUser := users.findByKey(err.Key,err.KeyValue)
			UserOk  := AlreadyExistedUser.CheckKeysPresence([string]{"Email","UserName"})
            //if Email or Username not exist, means is a integrity  Error, and these user can be removed
			if !UserOk  {
				err := AlreadyExistedUser.Remove()
				return err 
			}
	}
	return err
}

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	database := keep.NewDatabase(Props)
	users  := database.GetSchema("Users")

	createdUser,err := users.newItem(map[string]any{
        "Email": EmailToInsert,
        "UserName": UserNameToInsert,
		"Age":AgeToInsert,
	})

	if err != nil {
		err = FixIntegrity(database,err)
		if err != nil {
			fmt.Println(err)
			return
		}
		//try Again
        createdUser,err := users.newItem(map[string]any{
            "Email": EmailToInsert,
            "UserName": UserNameToInsert,
            "Age":AgeToInsert,
        })
        if err != nil {
            fmt.Println(err)
            return
        }
	}	
	
	fmt.Println("User created successfully")

}
