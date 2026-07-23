# TutorialDocs Specification

## Description
Defines the required shape of a **Tutorial** page — any `.md` file under `docs/Tutorials/`. A tutorial is a workflow guide for a single goal (e.g. `AddSample.md`, `AddDependency.md`), built from actionable numbered steps.

### Rules
- Every page must comply with [GeneralDoc](/docs/Reference/Meta/GeneralDoc/Specs.md).
- **One goal per page**: a tutorial has exactly one `## Workflow`. A guide covering more than one goal must be split into one page per goal, cross-linked by their steps.
- The title names that single goal as an action (e.g. `Add a Sample`, not `Handle Samples`).
- A tutorial page must contain a `## Description` and a `## Workflow` with **actionable**, numbered steps.
- An optional `### Rules` section states constraints specific to that tutorial.
- Steps prescribe actions, not descriptions; use fenced code blocks when a step involves writing or running code.
- Background explanations belong in `docs/Explanation/` — link to them instead of embedding them.
- Every new page must be registered in the [README.md](/README.md) Doc Index.

## Structure
1. **Title** (H1): the goal, phrased as an action.
2. **`## Description`**: one short paragraph on what the tutorial covers, linking to the neighbouring tutorials it does *not* cover.
3. **`### Rules`** *(optional)*: constraints for this tutorial.
4. **`---`**: horizontal rule separating the header from the workflow.
5. **`## Workflow`**: numbered, actionable steps, with fenced code blocks where a step involves code, and links to other tutorials for any step that is itself a separate goal.

> **Note**: For a concrete example, refer to [sample.md](/docs/Reference/Meta/TutorialDocs/sample.md).
