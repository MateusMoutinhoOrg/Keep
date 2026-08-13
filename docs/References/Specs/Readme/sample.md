# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

An OS-independent Go library template demonstrating **Dependency Injection** — the whole behavior lives inside a closed library.

---

## Overview

Agnos is a structured Go template that showcases how to build libraries that are fully decoupled from their runtime dependencies. It uses a **Dependency Injection** pattern in which:

- **`/sandbox/contracts/`** defines the `Deps` contract every adapter must fill and the `api` structs the library hands back.
- **`/adapters/`** contains opinionated, concrete implementations of the `Deps` contract.
- **`/sandbox/lib/`** contains the pure library logic as factories filling the `api` contract structs — it never imports concrete implementations.
- **`/sandbox/new.go`** is the entry point: it takes a `Deps` and returns an `api.Lib`, whose fields are every behavior the library offers.
- **`/examples/libraryExamples/`** are runnable programs: each wires an adapter into the lib and calls those fields.

This design ensures the library remains portable, testable, and easy to extend without modifying its core.

---

## Doc Index

Documentation is split into three themes, one index page each under `docs/Index/`, listing that theme's **Tutorials** — step-by-step workflows — and its **References** — explanations and lookups. Start from the theme index matching what you want to do.

| Theme | Description |
| --- | --- |
| [Library Usage](/docs/Index/LibUsage.md) | For library consumers: installing the module, creating deps, and calling the Go API. |
| [Development](/docs/Index/Development.md) | For contributors: the rules, the mechanics, the workflows, and the specifications. |
| [Templating](/docs/Index/Templating.md) | For template users: forking, renaming, and adapting this structure into a new library. |

New here? [Library Usage → LibInitialization.md](/docs/Tutorials/LibInitialization.md) initializes the library.

---

## License

This project is licensed under the [MIT License](./LICENSE).
