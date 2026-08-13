//go:build ignore

// This file is an illustrative sample, not part of the build.
package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

var Schemas = []keeptypes.Schema{
	{
		Name: "user",
		Itens: []keeptypes.Item{
			{Name: "email", Type: keeptypes.Key, Required: true},
			{Name: "age", Type: keeptypes.Int, Required: true},
		},
	},
}

var Props = keeptypes.Props{
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
