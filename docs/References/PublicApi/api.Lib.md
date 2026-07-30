# `api.Lib`

**Type:** Interface

## Definition

```go
type Lib interface {
	NewDatabase(props api.Props) KeepDatabase
}
```

## Description

The library entry point, handed back by [`lib.New`](./lib.New.md). It holds the injected dependency adapter and creates databases with it wired in. The struct implementing it lives in `sandbox/internal/lib/` and is unreachable from outside the sandbox — callers only ever see this interface.

## Methods

### `NewDatabase`

```go
func NewDatabase(props api.Props) api.KeepDatabase
```

Creates a [`KeepDatabase`](./api.KeepDatabase.md) from a [`Props`](./api.Props.md) description, with the lib's deps wired in.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `props` | [`api.Props`](./api.Props.md) | The database description: key prefix and schemas. |

| Returns | Description |
| :--- | :--- |
| [`api.KeepDatabase`](./api.KeepDatabase.md) | A database ready to hand out its collections via `GetSchema`. |

## Examples

```go
import (
	"github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	lib "github.com/MateusMoutinhoOrg/Keep/sandbox"
)

keep := lib.New(standard.New())
db := keep.NewDatabase(Props)
```
