# Template File Actions

## Description
Lists every file and directory of this template and the action it takes when the template is forked into a new library or an existing library is adapted to it. Each file falls into exactly one action: **Copy**, **Create**, **Rewrite**, or **Delete**. The workflows using this list are [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) and [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md).

To locate any file: find its exact path below; if it is not listed by name, it falls under the `*` row of its directory. Under `docs/`, pages are listed **by name** — the generic guides are copied, the library-shaped pages are rewritten, and the domain-specific pages are deleted.

---

## Copy

Taken as-is from the template. They describe the structure itself, not the library, so they carry over unchanged. Adapting them is allowed but never required.

Copying these files carries over the template's **generic** guides and specifications only. The new library must still **create** its own case-specific tutorials and reference pages — see [Create](#create).

| Path | Description |
|------|-------------|
| `docs/References/Specs/*` | The specifications every file of the new library must be shaped by |
| `docs/References/Specs.md` | The index locating each specification |
| `docs/References/RULES.md` | The binding contribution rules |
| `docs/Tutorials/ForkTemplate.md`, `docs/Tutorials/AdaptExistingLib.md`, `docs/Tutorials/RenameModule.md`, `docs/References/TemplateFileActions.md` | The template workflows and this page |
| `docs/References/SandboxIsolation.md`, `docs/References/StructContracts.md`, `docs/References/DepsMechanic.md` | The explanations of the structure's mechanics |
| `docs/Tutorials/HandleDependencies.md`, `docs/Tutorials/HandleLibElements.md`, `docs/Tutorials/HandleLibrarySamples.md`, `docs/Tutorials/HandleDocuments.md` | The generic workflow guides for extending any library built on this structure |
| `sandbox/new.go` | The `New` constructor delegating to the internal lib constructor |

---

## Create

Written from scratch for the library being built or adapted. Nothing of the template's content survives here — the example files occupying these paths are removed by **[Delete](#delete)**. Every created file must be shaped by the specification in its row.

| Path | Description | Specification |
|------|-------------|---------------|
| `sandbox/lib/publicfunctions/*` | One file per public function field of `api.Lib`, holding its `<Field>Factory` | [LibFunctions](/docs/References/Specs/LibFunctions/Specs.md) |
| `sandbox/lib/<object>/*` | One package per object the library hands back: its field factories and the `New` constructor running them all | [LibObjects](/docs/References/Specs/LibObjects/Specs.md) |
| `sandbox/config/*` | The new library's compile-time constants — its version, and any other fixed text it reports | |
| `adapters/<name>/<name>.go` | One adapter per additional opinionated implementation of the `Deps` contract | [Adapters](/docs/References/Specs/Adapters/Specs.md) |
| `examples/libraryExamples/<example>/<example>.go` | One runnable Go sample per demonstrated use case | [LibraryExamples](/docs/References/Specs/LibraryExamples/Specs.md) |
| `docs/References/PublicApi/<pkg>.<Symbol>.md` | One detail page per public API entry | [ReferenceDocs](/docs/References/Specs/ReferenceDocs/Specs.md) |
| `docs/Tutorials/<Subject>.md` | One tutorial per workflow specific to the new library — the generic guides carried over by **[Copy](#copy)** do **not** fulfil this | [TutorialDocs](/docs/References/Specs/TutorialDocs/Specs.md) |
| `docs/References/<Name>.md` | Any reference page the new library needs beyond the public API index | [ReferenceDocs](/docs/References/Specs/ReferenceDocs/Specs.md) |

---

## Rewrite

Kept in place, with their content replaced by the new library's. The file keeps its path and its shape; only what it declares or documents changes. Every rewritten file must be shaped by the specification in its row.

| Path | Rewrite with | Specification |
|------|--------------|---------------|
| `README.md` | The new library's overview, badges, and the Doc Index pointing at each theme index | [Readme](/docs/References/Specs/Readme/Specs.md) |
| `sandbox/contracts/deps/deps.go` | The `Deps` function fields the new library requires | [Deps](/docs/References/Specs/Deps/Specs.md) |
| `sandbox/contracts/api/api.go` | The `Lib` struct and one struct per object the new library hands back | [Outputs](/docs/References/Specs/Outputs/Specs.md) |
| `adapters/standard/standard.go` | The default adapter, filling the new `Deps` contract | [Adapters](/docs/References/Specs/Adapters/Specs.md) |
| `sandbox/lib/new.go` | The `New` constructor assigning every new lib factory's return value | [LibFunctions](/docs/References/Specs/LibFunctions/Specs.md) |
| `docs/References/Structure.md` | The layout of the new library | [Structure](/docs/References/Specs/Structure/Specs.md) |
| `docs/References/PublicApi.md` | The index of the new public API entries | [ReferenceDocs](/docs/References/Specs/ReferenceDocs/Specs.md) |
| `docs/References/Adapters.md` | The adapters the new library ships | [AdaptersDoc](/docs/References/Specs/AdaptersDoc/Specs.md) |
| `docs/References/ApiSamplesList.md` | The new library's own samples | [ReferenceDocs](/docs/References/Specs/ReferenceDocs/Specs.md) |
| `docs/Index/LibUsage.md`, `docs/Index/Development.md`, `docs/Index/Templating.md` | The new library's page list, one index per theme | [Index](/docs/References/Specs/Index/Specs.md) |
| `docs/Tutorials/LibInitialization.md` | Installing the new library and wiring its standard adapter | [TutorialDocs](/docs/References/Specs/TutorialDocs/Specs.md) |

---

## Delete

The template's example content — the schema database. Removed once the new library's own files exist. For `.md` files, follow [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md#delete-a-document) so the theme indexes stay in sync.

| Path | Description |
|------|-------------|
| `sandbox/lib/*` | The engine's lib factories, object packages, and dense helpers — replaced by **[Create](#create)** |
| `sandbox/config/*` | The engine's version constant — replaced by **[Create](#create)** |
| `adapters/*` — except `adapters/standard/` | The example alternative adapters |
| `examples/libraryExamples/*` | The engine's Go samples |
| `docs/References/PublicApi/*` | The engine's public API detail pages |
| `docs/Tutorials/DefineDatabase.md`, `docs/Tutorials/AddSchemaField.md`, `docs/Tutorials/AddNestedCollection.md`, `docs/Tutorials/AddDatabaseOperation.md` | The engine's domain tutorials |
| `docs/References/Schemas.md`, `docs/References/Records.md`, `docs/References/Errors.md`, `docs/References/DenseRecordPattern.md`, `docs/References/RequiredApi.md` | The engine's domain reference and explanation pages |
