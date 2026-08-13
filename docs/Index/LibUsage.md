# Library Usage

## Description
Index of the documentation for developers consuming Keep as a Go library: wiring an adapter into the sandbox, describing a database, operating on records, and looking up the public API. Changing the library is indexed by [Development.md](/docs/Index/Development.md); turning it into a library of your own is indexed by [Templating.md](/docs/Index/Templating.md).

The library is always built the same way: an adapter produces a `deps.Deps`, `lib.New` injects it into the closed sandbox, and the returned `api.Lib` carries every behavior.

---

## Tutorials

- [LibInitialization.md](/docs/Tutorials/LibInitialization.md)
  - **description:** Install the lib, create deps via an adapter, and run a first program
- [DefineDatabase.md](/docs/Tutorials/DefineDatabase.md)
  - **description:** Describe a database with its collections and open it in a program
- [AddSchemaField.md](/docs/Tutorials/AddSchemaField.md)
  - **description:** Add a field to a collection that already holds stored records
- [AddNestedCollection.md](/docs/Tutorials/AddNestedCollection.md)
  - **description:** Give a record its own nested collection of sub-records
- [HandleLibrarySamples.md](/docs/Tutorials/HandleLibrarySamples.md)
  - **description:** Run the shipped Go samples, and create one for a new use case
  - [Run a Library Sample](/docs/Tutorials/HandleLibrarySamples.md#run-a-library-sample)
  - [Add a Library Sample](/docs/Tutorials/HandleLibrarySamples.md#add-a-library-sample)

---

## References

- [PublicApi.md](/docs/References/PublicApi.md)
  - **description:** Index of every public-facing entry of the library, grouped by role
  - [Structs](/docs/References/PublicApi.md#structs)
  - [Functions](/docs/References/PublicApi.md#functions)
  - [Fields](/docs/References/PublicApi.md#fields)
- [Schemas.md](/docs/References/Schemas.md)
  - **description:** How collections, field types, and nested sub-databases are described
  - [Collections](/docs/References/Schemas.md#collections)
  - [Fields (`keeptypes.Item`)](/docs/References/Schemas.md#fields-keeptypesitem)
  - [Example](/docs/References/Schemas.md#example)
  - [Sub-databases](/docs/References/Schemas.md#sub-databases)
  - [Naming rules](/docs/References/Schemas.md#naming-rules)
- [Records.md](/docs/References/Records.md)
  - **description:** Creating, finding, reading, updating, deleting, and listing records
  - [Starting point](/docs/References/Records.md#starting-point)
  - [Create — `NewItem`](/docs/References/Records.md#create--newitem)
  - [Find — `FindByKey`](/docs/References/Records.md#find--findbykey)
  - [Read — `Get`](/docs/References/Records.md#read--get)
  - [Update — `Update`](/docs/References/Records.md#update--update)
  - [Delete — `Remove`](/docs/References/Records.md#delete--remove)
  - [List — `ListAll` and `List`](/docs/References/Records.md#list--listall-and-list)
  - [Sub-databases — `NewSubItem` and `ListAll(field)`](/docs/References/Records.md#sub-databases--newsubitem-and-listallfield)
  - [Other helpers](/docs/References/Records.md#other-helpers)
  - [Concurrency](/docs/References/Records.md#concurrency)
- [Errors.md](/docs/References/Errors.md)
  - **description:** The error types operations return, and how a caller reacts to each
  - [The Error struct](/docs/References/Errors.md#the-error-struct)
  - [Error types](/docs/References/Errors.md#error-types)
  - [Reacting to an error](/docs/References/Errors.md#reacting-to-an-error)
  - [Special cases](/docs/References/Errors.md#special-cases)
- [Adapters.md](/docs/References/Adapters.md)
  - **description:** Every shipped storage backend you can inject, and when to use each
  - [Available Adapters](/docs/References/Adapters.md#available-adapters)
- [DepsMechanic.md](/docs/References/DepsMechanic.md)
  - **description:** Choosing a backend, overriding one dep, or writing your own adapter
  - [Dependency Injection](/docs/References/DepsMechanic.md#dependency-injection)
  - [Creating Custom Dependencies](/docs/References/DepsMechanic.md#creating-custom-dependencies)
  - [Overriding Dependencies](/docs/References/DepsMechanic.md#overriding-dependencies)
- [RequiredApi.md](/docs/References/RequiredApi.md)
  - **description:** The contract each `Deps` field must honor to back the library
  - [Sentinel errors](/docs/References/RequiredApi.md#sentinel-errors)
  - [Fields](/docs/References/RequiredApi.md#fields)
- [ApiSamplesList.md](/docs/References/ApiSamplesList.md)
  - **description:** Every Go sample shipped in `examples/libraryExamples/`
  - [Examples](/docs/References/ApiSamplesList.md#examples)
