# LibExamples

Every example of Keep used as a Go module. Each one is a `package main` program that
runs with its own directory as the working directory and writes only into its own `test-dir`,
so it can be read as documentation and copied as a starting point. It ends by copying out of
`test-dir` into `assert-dir` the paths it asserts — `os.CopyFS(dst, os.DirFS(src))`, one call per
path, each keeping the place it holds in the tree.

`agnos run-examples` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `assert-dir` file, written by
`run-examples` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one.

| Example | Description | Source |
|---|---|---|
| `bytes-fields` | Store raw binary content in a Bytes field and read it back byte for byte | [example.go](../../examples/lib/bytes-fields/example.go) |
| `concurrent-inserts` | Insert from many goroutines at once, and see a contested unique key won exactly once | [example.go](../../examples/lib/concurrent-inserts/example.go) |
| `create-user` | Insert a record, and see a unique key and a required field refused | [example.go](../../examples/lib/create-user/example.go) |
| `delete-user` | Remove a record with its index entries and everything nested under it | [example.go](../../examples/lib/delete-user/example.go) |
| `find-nested-by-key` | Open a nested field as a collection of its own, and find a nested record by its key | [example.go](../../examples/lib/find-nested-by-key/example.go) |
| `find-user-by-id` | Resolve a record straight from its permanent id, with no index read and no reuse | [example.go](../../examples/lib/find-user-by-id/example.go) |
| `find-user-by-key` | Look a record up through a unique Key field, at a fixed cost, ignoring case | [example.go](../../examples/lib/find-user-by-key/example.go) |
| `in-memory-database` | Run the same code over the in-memory backend by importing another binding | [example.go](../../examples/lib/in-memory-database/example.go) |
| `link-records` | Point a record at a record of another collection with a Link field, and follow it with GetLink | [example.go](../../examples/lib/link-records/example.go) |
| `list-all-users` | Walk every record of a collection over a backend that cannot list keys | [example.go](../../examples/lib/list-all-users/example.go) |
| `list-users-paginated` | Read a collection one page at a time, paying only for the records returned | [example.go](../../examples/lib/list-users-paginated/example.go) |
| `nested-collections` | Declare a collection inside a record, link out of it, and see it removed with its owner | [example.go](../../examples/lib/nested-collections/example.go) |
| `nosync-writes` | Turn filestorage's NoSync on and off where the deps are built | [example.go](../../examples/lib/nosync-writes/example.go) |
| `plain-value-fields` | Store a decimal and a repeatable text field, and see String accept what Key refuses | [example.go](../../examples/lib/plain-value-fields/example.go) |
| `repair-collection` | Repair a collection after a crash, and index a field turned from String into Key | [example.go](../../examples/lib/repair-collection/example.go) |
| `retrieve-user-info` | Read fields back off a record, and tell an unset field from an undeclared one | [example.go](../../examples/lib/retrieve-user-info/example.go) |
| `schema-validation` | See Databases.New refuse a Props a database cannot be built from | [example.go](../../examples/lib/schema-validation/example.go) |
| `update-user` | Write a new value for a plain field, and see a wrong Go type refused | [example.go](../../examples/lib/update-user/example.go) |
| `update-user-key` | Write a new value for an indexed field, moving its entry in the unique index a String has none of | [example.go](../../examples/lib/update-user-key/example.go) |

