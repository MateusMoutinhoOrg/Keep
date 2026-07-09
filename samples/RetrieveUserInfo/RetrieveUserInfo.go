package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

const (
	EmailToSearch = "mateus@gmail.com"
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
	users := db.GetSchema("user")

	// Create user before retrieving info (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != database.KeyConflict {
			fmt.Println("Error creating user before retrieve:", err)
			return
		}
		fmt.Println("User already exists, retrieving the existing one")
	}

	// Find the user by email
	foundUser := users.FindByKey("email", EmailToSearch)
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Retrieve and print each field individually
	email, errEmail := foundUser.Get("email")
	if errEmail != nil {
		fmt.Println("Error retrieving email:", errEmail)
		return
	}

	userName, errUsername := foundUser.Get("username")
	if errUsername != nil {
		fmt.Println("Error retrieving username:", errUsername)
		return
	}

	age, errAge := foundUser.Get("age")
	if errAge != nil {
		fmt.Println("Error retrieving age:", errAge)
		return
	}

	fmt.Println("=== User Information ===")
	fmt.Println("Email:   ", email)
	fmt.Println("UserName:", userName)
	fmt.Println("Age:     ", age)
}
