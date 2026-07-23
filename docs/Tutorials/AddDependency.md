# Add a Dependency

## Description
Covers adding a requirement to the `Deps` contract in [pkg/deps/deps.go](../../pkg/deps/deps.go) and implementing it in every existing adapter.

### Rules
- `Deps` fields are function fields — the contract declares behavior, never a concrete implementation.
- A new field must be implemented by **every** adapter in [adapters/](../../adapters/) in the same commit; a partial adapter breaks all consumers.
- The `Deps` struct must follow its specification — locate it in [Specs.md](/docs/Reference/Specs.md).
- Every field of the contract must have its behavior documented in [RequiredApi.md](/docs/Reference/RequiredApi.md).

---

## Workflow
1. Add the field to the `Deps` struct in [pkg/deps/deps.go](../../pkg/deps/deps.go), named after the behavior it provides:
   ```go
   type Deps struct {
       Write      func(key string, value []byte) error
       ReadPrefix func(prefix string) ([]string, error) // new requirement
   }
   ```
2. If the function must distinguish an expected condition (e.g. a missing key) from a real failure, add a sentinel error next to the existing ones in the same file:
   ```go
   var ErrPrefixNotFound = errors.New("keep: prefix not found")
   ```
3. Implement the new field in every adapter under [adapters/](../../adapters/), following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
4. Document the function's contract — arguments, return values, and the sentinel it must return — in [RequiredApi.md](/docs/Reference/RequiredApi.md), and update the [deps.Deps](/docs/Reference/PublicApi/deps.Deps.md) detail page.
5. Use the dependency from the library through the injected `Deps`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md).
6. If the requirement changes how dependencies behave for consumers, update [DepsMechanic.md](/docs/Explanation/DepsMechanic.md).
7. Build the project and confirm no adapter, sample, or test breaks:
   ```bash
   go build ./... && go test ./...
   ```
