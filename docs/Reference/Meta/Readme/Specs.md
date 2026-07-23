# Readme Specification

## Description
Defines the required structure and layout for the project's root `README.md` file.

### Rules
- The `README.md` must strictly follow the section order defined below.
- Every table must follow the formatting defined in [GeneralDoc](/docs/Reference/Meta/GeneralDoc/Specs.md).

#### Doc Index Tables
- Each documentation section (Reference, Explanation, Tutorials) indexes its documents in a **Doc Index table** with `Name` and `Description` columns.
- The Tutorials section splits its index into `####` **context groups** (e.g. *Getting Started*, *Library Development*, *Documentation*, *Templating*), one table per group.
- Each row carries an anchor ID, a relative link to the file, and a description.
- Anchor IDs follow the pattern `<section>-<kebab-case-name>` (e.g. `reference-public-api`, `explanation-deps-mechanic`, `tutorial-add-sample`).
- Descriptions are prefixed with **Reference**, **Explanation**, or **Tutorial** and kept between 50–100 characters.
- Creating, renaming, or deleting a `.md` file requires updating its Doc Index row in the same commit — for a tutorial, in the table of the context group it belongs to.

## Structure

1. **Title**: The project's name (H1).
2. **Headers/Badges**: Links to relevant external resources.
3. **Short Description**: A brief, single-sentence summary of the project.
4. **Overview**: A high-level explanation of the project's design and purpose.
5. **Quick Start**: Step-by-step instructions to install and run a basic example.
6. **Reference**:
   - Description of the section.
   - Doc indexation table for listable documents (structures, specs, API indexes).
7. **Explanation**:
   - Description of the section.
   - Doc indexation table for documents explaining mechanics and features.
8. **Tutorials**:
   - Description of the section.
   - One `####` context group per theme (e.g. *Getting Started*, *Library Development*, *Documentation*, *Templating*), each with its own doc indexation table for the workflow guides of that theme.
   - Samples table for runnable examples.
9. **Additional Theme** *(optional)*:
   - Description of the section.
   - Doc indexation table for any additional topic.
10. **License Ref**: A reference link to the project's license file.

> **Note**: For a concrete example, refer to [sample.md](/docs/Reference/Meta/Readme/sample.md).