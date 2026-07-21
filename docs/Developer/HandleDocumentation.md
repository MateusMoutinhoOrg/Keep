# Handle Documentation

## Description
Covers creating, updating, renaming, and deleting documentation files, including standalone docs and use cases. Every `.md` file must comply with [DocumentationStandards.md](../DocumentationStandards.md).

### Rules
- Every `.md` file must follow all rules in [DocumentationStandards.md](../DocumentationStandards.md).
- Creating, deleting, or renaming a `.md` file requires updating the [README.md](../../README.md) Internal Docs section.
- Creating, deleting, or renaming a source file or directory requires updating [STRUCTURE.md](../Structure.md).

---

## Add or Update a Document

### Workflow
1. Identify where the document belongs:
   - **Usage-facing** docs → `docs/Usage/`.
   - **Developer-facing** docs → `docs/Developer/`.
2. Create the `.md` file in the correct directory with a descriptive, topic-based name (e.g., `HandleSamples.md`, `PublicApi.md`).
3. Write the content following all rules in [DocumentationStandards.md](../DocumentationStandards.md), paying special attention to:
   - **Topic-driven structure** — one concern per section.
   - **Conciseness** — short, direct sentences.
   - **Heading hierarchy** — never skip heading levels.
4. Add cross-references using **relative paths**, and add the reverse link in any document that should point back to this one.
5. Update the [README.md](../../README.md) Internal Docs section (Usage or Developer table).
6. If the document describes a file or directory, update [STRUCTURE.md](../Structure.md) to reflect the structural change.
7. Verify all impacted files are updated in the **same commit** (see [Keep Documentation in Sync](../DocumentationStandards.md#8-keep-documentation-in-sync)).

---

## Add a Use Case

### Workflow
1. Create the `.md` file inside the relevant directory with a topic-based name that groups related actions (e.g., `HandleDependencies.md`).
2. Follow the [Use Case Document Format](../DocumentationStandards.md#5-use-case-document-format): a `Description`, an optional `Rules` section, and one or more `Workflow` sections with **actionable** numbered steps.
3. Complete the shared steps in [Add or Update a Document](#add-or-update-a-document) to register the file in the README and keep the structure in sync.

---

## Rename or Delete a Document

### Workflow
1. Rename or delete the `.md` file.
2. Update every document that links to it, using the [DocumentationStandards](../DocumentationStandards.md#3-cross-reference-between-documents) cross-reference rules.
3. Update the [README.md](../../README.md) Internal Docs section and, if the change is structural, [STRUCTURE.md](../Structure.md).
