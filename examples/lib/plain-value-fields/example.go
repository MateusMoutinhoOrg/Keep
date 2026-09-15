package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// The two plain value types that carry no index: Float and String.
//
// Int, Float and String are all written as a value and read back typed.
// What separates String from Key is the index: a Key is unique across the
// collection and is what FindByKey reads, a String is neither. Two records
// may hold the same String, which is what makes it the type for a title, a
// bio or a display name.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "TestDir/database/",
	Schemas: []api.Schema{
		{
			Name: "product",
			Itens: []api.Item{
				{Name: "sku", Type: api.Key, Required: true},
				{Name: "title", Type: api.String, Required: true},
				{Name: "price", Type: api.Float, Required: true},
			},
		},
	},
}

func main() {

	deps := standard.New()
	lib := sandbox.New(&deps)

	db := lib.Databases.New(Props)
	products, ok := db.GetSchema("product")
	if !ok {
		panic(`the Props declares no "product" schema`)
	}

	first, failure := products.NewItem(map[string]any{
		"sku":   "kb-001",
		"title": "Mechanical Keyboard",
		"price": 249.90,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("first:", first.String())

	// Get hands each field back in its own Go type: a float64 for a Float
	// field, a string for a String one.
	price, failure := first.Get("price")
	if failure != nil {
		panic(failure.Message)
	}
	title, failure := first.Get("title")
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Printf("price: %.2f (%T), title: %s (%T)\n", price, price, title, title)

	// A String carries no unique index, so a second record may hold the
	// same title. The sku is a Key, so that one still has to differ.
	second, failure := products.NewItem(map[string]any{
		"sku":   "kb-002",
		"title": "Mechanical Keyboard",
		"price": 199,
	})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("second:", second.String())

	// 199 was written as an int and read back as a float64: a Float field
	// accepts either, and stores one canonical decimal form.
	price, failure = second.Get("price")
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Printf("second price: %v (%T)\n", price, price)

	// FindByKey only reads the index of a Key field, so it never finds a
	// record by a String — not even one holding that exact value.
	_, ok = products.FindByKey("title", "Mechanical Keyboard")
	fmt.Println("found by title:", ok)
	found, ok := products.FindByKey("sku", "kb-002")
	fmt.Println("found by sku:", ok, found.Id)

	// An Update to either type is a single write: neither moves an index.
	if failure := second.Update("price", 179.5); failure != nil {
		panic(failure.Message)
	}
	price, failure = second.Get("price")
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("second price after update:", price)

	// The wrong Go type is refused before anything is written.
	_, failure = products.NewItem(map[string]any{
		"sku": "kb-003", "title": "Numpad", "price": "cheap",
	})
	if failure == nil || failure.Type != api.InvalidField {
		panic("a string price should have been refused")
	}
	fmt.Printf("refused: %s\n", failure.Message)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
