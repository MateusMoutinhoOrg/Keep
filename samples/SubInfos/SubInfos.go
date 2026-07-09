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

			{
				Name:"Sessions",
				Type: database.Database,
				Itens: []database.Item{
					{
						Name:"Token",
						Type: database.Key,
						Required: true,
					},
					{
						Name:"Creation",
						Type: database.Int,
						Required: true,
					},
					{
						Name: "Expiration",
						Type: database.Int,
						Required: true,
					},
				},
			}
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

	// Find the user by email
	foundUser := users.FindByKey("Email", EmailToSearch)
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Retrieve and print each field individually
	email, err := foundUser.Get("Email")
	if err != nil {
		fmt.Println("Error retrieving Email:", err)
		return
	}

	userName, err := foundUser.Get("UserName")
	if err != nil {
		fmt.Println("Error retrieving UserName:", err)
		return
	}

	age, err := foundUser.Get("Age")
	if err != nil {
		fmt.Println("Error retrieving Age:", err)
		return
	}

	fmt.Println("=== User Information ===")
	fmt.Println("Email:   ", email)
	fmt.Println("UserName:", userName)
	fmt.Println("Age:     ", age)
}
