package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MateusMoutinhoOrg/Keep/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Keep/sandbox"
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
)

// A schema is checked once, where the database is declared.
//
// A mistake in a Props is a mistake in the code, so Databases.New refuses it
// with InvalidSchema before anything is built, rather than letting it
// surface later as a collection no insert can satisfy, a link that never
// resolves, or a nested collection writing over the record that owns it.
// Error.Field is the dotted path of what is wrong.

// user is a schema every Props below starts from.
func user(fields ...api.Field) api.Schema {
	return api.Schema{
		Name:   "user",
		Fields: append([]api.Field{{Name: "email", Type: api.Key, Required: true}}, fields...),
	}
}

func main() {

	deps := standard.New()
	lib := sandbox.New(&deps)

	refused := []struct {
		why   string
		props api.Props
	}{
		{"a field with no Type", api.Props{Path: "test-dir/database/", Schemas: []api.Schema{
			user(api.Field{Name: "city"}),
		}}},
		{"a field declared twice", api.Props{Path: "test-dir/database/", Schemas: []api.Schema{
			user(api.Field{Name: "email", Type: api.String}),
		}}},
		{"a link to no schema", api.Props{Path: "test-dir/database/", Schemas: []api.Schema{
			user(api.Field{Name: "team", Type: api.Link, Target: "teams"}),
		}}},
		{"a nested field under a reserved name", api.Props{Path: "test-dir/database/", Schemas: []api.Schema{
			user(api.Field{Name: "values", Type: api.Nested, Fields: []api.Field{{Name: "size", Type: api.Int}}}),
		}}},
		{"a schema declared twice", api.Props{Path: "test-dir/database/", Schemas: []api.Schema{
			user(), user(),
		}}},
		{"a Path leaving its base", api.Props{Path: "../outside/", Schemas: []api.Schema{
			user(),
		}}},
	}

	report := strings.Builder{}
	for _, refusal := range refused {
		_, failure := lib.Databases.New(refusal.props)
		if failure == nil || failure.Type != api.InvalidSchema {
			panic(refusal.why + " should have been refused")
		}
		line := fmt.Sprintf("%s: field %q, %s\n", refusal.why, failure.Field, failure.Message)
		fmt.Print(line)
		report.WriteString(line)
	}

	// A Props that passes builds a database, and still writes nothing.
	if _, failure := lib.Databases.New(api.Props{Path: "test-dir/database/", Schemas: []api.Schema{
		user(api.Field{Name: "city", Type: api.String}),
	}}); failure != nil {
		panic(failure.Message)
	}
	fmt.Println("a valid Props: accepted")

	// Nothing above wrote a key, so the report is what this example asserts.
	if err := os.MkdirAll("test-dir", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("test-dir/refusals.txt", []byte(report.String()), 0o644); err != nil {
		panic(err)
	}
	if err := os.CopyFS("assert-dir", os.DirFS("test-dir")); err != nil {
		panic(err)
	}
}
