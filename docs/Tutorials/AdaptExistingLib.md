# Adapt a Pre-Existing Library

## Description
Covers converting a library that already exists into this project's dependency-injected structure. To start a new library from scratch, follow [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) instead.

### Rules
- Read [RULES.md](/docs/References/RULES.md) and [Structure.md](/docs/References/Structure.md) before starting.
- Keep the separation defined in [Structure.md](/docs/References/Structure.md): pure logic in `pkg/`, concrete implementations in `adapters/`.
- Every file of the template has one action — **Copy**, **Create**, **Rewrite**, or **Delete**. Take it from [TemplateFileActions.md](/docs/References/TemplateFileActions.md); the steps below follow that order.
- The pre-existing package layout does **not** survive: all library logic ends up in `sandbox/internal/`, calling every OS-bound and third-party dependency through `l.deps`. Code left in its original packages, or still calling `os`/`net`/third-party APIs directly, is not adapted.
- Every file created or rewritten — code and `.md` alike — must follow its specification, located through [Specs.md](/docs/References/Specs.md). A file that ignores its specification is not adapted.
- The adaptation is not complete until the final checklist in the last workflow step passes.

---

## Workflow
1. Recreate this project's directory layout inside the library being converted, using [Structure.md](/docs/References/Structure.md) as reference.
2. Copy every **[Copy](/docs/References/TemplateFileActions.md#copy)** file into the library unchanged — the specifications, rules, tutorials, explanations, and [sandbox/new.go](../../sandbox/new.go).
3. Rewrite `sandbox/contracts/deps/deps.go` with the OS-bound and third-party calls the library must receive as dependencies, following [AddDependency.md](/docs/Tutorials/AddDependency.md).
4. Rewrite `adapters/standard/standard.go` so the default adapter satisfies that contract with the library's current behavior, following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
5. Declare every interface the library exchanges — inputs and outputs alike — in `sandbox/contracts/api/api.go`, add a constructor in `sandbox/description.go` for each input interface, then implement them under `sandbox/internal/<object>/`, following [AddLibObject.md](/docs/Tutorials/AddLibObject.md).
6. Rewrite the existing library code into `sandbox/internal/`: move each source file in, hang its public functions off `Lib` (or off objects it creates), and replace **every** OS-bound or third-party call with a call through `l.deps.<Field>()`, following [AddLibFunction.md](/docs/Tutorials/AddLibFunction.md) and [AddLibObject.md](/docs/Tutorials/AddLibObject.md). Do not keep the code in its original packages or leave direct calls in place.
7. Create any additional adapter in `adapters/`, following [AddAdapter.md](/docs/Tutorials/AddAdapter.md).
8. Create the samples in `examples/` demonstrating the converted entry points, following [AddSample.md](/docs/Tutorials/AddSample.md).
9. Create the detail pages in `docs/References/PublicApi/` and rewrite `docs/References/PublicApi.md`, following [ExposePublicApi.md](/docs/Tutorials/ExposePublicApi.md).
10. Delete every **[Delete](/docs/References/TemplateFileActions.md#delete)** file carried over from the template, plus the pre-existing code the converted lib replaced. For `.md` files, follow [DeleteDocument.md](/docs/Tutorials/DeleteDocument.md).
11. Rewrite `docs/References/Structure.md` to describe the library's actual layout.
12. Create the tutorials specific to the converted library — one page per workflow its maintainers will repeat (e.g. adding a domain object, extending a feature, releasing) — following [AddDocument.md](/docs/Tutorials/AddDocument.md) and the [TutorialDocs specification](/docs/References/Meta/TutorialDocs/Specs.md). The template tutorials copied in step 2 cover the structure only; they do not document the library's own use cases.
13. Create any reference page the library needs beyond the public API — following [AddDocument.md](/docs/Tutorials/AddDocument.md) and the [ReferenceDocs specification](/docs/References/Meta/ReferenceDocs/Specs.md).
14. Rewrite the `README.md`: overview, quick start, badges, Doc Index, and Samples section.
15. Verify the result:
```bash
go build ./...
```
Then confirm every item below — the adaptation is only done when all pass:
- All library logic lives in `sandbox/internal/`; no file there imports `os`, `net`, or a third-party implementation directly — every such call goes through `l.deps`.
- `sandbox/contracts/deps/deps.go` declares one method per injected call, and **every** adapter in `adapters/` implements every method.
- Tutorials and reference pages specific to this library exist under `docs/Tutorials/` and `docs/References/`.
- Every created or rewritten file matches its specification from [Specs.md](/docs/References/Specs.md).
- The `README.md` Doc Index lists every `.md` file and the Samples section lists every sample.
