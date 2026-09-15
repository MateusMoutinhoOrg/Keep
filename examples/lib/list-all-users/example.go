package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Walking every record of a collection.
//
// The backend offers no way to list keys, scan a prefix or query a range —
// and iteration still works, because the collection keeps its own position
// list: positions 1..size, with no gap, each holding one record id. ListAll
// reads that list.

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

	deps := standard.New()
	lib := sandbox.New(&deps)

	db := lib.Databases.New(Props)
	users, ok := db.GetSchema("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}

	// An empty collection lists as an empty slice, not as a failure.
	empty, failure := users.ListAll()
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("before any insert:", len(empty))

	for _, fields := range []map[string]any{
		{"email": "mateus@gmail.com", "username": "mateus", "age": 27},
		{"email": "ana@gmail.com", "username": "ana", "age": 31},
		{"email": "bruno@gmail.com", "username": "bruno", "age": 44},
	} {
		if _, failure := users.NewItem(fields); failure != nil {
			panic(failure.Message)
		}
	}

	all, failure := users.ListAll()
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("users:", len(all))
	for _, user := range all {
		fmt.Println(" -", user.String())
	}

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
