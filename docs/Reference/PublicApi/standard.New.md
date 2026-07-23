# `standard.New` / `standard.NewWithBase`

**Type:** Functions

## Signature

```go
func New() deps.Deps
func NewWithBase(base string) deps.Deps
```

## Description

Creates a [`deps.Deps`](./deps.Deps.md) backed by the filesystem: each key becomes a file (path segments split on `/`, each segment escaped so keys can contain arbitrary characters). Data survives across process restarts.

`New` stores files relative to the current working directory; `NewWithBase` stores all keys under the given directory, creating it as needed.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `base` | `string` | (`NewWithBase` only) The directory every key is stored under. |

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](./deps.Deps.md) | A fully populated dependency struct ready to be passed to [`lib.New`](./lib.New.md). |

## Examples

```go
import "github.com/MateusMoutinhoOrg/Keep/adapters/standard"

deps := standard.New()               // relative to the working directory
deps  = standard.NewWithBase("/srv") // under a specific directory
```
