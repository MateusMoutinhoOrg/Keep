package main

import (
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/adapters/impls/filestorage"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Turning NoSync on and off.
//
// standard.New() binds filestorage durable: every write is flushed to the
// disk, file and directory, before it returns — on macOS that is an
// F_FULLFSYNC, and it dominates what an insert costs. NoSync skips the
// flush. A write stays atomic, so a crashed process still sees the old value
// or the new one, but a power loss may undo writes already reported done.
// Turn it on only for data that can be rebuilt: a cache, a test, a bulk
// import that is rerun if the machine goes down.
//
// NoSync is not a Props field, it is how the storage is built, so it is
// chosen where the deps are: replace deps.StorageDeps after standard.New()
// and before sandbox.New(&deps). Nothing below sandbox.New changes, and the
// files on disk are the same either way.

// schemas is what both databases below hold.
var schemas = []api.Schema{
	{
		Name: "user",
		Fields: []api.Field{
			{Name: "email", Type: api.Key, Required: true},
			{Name: "age", Type: api.Int, Required: true},
		},
	},
}

// open builds a database under path, flushing every write unless noSync is
// set.
func open(path string, noSync bool) api.Database {
	deps := standard.New() // filestorage, durable
	deps.StorageDeps = filestorage.NewWithOptions(".", filestorage.Options{
		NoSync: noSync, // true: on, false: off — the same as filestorage.New(".")
	})
	lib := sandbox.New(&deps)

	db, failure := lib.Databases.New(api.Props{Path: path, Schemas: schemas})
	if failure != nil {
		panic(failure.Message)
	}
	return db
}

func fill(db api.Database) {
	users, ok := db.Collection("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}
	for _, fields := range []map[string]any{
		{"email": "mateus@gmail.com", "age": 27},
		{"email": "ana@gmail.com", "age": 31},
	} {
		created, failure := users.Insert(fields)
		if failure != nil {
			panic(failure.Message)
		}
		fmt.Println("  inserted:", created.String())
	}
}

func main() {

	fmt.Println("NoSync off (durable, the default):")
	fill(open("test-dir/durable/", false))

	fmt.Println("NoSync on (fast, lost on a power loss):")
	fill(open("test-dir/nosync/", true))

	// Both trees hold the same files with the same bytes: NoSync changes
	// when a write reaches the disk, never what is written.
	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
