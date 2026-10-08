# `deps.StorageDeps`

`sandbox/deps/storagedeps`

## `Contract`

Contract is the storage library injected whole as the Deps.StorageDeps field. The first group of fields writes, the second reads, and the last two are the optional advisory lease a multi-writer backend can offer.

| Field | Type | Description |
| --- | --- | --- |
| `Write` | `func(key []string, value []byte) error` | Write stores value under key, overwriting any current value and creating the key when it is absent. |
| `WriteIfAbsent` | `func(key []string, value []byte) (written bool, err error)` | WriteIfAbsent stores value only when key holds nothing. written is false when the key already exists, which is not an error. |
| `WriteIfValueEquals` | `func(key []string, value []byte, oldValue []byte) (written bool, err error)` | WriteIfValueEquals stores value only when the current value of key is exactly oldValue. written is false when the key is absent or holds something else, which is not an error. |
| `Append` | `func(key []string, value []byte) error` | Append adds value to the end of the current value of key, creating the key when it is absent. |
| `InsertAt` | `func(key []string, position int64, value []byte) error` | InsertAt splices value into the current value of key at position, counted in bytes from the start. A position beyond the current length is an error: it would leave a hole. |
| `Exists` | `func(key []string) (exists bool, err error)` | Exists reports whether key currently holds a value. |
| `Read` | `func(key []string) (value []byte, found bool, err error)` | Read returns the whole value of key. found is false when the key holds nothing, and value is then nil. |
| `ReadAt` | `func(key []string, position int64, size int64) (value []byte, found bool, err error)` | ReadAt returns at most size bytes of the value of key, starting at position, counted in bytes from the start. A range reaching past the end is truncated rather than refused. found is false when the key holds nothing. |
| `Delete` | `func(key []string) error` | Delete removes key. Removing a key that holds nothing is not an error: absent before and absent after is the same outcome. |
| `Lock` | `func(key []string, seconds int) (locked bool, err error)` | Lock takes an advisory lease on key for seconds seconds. locked is false when someone else already holds a lease that has not expired, which is not an error. Keep never calls it itself — a database is written by one writer at a time — so a backend with no leases may report locked == true and do nothing. |
| `Unlock` | `func(key []string) error` | Unlock releases a lease taken by Lock. Releasing a lease nobody holds is not an error. |

[every contract](doc.md)
