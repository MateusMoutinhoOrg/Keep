package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Referencing a record of another collection with a Link field.
//
// A Link is an id with a Target: the schema says which collection the id
// belongs to, so the record can be followed in one call instead of the
// caller carrying the id back to the other collection itself. It is stored
// exactly as an Int is — a bare decimal id, no index — so it costs nothing
// over the raw-id form in find-user-by-id, and it dangles the same way: a
// link to a removed record resolves to nothing, never to a different one.

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
				{Name: "author", Type: api.Link, Target: "user", Required: true},
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

	// A Link takes the record itself, or its Id — both store the same id.
	post, failure := posts.NewItem(map[string]any{
		"slug":   "storage-independent-databases",
		"author": author,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("post:", post.String())

	// GetLink reads the id and resolves it in the collection Target names,
	// with no GetSchema and no FindById at the call site.
	resolved, ok := post.GetLink("author")
	if !ok {
		panic("the author should have resolved")
	}
	fmt.Println("author:", resolved.String())

	// Get on the same field still hands back the bare id, as an int64.
	raw, failure := post.Get("author")
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Printf("author id: %v (%T)\n", raw, raw)

	// A Link may be repointed like any other field: it carries no index,
	// so two posts may name the same author.
	other, failure := users.NewItem(map[string]any{
		"email": "other@gmail.com",
		"age":   31,
	})
	if failure != nil {
		panic(failure.Message)
	}
	if failure := post.Update("author", other.Id); failure != nil {
		panic(failure.Message)
	}
	resolved, _ = post.GetLink("author")
	fmt.Println("author after update:", resolved.String())

	// Remove the author and the link stops resolving. It never resolves to
	// a different user: ids are allocated from a counter that only grows,
	// so nothing is ever handed id 2 again.
	if failure := resolved.Remove(); failure != nil {
		panic(failure.Message)
	}
	_, ok = post.GetLink("author")
	fmt.Println("author after removal:", ok)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
