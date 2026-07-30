# `deps.Deps`

**Type:** Interface

## Definition

```go
type Deps interface {
	Write(key string, value []byte) error
	WriteIfKeyNotExists(key string, value []byte) error
	WriteIfValueEquals(key string, value []byte, oldValue []byte) error
	Append(key string, value []byte) error
	InsertAt(key string, position int64, value []byte) error
	Exists(key string) (bool, error)
	Read(key string) ([]byte, error)
	ReadAt(key string, position int64, size int64) ([]byte, error)
	Delete(key string) error
	Lock(key string, time int) error
	UnLock(key string) error
}
```

## Description

The contract every storage backend must satisfy. The library performs all storage access through these methods and never touches storage directly — that is what keeps the engine inside its [closed sandbox](/docs/Explanations/SandboxIsolation.md). The behavior each method must honor — including the sentinel errors `ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, and `ErrKeyLocked` — is specified in [Required API](../RequiredApi.md).

Obtain an implementation from an adapter ([`standard.New`](./standard.New.md), [`native.New`](./native.New.md)) — the shipped ones are listed in [Adapters](../Adapters.md) — or write your own. Since `Deps` is an interface, individual behaviors are overridden by **embedding** an existing implementation and shadowing the methods you want to change (see [Dependency Mechanic](/docs/Explanations/DepsMechanic.md)).

## Examples

```go
// Embed an adapter and shadow one method to change its behavior.
type readOnly struct {
	deps.Deps // every other method is inherited
}

func (readOnly) Delete(key string) error {
	return fmt.Errorf("deletes are disabled")
}

keep := lib.New(readOnly{Deps: standard.New()})
```
