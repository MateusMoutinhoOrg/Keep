# Handle Module

## Description
Covers renaming the Go module, typically needed when forking this repository.

---

## Rename the Module

### Workflow
1. Open the `go.mod` file at the project root.
2. Change the module path on the first line (`module github.com/MateusMoutinhoOrg/Keep`) to your new module path.
3. Update all internal import paths. Use your IDE's global "Find and Replace" to replace `github.com/MateusMoutinhoOrg/Keep` with your new module path, or run:
   - **macOS:**
     ```sh
     find . -type f \( -name '*.go' -o -name '*.md' \) -exec sed -i '' 's|github.com/MateusMoutinhoOrg/Keep|<your-new-module-path>|g' {} +
     ```
   - **Linux:**
     ```sh
     find . -type f \( -name '*.go' -o -name '*.md' \) -exec sed -i 's|github.com/MateusMoutinhoOrg/Keep|<your-new-module-path>|g' {} +
     ```
4. Run `go mod tidy` to verify and clean up dependencies.
