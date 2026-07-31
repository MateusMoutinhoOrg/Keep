//go:build ignore

// This file is an illustrative sample, not part of the build.
package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	database "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

var Schemas = []database.Schema{
	{
		Name: "user",
		Itens: []database.Item{
			{Name: "email", Type: database.Key, Required: true},
			{Name: "age", Type: database.Int, Required: true},
		},
	},
}

var Props = database.Props{
	Path:    "testDatabase/",
	Schemas: Schemas,
}

func main() {
	// 1. Build deps through an adapter (the opinionated layer).
	deps := keepadapter.New()

	// 2. Inject deps into the pure library.
	keep := keeplib.New(deps)

	// 3. Exercise the library — it never knows which adapter is behind it.
	db := keep.NewDatabase(Props)
	users, _ := db.GetSchema("user")

	created, err := users.NewItem(map[string]any{"email": "a@x.com", "age": 30})
	if err != nil {
		fmt.Println("error creating user:", err.Message)
		return
	}
	fmt.Println("created:", created.String())
}
