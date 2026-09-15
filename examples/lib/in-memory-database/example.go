package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/adapters/availables/native"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// The same database over a different backend.
//
// Every other example builds its deps from adapters/availables/standard,
// which binds filestorage and writes one file per key. This one imports
// native instead, which binds memstorage and writes nothing anywhere. Not a
// line below the import changes: the sandbox calls the same eleven
// single-key functions either way, and which implementation stands behind
// them is decided by the available a program picks.

// Props describes the database this example builds. Path is still a prefix
// — the in-memory backend just keeps it as part of the key.
var Props = api.Props{
	Path: "users/",
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

	deps := native.New()      // memstorage in place of filestorage
	lib := sandbox.New(&deps) // *api.Sandbox, the same type

	fmt.Printf("%s %s\n", lib.Info.Name(), lib.Info.Version())

	db := lib.Databases.New(Props)
	users, ok := db.GetSchema("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}

	for _, fields := range []map[string]any{
		{"email": "mateus@gmail.com", "age": 27},
		{"email": "ana@gmail.com", "age": 31},
	} {
		if _, failure := users.NewItem(fields); failure != nil {
			panic(failure.Message)
		}
	}

	found, ok := users.FindByKey("email", "ana@gmail.com")
	if !ok {
		panic("ana should have been found")
	}
	fmt.Println("found:", found.String())

	all, failure := users.ListAll()
	if failure != nil {
		panic(failure.Message)
	}

	// Nothing above touched the filesystem, so this example has to write
	// what it asserts itself: the report is the result, and it is the only
	// file the run leaves behind.
	report := strings.Builder{}
	for _, user := range all {
		report.WriteString(user.String())
		report.WriteString("\n")
	}
	if err := os.MkdirAll("TestDir", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("TestDir/users.txt", []byte(report.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Println("wrote TestDir/users.txt with", len(all), "users")

	if err := os.CopyFS("AssertDir", os.DirFS("TestDir")); err != nil {
		panic(err)
	}
}
