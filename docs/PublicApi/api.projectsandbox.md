# `sandbox/api/projectsandbox.go`

## `ProjectSandbox`

ProjectSandbox is the part of the Sandbox Keep declares itself: api.Sandbox embeds it, so every field typed here is a field of the Sandbox, reached as sandbox.<Field> by every function handed it. Each field is filled by the Constructor of the package of the same name under sandbox/constructors/, which calls the New<Field> of sandbox/internal/<field>/new.go. Written once by `agnos start` and then Keep's: no build rewrites this file.

| Field | Type | Description |
| --- | --- | --- |
| `Databases` | `Databases` | Databases is the entry point of the library: it binds a Props description to the storage the sandbox was built over. |
| `Info` | `Info` | Info reports which Keep the caller linked against. |

[every contract](doc.md)
