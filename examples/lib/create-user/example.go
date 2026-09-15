package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Inserting a record, and what happens when a unique key is already taken.
//
// A database is described by a value: Props names the prefix every key is
// written under and the collections it holds. Nothing is created until the
// first record is written, so building the handle is free.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "TestDir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Itens: []api.Item{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "username", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
			},
		},
	},
}

func main() {

	deps := standard.New()    // filestorage, hashdeps, std, stringsdeps
	lib := sandbox.New(&deps) // *api.Sandbox

	db := lib.Databases.New(Props)
	users, ok := db.GetSchema("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}

	created, failure := users.NewItem(map[string]any{
		"email":    "mateus@gmail.com",
		"username": "mateus",
		"age":      27,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("created:", created.String())

	// "email" and "username" are Key fields, so the value of each is unique
	// across every live record of the collection. Inserting the same email
	// again is refused before anything is written.
	_, failure = users.NewItem(map[string]any{
		"email":    "mateus@gmail.com",
		"username": "other",
		"age":      31,
	})
	switch {
	case failure == nil:
		panic("the duplicate email should have been refused")
	case failure.Type == api.KeyConflict:
		fmt.Printf("refused: %s (field %q, value %v)\n", failure.Message, failure.Key, failure.KeyValue)
	default:
		panic(failure.Message)
	}

	// A required field left out is refused the same way, with its own Type.
	_, failure = users.NewItem(map[string]any{"email": "other@gmail.com"})
	if failure == nil || failure.Type != api.MissingField {
		panic("the missing field should have been refused")
	}
	fmt.Printf("refused: %s\n", failure.Message)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
