# Structure Specification

## Description
Defines the required shape of `docs/Reference/Structure.md`. Structure.md describes the project's **schema** — the *kinds* of files a project is built from — not an exhaustive listing of every concrete file.

### Rules
- Structure.md documents **schema slots**, not individual files. One row represents a kind of file (e.g. "an adapter", "a lib function"), not a specific instance.
- Each top-level directory of the project gets its own `##` section, separated by `---`.
- Directory contents are described with Markdown tables (see [GeneralDoc](/docs/Reference/Meta/GeneralDoc/Specs.md#use-file-tables-for-directory-descriptions)).
- Structure.md must **not** link to individual specifications. A slot governed by a specification carries a **Spec** column naming it; the reader resolves that name through [Specs.md](/docs/Reference/Specs.md).
- The document must let a reader construct the project from scratch — including creating `pkg/`, `adapters/`, and `docs/` — using the schema it maps and the specifications it names.
- Slots with no meaningful shape to specify (e.g. `LICENSE`, `go.mod`) leave the Spec cell empty.

## Structure
1. **Title** (H1): `# Project Structure`.
2. **Intro**: one short paragraph explaining that the document is a schema map and that named specs are resolved through `docs/Reference/Specs.md`.
3. **One `##` section per top-level directory** (`Root`, `/adapters/`, `/docs/`, `/examples/`, `/pkg/`), each with a `File | Description | Spec` table (or nested subsection tables), separated by `---`.

> **Note**: For a concrete example, refer to [sample.md](/docs/Reference/Meta/Structure/sample.md).
