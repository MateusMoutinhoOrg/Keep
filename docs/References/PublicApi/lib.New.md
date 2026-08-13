# `lib.New`

**Type:** Function

## Signature

```go
func New(d deps.Deps) api.Lib
```

## Description

Injects a [`deps.Deps`](./deps.Deps.md) implementation into the library and returns the [`api.Lib`](./api.Lib.md) entry point. It is the single wiring point consumers touch, and the only exported symbol of the `sandbox` package. It delegates to `sandbox/lib/lib.New`, which stores the deps on the struct and runs every field factory over it.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `d` | [`deps.Deps`](./deps.Deps.md) | A filled storage contract, usually from an [adapter](../Adapters.md). |

## Returns

| Type | Description |
| :--- | :--- |
| [`api.Lib`](./api.Lib.md) | A fully initialized, ready-to-use library instance. |

## Examples

```go
package main

import (
	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

func main() {
	// 1. Initialize the desired adapter
	deps := keepadapter.New()

	// 2. Instantiate the library using the configured dependencies
	keep := keeplib.New(deps)

	// The library instance is now ready for use.
	_ = keep
}
```
