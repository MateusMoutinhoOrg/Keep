

# Contribution Rules

Rules to follow when contributing to this project. Every file must also be shaped by the specification that governs it — locate it in [Specs.md](/docs/Reference/Specs.md).

---

## Tutorials Guide
Before making anything, Read the [README.md](/README.md) and search if there is a tutorial about what you want to do, if there is, Follow it, if there isn't,you need to create one following the spec defined in [TutorialDocs](./Meta/TutorialDocs/)


## Specification Compliance

Before creating or editing any file, read [Specs.md](/docs/Reference/Specs.md) and check whether the file matches an **Applies To** entry. If it does, create or edit it following the specification that entry points to — reproduce the shape it requires, using its `sample` as reference.

---

## File Changes

Before creating, deleting, or renaming any file or directory, read [Structure.md](/docs/Reference/Structure.md) and check whether the change affects the project structure. If it does, update [Structure.md](/docs/Reference/Structure.md) in the same commit.

---

## Specification Changes

When you create, delete, or rename a specification inside [Meta/](./Meta),  you MUST adapt all the files that match the spec's Applies To rule, and update the index in [Specs.md](/docs/Reference/Specs.md).

---

## Documentation Changes

When you create, delete, or rename a `.md` file, update the Doc Index of [README.md](/README.md).

---

## Sample Changes

When you create, delete, or rename a sample (any file inside [examples/](../../examples)), update the Samples section of [README.md](/README.md).
