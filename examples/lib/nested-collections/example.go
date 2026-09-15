package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// A collection nested inside a record.
//
// A field of type Database is a collection of its own, rooted at the
// record that owns it. It behaves like any top-level collection — its Key
// fields are unique within it, its records carry their own ids, it lists
// the same way — and it is removed with the record that owns it, at any
// depth. A record of one links out to another collection like any other
// record does.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "TestDir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Itens: []api.Item{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
				{
					Name: "sessions",
					Type: api.Database,
					Itens: []api.Item{
						{Name: "token", Type: api.Key, Required: true},
						{Name: "creation", Type: api.Int, Required: true},
						{Name: "expiration", Type: api.Int, Required: true},
						// A nested record links out like any other.
						{Name: "device", Type: api.Link, Target: "device"},
					},
				},
			},
		},
		{
			Name: "device",
			Itens: []api.Item{
				{Name: "serial", Type: api.Key, Required: true},
				{Name: "label", Type: api.String, Required: true},
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
	ana, failure := users.NewItem(map[string]any{"email": "ana@gmail.com", "age": 31})
	if failure != nil {
		panic(failure.Message)
	}

	// "device" is a top-level collection, not a nested one: a laptop
	// outlives any session opened from it.
	devices, _ := db.GetSchema("device")
	laptop, failure := devices.NewItem(map[string]any{"serial": "SN-1", "label": "work laptop"})
	if failure != nil {
		panic(failure.Message)
	}

	for _, session := range []map[string]any{
		{"token": "token-1", "creation": 1000, "expiration": 2000, "device": laptop},
		{"token": "token-2", "creation": 1500, "expiration": 2500},
	} {
		if _, failure := mateus.NewSubItem("sessions", session); failure != nil {
			panic(failure.Message)
		}
	}

	for _, session := range mateus.ListAll("sessions") {
		token, _ := session.Get("token")
		creation, _ := session.Get("creation")
		expiration, _ := session.Get("expiration")
		fmt.Printf("session %d: %v, %v -> %v\n", session.Id, token, creation, expiration)
	}

	// A record of a nested collection follows a Link exactly the way a
	// top-level one does: the Target names a schema of the same Props, at
	// any depth.
	first := mateus.ListAll("sessions")[0]
	device, ok := first.GetLink("device")
	if !ok {
		panic("the session should have resolved its device")
	}
	label, _ := device.Get("label")
	fmt.Printf("session %d opened from: %v\n", first.Id, label)

	// The second session set no device, so there is nothing to follow.
	_, ok = mateus.ListAll("sessions")[1].GetLink("device")
	fmt.Println("second session has a device:", ok)

	// A Key of a nested collection is unique inside that collection only,
	// so the same token can live under another user.
	_, failure = mateus.NewSubItem("sessions", map[string]any{
		"token": "token-1", "creation": 3000, "expiration": 4000,
	})
	if failure == nil || failure.Type != api.KeyConflict {
		panic("a duplicate token under the same user should be refused")
	}
	fmt.Println("refused under mateus:", failure.Message)

	if _, failure := ana.NewSubItem("sessions", map[string]any{
		"token": "token-1", "creation": 3000, "expiration": 4000,
	}); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("accepted under ana: token-1")

	// A nested field is not a value: Get refuses it, ListAll is the way in.
	_, failure = mateus.Get("sessions")
	fmt.Println("Get on a nested field:", failure.Message)

	// Removing the owner removes the collection under it.
	if failure := mateus.Remove(); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("ana's sessions after removing mateus:", len(ana.ListAll("sessions")))

	// The device is a collection of its own, so nothing nested under mateus
	// took it with it. A link is a reference, never ownership.
	_, ok = devices.FindById(laptop.Id)
	fmt.Println("device still there:", ok)

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
