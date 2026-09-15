package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Looking a record up through a unique key.
//
// Every Key field carries a unique index, so FindByKey is a single read
// whatever the size of the collection — the value is hashed and the hash is
// the key the index lives under. Lookups are case-insensitive because the
// value is lower-cased before it is hashed.

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

	for _, fields := range []map[string]any{
		{"email": "mateus@gmail.com", "username": "mateus", "age": 27},
		{"email": "ana@gmail.com", "username": "ana", "age": 31},
	} {
		if _, failure := users.NewItem(fields); failure != nil {
			panic(failure.Message)
		}
	}

	found, ok := users.FindByKey("email", "mateus@gmail.com")
	if !ok {
		panic("mateus@gmail.com should have been found")
	}
	fmt.Println("by email:", found.String())

	// Any Key field of the schema indexes its own values.
	found, ok = users.FindByKey("username", "ana")
	if !ok {
		panic("ana should have been found")
	}
	fmt.Println("by username:", found.String())

	// The index is case-insensitive.
	found, ok = users.FindByKey("email", "MATEUS@GMAIL.COM")
	fmt.Println("by upper-cased email:", ok, found.Id)

	// A value nobody holds is a miss, not a failure.
	_, ok = users.FindByKey("email", "nobody@gmail.com")
	fmt.Println("unknown email found:", ok)

	// So is a field that carries no index: only a Key field does.
	_, ok = users.FindByKey("age", 27)
	fmt.Println("non-key field found:", ok)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
