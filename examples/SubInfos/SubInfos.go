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

	// Create user before searching (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != lib.KeyConflict {
			fmt.Println("Error creating user:", err)
			return
		}
		fmt.Println("User already exists, reusing it")
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
			if errSession.Type == lib.KeyConflict {
				// Already created by a previous run, keep going
				fmt.Printf("Session %v already exists, skipping\n", s["token"])
				continue
			}
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
