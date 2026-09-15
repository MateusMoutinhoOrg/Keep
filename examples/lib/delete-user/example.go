package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Removing a record, and what goes with it.
//
// Remove costs the same whatever the size of the collection: the record at
// the last position moves into the hole the removed one leaves, so the
// position list stays dense with no shifting. That is also why list order
// is not stable across removals.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "TestDir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Itens: []api.Item{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
				{
					Name: "sessions",
					Type: api.Database,
					Itens: []api.Item{
						{Name: "token", Type: api.Key, Required: true},
					},
				},
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

	for _, fields := range []map[string]any{
		{"email": "mateus@gmail.com", "age": 27},
		{"email": "ana@gmail.com", "age": 31},
		{"email": "bruno@gmail.com", "age": 44},
	} {
		if _, failure := users.NewItem(fields); failure != nil {
			panic(failure.Message)
		}
	}

	mateus, _ := users.FindByKey("email", "mateus@gmail.com")
	if _, failure := mateus.NewSubItem("sessions", map[string]any{"token": "token-1"}); failure != nil {
		panic(failure.Message)
	}

	// Removing the record removes its values, its index entries and every
	// record of every collection nested under it.
	if failure := mateus.Remove(); failure != nil {
		panic(failure.Message)
	}

	_, ok = users.FindByKey("email", "mateus@gmail.com")
	fmt.Println("found after removal:", ok)

	// Removing a record that is already gone is not a failure.
	if failure := mateus.Remove(); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("removing twice: no error")

	remaining, failure := users.ListAll()
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("remaining:", len(remaining))
	for _, user := range remaining {
		fmt.Println(" -", user.String())
	}

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
