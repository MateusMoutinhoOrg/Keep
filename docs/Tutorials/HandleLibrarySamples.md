# Handle Library Samples

## Description
Covers running the executable Go samples in [examples/libraryExamples/](/examples/libraryExamples/) and creating new ones. There is one sample per database operation, and each is the only kind of file where an adapter and the sandbox meet. Every shipped sample is listed in [ApiSamplesList.md](/docs/References/ApiSamplesList.md).

---

## Run a Library Sample

### Workflow
1. Browse [ApiSamplesList.md](/docs/References/ApiSamplesList.md) and pick a sample (e.g. `CreateUserSample`).
2. Run it from the project root with the Go toolchain:
   ```bash
   go run ./examples/libraryExamples/CreateUserSample/CreateUserSample.go
   ```
3. Inspect the `testDatabase/` directory the samples create — with the standard adapter, each key becomes a file.
4. Re-run the sample to exercise the paths that react to already-existing data; the samples never reset the lib.

---

## Add a Library Sample

### Rules
- Creating a sample requires updating [ApiSamplesList.md](/docs/References/ApiSamplesList.md) and [Structure.md](/docs/References/Structure.md) in the same commit.
- A sample must be self-contained and runnable with a single `go run` command.
- Samples must not reset the database on start — data persists across runs, so re-running a sample exercises the "already exists" paths.
- The sample file must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).

### Workflow
1. Create a directory inside [examples/libraryExamples/](/examples/libraryExamples/) named after the operation being demonstrated, suffixed `Sample` (e.g. `examples/libraryExamples/CountUsersSample/`).
2. Inside it, create the sample file with the same name as the directory (e.g. `CountUsersSample.go`).
3. Write a runnable `package main` program that describes the data as package-level `Schemas` and `Props` values, builds deps through an adapter, injects them with `lib.New`, and exercises the feature. Store data under `testDatabase/`:
   ```go
   var Schemas = []keeptypes.Schema{
       {
           Name: "user",
           Itens: []keeptypes.Item{
               {Name: "email", Type: keeptypes.Key, Required: true},
               {Name: "age", Type: keeptypes.Int, Required: true},
           },
       },
   }

   var Props = keeptypes.Props{
       Path:    "testDatabase/",
       Schemas: Schemas,
   }
   ```
   ```go
   deps := keepadapter.New()             // filesystem adapter
   keep := keeplib.New(deps)
   db := keep.NewDatabase(Props)
   users, _ := db.GetSchema("user")
   ```
4. If the sample needs setup instructions, add a `README.md` in the sample's directory.
5. Add the sample to [ApiSamplesList.md](/docs/References/ApiSamplesList.md).
6. Register the new directory and file in [Structure.md](/docs/References/Structure.md).
7. Verify the sample runs, following [Run a Library Sample](#run-a-library-sample).
