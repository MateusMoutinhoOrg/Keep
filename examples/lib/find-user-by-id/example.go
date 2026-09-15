package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Pointing one record at another through its id.
//
// Every record carries a permanent Id, and an id is never reused. Storing
// one in an Int field of another collection is how a record references a
// record: FindById resolves it in a single read, and a reference to a
// record that has been removed resolves to nothing rather than to whatever
// took its place.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "TestDir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Itens: []api.Item{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
			},
		},
		{
			Name: "post",
			Itens: []api.Item{
				{Name: "slug", Type: api.Key, Required: true},
				{Name: "author", Type: api.Int, Required: true}, // a user id
			},
		},
	},
}

func main() {

	deps := standard.New()
	lib := sandbox.New(&deps)

	db := lib.Databases.New(Props)
	users, _ := db.GetSchema("user")
	posts, _ := db.GetSchema("post")

	author, failure := users.NewItem(map[string]any{
		"email": "mateus@gmail.com",
		"age":   27,
	})
	if failure != nil {
		panic(failure.Message)
	}

	post, failure := posts.NewItem(map[string]any{
		"slug":   "storage-independent-databases",
		"author": author.Id,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("post:", post.String())

	// Read the reference back and resolve it.
	authorId, failure := post.Get("author")
	if failure != nil {
		panic(failure.Message)
	}
	resolved, ok := users.FindById(authorId.(int64))
	if !ok {
		panic("the author should have resolved")
	}
	fmt.Println("author:", resolved.String())

	// Remove the author, and the reference stops resolving. It never
	// resolves to a different user: ids are allocated from a counter that
	// only grows, so nothing is ever handed id 1 again.
	if failure := resolved.Remove(); failure != nil {
		panic(failure.Message)
	}
	_, ok = users.FindById(authorId.(int64))
	fmt.Println("author after removal:", ok)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
