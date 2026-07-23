package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

const (
	EmailToSearch = "mateus@gmail.com"
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

	// Create user before retrieving info (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != lib.KeyConflict {
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
