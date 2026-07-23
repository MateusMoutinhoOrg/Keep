# `deps.Deps`

**Type:** Struct

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

The struct of injectable storage functions every backend must populate. The library performs all storage access through these functions and never touches storage directly. The behavior each function must honor — including the sentinel errors `ErrKeyNotFound`, `ErrKeyAlreadyExists`, `ErrValueMismatch`, and `ErrKeyLocked` — is specified in [Required API](../RequiredApi.md).

Obtain a populated `Deps` from an adapter ([`standard.New`](./standard.New.md), [`native.New`](./native.New.md)), build one by hand, or overwrite individual fields of an adapter's result (see [Dependency Mechanic](../DepsMechanic.md)).

## Examples

```go
myDeps := standard.New()

// Individual functions are plain values and can be replaced:
myDeps.Delete = func(key string) error {
	return fmt.Errorf("deletes are disabled")
}

keep := keep_lib.New(myDeps)
```
