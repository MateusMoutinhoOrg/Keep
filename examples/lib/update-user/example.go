package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Writing a new value for a plain field.
//
// Update writes one field of one record. For a field that carries no index
// — anything that is not a Key — it is a single write, and the value has to
// match the Go type the schema declares for it.

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

	user, failure := users.NewItem(map[string]any{
		"email":    "mateus@gmail.com",
		"username": "mateus",
		"age":      27,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("before:", user.String())

	if failure := user.Update("age", 28); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("after:", user.String())

	// The record handed back by a lookup reads the same value: Update wrote
	// the key, it did not change a copy held in memory.
	reloaded, _ := users.FindByKey("email", "mateus@gmail.com")
	age, _ := reloaded.Get("age")
	fmt.Println("reloaded age:", age)

	// "age" is an Int field, so a string is refused before anything is
	// written.
	failure = user.Update("age", "twenty-eight")
	if failure == nil || failure.Type != api.InvalidField {
		panic("a string written to an Int field should be refused")
	}
	fmt.Println("refused:", failure.Message)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
