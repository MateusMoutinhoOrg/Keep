package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
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

func main() {
	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("Users")

	// Iterate all records in the dense list (positions 1..size)
	allUsers, err := users.ListAll()
	if err != nil {
		fmt.Println("Error listing users", err)
		return
	}
	for _, user := range allUsers {
		fmt.Println("User:", user)
	}
}
