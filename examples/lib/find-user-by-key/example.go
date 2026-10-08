package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Looking a record up through a unique key.
//
// Every Key field carries a unique index, so FindByKey costs the same
// whatever the size of the collection — the value is hashed and the hash is
// the key the index lives under. Lookups ignore case because the value is
// case-folded before it is hashed.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
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

	db, failure := lib.Databases.New(Props)
	if failure != nil {
		panic(failure.Message)
	}
	users, ok := db.Collection("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}

	for _, fields := range []map[string]any{
		{"email": "mateus@gmail.com", "username": "mateus", "age": 27},
		{"email": "ana@gmail.com", "username": "ana", "age": 31},
	} {
		if _, failure := users.Insert(fields); failure != nil {
			panic(failure.Message)
		}
	}

	found, ok, failure := users.FindByKey("email", "mateus@gmail.com")
	if failure != nil {
		panic(failure.Message)
	}
	if !ok {
		panic("mateus@gmail.com should have been found")
	}
	fmt.Println("by email:", found.String())

	// Any Key field of the schema indexes its own values.
	found, ok, failure = users.FindByKey("username", "ana")
	if failure != nil || !ok {
		panic("ana should have been found")
	}
	fmt.Println("by username:", found.String())

	// The index ignores case.
	found, ok, _ = users.FindByKey("email", "MATEUS@GMAIL.COM")
	fmt.Println("by upper-cased email:", ok, found.ID)

	// A value nobody holds is a miss, not a failure: ok is false and the
	// failure is nil.
	_, ok, failure = users.FindByKey("email", "nobody@gmail.com")
	fmt.Println("unknown email found:", ok, "failure:", failure)

	// A field that carries no index is a mistake in the call, so it is a
	// failure: only a Key field can be looked up.
	_, _, failure = users.FindByKey("age", 27)
	fmt.Println("non-key field refused:", failure != nil && failure.Type == api.InvalidField)

	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
