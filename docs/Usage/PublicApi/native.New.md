# `native.New`

**Type:** Function

## Signature

```go
func New() deps.Deps
```

## Description

Creates a [`deps.Deps`](./deps.Deps.md) backed by process memory. Data lives only for the lifetime of the process, which makes it ideal for tests and prototypes. All operations are safe for concurrent use through an internal mutex.

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](./deps.Deps.md) | A fully populated dependency struct ready to be passed to [`keep.New`](./keep.New.md). |

## Examples

```go
import keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/native"

deps := keep_deps.New() // zero-setup database, gone when the process exits
```
