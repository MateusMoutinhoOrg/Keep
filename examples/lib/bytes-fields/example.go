package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// The plain value type for binary content: Bytes.
//
// A Bytes field is written as a []byte and stored exactly as given — no
// text encoding, no escaping — so a file, an image or a hash comes back byte
// for byte. Like a String it carries no index, which is what makes it the
// type for content rather than for a lookup.

// Props describes the database this example writes.
var Props = api.Props{
	Path: "test-dir/database/",
	Schemas: []api.Schema{
		{
			Name: "file",
			Fields: []api.Field{
				{Name: "name", Type: api.Key, Required: true},
				{Name: "content", Type: api.Bytes, Required: true},
			},
		},
	},
}

func main() {

	deps := standard.New()
	lib := sandbox.New(&deps)

	db := lib.Databases.New(Props)
	files, ok := db.Collection("file")
	if !ok {
		panic(`the Props declares no "file" schema`)
	}

	// The signature a PNG file starts with: not text, and a zero byte at the
	// end besides — exactly what a string field could not be trusted with.
	header := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0x00}

	logo, failure := files.Insert(map[string]any{
		"name":    "logo.png",
		"content": header,
	})
	if failure != nil {
		panic(failure.Message)
	}
	// String prints a Bytes field by its length, never by its contents.
	fmt.Println("logo:", logo.String())

	// Get hands the field back as a []byte holding the very bytes written.
	content, failure := logo.Get("content")
	if failure != nil {
		panic(failure.Message)
	}
	raw := content.([]byte)
	fmt.Printf("content: % x (%T)\n", raw, content)
	fmt.Println("same bytes:", bytes.Equal(raw, header))

	// An Update is a single write, like any field without an index.
	if failure := logo.Update("content", []byte("GIF89a")); failure != nil {
		panic(failure.Message)
	}
	content, failure = logo.Get("content")
	if failure != nil {
		panic(failure.Message)
	}
	fmt.Printf("after update: %q, %d bytes\n", content, len(content.([]byte)))

	// A Bytes field takes a []byte and nothing else: a string is refused
	// before anything is written.
	_, failure = files.Insert(map[string]any{
		"name": "notes.txt", "content": "plain text",
	})
	if failure == nil || failure.Type != api.InvalidField {
		panic("a string content should have been refused")
	}
	fmt.Printf("refused: %s\n", failure.Message)

	// database/file/1/values/content holds the bytes of the last Update and
	// nothing else — its sha in result.yaml is the sha of "GIF89a".
	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
