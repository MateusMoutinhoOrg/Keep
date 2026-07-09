package main

import (
	"fmt"
	"os"

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
	// Start from a clean database directory so the sample is deterministic
	os.RemoveAll("testDatabase")

	deps := keep_deps.New()
	keep := keep_lib.New(deps)
	db := keep.NewDatabase(Props)
	users := db.GetSchema("user")

	// Create user before searching
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}

	// Find the user by email
	foundUser := users.FindByKey("email", EmailToSearch)
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Create two sessions in the user's "sessions" sub-database
	sessionsToCreate := []map[string]any{
		{"token": "token-1", "creation": 1000, "expiration": 2000},
		{"token": "token-2", "creation": 1500, "expiration": 2500},
	}
	for _, s := range sessionsToCreate {
		_, errSession := foundUser.NewSubItem("sessions", s)
		if errSession != nil {
			fmt.Println("Error creating session:", errSession)
			return
		}
	}

	sessions := foundUser.ListAll("sessions")
	for _, session := range sessions {
		token, err := session.Get("token")
		if err != nil {
			fmt.Println("Error getting token", err)
			continue
		}

		creation, err := session.Get("creation")
		if err != nil {
			fmt.Println("Error getting creation", err)
			continue
		}

		expiration, err := session.Get("expiration")
		if err != nil {
			fmt.Println("Error getting expiration", err)
			continue
		}

		fmt.Println("Token:", token)
		fmt.Println("Creation:", creation)
		fmt.Println("Expiration:", expiration)
	}

}
