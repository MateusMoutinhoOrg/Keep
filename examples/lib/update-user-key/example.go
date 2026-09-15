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
//
// A String field is the same value written with none of that: no index to
// move, no conflict to refuse, one write. It is the whole difference
// between the two types.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "TestDir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Itens: []api.Item{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "handle", Type: api.String, Required: true},
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

	mateus, failure := users.NewItem(map[string]any{
		"email": "mateus@gmail.com", "handle": "mateus", "age": 27,
	})
	if failure != nil {
		panic(failure.Message)
	}
	ana, failure := users.NewItem(map[string]any{
		"email": "ana@gmail.com", "handle": "ana", "age": 31,
	})
	if failure != nil {
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

	// "handle" holds the same kind of value and is declared a String, so it
	// carries no index: the value another record already holds is written
	// without complaint, and nothing indexes either of them.
	if failure := mateus.Update("handle", "ana"); failure != nil {
		panic(failure.Message)
	}
	mine, _ := mateus.Get("handle")
	theirs, _ := ana.Get("handle")
	fmt.Printf("handles: %v and %v\n", mine, theirs)

	// Which is why FindByKey cannot read it back: only a Key is indexed.
	_, ok = users.FindByKey("handle", "ana")
	fmt.Println("found by handle:", ok)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
