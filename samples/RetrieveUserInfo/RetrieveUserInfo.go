package main

import (
	"fmt"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"
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
	users, err := db.GetSchema("users")
	if err != nil {
		fmt.Println("Error getting schema:", err)
		return
	}

	// Find the user by email
	foundUser, err := users.FindByKey("email", EmailToSearch)
	if err != nil {
		fmt.Println("Error finding user:", err)
		return
	}
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Retrieve and print each field individually
	email, err := foundUser.Get("email")
	if err != nil {
		fmt.Println("Error retrieving email:", err)
		return
	}

	userName, err := foundUser.Get("username")
	if err != nil {
		fmt.Println("Error retrieving username:", err)
		return
	}

	age, err := foundUser.Get("age")
	if err != nil {
		fmt.Println("Error retrieving age:", err)
		return
	}

	fmt.Println("=== User Information ===")
	fmt.Println("Email:   ", email)
	fmt.Println("UserName:", userName)
	fmt.Println("Age:     ", age)
}
