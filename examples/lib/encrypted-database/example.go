package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// An encrypted database: the same Props, plus an EncryptKey.
//
// Every value the database writes passes through Deps.CryptoDeps.Encrypt
// under the key, and every value it reads through Decrypt, so a caller reads
// and writes plain values while the backend only ever holds ciphertext. The
// unique index is addressed by a digest keyed by the same key, so not even
// the name of an index entry gives a value away.

// Props describes the database this example writes. EncryptKey is the one
// line that differs from an unencrypted database.
var Props = api.Props{
	Path:       "test-dir/database/",
	EncryptKey: "a long random secret, kept out of the source in a real program",
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

	deps := standard.New()    // aesgcmcrypto fills Deps.CryptoDeps
	lib := sandbox.New(&deps) // *api.Sandbox

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

	// Reads hand back plain values, and a lookup by key still ignores case.
	ana, ok, failure := users.FindByKey("email", "ANA@gmail.com")
	if failure != nil || !ok {
		panic("ana should have been found")
	}
	fmt.Println("found:", ana.String())

	// What the backend holds is ciphertext: the file of a value never
	// contains the value, and is 28 bytes longer than it — a 12-byte nonce in
	// front, a 16-byte tag behind.
	raw, err := os.ReadFile("test-dir/database/user/2/values/email")
	if err != nil {
		panic(err)
	}
	fmt.Println("stored email holds the plain value:", strings.Contains(string(raw), "ana@gmail.com"))
	fmt.Println("stored email is", len(raw)-len("ana@gmail.com"), "bytes longer than the value")

	// The same Path opened with another key reads nothing it can decrypt.
	wrong := Props
	wrong.EncryptKey = "another secret"
	other, failure := lib.Databases.New(wrong)
	if failure != nil {
		panic(failure.Message)
	}
	otherUsers, _ := other.Collection("user")
	_, _, failure = otherUsers.FindByID(ana.ID)
	fmt.Println("another key fails with Internal:", failure != nil && failure.Type == api.Internal)

	// Every value file differs on every run — the nonce is random — so the
	// golden records what does not: each key of the tree and the size of what
	// it holds.
	report := strings.Builder{}
	err = filepath.WalkDir("test-dir/database", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(&report, "%s %d\n", filepath.ToSlash(path), info.Size())
		return nil
	})
	if err != nil {
		panic(err)
	}
	if err := os.MkdirAll("test-dir/report", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("test-dir/report/tree.txt", []byte(report.String()), 0o644); err != nil {
		panic(err)
	}

	if err := os.CopyFS("assert-dir/report", os.DirFS("test-dir/report")); err != nil {
		panic(err)
	}
}
