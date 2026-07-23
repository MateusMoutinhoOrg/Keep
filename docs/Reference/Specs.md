# Specifications Index

## Description
Entry point for every specification in this project. A specification is a **description of how a file, or a kind of file, must be shaped** — its required sections, in the required order, plus the rules it must respect. Each specification pairs a `Specs.md` (the description) with a `sample` (a concrete file that satisfies it).

This index is the **only** place a specification is located from. Never browse `docs/Reference/Meta/` looking for one: find the file you are about to touch in an **Applies To** column below and follow the link.

### Rules
- Before creating or editing a file, look it up in the **Applies To** columns below. If a row matches, the file must follow that specification — see [RULES.md](/docs/Reference/RULES.md#specification-compliance).
- Every specification lives in its own directory under `docs/Reference/Meta/`, containing a `Specs.md` and a `sample` file.
- Creating, renaming, or deleting a specification requires updating this index in the same commit.

---

## Documentation Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| GeneralDoc | **Every** `.md` file in the project | [Specs](/docs/Reference/Meta/GeneralDoc/Specs.md) · [sample](/docs/Reference/Meta/GeneralDoc/sample.md) |
| Readme | Root `README.md` | [Specs](/docs/Reference/Meta/Readme/Specs.md) · [sample](/docs/Reference/Meta/Readme/sample.md) |
| Rules | `docs/Reference/RULES.md` | [Specs](/docs/Reference/Meta/Rules/Specs.md) · [sample](/docs/Reference/Meta/Rules/sample.md) |
| Structure | `docs/Reference/Structure.md` | [Specs](/docs/Reference/Meta/Structure/Specs.md) · [sample](/docs/Reference/Meta/Structure/sample.md) |
| ReferenceDocs | Any other `.md` under `docs/Reference/`, except this index and `docs/Reference/Meta/` | [Specs](/docs/Reference/Meta/ReferenceDocs/Specs.md) · [sample](/docs/Reference/Meta/ReferenceDocs/sample.md) |
| ExplanationDocs | Any `.md` under `docs/Explanation/` | [Specs](/docs/Reference/Meta/ExplanationDocs/Specs.md) · [sample](/docs/Reference/Meta/ExplanationDocs/sample.md) |
| TutorialDocs | Any `.md` under `docs/Tutorials/` | [Specs](/docs/Reference/Meta/TutorialDocs/Specs.md) · [sample](/docs/Reference/Meta/TutorialDocs/sample.md) |

GeneralDoc applies on top of the others: a tutorial follows **both** GeneralDoc and TutorialDocs.

---

## Code Specifications

| Spec | Applies To | Links |
|------|------------|-------|
| Deps | `pkg/deps/deps.go` | [Specs](/docs/Reference/Meta/Deps/Specs.md) · [sample](./Meta/Deps/sample.go) |
| Adapters | `adapters/<name>/<name>.go` | [Specs](/docs/Reference/Meta/Adapters/Specs.md) · [sample](./Meta/Adapters/sample.go) |
| LibFunctions | Public functions in `pkg/lib/` | [Specs](/docs/Reference/Meta/LibFunctions/Specs.md) · [sample](./Meta/LibFunctions/sample.go) |
| LibObjects | Objects created by the lib in `pkg/lib/` | [Specs](/docs/Reference/Meta/LibObjects/Specs.md) · [sample](./Meta/LibObjects/sample.go) |
| Examples | `examples/<example>/<example>.go` | [Specs](/docs/Reference/Meta/Examples/Specs.md) · [sample](./Meta/Examples/sample.go) |

---

## Workflow

1. Locate the file you are about to create or edit in an **Applies To** column above.
2. If no row matches, no specification governs the file — follow [Structure.md](/docs/Reference/Structure.md) and, for `.md` files, [GeneralDoc](/docs/Reference/Meta/GeneralDoc/Specs.md).
3. If a row matches, read its `Specs.md` and reproduce the required **Structure** section by section.
4. Use the linked `sample` as the reference implementation.
5. Apply the companion updates required by [RULES.md](/docs/Reference/RULES.md).
