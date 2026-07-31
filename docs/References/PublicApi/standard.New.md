# `standard.New` / `standard.NewWithBase`

**Type:** Functions

## Signature

```go
func New() deps.Deps
func NewWithBase(base string) deps.Deps
```

## Description

Creates a [`deps.Deps`](./deps.Deps.md) backed by the filesystem: each key becomes a file (path segments split on `/`, each segment escaped so keys can contain arbitrary characters). Data survives across process restarts. Each builds the adapter instance and runs every field factory over it, so each closure reads the adapter's state at call time. See [Adapters](../Adapters.md) for how it compares with the other shipped backends.

`New` stores files relative to the current working directory; `NewWithBase` stores all keys under the given directory, creating it as needed.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `base` | `string` | (`NewWithBase` only) The directory every key is stored under. |

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](./deps.Deps.md) | A filled storage contract ready to be passed to [`lib.New`](./lib.New.md). |

## Examples

```go
import keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"

deps := keepadapter.New()               // relative to the working directory
deps  = keepadapter.NewWithBase("/srv") // under a specific directory
```
