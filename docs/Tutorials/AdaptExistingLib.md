# Adapt a Pre-Existing Library

## Description
Covers converting a library that already exists into this project's dependency-injected structure. To start a new library from scratch, follow [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) instead.

### Rules
- Read [RULES.md](/docs/Reference/RULES.md) and [Structure.md](/docs/Reference/Structure.md) before starting.
- Keep the separation defined in [Structure.md](/docs/Reference/Structure.md): pure logic in `pkg/`, concrete implementations in `adapters/`.
- Every file has one action — **Copy**, **Create**, **Rewrite**, or **Delete**. Take it from [TemplateFileActions.md](/docs/Reference/TemplateFileActions.md); the steps below follow that order.

---

## Workflow
1. Recreate this project's directory layout inside the library being converted, using [Structure.md](/docs/Reference/Structure.md) as reference.
2. Copy every **[Copy](/docs/Reference/TemplateFileActions.md#copy)** file into the library unchanged — the specifications, rules, tutorials, and the lib's `New` constructor.
3. Rewrite `pkg/deps/deps.go` with the OS-bound and third-party calls the library must receive as dependencies, following [AddDependency.md](/docs/Tutorials/AddDependency.md).
4. Rewrite `adapters/standard/standard.go` so the default adapter satisfies that contract with the library's current behavior, following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
5. Rewrite the entry-point package with the library's own lib type and constructor, following [AddLibObject.md](/docs/Tutorials/AddLibObject.md).
6. Create the library logic under `pkg/` by moving the existing code in and routing every dependency call through the injected `Deps`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) and [AddLibObject.md](/docs/Tutorials/AddLibObject.md).
7. Create any additional adapter in `adapters/`, following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
8. Create the samples in `examples/` demonstrating the converted entry points, following [AddSample.md](/docs/Tutorials/AddSample.md).
9. Create the detail pages in `docs/Reference/PublicApi/` and rewrite `docs/Reference/PublicApi.md`, following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
10. Delete every **[Delete](/docs/Reference/TemplateFileActions.md#delete)** file carried over, plus the pre-existing code the converted lib replaced. For `.md` files, follow [DeleteDocument.md](/docs/Tutorials/DeleteDocument.md).
11. Rewrite `docs/Reference/Structure.md` to describe the library's actual layout.
12. Write any remaining documentation, following [AddDocument.md](/docs/Tutorials/AddDocument.md).
13. Rewrite the `README.md`: overview, quick start, badges, Doc Index, and Samples section.
14. Verify the result:
    ```bash
    go build ./... && go test ./...
    ```
