package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Reading a collection one page at a time.
//
// List(position, chunk) reads the position list from position, counted from
// 1, and stops after chunk records. A page costs one read per record it
// returns and nothing for the records it skips, so paging through a large
// collection never reads the whole of it.

// ChunkSize is how many records each page holds.
const ChunkSize = 2

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
		{"email": "mateus1@gmail.com", "age": 21},
		{"email": "mateus2@gmail.com", "age": 22},
		{"email": "mateus3@gmail.com", "age": 23},
		{"email": "mateus4@gmail.com", "age": 24},
		{"email": "mateus5@gmail.com", "age": 25},
	} {
		if _, failure := users.NewItem(fields); failure != nil {
			panic(failure.Message)
		}
	}

	for position := 1; ; position += ChunkSize {
		page, failure := users.List(position, ChunkSize)
		if failure != nil {
			panic(failure.Message)
		}
		if len(page) == 0 {
			// Past the end of the list: an empty page, not a failure.
			fmt.Printf("page at %d: empty, done\n", position)
			break
		}
		fmt.Printf("page at %d:\n", position)
		for _, user := range page {
			fmt.Println(" -", user.String())
		}
	}

	// A chunk of 0 means "to the end of the collection", so List(3, 0) is
	// everything from the third record on.
	tail, failure := users.List(3, 0)
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("from position 3 to the end:", len(tail))

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
