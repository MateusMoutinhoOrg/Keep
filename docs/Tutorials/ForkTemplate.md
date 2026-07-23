# Fork This Repository as a Template

## Description
Covers using this repository's structure as the starting point for a **new** dependency-injected library. To convert a library that already exists, follow [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md) instead.

### Rules
- Read [RULES.md](/docs/Reference/RULES.md) and [Structure.md](/docs/Reference/Structure.md) before starting.
- Keep the separation defined in [Structure.md](/docs/Reference/Structure.md): pure logic in `pkg/`, concrete implementations in `adapters/`.
- Every file has one action — **Copy**, **Create**, **Rewrite**, or **Delete**. Take it from [TemplateFileActions.md](/docs/Reference/TemplateFileActions.md); the steps below follow that order.

---

## Workflow
1. On the GitHub repository page, click **"Use this template"** and create the new repository.
2. Rename the module to the new GitHub path, following [RenameModule.md](/docs/Tutorials/RenameModule.md).
3. Leave every **[Copy](/docs/Reference/TemplateFileActions.md#copy)** file untouched — they describe the structure, not the library.
4. Rewrite [pkg/deps/deps.go](../../pkg/deps/deps.go) with the dependencies the new library requires, following [AddDependency.md](/docs/Tutorials/AddDependency.md).
5. Rewrite [adapters/standard/standard.go](../../adapters/standard/standard.go) so the default adapter satisfies the new contract, following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
6. Rename [pkg/keep/](../../pkg/keep/) to the new library's entry-point package and rewrite its `New` constructor to return the new lib type, following [AddLibObject.md](/docs/Tutorials/AddLibObject.md).
7. Create the new library logic under [pkg/](../../pkg/), following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) and [AddLibObject.md](/docs/Tutorials/AddLibObject.md).
8. Create any additional adapter in [adapters/](../../adapters/), following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
9. Create the new samples in [examples/](../../examples/), following [AddSample.md](/docs/Tutorials/AddSample.md).
10. Create the new detail pages in [docs/Reference/PublicApi/](../Reference/PublicApi/) and rewrite [PublicApi.md](/docs/Reference/PublicApi.md), following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
11. Delete every remaining **[Delete](/docs/Reference/TemplateFileActions.md#delete)** file — the lib logic, adapters, samples, and documentation the new library replaced. For `.md` files, follow [DeleteDocument.md](/docs/Tutorials/DeleteDocument.md).
12. Rewrite [docs/Reference/Structure.md](/docs/Reference/Structure.md) to describe the resulting layout.
13. Write the remaining documentation, following [AddDocument.md](/docs/Tutorials/AddDocument.md).
14. Rewrite the [README.md](/README.md): overview, quick start, badges, Doc Index, and Samples section.
15. Verify the result:
    ```bash
    go build ./... && go test ./...
    ```
