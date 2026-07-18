# Run Sample

## Description
This guide describes how to run the provided executable samples to see the library in action. There is one sample per database operation.

## Workflow

1. Browse to the [examples/](../../../examples/) directory in the project repository.
2. Select a sample you want to run (e.g., `CreateUser/`).
3. Run the sample from the project root using the Go toolchain:
```bash
go run examples/CreateUser/CreateUser.go
```
4. Inspect the `testDatabase/` directory created by the samples that use the filesystem adapter — each key becomes a file.
