# Public API

## Description
Index of all public-facing components of the library, grouped by kind, with links to their detail pages.

---

## Structs

### [api.Lib](/docs/References/PublicApi/api.Lib.md)
The main library entry point. Constructed via `lib.New`; exposes all library methods.

### [lib.ExampleLibObject](/docs/References/PublicApi/lib.ExampleLibObject.md)
An object created by the library with its dependencies automatically wired in.

---

## Functions

### [lib.New](/docs/References/PublicApi/lib.New.md)
Initializes and returns a new `Lib` instance configured with the provided dependency adapter.

### [standard.New](/docs/References/PublicApi/standard.New.md)
Creates a `deps.Deps` instance using the standard library adapter.
