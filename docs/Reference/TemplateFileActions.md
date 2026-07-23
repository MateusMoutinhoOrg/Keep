# Template File Actions

## Description
Lists every file and directory of this template and the action it takes when the template is forked into a new library or an existing library is adapted to it. Each file falls into exactly one action: **Copy**, **Create**, **Rewrite**, or **Delete**. The workflows using this list are [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) and [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md).

---

## Copy

Taken as-is from the template. They describe the structure itself, not the library, so they carry over unchanged. Adapting them is allowed but never required.

Copying these files carries over the template's **generic** guides and specifications only. The new library must still **create** its own case-specific tutorials and reference pages — see [Create](#create).

| Path | Description |
|------|-------------|
| `docs/Reference/Meta/*` | The specifications every file of the new library must be shaped by |
| `docs/Reference/RULES.md` | The binding contribution rules |
| `docs/Reference/Specs.md` | The index locating each specification |
| `docs/Reference/TemplateFileActions.md` | This page |
| `docs/Tutorials/*` | The workflow guides |
| `docs/Explanation/*` | The mechanics of the dependency injection system |
| `pkg/lib/new.go` | The `New` constructor that injects `Deps` into `Lib` |

---

## Create

Written from scratch for the library being built or adapted. Nothing of the template's content survives here — the example files occupying these paths are removed by **[Delete](#delete)**. Every created file must be shaped by the specification in its row.

| Path | Description | Specification |
|------|-------------|---------------|
| `adapters/<name>/<name>.go` | One adapter per opinionated implementation of the `Deps` contract | [Adapters](/docs/Reference/Meta/Adapters/Specs.md) |
| `docs/Reference/PublicApi/*` | One detail page per public API entry | [ReferenceDocs](/docs/Reference/Meta/ReferenceDocs/Specs.md) |
| `docs/Reference/<Name>.md` | Any reference page the new library needs beyond the public API index | [ReferenceDocs](/docs/Reference/Meta/ReferenceDocs/Specs.md) |
| `docs/Tutorials/<Goal>.md` | One tutorial per workflow specific to the new library — the template tutorials carried over by **[Copy](#copy)** do **not** fulfil this | [TutorialDocs](/docs/Reference/Meta/TutorialDocs/Specs.md) |
| `examples/<example>/<example>.go` | One runnable sample per demonstrated use case | [Examples](/docs/Reference/Meta/Examples/Specs.md) |
| `pkg/lib/*` | The library logic, calling every dependency through `l.deps` | [LibFunctions](/docs/Reference/Meta/LibFunctions/Specs.md) · [LibObjects](/docs/Reference/Meta/LibObjects/Specs.md) |

---

## Rewrite

Kept in place, with their content replaced by the new library's. The file keeps its path and its shape; only what it declares or documents changes. Every rewritten file must be shaped by the specification in its row.

| Path | Rewrite with | Specification |
|------|--------------|---------------|
| `README.md` | The new library's overview, quick start, badges, Doc Index, and Samples section | [Readme](/docs/Reference/Meta/Readme/Specs.md) |
| `adapters/standard/standard.go` | The default adapter, satisfying the new `Deps` contract | [Adapters](/docs/Reference/Meta/Adapters/Specs.md) |
| `docs/Reference/PublicApi.md` | The index of the new public API entries | [ReferenceDocs](/docs/Reference/Meta/ReferenceDocs/Specs.md) |
| `docs/Reference/Structure.md` | The layout of the new library | [Structure](/docs/Reference/Meta/Structure/Specs.md) |
| `pkg/deps/deps.go` | The `Deps` fields the new library requires | [Deps](/docs/Reference/Meta/Deps/Specs.md) |
| `pkg/lib/types.go` | The `Lib` struct and the types the new library creates | [LibObjects](/docs/Reference/Meta/LibObjects/Specs.md) |

---

## Delete

The template's example content. Removed once the new library's own files exist.

| Path | Description |
|------|-------------|
| `adapters/*` — except `adapters/standard/` | The example alternative adapters |
| `docs/Reference/PublicApi/*` | The example API detail pages |
| `examples/*` | The example samples |
| `pkg/lib/*` — except `pkg/lib/types.go` and `pkg/lib/new.go` | The example library logic |
