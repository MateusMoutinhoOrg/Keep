# Template File Actions

## Description
Lists every file and directory of this project and the action it takes when the structure is reused for another library — forked into a new one, or adapted onto a library that already exists. Each file falls into exactly one action: **Copy**, **Create**, **Rewrite**, or **Delete**. The workflows using this list are [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) and [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md).

---

## Copy

Taken as-is. They describe the structure itself, not the library, so they carry over unchanged. Adapting them is allowed but never required.

| Path | Description |
|------|-------------|
| `docs/Reference/Meta/*` | The specifications every file of the new library must be shaped by |
| `docs/Reference/RULES.md` | The binding contribution rules |
| `docs/Reference/Specs.md` | The index locating each specification |
| `docs/Reference/TemplateFileActions.md` | This page |
| `docs/Tutorials/*` | The workflow guides |
| `pkg/keep/new.go` | The `New` constructor that injects `Deps` into the lib entry point |

---

## Create

Written from scratch for the library being built or adapted. Nothing of this project's content survives here — the files occupying these paths are removed by **[Delete](#delete)**.

| Path | Description | Specification |
|------|-------------|---------------|
| `adapters/<name>/<name>.go` | One adapter per opinionated implementation of the `Deps` contract | [Adapters](/docs/Reference/Meta/Adapters/Specs.md) |
| `docs/Explanation/*` | One page per mechanic of the new library | [ExplanationDocs](/docs/Reference/Meta/ExplanationDocs/Specs.md) |
| `docs/Reference/PublicApi/*` | One detail page per public API entry | [ReferenceDocs](/docs/Reference/Meta/ReferenceDocs/Specs.md) |
| `examples/<example>/<example>.go` | One runnable sample per demonstrated use case | [Examples](/docs/Reference/Meta/Examples/Specs.md) |
| `pkg/<lib>/*` | The library logic, reaching storage only through the injected `Deps` | [LibFunctions](/docs/Reference/Meta/LibFunctions/Specs.md) · [LibObjects](/docs/Reference/Meta/LibObjects/Specs.md) |

---

## Rewrite

Kept in place, with their content replaced by the new library's. The file keeps its path and its shape; only what it declares or documents changes.

| Path | Rewrite with | Specification |
|------|--------------|---------------|
| `README.md` | The new library's overview, quick start, badges, Doc Index, and Samples section | [Readme](/docs/Reference/Meta/Readme/Specs.md) |
| `adapters/standard/standard.go` | The default adapter, satisfying the new `Deps` contract | [Adapters](/docs/Reference/Meta/Adapters/Specs.md) |
| `docs/Explanation/DepsMechanic.md` | How the new library's dependencies are injected, overwritten, and implemented | [ExplanationDocs](/docs/Reference/Meta/ExplanationDocs/Specs.md) |
| `docs/Reference/PublicApi.md` | The index of the new public API entries | [ReferenceDocs](/docs/Reference/Meta/ReferenceDocs/Specs.md) |
| `docs/Reference/Structure.md` | The layout of the new library | [Structure](/docs/Reference/Meta/Structure/Specs.md) |
| `pkg/deps/deps.go` | The `Deps` fields the new library requires | [Deps](/docs/Reference/Meta/Deps/Specs.md) |
| `pkg/keep/new.go` | The lib entry point, when the new library renames it | [LibObjects](/docs/Reference/Meta/LibObjects/Specs.md) |

---

## Delete

This project's own content. Removed once the new library's files exist.

| Path | Description |
|------|-------------|
| `adapters/*` — except `adapters/standard/` | The alternative adapters |
| `docs/Explanation/*` — except `DepsMechanic.md` | The pages explaining this library's mechanics |
| `docs/Reference/Errors.md`, `docs/Reference/RequiredApi.md` | The reference pages specific to this library |
| `docs/Reference/PublicApi/*` | The API detail pages |
| `examples/*` | The samples |
| `pkg/database/*` | The storage engine this library replaced |
