//go:build ignore

// This file is an illustrative sample, not part of the build.
package main

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	"github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

func createProps() api.Props {

	//========================User==========================
	email := lib.NewKeyItem("email", true)
	age := lib.NewIntItem("age", true)
	user := lib.NewSchema("user", email, age)

	//========================Props==========================
	return lib.NewProps("testDatabase/", user)
}

func main() {
	// 1. Build deps through an adapter (the opinionated layer).
	deps := standard.New()

	// 2. Inject deps into the pure library.
	keep := lib.New(deps)

	// 3. Exercise the library — it never knows which adapter is behind it.
	props := createProps()
	db := keep.NewDatabase(props)
	users := db.GetSchema("user")

	created, err := users.NewItem(map[string]any{"email": "a@x.com", "age": 30})
	if err != nil {
		fmt.Println("error creating user:", err)
		return
	}
	fmt.Println("created:", created)
}
