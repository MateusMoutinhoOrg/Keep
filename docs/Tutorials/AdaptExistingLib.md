# Adapt a Pre-Existing Library

## Description
Covers converting a library that already exists into this project's dependency-injected structure. To start a new library from scratch, follow [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) instead.

### Rules
- Read [Structure.md](/docs/References/Structure.md) and [Specs.md](/docs/References/Specs.md) before starting.
- Keep the separation defined in [Structure.md](/docs/References/Structure.md): pure logic in `sandbox/`, concrete implementations in `adapters/`.
- Every file of the template has one action — **Copy**, **Create**, **Rewrite**, or **Delete**. Take it from [TemplateFileActions.md](/docs/References/TemplateFileActions.md); the steps below follow that order.
- The pre-existing package layout does **not** survive: all library logic ends up in `sandbox/lib/` as factories, calling every OS-bound and third-party dependency through the carrier's `Deps` field. Code left in its original packages, or still calling `os`/`net`/third-party APIs directly, is not adapted.
- Every file created or rewritten — code and `.md` alike — must follow its specification, located through [Specs.md](/docs/References/Specs.md). A file that ignores its specification is not adapted.
- The adaptation is not complete until the final checklist in the last workflow step passes.

---

## Workflow
1. Recreate this project's directory layout inside the library being converted, using [Structure.md](/docs/References/Structure.md) as reference.
2. Copy every **[Copy](/docs/References/TemplateFileActions.md#copy)** file into the library unchanged — the specifications, rules, tutorials, explanations, and [sandbox/new.go](../../sandbox/new.go).
3. Rewrite `sandbox/contracts/deps/deps.go` with the OS-bound and third-party calls the library must receive as dependencies, following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#add-a-dependency).
4. Rewrite `adapters/standard/standard.go` so the default adapter satisfies that contract with the library's current behavior, following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#create-an-adapter-in-this-repository).
5. Declare every struct the library exchanges — inputs and outputs alike — in `sandbox/contracts/api/api.go`, plain data for the ones with no behavior and a struct of function fields (leading with `Deps`) for the ones with behavior, then fill the latter with factories under `sandbox/lib/<object>/`, following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#add-a-library-object).
6. Rewrite the existing library code into `sandbox/lib/`: move each source file in, turn its public functions into `<Field>Factory` functions assigned from the object's constructor, and replace **every** OS-bound or third-party call with a call through the carrier's `Deps.<Field>()`, following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#add-a-library-function) and [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#add-a-library-object). Do not keep the code in its original packages or leave direct calls in place.
7. Create any additional adapter in `adapters/`, following [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md#create-an-adapter-in-this-repository).
8. Create the samples in `examples/libraryExamples/` demonstrating the converted entry points, following [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample).
9. Create the detail pages in `docs/References/PublicApi/` and rewrite `docs/References/PublicApi.md`, following [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md#expose-in-the-public-api).
10. Delete every **[Delete](/docs/References/TemplateFileActions.md#delete)** file carried over from the template, plus the pre-existing code the converted lib replaced. For `.md` files, follow [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md#delete-a-document).
11. Rewrite `docs/References/Structure.md` to describe the library's actual layout.
12. Create the tutorials specific to the converted library — one page per workflow its maintainers will repeat (e.g. adding a domain object, extending a feature, releasing) — following [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md#add-a-document) and the [TutorialDocs specification](/docs/References/Specs/TutorialDocs/Specs.md). The template tutorials copied in step 2 cover the structure only; they do not document the library's own use cases.
13. Create any reference page the library needs beyond the public API — following [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md#add-a-document) and the [ReferenceDocs specification](/docs/References/Specs/ReferenceDocs/Specs.md).
14. Rewrite the `README.md`: overview, quick start, badges, Doc Index, and Samples section.
15. Verify the result:
```bash
go build ./...
```
Then confirm every item below — the adaptation is only done when all pass:
- All library logic lives in `sandbox/lib/` as factories; no file there imports `os`, `net`, or a third-party implementation directly — every such call goes through the carrier's `Deps` field.
- `sandbox/contracts/deps/deps.go` declares one function field per injected call, and **every** adapter in `adapters/` fills every field with a factory.
- Tutorials and reference pages specific to this library exist under `docs/Tutorials/` and `docs/References/`.
- Every created or rewritten file matches its specification from [Specs.md](/docs/References/Specs.md).
- The `README.md` Doc Index lists every `.md` file and the Samples section lists every sample.
