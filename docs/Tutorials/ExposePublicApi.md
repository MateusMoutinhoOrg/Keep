# Expose in the Public API

## Description
Covers publishing a library function, object, or method in the public API index at [PublicApi.md](/docs/Reference/PublicApi.md).

### Rules
- Every public-facing entry must be listed in [PublicApi.md](/docs/Reference/PublicApi.md).
- Detail pages live in [docs/Reference/PublicApi/](../Reference/PublicApi/) and are named `<pkg>.<Symbol>.md`.
- Adding a detail page requires updating [Structure.md](/docs/Reference/Structure.md) and the [README.md](/README.md) Doc Index.

---

## Workflow
1. Open [PublicApi.md](/docs/Reference/PublicApi.md).
2. Add the function, struct, or method to the section matching its kind, with a one-line description.
3. Create the detail page under [docs/Reference/PublicApi/](../Reference/PublicApi/), named `<pkg>.<Symbol>.md` (e.g., `lib.NewExampleObject.md`), following [AddDocument.md](/docs/Tutorials/AddDocument.md).
4. Link the new detail page from its entry in [PublicApi.md](/docs/Reference/PublicApi.md).
5. Register the detail page in [Structure.md](/docs/Reference/Structure.md).
