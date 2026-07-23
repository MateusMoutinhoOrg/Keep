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
| [`deps.Deps`](./deps.Deps.md) | A fully populated dependency struct ready to be passed to [`lib.New`](./lib.New.md). |

## Examples

```go
import "github.com/MateusMoutinhoOrg/Keep/adapters/native"

deps := standard.New() // zero-setup database, gone when the process exits
```
