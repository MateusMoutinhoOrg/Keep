package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

const (
	UserNameToInsert = "mateus"
	EmailToInsert    = "mateus@gmail.com"
	AgeToInsert      = 27
)

var Schemas = []database.Schema{
	{
		Name: "user",
		Itens: []database.Item{
			{
				Type:     database.Key,
				Required: true,
				Name:     "email",
			},
			{
				Type:     database.Key,
				Required: true,
				Name:     "username",
			},
			{
				Name:     "age",
				Required: true,
				Type:     database.Int,
			},

			{
				Name: "sessions",
				Type: database.Database,
				Itens: []database.Item{
					{
						Name:     "token",
						Type:     database.Key,
						Required: true,
					},
					{
						Name:     "creation",
						Type:     database.Int,
						Required: true,
					},
					{
						Name:     "expiration",
						Type:     database.Int,
						Required: true,
					},
				},
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
	users, err := db.GetSchema("users")
	if err != nil {
		fmt.Println("Error getting schema:", err)
		return
	}

	createdUser, err := users.NewItem(map[string]any{
		"email":    EmailToInsert,
		"username": UserNameToInsert,
		"age":      AgeToInsert,
	})
	if err != nil {
		fmt.Println("Error creating user", err)
		return
	}
	fmt.Println("User created successfully", createdUser)
}
