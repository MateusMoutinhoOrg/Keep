# Add a Sample

## Description
Covers creating a runnable sample in [examples/](../../examples/) that demonstrates a library feature. To run an existing one, follow [RunSample.md](/docs/Tutorials/RunSample.md) instead.

### Rules
- Creating a sample requires updating the Samples section of the [README.md](/README.md) and [Structure.md](/docs/Reference/Structure.md).
- A sample must be self-contained and runnable with a single `go run` command.
- Samples must not reset the database on start — data persists across runs, so re-running a sample exercises the "already exists" paths.
- The sample file must follow its specification — locate it in [Specs.md](/docs/Reference/Specs.md).

---

## Workflow
1. Create a directory inside [examples/](../../examples/) named after the operation being demonstrated (e.g. `examples/CountUsers/`).
2. Inside it, create the sample file with the same name as the directory (e.g. `CountUsers.go`).
3. Write a runnable `package main` program that builds deps through an adapter, injects them with `lib.New`, and exercises the feature. Store data under `testDatabase/` and comment the key parts:
   ```go
   deps := standard.New()             // filesystem adapter
   db := lib.New(deps).NewDatabase(Props)
   users := db.GetSchema("user")
   ```
4. If the sample needs setup instructions, add a `README.md` in the sample's directory.
5. Add the sample to the Samples section of the [README.md](/README.md).
6. Register the new directory and file in [Structure.md](/docs/Reference/Structure.md).
7. Verify the sample runs, following [RunSample.md](/docs/Tutorials/RunSample.md).
