# `native.New`

**Type:** Function

## Signature

```go
func New() deps.Deps
```

## Description

Creates a [`deps.Deps`](./deps.Deps.md) backed by process memory. Data lives only for the lifetime of the process, which makes it ideal for tests and prototypes. All operations are safe for concurrent use through an internal mutex. See [Adapters](../Adapters.md) for how it compares with the other shipped backends.

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](./deps.Deps.md) | An implementation of the storage contract ready to be passed to [`lib.New`](./lib.New.md). |

## Examples

```go
import "github.com/MateusMoutinhoOrg/Keep/adapters/native"

deps := native.New() // zero-setup database, gone when the process exits
```
