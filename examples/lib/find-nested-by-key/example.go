package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Looking a record of a nested collection up by key.
//
// Record.Nested hands a nested field out as a Collection of its own, rooted
// at the record that owns it: Insert, FindByKey, FindByID, List and Repair
// work on it exactly as on a top-level collection. Its unique index is its
// own, so a session is found by its token in one lookup, with no walk over
// the sessions of the user.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
				{
					Name: "sessions",
					Type: api.Nested,
					Fields: []api.Field{
						{Name: "token", Type: api.Key, Required: true},
						{Name: "agent", Type: api.String},
					},
				},
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

	mateus, failure := users.Insert(map[string]any{"email": "mateus@gmail.com"})
	if failure != nil {
		panic(failure.Message)
	}
	ana, failure := users.Insert(map[string]any{"email": "ana@gmail.com"})
	if failure != nil {
		panic(failure.Message)
	}

	sessions, failure := mateus.Nested("sessions")
	if failure != nil {
		panic(failure.Message)
	}
	for _, session := range []map[string]any{
		{"token": "token-1", "agent": "firefox"},
		{"token": "token-2", "agent": "curl"},
	} {
		if _, failure := sessions.Insert(session); failure != nil {
			panic(failure.Message)
		}
	}

	// One lookup in mateus's own index, ignoring case like any Key.
	found, ok, failure := sessions.FindByKey("token", "TOKEN-2")
	if failure != nil || !ok {
		panic("token-2 should have been found")
	}
	fmt.Println("by token:", found.String())

	// Ids, pages and lookups by id work the same way.
	byID, ok, _ := sessions.FindByID(found.ID)
	fmt.Println("by id:", ok, byID.String())
	page, failure := sessions.List(1, 1)
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("first page:", len(page), page[0].String())

	// Ana's sessions are another collection, with another index.
	anaSessions, failure := ana.Nested("sessions")
	if failure != nil {
		panic(failure.Message)
	}
	_, ok, _ = anaSessions.FindByKey("token", "token-1")
	fmt.Println("token-1 under ana:", ok)

	// Only a Nested field is a collection.
	_, failure = mateus.Nested("email")
	fmt.Println("Nested on a plain field refused:", failure != nil && failure.Type == api.InvalidField)

	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
