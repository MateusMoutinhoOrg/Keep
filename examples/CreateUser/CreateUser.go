package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

const (
	UserNameToInsert = "mateus"
	EmailToInsert    = "mateus@gmail.com"
	AgeToInsert      = 27
)

var Schemas = []lib.Schema{
	{
		Name: "user",
		Itens: []lib.Item{
			{
				Type:     lib.Key,
				Required: true,
				Name:     "email",
			},
			{
				Type:     lib.Key,
				Required: true,
				Name:     "username",
			},
			{
				Name:     "age",
				Required: true,
				Type:     lib.Int,
			},

			{
				Name: "sessions",
				Type: lib.Database,
				Itens: []lib.Item{
					{
						Name:     "token",
						Type:     lib.Key,
						Required: true,
					},
					{
						Name:     "creation",
						Type:     lib.Int,
						Required: true,
					},
					{
						Name:     "expiration",
						Type:     lib.Int,
						Required: true,
					},
				},
			},
		},
	},
}

var Props = lib.Props{
	Path:    "testDatabase/",
	Schemas: Schemas,
}

func main() {
	deps := standard.New()
	keep := lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("user")

	createdUser, err := users.NewItem(map[string]any{
		"email":    EmailToInsert,
		"username": UserNameToInsert,
		"age":      AgeToInsert,
	})
	if err != nil {
		if err.Type == lib.KeyConflict {
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
