# `deps.Deps`

**Type:** Struct (struct of function fields)

## Definition

```go
type Deps struct {
	Write               func(key string, value []byte) error
	WriteIfKeyNotExists func(key string, value []byte) error
	WriteIfValueEquals  func(key string, value []byte, oldValue []byte) error
	Append              func(key string, value []byte) error
	InsertAt            func(key string, position int64, value []byte) error
	Exists              func(key string) (bool, error)
	Read                func(key string) ([]byte, error)
	ReadAt              func(key string, position int64, size int64) ([]byte, error)
	Delete              func(key string) error
	Lock                func(key string, time int) error
	UnLock              func(key string) error
}
```

## Description

The contract every storage backend must satisfy. It is a struct of function fields, not an interface: an adapter fills every field with the behavior it provides, and the library calls those fields directly. The library performs all storage access through them and never touches storage itself — that is what keeps the engine inside its [closed sandbox](/docs/Explanations/SandboxIsolation.md). The behavior each field must honor — including the sentinel errors `ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, and `ErrKeyLocked` — is specified in [Required API](../RequiredApi.md).

Obtain a filled `Deps` from an adapter ([`standard.New`](./standard.New.md), [`native.New`](./native.New.md)) — the shipped ones are listed in [Adapters](../Adapters.md) — or write your own. Since `Deps` is a struct, overriding one behavior is a plain field assignment rather than an embedding trick (see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md)).

## Examples

```go
// Take an adapter's Deps and override one field.
myDeps := standard.New()
myDeps.Delete = func(key string) error {
	return fmt.Errorf("deletes are disabled")
}

keep := keeplib.New(myDeps)
```
