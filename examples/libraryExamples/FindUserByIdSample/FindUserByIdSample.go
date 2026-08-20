package main

import (
	"fmt"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keeptypes "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

const (
	AuthorEmail = "mateus@gmail.com"
	PostTitle   = "keep-by-id"
)

// Two collections pointing at each other: a post stores the id of the
// user that wrote it in a plain Int field, and that id is resolved back
// into the record with users.FindById.
var Schemas = []keeptypes.Schema{
	{
		Name: "user",
		Itens: []keeptypes.Item{
			{Name: "email", Type: keeptypes.Key, Required: true},
			{Name: "username", Type: keeptypes.Key, Required: true},
			{Name: "age", Type: keeptypes.Int, Required: true},
		},
	},
	{
		Name: "post",
		Itens: []keeptypes.Item{
			{Name: "title", Type: keeptypes.Key, Required: true},
			// The foreign key: it holds a user record's permanent Id.
			{Name: "author", Type: keeptypes.Int, Required: true},
		},
	},
}

var Props = keeptypes.Props{
	Path:    "testDatabase/",
	Schemas: Schemas,
}

func main() {
	deps := keepadapter.New()
	keep := keeplib.New(deps)
	db := keep.NewDatabase(Props)
	users, _ := db.GetSchema("user")
	posts, _ := db.GetSchema("post")

	// Create the author, or reuse the one a previous run left behind.
	author, err := users.NewItem(map[string]any{
		"email":    AuthorEmail,
		"username": "mateus",
		"age":      27,
	})
	if err != nil {
		if err.Type != keeptypes.KeyConflict {
			fmt.Println("Error creating the author:", err.Message)
			return
		}
		author, _ = users.FindByKey("email", AuthorEmail)
	}

	// Point the post at the author by storing the author's id.
	post, err := posts.NewItem(map[string]any{
		"title":  PostTitle,
		"author": author.Id,
	})
	if err != nil {
		if err.Type != keeptypes.KeyConflict {
			fmt.Println("Error creating the post:", err.Message)
			return
		}
		post, _ = posts.FindByKey("title", PostTitle)
	}

	// Follow the pointer: read the stored id, then resolve it.
	authorId, err := post.Get("author")
	if err != nil {
		fmt.Println("Error reading the author field:", err.Message)
		return
	}
	foundAuthor, ok := users.FindById(authorId.(int64))
	if !ok {
		fmt.Println("Author not found for id", authorId)
		return
	}
	fmt.Println("Post:", post.String())
	fmt.Println("Author found by id:", foundAuthor.String())

	// An id that was never allocated simply reports ok == false.
	if _, ok := users.FindById(999999); !ok {
		fmt.Println("Id 999999 holds no record, as expected")
	}
}
