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

	// Create user before searching (skip if it already exists from a previous run)
	_, err := users.NewItem(map[string]any{
		"email":    EmailToSearch,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != database.KeyConflict {
			fmt.Println("Error creating user:", err.Message)
			return
		}
		fmt.Println("User already exists, reusing it")
	}

	// Find the user by email
	foundUser, ok := users.FindByKey("email", EmailToSearch)
	if !ok {
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
			if errSession.Type == database.KeyConflict {
				// Already created by a previous run, keep going
				fmt.Printf("Session %v already exists, skipping\n", s["token"])
				continue
			}
			fmt.Println("Error creating session:", errSession.Message)
			return
		}
	}

	sessions := foundUser.ListAll("sessions")
	for _, session := range sessions {
		token, err := session.Get("token")
		if err != nil {
			fmt.Println("Error getting token", err.Message)
			continue
		}

		creation, err := session.Get("creation")
		if err != nil {
			fmt.Println("Error getting creation", err.Message)
			continue
		}

		expiration, err := session.Get("expiration")
		if err != nil {
			fmt.Println("Error getting expiration", err.Message)
			continue
		}

		fmt.Println("Token:", token)
		fmt.Println("Creation:", creation)
		fmt.Println("Expiration:", expiration)
	}

}
