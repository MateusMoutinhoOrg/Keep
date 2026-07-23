# `lib.New`

**Type:** Function

## Signature

```go
func New(d deps.Deps) Lib
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
| [`Lib`](./lib.Lib.md) | A fully initialized, ready-to-use library instance. |

## Examples

```go
package main

import (
	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	"github.com/MateusMoutinhoOrg/Keep/pkg/lib"
)

func main() {
	// 1. Initialize the desired adapter
	deps := standard.New()

	// 2. Instantiate the library using the configured dependencies
	keep := lib.New(deps)

	// The library instance is now ready for use.
	_ = keep
}
```
