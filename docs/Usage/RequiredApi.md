# Required API

The functions a storage backend must provide to power Keep — one per field of the `Deps` struct in [pkg/deps/deps.go](../../pkg/deps/deps.go); see [Dependency Injection](DependencyInjection.md) for how to plug an implementation in.

Every function returns an `error` (`nil` on success). Where a specific condition is expected, return the matching sentinel from `pkg/deps`, wrapped so `errors.Is` matches:

| Sentinel | When |
|---|---|
| `deps.ErrKeyNotFound` | `Read`/`ReadAt`/`WriteIfValueEquals` on a missing key |
| `deps.ErrKeyAlreadyExists` | `WriteIfKeyNotExists` on an existing key |
| `deps.ErrValueMismatch` | `WriteIfValueEquals` when the current value differs |
| `deps.ErrKeyLocked` | `Lock` on a key already held |

---

### `Write(key string, value []byte) error`
Store `value` under `key`, overwriting any existing value.

### `WriteIfKeyNotExists(key string, value []byte) error`
Store `value` only if `key` does not exist yet. Returns `ErrKeyAlreadyExists` otherwise. Must be atomic.

### `WriteIfValueEquals(key string, value []byte, oldValue []byte) error`
Store `value` only if the current value equals `oldValue` (compare-and-swap). Returns `ErrKeyNotFound` if the key is missing, `ErrValueMismatch` if the value differs. Must be atomic.

### `Append(key string, value []byte) error`
Append `value` to the end of the key's current value, creating the key if needed.

### `InsertAt(key string, position int64, value []byte) error`
Insert `value` into the key's current value at byte offset `position`. Fails if `position` is out of range.

### `Exists(key string) (bool, error)`
Report whether `key` exists. A missing key is `(false, nil)`, not an error.

### `Read(key string) ([]byte, error)`
Return the value stored under `key`. Returns `ErrKeyNotFound` if missing.

### `ReadAt(key string, position int64, size int64) ([]byte, error)`
Return up to `size` bytes of the value starting at byte offset `position`.

### `Delete(key string) error`
Remove `key`. Deleting a missing key is a no-op, not an error.

### `Lock(key string, time int) error`
Acquire a lock on `key` that expires after `time` seconds. Returns `ErrKeyLocked` if the key is already locked and the lock has not expired.

### `UnLock(key string) error`
Release the lock on `key`. Unlocking a key that is not locked is a no-op.
