# Development

## Description
Index of the documentation for contributors changing this repository: the mechanics every change runs into, the per-goal workflows, and the specifications every file must satisfy. Using the library is indexed by [LibUsage.md](/docs/Index/LibUsage.md); turning the project into a new library is indexed by [Templating.md](/docs/Index/Templating.md).

> [!IMPORTANT]
> **Read before contributing.** [Structure.md](/docs/References/Structure.md) and [Specs.md](/docs/References/Specs.md) are required reading: they say **where** a change belongs and **how** the file you touch must be shaped.

---

## Tutorials

- [HandleLibElements.md](/docs/Tutorials/HandleLibElements.md)
  - **description:** Add a function field or a whole object to the lib, then publish it
  - [Add a Library Function](/docs/Tutorials/HandleLibElements.md#add-a-library-function)
  - [Add a Library Object](/docs/Tutorials/HandleLibElements.md#add-a-library-object)
  - [Expose in the Public API](/docs/Tutorials/HandleLibElements.md#expose-in-the-public-api)
- [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md)
  - **description:** How injected deps travel the object graph, and how to add or back one
  - [Find the Dependencies You Can Use](/docs/Tutorials/HandleDependencies.md#find-the-dependencies-you-can-use)
  - [Add a Dependency](/docs/Tutorials/HandleDependencies.md#add-a-dependency)
  - [Overwrite an Adapter Function](/docs/Tutorials/HandleDependencies.md#overwrite-an-adapter-function)
  - [Create an Adapter in This Repository](/docs/Tutorials/HandleDependencies.md#create-an-adapter-in-this-repository)
  - [Create an Adapter in Your Project](/docs/Tutorials/HandleDependencies.md#create-an-adapter-in-your-project)
- [AddDatabaseOperation.md](/docs/Tutorials/AddDatabaseOperation.md)
  - **description:** Add an engine operation without breaking the dense key layout
- [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md)
  - **description:** Run the shipped Go samples, and create one for a new use case
  - [Run a Library Sample](/docs/Tutorials/HandleLibrarySamples.md#run-a-library-sample)
  - [Add a Library Sample](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample)
- [HandleDocuments.md](/docs/Tutorials/HandleDocuments.md)
  - **description:** Create, rename, move, or delete a `.md` file without breaking references
  - [Add a Document](/docs/Tutorials/HandleDocuments.md#add-a-document)
  - [Rename or Move a Document](/docs/Tutorials/HandleDocuments.md#rename-or-move-a-document)
  - [Delete a Document](/docs/Tutorials/HandleDocuments.md#delete-a-document)

---

## References

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
- [Specs/](/docs/References/Specs/)
  - **description:** The specifications themselves — always reached through `Specs.md`, never browsed
- [SandboxIsolation.md](/docs/References/SandboxIsolation.md)
  - **description:** The sandbox wall: what `sandbox/` may not import, and why effects are deps
  - [The Four Trees](/docs/References/SandboxIsolation.md#the-four-trees)
  - [What the Wall Forbids](/docs/References/SandboxIsolation.md#what-the-wall-forbids)
  - [What the Wall Forbids in the Other Direction](/docs/References/SandboxIsolation.md#what-the-wall-forbids-in-the-other-direction)
  - [Why the Entry Point Lives Inside](/docs/References/SandboxIsolation.md#why-the-entry-point-lives-inside)
- [StructContracts.md](/docs/References/StructContracts.md)
  - **description:** Why every contract is a struct of function fields, and how factories fill them
  - [The Shape](/docs/References/StructContracts.md#the-shape)
  - [Factories Fill the Fields](/docs/References/StructContracts.md#factories-fill-the-fields)
  - [Adapters Fill Their Contract the Same Way](/docs/References/StructContracts.md#adapters-fill-their-contract-the-same-way)
  - [Replacing One Behavior](/docs/References/StructContracts.md#replacing-one-behavior)
  - [No Nil Form for a Struct](/docs/References/StructContracts.md#no-nil-form-for-a-struct)
  - [What It Costs](/docs/References/StructContracts.md#what-it-costs)
- [DepsMechanic.md](/docs/References/DepsMechanic.md)
  - **description:** How storage dependencies are injected, propagated, and implemented
  - [Dependency Injection](/docs/References/DepsMechanic.md#dependency-injection)
  - [Creating Custom Dependencies](/docs/References/DepsMechanic.md#creating-custom-dependencies)
  - [Overriding Dependencies](/docs/References/DepsMechanic.md#overriding-dependencies)
- [DenseRecordPattern.md](/docs/References/DenseRecordPattern.md)
  - **description:** The key layout and write orderings the storage engine is built on
  - [The pattern](/docs/References/DenseRecordPattern.md#the-pattern)
  - [Design Goals](/docs/References/DenseRecordPattern.md#design-goals)
  - [Key Layout](/docs/References/DenseRecordPattern.md#key-layout)
  - [Normalization Rules](/docs/References/DenseRecordPattern.md#normalization-rules)
  - [Insertion](/docs/References/DenseRecordPattern.md#insertion)
  - [Deletion (Swap-With-Last)](/docs/References/DenseRecordPattern.md#deletion-swap-with-last)
  - [Updating an Indexed Field](/docs/References/DenseRecordPattern.md#updating-an-indexed-field)
  - [Lookup](/docs/References/DenseRecordPattern.md#lookup)
  - [Concurrency and Atomicity](/docs/References/DenseRecordPattern.md#concurrency-and-atomicity)
  - [Recovery](/docs/References/DenseRecordPattern.md#recovery)
  - [Invariants](/docs/References/DenseRecordPattern.md#invariants)
- [RequiredApi.md](/docs/References/RequiredApi.md)
  - **description:** The contract each `Deps` field must honor to back the library
  - [Sentinel errors](/docs/References/RequiredApi.md#sentinel-errors)
  - [Fields](/docs/References/RequiredApi.md#fields)
- [PublicApi.md](/docs/References/PublicApi.md)
  - **description:** Index of every public-facing entry of the library, grouped by role
  - [Structs](/docs/References/PublicApi.md#structs)
  - [Functions](/docs/References/PublicApi.md#functions)
  - [Fields](/docs/References/PublicApi.md#fields)
- [Adapters.md](/docs/References/Adapters.md)
  - **description:** Every shipped storage backend you can inject, and when to use each
  - [Available Adapters](/docs/References/Adapters.md#available-adapters)
- [ApiSamplesList.md](/docs/References/ApiSamplesList.md)
  - **description:** Every Go sample shipped in `examples/libraryExamples/`
  - [Examples](/docs/References/ApiSamplesList.md#examples)
