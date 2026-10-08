package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// Writing to one collection from many goroutines at once.
//
// Every write to a collection — an insert, an update, a removal, anything
// written under one of its records — holds the write lock of that
// collection: an in-process lock, then a lease the storage backend keeps
// beside the keys, so a writer in another process sharing the same files
// waits too. A handler that inserts from every request needs nothing of its
// own: no insert is lost, and a unique key is held by exactly one record
// however many writers race for it.
//
// Which goroutine gets which id is up to the scheduler, so this example
// asserts the two counters only: every insert that succeeded took one id and
// one position.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "user",
			Fields: []api.Field{
				{Name: "email", Type: api.Key, Required: true},
			},
		},
	},
}

const (
	writers          = 8
	insertsPerWriter = 10
)

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

	var inserted, winners, conflicts, failed atomic.Int64
	var group sync.WaitGroup
	for writer := 0; writer < writers; writer++ {
		group.Add(2)

		// One goroutine inserts emails nobody else writes...
		go func() {
			defer group.Done()
			for index := 0; index < insertsPerWriter; index++ {
				email := fmt.Sprintf("user-%d-%d@gmail.com", writer, index)
				if _, failure := users.Insert(map[string]any{"email": email}); failure != nil {
					failed.Add(1)
					continue
				}
				inserted.Add(1)
			}
		}()

		// ...and another races every other writer for the same one.
		go func() {
			defer group.Done()
			_, failure := users.Insert(map[string]any{"email": "everyone@gmail.com"})
			switch {
			case failure == nil:
				winners.Add(1)
			case failure.Type == api.KeyConflict:
				conflicts.Add(1)
			default:
				failed.Add(1)
			}
		}()
	}
	group.Wait()

	listed, failure := users.ListAll()
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Println("distinct inserts:", inserted.Load())
	fmt.Println("same email: won", winners.Load(), "refused", conflicts.Load())
	fmt.Println("failures:", failed.Load())
	fmt.Println("listed:", len(listed))

	for _, name := range []string{"size", "last-id"} {
		target := filepath.Join("assert-dir", "database", "user", name)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			panic(err)
		}
		content, err := os.ReadFile(filepath.Join("test-dir", "database", "user", name))
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(target, content, 0o644); err != nil {
			panic(err)
		}
	}
}
