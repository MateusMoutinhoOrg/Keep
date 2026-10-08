package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Resolving a record by its permanent id.
//
// Every record carries an ID the moment it is inserted. FindByID resolves
// one with no index read at all — it checks the back-pointer that marks a
// record live and hands the record back — so it costs less than FindByKey
// and works for a collection that declares no Key field.
//
// An id is allocated from a counter that only grows, so it is never handed
// out twice. That is what makes it safe to hold on to: see link-records for
// a field that stores one and follows it.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "bio", Type: api.String},
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

	mateus, failure := users.Insert(map[string]any{
		"email": "mateus@gmail.com",
		"bio":   "writes databases",
		"age":   27,
	})
	if failure != nil {
		panic(failure.Message)
	}
	ana, failure := users.Insert(map[string]any{
		"email": "ana@gmail.com",
		"bio":   "writes databases",
		"age":   31,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("ids:", mateus.ID, ana.ID)

	// The id is all a lookup needs. Nothing about the record's values is
	// read to find it, so a field no index covers — "bio" is a String —
	// still comes back with it.
	resolved, ok, failure := users.FindByID(mateus.ID)
	if failure != nil || !ok {
		panic("the record should have resolved")
	}
	bio, _ := resolved.Get("bio")
	fmt.Printf("by id %d: %s (bio: %v)\n", mateus.ID, resolved.String(), bio)

	// An id that was never allocated resolves to nothing.
	_, ok, _ = users.FindByID(99)
	fmt.Println("id 99:", ok)

	// Remove the record and its id stops resolving — permanently. The
	// counter is never rewound, so the next insert is id 3, not id 1, and
	// id 1 resolves to nothing rather than to whoever came after.
	if failure := resolved.Remove(); failure != nil {
		panic(failure.Message)
	}
	_, ok, _ = users.FindByID(mateus.ID)
	fmt.Println("id 1 after removal:", ok)

	fresh, failure := users.Insert(map[string]any{
		"email": "bruno@gmail.com",
		"age":   22,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("next id allocated:", fresh.ID)
	_, ok, _ = users.FindByID(mateus.ID)
	fmt.Println("id 1 still:", ok)

	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
