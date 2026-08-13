# TutorialDocs Specification

## Description
Defines the required shape of a **Tutorial** page — any `.md` file under `docs/Tutorials/`. A tutorial guides workflows: one page per **subject**, built from actionable numbered steps. A page covering a single workflow (e.g. `RenameModule.md`) carries one `## Workflow`; a page covering several closely related workflows on one subject (e.g. `HandleDocuments.md`, `HandleDependencies.md`) carries one `##` section per workflow instead.

### Rules
- Every page must comply with [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md).
- **One subject per page.** A page groups the workflows a reader performs on the same subject — documents, dependencies, library elements, samples — and nothing else. Two unrelated subjects are two pages.
- The title names the subject; each workflow section names its goal as an action (`Add a Document`, not `Documents`).
- **Single-workflow page**: a `## Description`, an optional `### Rules`, and exactly one `## Workflow` of numbered, actionable steps.
- **Multi-workflow page**: a `## Description`, an optional `### Rules` holding the constraints shared by every workflow, then one `##` section per workflow, each with an optional `### Rules` of its own and a `### Workflow` of numbered steps. Workflow sections are the page's topics, and the theme index lists them one by one — see the [Index](/docs/References/Specs/Index/Specs.md) specification.
- Steps prescribe actions, not descriptions; use fenced code blocks when a step involves writing or running code.
- Sections of the same page cross-link by anchor (`[Add a Document](#add-a-document)`) rather than repeating each other's steps.
- **Full Code at the end**: when a workflow builds a program across multiple steps (code snippets the reader assembles into one file), the section must end with a `## Full Code` — or `### Full Code`, inside a multi-workflow section — containing the complete, copy-pasteable program the steps produce. A workflow whose code already appears whole in a single step is exempt.
- Background explanations belong in `docs/References/` — link to them instead of embedding them.
- Every new page must be registered in its theme index under `docs/Index/`.

## Structure
1. **Title** (H1): the subject of the tutorial.
2. **`## Description`**: one short paragraph on what the page covers, linking to the neighbouring tutorials it does *not* cover.
3. **`### Rules`** *(optional)*: the constraints every workflow on the page shares.
4. **`---`**: horizontal rule separating the header from the workflows.
5. **`## Workflow`** — for a single-workflow page: numbered, actionable steps, with fenced code blocks where a step involves code, and links to other tutorials for any step that is itself a separate goal.
6. **`## <Goal>` sections** — for a multi-workflow page: one per workflow, separated by `---`, each holding an optional `### Rules` and a `### Workflow` shaped as above.
7. **`## Full Code`** *(required when a workflow assembles a program across steps)*: the complete resulting code in a single fenced block, ready to copy and run.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/TutorialDocs/sample.md).
