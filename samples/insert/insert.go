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
	


}
