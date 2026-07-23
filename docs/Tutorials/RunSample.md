# Run a Sample

## Description
Covers running the executable samples in [examples/](../../examples/) to see the library in action. There is one sample per database operation. To create a new one, follow [AddSample.md](/docs/Tutorials/AddSample.md) instead.

---

## Workflow
1. Browse the [examples/](../../examples/) directory and pick a sample (e.g. `CreateUser/`).
2. Run it from the project root with the Go toolchain:
   ```bash
   go run ./examples/CreateUser/CreateUser.go
   ```
3. Inspect the `testDatabase/` directory the samples create — with the standard adapter, each key becomes a file.
4. Re-run the sample to exercise the paths that react to already-existing data; the samples never reset the database.
