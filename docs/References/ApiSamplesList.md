# API Samples List

## Description
A reference list of every Go sample shipped in [examples/libraryExamples/](/examples/libraryExamples/). Each one is a self-contained `package main` that wires an adapter into the library and drives it from code, one sample per database operation. Running or adding one is [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md).

---

## Examples

| Sample | Description |
| --- | --- |
| [CreateUserSample](/examples/libraryExamples/CreateUserSample/CreateUserSample.go) | Insert a record with unique keys |
| [FindUserByKeySample](/examples/libraryExamples/FindUserByKeySample/FindUserByKeySample.go) | Look a record up by a unique field |
| [RetrieveUserInfoSample](/examples/libraryExamples/RetrieveUserInfoSample/RetrieveUserInfoSample.go) | Read individual fields of a record |
| [UpdateUserSample](/examples/libraryExamples/UpdateUserSample/UpdateUserSample.go) | Update a plain field |
| [UpdateUserKeySample](/examples/libraryExamples/UpdateUserKeySample/UpdateUserKeySample.go) | Update a unique indexed field, re-indexing it |
| [DeleteUserSample](/examples/libraryExamples/DeleteUserSample/DeleteUserSample.go) | Remove a record and its index entries |
| [ListAllUsersSample](/examples/libraryExamples/ListAllUsersSample/ListAllUsersSample.go) | Iterate every record of a collection |
| [ListUsersPaginatedSample](/examples/libraryExamples/ListUsersPaginatedSample/ListUsersPaginatedSample.go) | Paginate through a collection |
| [SubInfosSample](/examples/libraryExamples/SubInfosSample/SubInfosSample.go) | Manage nested sub-database records |

Samples using the standard adapter write their data under `testDatabase/` (gitignored) and never reset it — re-running one exercises the "already exists" paths.
