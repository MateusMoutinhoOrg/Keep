package main

import (
	"fmt"
	"os"

	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/database"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

const (
	EmailToDelete = "mateus@gmail.com"
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

	// Create the user first before deleting
	_, err := users.NewItem(map[string]any{
		"email":    EmailToDelete,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		fmt.Println("Error creating user before delete:", err)
		return
	}

	// First, find the user by key
	foundUser := users.FindByKey("email", EmailToDelete)
	if foundUser == nil {
		fmt.Println("User not found")
		return
	}

	// Then, remove the user (swap-with-last deletion)
	errRemove := foundUser.Remove()
	if errRemove.Msg != "" {
		fmt.Println("Error deleting user:", errRemove.Error())
		return
	}
	fmt.Println("User deleted successfully")
}
