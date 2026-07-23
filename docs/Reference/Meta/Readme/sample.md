# Agnos

[![Go Reference](https://pkg.go.dev/badge/github.com/MateusMoutinhoOrg/Agnos.svg)](https://pkg.go.dev/github.com/MateusMoutinhoOrg/Agnos)
[![Release](https://img.shields.io/github/v/release/MateusMoutinhoOrg/Agnos)](https://github.com/MateusMoutinhoOrg/Agnos/releases/latest)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

An OS-independent Go library template demonstrating **Dependency Injection** with a clean separation between pure library logic and adapter implementations.

---

## Overview

Agnos is a structured Go template that showcases how to build libraries that are fully decoupled from their runtime dependencies. It uses a **Dependency Injection** pattern in which:

- **`/pkg/`** contains the pure library logic — it never imports concrete implementations.
- **`/adapters/`** contains opinionated, concrete implementations of the dependency interfaces.
- **`/pkg/deps/`** defines the `Deps` contract that all adapters must satisfy.

This design ensures the library remains portable, testable, and easy to extend without modifying its core.

---

## Quick Start

**1. Install the library:**
```bash
go get github.com/MateusMoutinhoOrg/Agnos
```

**2. Create a `main.go` file:**
```go
package main

import (
    "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
    "github.com/MateusMoutinhoOrg/Agnos/pkg/lib"
)

func main() {
    // 1. Create deps via an adapter (the "opinionated" part)
    deps := standard.New(3)

    // 2. Inject deps into the pure library
    l := lib.New(deps)

    // 3. Use the library — it never knows which adapter is behind the scenes
    obj := l.NewExampleObject(1, "2")
    println(obj.ExampleObjectMethod())
}
```

**3. Run:**
```bash
go run main.go
```

---

> [!IMPORTANT]
> **Must Read before contributing.** The following documents are **required reading** for every developer. Do not open a pull request or make changes without first reading them:
>
> | Document | Why it's required |
> |----------|-------------------|
> | [Rules](/docs/Reference/Meta/Readme/docs/Reference/RULES.md) | The contribution rules and guidelines that **must** be followed for any change to be accepted. |
> | [Structure](/docs/Reference/Meta/Readme/docs/Reference/Structure.md) | The project's directory layout and the purpose of each component — needed to know **where** changes belong. |
> | [Specs](/docs/Reference/Meta/Readme/docs/Reference/Specs.md) | The index of every specification — needed to know **how** the file you are about to touch must be shaped. |

### Reference Documentation

> Listable material — structures, rules, specifications, and the public API.

| Name | Description |
|:-|:-|
| <a id="reference-structure"></a>[Structure.md](/docs/Reference/Meta/Readme/docs/Reference/Structure.md) | **Reference** — The project's directory layout and the purpose of each component. |
| <a id="reference-rules"></a>[RULES.md](/docs/Reference/Meta/Readme/docs/Reference/RULES.md) | **Reference** — The binding contribution rules and their required companion updates. |
| <a id="reference-specs"></a>[Specs.md](/docs/Reference/Meta/Readme/docs/Reference/Specs.md) | **Reference** — Lists every specification and the files each one governs. |
| <a id="reference-public-api"></a>[PublicApi.md](/docs/Reference/Meta/Readme/docs/Reference/PublicApi.md) | **Reference** — Index of all public structs, functions, and methods with detail links. |

---

### Explanation Documentation

> How the project's mechanics and features work.

| Name | Description |
|:-|:-|
| <a id="explanation-deps-mechanic"></a>[DepsMechanic.md](/docs/Reference/Meta/Readme/docs/Explanation/DepsMechanic.md) | **Explanation** — Manage and customize dependencies, including overwriting and custom setups. |

---

### Tutorials

> Workflow guides, grouped by context. Each tutorial covers a single goal.

#### Getting Started

| Name | Description |
|:-|:-|
| <a id="tutorial-lib-initialization"></a>[LibInitialization.md](/docs/Reference/Meta/Readme/docs/Tutorials/LibInitialization.md) | **Tutorial** — Install the lib, create deps via an adapter, and run a first program. |
| <a id="tutorial-run-sample"></a>[RunSample.md](/docs/Reference/Meta/Readme/docs/Tutorials/RunSample.md) | **Tutorial** — Browse and run the executable samples in the examples/ directory. |

#### Library Development

| Name | Description |
|:-|:-|
| <a id="tutorial-add-lib-function"></a>[AddLibFunction.md](/docs/Reference/Meta/Readme/docs/Tutorials/AddLibFunction.md) | **Tutorial** — Add a function to pkg/lib/ and wire it to the injected deps. |
| <a id="tutorial-add-lib-object"></a>[AddLibObject.md](/docs/Reference/Meta/Readme/docs/Tutorials/AddLibObject.md) | **Tutorial** — Add an object created by the lib, with its deps wired in by the constructor. |
| <a id="tutorial-add-dependency"></a>[AddDependency.md](/docs/Reference/Meta/Readme/docs/Tutorials/AddDependency.md) | **Tutorial** — Add a field to the Deps contract and implement it in every adapter. |
| <a id="tutorial-add-adapter"></a>[AddAdapter.md](/docs/Reference/Meta/Readme/docs/Tutorials/AddAdapter.md) | **Tutorial** — Create a new opinionated implementation of the Deps contract. |
| <a id="tutorial-add-sample"></a>[AddSample.md](/docs/Reference/Meta/Readme/docs/Tutorials/AddSample.md) | **Tutorial** — Create a runnable sample in examples/ and register it in the README. |

#### Documentation

| Name | Description |
|:-|:-|
| <a id="tutorial-add-document"></a>[AddDocument.md](/docs/Reference/Meta/Readme/docs/Tutorials/AddDocument.md) | **Tutorial** — Create or update a .md file and register it in README and Structure. |
| <a id="tutorial-rename-document"></a>[RenameDocument.md](/docs/Reference/Meta/Readme/docs/Tutorials/RenameDocument.md) | **Tutorial** — Rename or move a .md file without leaving broken references behind. |
| <a id="tutorial-delete-document"></a>[DeleteDocument.md](/docs/Reference/Meta/Readme/docs/Tutorials/DeleteDocument.md) | **Tutorial** — Remove a .md file and clear every reference pointing to it. |
| <a id="tutorial-expose-public-api"></a>[ExposePublicApi.md](/docs/Reference/Meta/Readme/docs/Tutorials/ExposePublicApi.md) | **Tutorial** — Publish a lib function, object, or method in the public API index. |

#### Templating

| Name | Description |
|:-|:-|
| <a id="tutorial-rename-module"></a>[RenameModule.md](/docs/Reference/Meta/Readme/docs/Tutorials/RenameModule.md) | **Tutorial** — Rename the Go module path and update all internal imports. |
| <a id="tutorial-fork-template"></a>[ForkTemplate.md](/docs/Reference/Meta/Readme/docs/Tutorials/ForkTemplate.md) | **Tutorial** — Use this repo as a GitHub template to start a new DI library. |
| <a id="tutorial-adapt-existing-lib"></a>[AdaptExistingLib.md](/docs/Reference/Meta/Readme/docs/Tutorials/AdaptExistingLib.md) | **Tutorial** — Convert a pre-existing library to this DI structure. |

---

#### Samples
| Sample | Description |
|----------|-------------|
| [ExampleSample](/examples/ExampleSample/ExampleSample.go) | How to use the library |

---

## License

This project is licensed under the [MIT License](./LICENSE).
