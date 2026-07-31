package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	database "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
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
			{Name: "email", Type: database.Key, Required: true},
			{Name: "username", Type: database.Key, Required: true},
			{Name: "age", Type: database.Int, Required: true},
			{
				Name: "sessions",
				Type: database.Database,
				Itens: []database.Item{
					{Name: "token", Type: database.Key, Required: true},
					{Name: "creation", Type: database.Int, Required: true},
					{Name: "expiration", Type: database.Int, Required: true},
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
	deps := keepadapter.New()
	keep := keeplib.New(deps)
	db := keep.NewDatabase(Props)
	users, _ := db.GetSchema("user")

	createdUser, err := users.NewItem(map[string]any{
		"email":    EmailToInsert,
		"username": UserNameToInsert,
		"age":      AgeToInsert,
	})
	if err != nil {
		if err.Type == database.KeyConflict {
			// Second run: the unique index already holds this key
			fmt.Println("User already exists:", err.Message)
			if existing, ok := users.FindByKey("email", EmailToInsert); ok {
				fmt.Println("Existing user:", existing.String())
			}
			return
		}
		fmt.Println("Error creating user", err.Message)
		return
	}
	fmt.Println("User created successfully", createdUser.String())
}
