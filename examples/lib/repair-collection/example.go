package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Repairing a collection after a crash.
//
// Every write sequence of Keep is ordered so that stopping part-way — a
// crash, a power cut, a backend failing between two writes — only leaves
// keys behind that no read mistakes for data: a record whose insert never
// committed is neither listed nor found, and a removal that stopped in the
// middle of its swap leaves a slot every listing skips. Repair deletes that
// debris and closes the hole, with nothing but single-key reads and writes.
//
// The same pass indexes what is not indexed yet, which is how a field turned
// from a String into a Key gets its unique index.
//
// This example fakes two crashes by writing the keys they leave behind
// straight into the files of the standard backend.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "handle", Type: api.String},
			},
		},
	},
}

// Indexed is the same database after "handle" was turned into a Key.
var Indexed = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
				{Name: "handle", Type: api.Key},
			},
		},
	},
}

// write puts one key of the standard backend in place by hand.
func write(path string, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		panic(err)
	}
}

// indexPath is where the unique index keeps the owner of one value: under
// the SHA-256 of the case-folded value, which for ASCII is the lower-cased one.
func indexPath(field string, value string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(value)))
	return "test-dir/database/user/keys/" + field + "/" + hex.EncodeToString(sum[:])
}

// emails lists the email of every record ListAll returns.
func emails(users api.Collection) string {
	records, failure := users.ListAll()
	if failure != nil {
		panic(failure.Message)
	}
	listed := make([]string, 0, len(records))
	for _, record := range records {
		email, _ := record.Get("email")
		listed = append(listed, fmt.Sprint(email))
	}
	return strings.Join(listed, ", ")
}

func main() {

	deps := standard.New()
	lib := sandbox.New(&deps)

	db, failure := lib.Databases.New(Props)
	if failure != nil {
		panic(failure.Message)
	}
	users, ok := db.Collection("user")
	if !ok {
		panic(`the Props declares no "user" schema`)
	}

	for _, handle := range []string{"mateus", "ana", "bruno"} {
		if _, failure := users.Insert(map[string]any{"email": handle + "@gmail.com", "handle": handle}); failure != nil {
			panic(failure.Message)
		}
	}

	// Crash 1: removing mateus — id 1, at position 1 — stopped right after
	// writing the last record, bruno (id 3), into position 1.
	write("test-dir/database/user/list/1", "3")

	// Crash 2: inserting ghost@gmail.com stopped before its commit point. It
	// had taken id 4 and written its value, its back-pointer and its index
	// entry, but never grew the size.
	write("test-dir/database/user/last-id", "4")
	write("test-dir/database/user/4/values/email", "ghost@gmail.com")
	write("test-dir/database/user/4/position", "4")
	write(indexPath("email", "ghost@gmail.com"), "4")

	// Neither crash shows: the half-moved slot is skipped, and neither the
	// half-removed nor the never-committed record is live.
	fmt.Println("listed before repair:", emails(users))
	_, ok, _ = users.FindByID(1)
	fmt.Println("mateus found by id:", ok)
	_, ok, _ = users.FindByKey("email", "ghost@gmail.com")
	fmt.Println("ghost found by email:", ok)

	// Repair closes the hole and deletes what ids 1 and 4 left behind.
	if failure := users.Repair(); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("listed after repair:", emails(users))

	// A failed insert still spent its id: ids are never reused.
	ghost, failure := users.Insert(map[string]any{"email": "ghost@gmail.com"})
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("ghost inserted as id:", ghost.ID)

	// The schema now declares "handle" a Key. The records written before
	// were never indexed under it, until Repair indexes them.
	indexed, failure := lib.Databases.New(Indexed)
	if failure != nil {
		panic(failure.Message)
	}
	migrated, _ := indexed.Collection("user")
	_, ok, _ = migrated.FindByKey("handle", "ana")
	fmt.Println("ana found by handle before repair:", ok)
	if failure := migrated.Repair(); failure != nil {
		panic(failure.Message)
	}
	found, ok, _ := migrated.FindByKey("handle", "ANA")
	fmt.Println("ana found by handle after repair:", ok, found.ID)

	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
