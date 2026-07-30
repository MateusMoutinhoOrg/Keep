# Add a Dependency

## Description
Covers adding a requirement to the `Deps` contract in [sandbox/contracts/deps/deps.go](../../sandbox/contracts/deps/deps.go) and implementing it in every existing adapter.

### Rules
- `Deps` is an interface and each requirement is a method — the contract declares behavior, never a concrete implementation.
- A new method must be implemented by **every** adapter in [adapters/](../../adapters/) in the same commit; a partial adapter no longer satisfies the interface and breaks all consumers at compile time.
- The `Deps` interface must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).
- Every method of the contract must have its behavior documented in [RequiredApi.md](/docs/References/RequiredApi.md).
- Adding a method is the **only** way to give the sandbox a capability it lacks; reaching for `os`, `net`, or a third-party module inside `sandbox/` is forbidden by [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).

---

## Workflow
1. Add the method to the `Deps` interface in [sandbox/contracts/deps/deps.go](../../sandbox/contracts/deps/deps.go), named after the behavior it provides:
   ```go
   type Deps interface {
       Write(key string, value []byte) error
       ReadPrefix(prefix string) ([]string, error) // new requirement
   }
   ```
2. If the method must distinguish an expected condition (e.g. a missing key) from a real failure, add a sentinel error next to the existing ones in the same file:
   ```go
   var ErrPrefixNotFound = errors.New("keep: prefix not found")
   ```
3. Implement the new method in every adapter under [adapters/](../../adapters/), following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
4. Document the method's contract — arguments, return values, and the sentinel it must return — in [RequiredApi.md](/docs/References/RequiredApi.md), and update the [deps.Deps](/docs/References/PublicApi/deps.Deps.md) detail page.
5. Use the dependency from the library through the injected `Deps`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md).
6. If the requirement changes how dependencies behave for consumers, update [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).
7. Build the project and confirm no adapter, sample, or test breaks:
   ```bash
   go build ./... && go test ./...
   ```
