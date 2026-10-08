package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/memory"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// The same database over a different backend.
//
// Every other example builds its deps from adapters/bindings/standard,
// which binds filestorage and writes one file per key. This one imports
// memory instead, which binds memstorage and writes nothing anywhere. Not a
// line below the import changes: the sandbox calls the same single-key
// functions either way, and which implementation stands behind them is
// decided by the binding a program picks.

// Props describes the database this example builds. Path is still a prefix
// — the in-memory backend just keeps it as part of the key.
var Props = api.Props{
	Path: "users/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "age", Type: api.Int, Required: true},
			},
		},
	},
}

func main() {

	deps := memory.New()      // memstorage in place of filestorage
	lib := sandbox.New(&deps) // *api.Sandbox, the same type

	fmt.Printf("%s %s\n", lib.Info.Name(), lib.Info.Version())

	db, failure := lib.Databases.New(Props)
	if failure != nil {
		panic(failure.Message)
	}
	users, ok := db.Collection("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}

	for _, fields := range []map[string]any{
		{"email": "mateus@gmail.com", "age": 27},
		{"email": "ana@gmail.com", "age": 31},
	} {
		if _, failure := users.Insert(fields); failure != nil {
			panic(failure.Message)
		}
	}

	found, ok, failure := users.FindByKey("email", "ana@gmail.com")
	if failure != nil || !ok {
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
	if err := os.MkdirAll("test-dir", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("test-dir/users.txt", []byte(report.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Println("wrote test-dir/users.txt with", len(all), "users")

	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
