package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Reading fields back off a record.
//
// Get returns the value in the Go type the field's schema declares: a
// string for a Key or String field, an int64 for an Int or Link field, a
// float64 for a Float one. Asking for a field the schema does not declare
// is a failure, not a nil.

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
				{Name: "height", Type: api.Float, Required: true},
				{Name: "bio", Type: api.String, Required: true},
				{Name: "nickname", Type: api.String},
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

	user, failure := users.NewItem(map[string]any{
		"email":    "mateus@gmail.com",
		"username": "mateus",
		"age":      27,
		"height":   1.82,
		"bio":      "writes databases",
	})
	if failure != nil {
		panic(failure.Message)
	}

	// One field per type the schema can declare, each printed with the Go
	// type Get handed back.
	for _, field := range []string{"email", "username", "age", "height", "bio"} {
		value, failure := user.Get(field)
		if failure != nil {
			panic(failure.Message)
		}
		fmt.Printf("%s = %v (%T)\n", field, value, value)
	}

	// "nickname" is declared but not required, and this record never set
	// it: the field exists in the schema and holds nothing. It is a String
	// rather than a Key, so leaving it unset costs no index entry either.
	_, failure = user.Get("nickname")
	if failure == nil || failure.Type != api.NotFound {
		panic("an unset field should report NotFound")
	}
	fmt.Println("nickname:", failure.Message)

	// A field the schema does not declare at all is a different failure.
	_, failure = user.Get("phone")
	if failure == nil || failure.Type != api.InvalidField {
		panic("an undeclared field should report InvalidField")
	}
	fmt.Println("phone:", failure.Message)

	// CheckKeysPresence answers the same question for several fields at
	// once, without reading any value.
	fmt.Println("has email and age:", user.CheckKeysPresence([]string{"email", "age"}))
	fmt.Println("has email and nickname:", user.CheckKeysPresence([]string{"email", "nickname"}))

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
