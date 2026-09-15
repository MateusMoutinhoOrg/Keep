package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Writing a new value for an indexed field.
//
// Updating a Key field moves the record's entry in the unique index. Keep
// writes the new index entry first, then the value, then deletes the old
// entry — so a crash part-way through can leave a stale entry pointing at a
// record that no longer holds that value, which resolves to nothing, but
// never leaves the record unreachable.

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

	mateus, failure := users.NewItem(map[string]any{"email": "mateus@gmail.com", "age": 27})
	if failure != nil {
		panic(failure.Message)
	}
	if _, failure := users.NewItem(map[string]any{"email": "ana@gmail.com", "age": 31}); failure != nil {
		panic(failure.Message)
	}

	if failure := mateus.Update("email", "mateus@keep.dev"); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("updated:", mateus.String())

	// The index followed the value.
	_, ok = users.FindByKey("email", "mateus@keep.dev")
	fmt.Println("found under the new email:", ok)
	_, ok = users.FindByKey("email", "mateus@gmail.com")
	fmt.Println("found under the old email:", ok)

	// The record is still the same record: its id never changed.
	byId, _ := users.FindById(mateus.Id)
	fmt.Println("by id:", byId.String())

	// A value another live record already holds is refused, and nothing is
	// written — the old value stays indexed.
	failure = mateus.Update("email", "ana@gmail.com")
	if failure == nil || failure.Type != api.KeyConflict {
		panic("taking another record's key should be refused")
	}
	fmt.Printf("refused: %s (field %q, value %v)\n", failure.Message, failure.Key, failure.KeyValue)

	current, _ := mateus.Get("email")
	fmt.Println("email after the refusal:", current)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
