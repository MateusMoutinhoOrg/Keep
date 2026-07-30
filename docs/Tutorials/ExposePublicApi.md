# Expose in the Public API

## Description
Covers publishing a library function, object, or method in the public API index at [PublicApi.md](/docs/References/PublicApi.md).

### Rules
- Every public-facing entry must be listed in [PublicApi.md](/docs/References/PublicApi.md).
- Detail pages live in [docs/References/PublicApi/](/docs/References/PublicApi/) and are named `<pkg>.<Symbol>.md`.
- Adding a detail page requires updating [Structure.md](/docs/References/Structure.md) and the [README.md](/README.md) Doc Index.

---

## Workflow
1. Open [PublicApi.md](/docs/References/PublicApi.md).
2. Add the function, struct, or method to the section matching its kind, with a one-line description.
3. Create the detail page under [docs/References/PublicApi/](/docs/References/PublicApi/), named `<pkg>.<Symbol>.md` (e.g., `lib.NewExampleObject.md`), following [AddDocument.md](/docs/Tutorials/AddDocument.md).
4. Link the new detail page from its entry in [PublicApi.md](/docs/References/PublicApi.md).
5. Register the detail page in [Structure.md](/docs/References/Structure.md).
