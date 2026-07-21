# `keep.New`

**Type:** Function

## Signature

```go
func New(d deps.Deps) *KeepLib
```

## Description

Initializes and returns a new library instance configured with the specified dependency adapter.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `d` | [`deps.Deps`](./deps.Deps.md) | A populated `Deps` struct providing the storage functions for the library. |

## Returns

| Type | Description |
| :--- | :--- |
| [`*KeepLib`](./keep.KeepLib.md) | A fully initialized, ready-to-use library instance. |

## Examples

```go
package main

import (
	keep_deps "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keep_lib "github.com/MateusMoutinhoOrg/Keep/pkg/keep"
)

func main() {
	// 1. Initialize the desired adapter
	deps := keep_deps.New()

	// 2. Instantiate the library using the configured dependencies
	keep := keep_lib.New(deps)

	// The library instance is now ready for use.
	_ = keep
}
```
