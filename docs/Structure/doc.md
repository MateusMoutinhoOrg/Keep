# Structure

`(gen)` = written by `build`, never edited — the full list is in
[GeneratedFiles](../GeneratedFiles/doc.md).

```
adapters/  -->  sandbox/  <--  cmd/
(reaches OS)    (closed)       (wires)
```

Every line below is one entry of `AgnosConfig/structure.yaml` — add `<path>:
{description: "..."}` there, nested under `children:` of its parent, with `dir: true` on a
directory, `gen: true` on a file `build` rewrites, and `order:` to place it among its siblings
(unordered siblings follow, alphabetically).

```
AgnosConfig/         written once by `start`, read by every `build`
sandbox/             closed: imports nothing outside sandbox/, no OS packages
  new.go             (gen) New(deps) *api.Sandbox, one <x>.New<X> per api/ file
  api/               contracts only; imports nothing but sandbox/deps
    sandbox.go       (gen) one field per other file of api/, plus Deps
    databases.go     the Databases contract, plus Props, Schema, Item, Error and the objects a database hands back
    info.go          the Info contract - the library's own name and version
  deps/              one dir per capability the sandbox needs from outside; each imports nothing at all
    deps.go          (gen) one field per sub-directory, title-cased
    storagedeps/     the single-key storage backend - the whole of what Keep asks of one
  internal/          the logic; unreachable from outside the sandbox
    config/          (gen) ProjectName and Version, from AgnosConfig/project.yaml
    databases/       NewDatabases, and the DatabaseHandle a Props is bound to
    schemainstance/  one collection - insert, find, list
    schemaitem/      one record - read, update, remove, and the nested collections under it
    dense/           the key layout and value encoding of the dense record pattern
    liberror/        the constructors every *api.Error is built by
    info/            NewInfo
adapters/            the only place OS-bound and third-party code lives
  libs/              one dir per adapter, each exporting Bind and carrying an adapter.yaml
    filestorage/     storagedeps over one file per key
    memstorage/      storagedeps over a map, lost when the process exits
  availables/        one dir per selection of adapters; new.go beside each is generated from its available.yaml
    standard/        filestorage - what a program that wants a database on disk imports
    native/          memstorage - the same library with nothing written anywhere
examples/            one dir per example under lib/, each a package main checked against its golden
  lib/               <name>/example.go + props.yaml, and the result.yaml exec-test writes
docs/                one dir per doc, holding doc.md + props.yaml. README.md indexes them all
go.mod               written by `start`; add-dep and remove-dep edit its require block
README.md            (gen) `render AgnosConfig/docs/ReadmeHeader.md` + the documentation index
```

Every rule this shape has to hold to — layers, naming, generated files, docs — is in
[Rules](../Rules/doc.md); the command that makes each change is in
[Workflow](../Workflow/doc.md).
