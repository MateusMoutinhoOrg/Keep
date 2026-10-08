package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Writing a new value for a plain field.
//
// Update writes one field of one record. For a field that carries no index
// — an Int, a Float, a String, a Link: anything that is not a Key — it is a
// single write, and the value has to match the Go type the schema declares
// for it. A nil value clears a field the schema does not require.

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

	db, failure := lib.Databases.New(Props)
	if failure != nil {
		panic(failure.Message)
	}
	users, ok := db.Collection("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}

	user, failure := users.Insert(map[string]any{
		"email":    "mateus@gmail.com",
		"username": "mateus",
		"age":      27,
		"height":   1.82,
		"bio":      "writes databases",
		"nickname": "teu",
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("before:", user.String())

	// Each of the three is one write, whatever the type.
	if failure := user.Update("age", 28); failure != nil {
		panic(failure.Message)
	}
	if failure := user.Update("height", 1.83); failure != nil {
		panic(failure.Message)
	}
	if failure := user.Update("bio", "writes storage-independent databases"); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("after:", user.String())

	// The record handed back by a lookup reads the same value: Update wrote
	// the key, it did not change a copy held in memory.
	reloaded, _, _ := users.FindByKey("email", "mateus@gmail.com")
	age, _ := reloaded.Get("age")
	fmt.Println("reloaded age:", age)

	// Each type refuses what it cannot hold, before anything is written.
	failure = user.Update("age", "twenty-eight")
	if failure == nil || failure.Type != api.InvalidField {
		panic("a string written to an Int field should be refused")
	}
	fmt.Println("refused:", failure.Message)

	failure = user.Update("height", "tall")
	if failure == nil || failure.Type != api.InvalidField {
		panic("a string written to a Float field should be refused")
	}
	fmt.Println("refused:", failure.Message)

	failure = user.Update("bio", 42)
	if failure == nil || failure.Type != api.InvalidField {
		panic("an int written to a String field should be refused")
	}
	fmt.Println("refused:", failure.Message)

	// A Float field takes a whole number too, and stores it as one.
	if failure := user.Update("height", 2); failure != nil {
		panic(failure.Message)
	}
	height, _ := user.Get("height")
	fmt.Printf("height: %v (%T)\n", height, height)

	// An Int field takes any Go integer, and a float holding a whole number
	// — which is what encoding/json hands back for every number.
	if failure := user.Update("age", float64(29)); failure != nil {
		panic(failure.Message)
	}
	age, _ = user.Get("age")
	fmt.Printf("age from a float64: %v (%T)\n", age, age)

	// nil clears a field the schema does not require...
	if failure := user.Update("nickname", nil); failure != nil {
		panic(failure.Message)
	}
	_, failure = user.Get("nickname")
	fmt.Println("nickname after clearing:", failure.Message)

	// ...and is refused for one it does.
	failure = user.Update("bio", nil)
	if failure == nil || failure.Type != api.MissingField {
		panic("clearing a required field should be refused")
	}
	fmt.Println("refused:", failure.Message)

	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
