//go:build ignore

// This file is an illustrative sample, not part of the build.
package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

var Props = lib.Props{
	Path: "testDatabase/",
	Schemas: []lib.Schema{
		{
			Name: "user",
			Itens: []lib.Item{
				{Name: "email", Type: lib.Key, Required: true},
				{Name: "age", Type: lib.Int, Required: true},
			},
		},
	},
}

func main() {
	// 1. Build deps through an adapter (the opinionated layer).
	deps := standard.New()

	// 2. Inject deps into the pure library.
	keep := lib.New(deps)

	// 3. Exercise the library — it never knows which adapter is behind it.
	db := keep.NewDatabase(Props)
	users := db.GetSchema("user")

	created, err := users.NewItem(map[string]any{"email": "a@x.com", "age": 30})
	if err != nil {
		fmt.Println("error creating user:", err)
		return
	}
	fmt.Println("created:", created)
}
