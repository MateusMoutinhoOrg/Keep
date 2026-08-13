# Templating

## Description
Index of the documentation for people using Keep as a **template**: turning this repository into their own library, or converting an existing library to the same dependency-injection structure. Contributing to Keep itself is indexed by [Development.md](/docs/Index/Development.md); consuming it as a library is indexed by [LibUsage.md](/docs/Index/LibUsage.md).

Pick **one** of the two workflows and follow it end to end. Both take each file's fate from the same per-file action table.

---

## Tutorials

- [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md)
  - **description:** **Start here for a new library**: use this repo as a GitHub template
- [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md)
  - **description:** **Start here for an existing library**: convert it to this DI structure
- [RenameModule.md](/docs/Tutorials/RenameModule.md)
  - **description:** Rename the Go module path and update every internal import

---

## References

- [TemplateFileActions.md](/docs/References/TemplateFileActions.md)
  - **description:** The per-file action both workflows follow: copy, create, rewrite, or delete
  - [Copy](/docs/References/TemplateFileActions.md#copy)
  - [Create](/docs/References/TemplateFileActions.md#create)
  - [Rewrite](/docs/References/TemplateFileActions.md#rewrite)
  - [Delete](/docs/References/TemplateFileActions.md#delete)
- [Structure.md](/docs/References/Structure.md)
  - **description:** The project's schema: which kind of file lives where, and its spec
  - [Root](/docs/References/Structure.md#root)
  - [`/sandbox/`](/docs/References/Structure.md#sandbox)
  - [`/adapters/`](/docs/References/Structure.md#adapters)
  - [`/examples/libraryExamples/`](/docs/References/Structure.md#exampleslibraryexamples)
  - [`/tests/`](/docs/References/Structure.md#tests)
  - [`/docs/`](/docs/References/Structure.md#docs)
- [Specs.md](/docs/References/Specs.md)
  - **description:** Index of every specification and the files each one governs
  - [Documentation Specifications](/docs/References/Specs.md#documentation-specifications)
  - [Code Specifications](/docs/References/Specs.md#code-specifications)
