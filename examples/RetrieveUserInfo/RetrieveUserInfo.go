package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	database "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	EmailToSearch = "mateus@gmail.com"
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

	// Create user before retrieving info (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != database.KeyConflict {
			fmt.Println("Error creating user before retrieve:", err.Message)
			return
		}
		fmt.Println("User already exists, retrieving the existing one")
	}

	// Find the user by email
	foundUser, ok := users.FindByKey("email", EmailToSearch)
	if !ok {
		fmt.Println("User not found")
		return
	}

	// Retrieve and print each field individually
	email, errEmail := foundUser.Get("email")
	if errEmail != nil {
		fmt.Println("Error retrieving email:", errEmail.Message)
		return
	}

	userName, errUsername := foundUser.Get("username")
	if errUsername != nil {
		fmt.Println("Error retrieving username:", errUsername.Message)
		return
	}

	age, errAge := foundUser.Get("age")
	if errAge != nil {
		fmt.Println("Error retrieving age:", errAge.Message)
		return
	}

	fmt.Println("=== User Information ===")
	fmt.Println("Email:   ", email)
	fmt.Println("UserName:", userName)
	fmt.Println("Age:     ", age)
}
