# LibExamples

Every example of Keep used as a Go module. Each one is a `package main` program that
runs with its own directory as the working directory and writes only into its own `TestDir`,
so it can be read as documentation and copied as a starting point. It ends by copying out of
`TestDir` into `AssertDir` the paths it asserts — `os.CopyFS(dst, os.DirFS(src))`, one call per
path, each keeping the place it holds in the tree.

`agnos exec-test` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `AssertDir` file, written by
`exec-test` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one.

| Example | Description | Source |
|---|---|---|
| `create-user` | Insert a record, and see a unique key and a required field refused | [example.go](../../examples/lib/create-user/example.go) |
| `delete-user` | Remove a record with its index entries and everything nested under it | [example.go](../../examples/lib/delete-user/example.go) |
| `find-user-by-id` | Point one collection at another through a record id that is never reused | [example.go](../../examples/lib/find-user-by-id/example.go) |
| `find-user-by-key` | Look a record up through a unique Key field, in one read, case-insensitively | [example.go](../../examples/lib/find-user-by-key/example.go) |
| `in-memory-database` | Run the same code over the in-memory backend by importing another available | [example.go](../../examples/lib/in-memory-database/example.go) |
| `list-all-users` | Walk every record of a collection over a backend that cannot list keys | [example.go](../../examples/lib/list-all-users/example.go) |
| `list-users-paginated` | Read a collection one page at a time, paying only for the records returned | [example.go](../../examples/lib/list-users-paginated/example.go) |
| `nested-collections` | Declare a collection inside a record, and see it removed with its owner | [example.go](../../examples/lib/nested-collections/example.go) |
| `retrieve-user-info` | Read fields back off a record, and tell an unset field from an undeclared one | [example.go](../../examples/lib/retrieve-user-info/example.go) |
| `update-user` | Write a new value for a plain field, and see a wrong Go type refused | [example.go](../../examples/lib/update-user/example.go) |
| `update-user-key` | Write a new value for an indexed field, moving its entry in the unique index | [example.go](../../examples/lib/update-user-key/example.go) |

