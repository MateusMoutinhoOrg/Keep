# Keep

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Keep.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Keep)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Keep)](https://github.com/MateusMoutinhoOrg/Keep/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A storage-independent database built on top of plain key-value operations.

---

## Overview

Keep lets you define schemas with typed fields, unique indexed keys, and nested collections — and runs them over **any** backend that can read, write, and delete a key. It needs no key listing, no prefix scans, and no range queries, so it works the same over the local filesystem, memory, S3-like blob stores, or anything you can wrap in a small struct.

The core of the library lives in **`/sandbox/`**: a **closed sandbox** that reaches nothing outside itself. Everything it can do arrives through an injected `Deps`, wired entirely through **structs of function fields** rather than interfaces — see [StructContracts.md](/docs/References/StructContracts.md).

```
adapters/  ──▶  sandbox/  ◀──  examples/libraryExamples/ , tests/
(reaches the OS)  (closed)     (wire the two together)
```

Nothing is exported but the library itself. There is no binary and no command of its own: a consumer wires an adapter into the sandbox and calls the fields of the `api.Lib` it gets back.

- **`/sandbox/`**: The closed database engine taking a `Deps` and returning an `api.Lib`. It may not import an adapter, a third-party module, or any OS-bound stdlib package — see [SandboxIsolation.md](/docs/References/SandboxIsolation.md).
- **`/sandbox/contracts/deps/`**: The `Deps` struct of function fields every adapter must fill.
- **`/sandbox/contracts/api/`**: Every type the library exchanges — structs only, never interfaces, so nothing crossing the boundary is hidden behind a method set.
- **`/adapters/`**: The opinionated, concrete backends — the only place OS-bound code is allowed.
- **`/examples/libraryExamples/`**: Places where an adapter and the library are wired together.

What you get on top of that:

- **Storage independent** — bring your own backend by filling a small struct of function fields, or use the built-in ones ([filesystem](adapters/standard/), [in-memory](adapters/native/)).
- **Constant-time operations** — create, lookup by key, and delete each touch a fixed number of keys, no matter how many records exist.
- **Unique keys** — fields of type `Key` are indexed and enforced unique (case-insensitive).
- **Nested collections** — a record can own sub-databases (e.g. a user owning its sessions).

---

## Doc Index

Documentation is split into three themes, one index page each under `docs/Index/`, listing that theme's **Tutorials** — step-by-step workflows — and its **References** — explanations and lookups. Start from the theme index matching what you want to do.

| Theme | Description |
| --- | --- |
| [Library Usage](/docs/Index/LibUsage.md) | For library consumers: installing the module, describing data, and calling the Go API. |
| [Development](/docs/Index/Development.md) | For contributors: the rules, the mechanics, the per-goal workflows, and the specifications. |
| [Templating](/docs/Index/Templating.md) | For template users: forking, renaming, and adapting this structure into a new library. |

New here? [Library Usage → LibInitialization.md](/docs/Tutorials/LibInitialization.md) installs the module and runs a first program.

---

## License

This project is licensed under the [MIT License](./LICENSE).
