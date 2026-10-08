# `sandbox/api/info.go`

## `Info`

Info is the contract reporting the library's own identity, carried by the Sandbox as the field of the same name. Both values are read from Sandbox.Config, which the build renders from AgnosConfig/project.yaml into sandbox/internal/generated/config, so a release bump is a one-line edit touching no logic — and a caller can report which Keep it linked against without importing anything but this package.

| Field | Type | Description |
| --- | --- | --- |
| `Name` | `func() string` | Name is the library's name, "Keep". |
| `Version` | `func() string` | Version is the release the caller linked against, in the "v0.0.0" spelling the repository tags with. |

[every contract](doc.md)
