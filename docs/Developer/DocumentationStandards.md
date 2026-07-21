# Documentation Standards

General rules and conventions for writing and maintaining documentation in this project.
For the step-by-step workflow on adding new documentation, see [HandleDocumentation](./HandleDocumentation.md).

---

## 1. Topic-Driven Structure

- Organize every document around **topics** and **subtopics** using Markdown headings.
- Each section must cover a single concern.
- Use horizontal rules (`---`) to visually separate top-level sections.

---

## 2. Be Concise

- Write short, direct sentences.
- Avoid filler words, redundant explanations, and unnecessary qualifiers.

| ❌ Avoid | ✅ Prefer |
|----------|----------|
| "This file is responsible for the implementation of the authentication logic." | "Implements the authentication logic." |
| "In order to be able to run the project…" | "To run the project…" |
| "It should be noted that this function returns an error." | "Returns an error on failure." |

---

## 3. Cross-Reference Between Documents

- Always use **relative paths** (e.g., `../../Usage/PublicApi.md`), never absolute filesystem paths.
- Link to the **most specific section** possible using anchors (e.g., `Structure.md#adapters`).
- When referencing content that is explained elsewhere, **link** to it instead of duplicating the content.

---

## 4. Use File Tables for Directory Descriptions

- When describing the contents of a directory, use a Markdown table with `File` and `Description` columns.

**Example:**

```markdown
### `/config/`

| File | Description |
|------|-------------|
| `config.go` | Loads and validates environment configuration |
| `defaults.go` | Defines default values for all settings |
```

---

## 5. Use Case Document Format

- Every use case document (in `docs/Developer/` or `docs/Usage/`) must follow this structure:

```markdown
## Title

## Description
Brief explanation of what this guide covers.

### Rules (optional)
- Specific constraints for this use case.

## Workflow
1. First step
2. Second step
3. Third step
```

- Each workflow step must be an **actionable instruction**, not a description.
- Use code blocks inside steps when the step involves writing or running code.
- Link to sub-sections (e.g., `[details](#section-name)`) for lengthy explanations instead of embedding them inline.

---

## 6. Code Examples

- Always specify the language in fenced code blocks (e.g., ` ```go `, ` ```bash `).
- Code examples should be **runnable** whenever possible — prefer complete snippets over fragments.
- Add inline comments to highlight the important parts of the snippet.
- Do **not** include unrelated boilerplate that distracts from the point being demonstrated.

---

## 7. Heading Hierarchy

| Level | Usage |
|-------|-------|
| `#` (H1) | Document title — **one per file** |
| `##` (H2) | Major sections (e.g., `Description`, `Workflow`, a top-level directory) |
| `###` (H3) | Subsections (e.g., a subdirectory, a grouped set of steps) |
| `####` (H4) | Minor details within a subsection |

- Never skip heading levels (e.g., don't jump from `##` to `####`).

---

## 8. Keep Documentation in Sync

- Documentation must reflect the current state of the code.
- When you make a change that affects documentation, update all impacted files in the **same commit**.

> See [Rules.md](./RULES.md) for the full set of contribution rules.

---

## 9. Avoid Duplication

- If the same information is needed in multiple places, write it **once** and link to it from everywhere else.
- Duplicated content drifts out of sync over time.

---

## 10. Use Consistent Terminology

- Use the same term for the same concept throughout all documentation.
- Define project-specific terms on first use if they are not obvious.

| Concept | Preferred Term |
|---------|---------------|
| The object that satisfies the dependency contract | **adapter** |
| The struct holding injectable functions | **Deps** |
| A runnable example in `/examples/` | **sample** |
| A step-by-step guide | **use case** |
