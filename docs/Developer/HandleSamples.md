# Handle Samples

## Description
Covers creating and running executable samples in [examples/](../../examples/) that demonstrate how to use the library.

### Rules
- Creating, deleting, or renaming a sample requires updating the [README.md](../../README.md) Samples section.
- Samples must not reset the database deterministically — data persists across runs so re-running a sample exercises the "already exists" paths.

---

## Add a Sample

### Workflow
1. Create a directory inside [examples/](../../examples/) named after the feature being demonstrated (e.g., `examples/CountUsers/`).
2. Inside it, create the sample file (e.g., `CountUsers.go`).
3. Write a runnable snippet that imports the library and instantiates it through the standard adapter, storing data under `testDatabase/`. Add comments explaining the key parts.
4. If the sample needs setup instructions, add a `README.md` in the sample's directory.
5. Update the [README.md](../../README.md) Samples section to list the new sample.
6. Verify the sample runs (see [Run a Sample](#run-a-sample)).

---

## Run a Sample

### Workflow
1. Run the sample from the project root:
   ```bash
   go run ./examples/CreateUser/CreateUser.go
   ```
