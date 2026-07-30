# Contribution Rules

Rules to follow when contributing to this project. Every file must also be shaped by the specification that governs it — locate it in [Specs.md](/docs/References/Specs.md).

---

## Tutorials Guide
Before making anything, read the [README.md](/README.md) and search for a tutorial about what you want to do. If there is one, follow it; if there isn't, you need to create one following the spec defined in [TutorialDocs](./Meta/TutorialDocs/).


## Specification Compliance

Before creating or editing any file, read [Specs.md](/docs/References/Specs.md) and check whether the file matches an **Applies To** entry. If it does, create or edit it following the specification that entry points to — reproduce the shape it requires, using its `sample` as reference.

---

## Sandbox Isolation

[sandbox/](/sandbox/) is a closed sandbox. No file inside it may import [adapters/](/adapters/), [examples/](/examples/), [tests/](/tests/), a third-party module, or an OS-bound standard-library package (`os`, `net`, `os/exec`, `syscall`, …). Every such effect must be declared as a method on the `Deps` contract and reached through the object's `Deps` field, following [AddDependency.md](/docs/Tutorials/AddDependency.md). The mechanic is explained in [SandboxIsolation.md](/docs/Explanations/SandboxIsolation.md).

Conversely, nothing outside the sandbox may reach into it beyond its three public packages: `sandbox` (package `lib`), `sandbox/contracts/deps`, and `sandbox/contracts/api`.

---

## Adapter Changes

When you create, delete, or rename an adapter inside [adapters/](/adapters/), update the table in [Adapters.md](/docs/References/Adapters.md) in the same commit. When you add a method to the `Deps` contract, every adapter must implement it — a missing method is a compile error for every consumer.

---

## File Changes

Before creating, deleting, or renaming any file or directory, read [Structure.md](/docs/References/Structure.md) and check whether the change affects the project structure. If it does, update [Structure.md](/docs/References/Structure.md) in the same commit.

---

## Specification Changes

When you create, delete, or rename a specification inside [Meta/](./Meta), you MUST adapt all the files that match the spec's Applies To rule, and update the index in [Specs.md](/docs/References/Specs.md).

---

## Documentation Changes

When you create, delete, or rename a `.md` file, update the Doc Index of [README.md](/README.md).

---

## Sample Changes

When you create, delete, or rename a sample (any file inside [examples/](/examples/)), update the Samples section of [README.md](/README.md).
