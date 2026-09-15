package storagedeps

// This package is the sandbox's *copy* of the api a single-key storage
// backend exposes — the same mechanic as hashdeps, std and stringsdeps, for
// the same reason: reaching a file, a socket or a remote store is an
// OS-bound effect, so `os`, `net` and every driver package may not appear
// inside the sandbox. The contract is restated here, and the adapter — which
// lives outside the sandbox — is what fills it.
//
// It is also the whole of what Keep asks of a backend. Every field addresses
// exactly one key: nothing here lists keys, scans a prefix or queries a
// range, which is what lets a schema database run over a plain key-value
// store, a bucket, a cache or a directory of files. The key layout that
// makes that possible is the Dense Record Pattern, documented in
// docs/DenseRecordPattern.
//
// A key is an opaque, slash-separated string built by the sandbox. An
// adapter may escape it however its backend requires, as long as the
// escaping is injective: two different keys must never collide.
//
// No field reports an expected condition as an error. A key that is absent
// comes back as found == false, a conditional write that did not apply as
// written == false, and a lease someone else holds as locked == false — the
// error result is reserved for a backend that actually failed. This is why
// the contract needs no sentinel error values, and so no import at all.

// Sandbox is the storage library injected whole as the Deps.Storagedeps
// field. The first group of fields writes, the second reads, and the last
// two are the optional advisory lease a multi-writer backend can offer.
type Sandbox struct {
	// Write stores value under key, overwriting any current value and
	// creating the key when it is absent.
	Write func(key string, value []byte) error

	// WriteIfKeyNotExists stores value only when key holds nothing.
	// written is false when the key already exists, which is not an error.
	WriteIfKeyNotExists func(key string, value []byte) (written bool, err error)

	// WriteIfValueEquals stores value only when the current value of key is
	// exactly old_value. written is false when the key is absent or holds
	// something else, which is not an error.
	WriteIfValueEquals func(key string, value []byte, old_value []byte) (written bool, err error)

	// Append adds value to the end of the current value of key, creating
	// the key when it is absent.
	Append func(key string, value []byte) error

	// InsertAt splices value into the current value of key at position,
	// counted in bytes from the start. A position beyond the current length
	// is an error: it would leave a hole.
	InsertAt func(key string, position int64, value []byte) error

	// Exists reports whether key currently holds a value.
	Exists func(key string) (exists bool, err error)

	// Read returns the whole value of key. found is false when the key
	// holds nothing, and value is then nil.
	Read func(key string) (value []byte, found bool, err error)

	// ReadAt returns at most size bytes of the value of key, starting at
	// position, counted in bytes from the start. A range reaching past the
	// end is truncated rather than refused. found is false when the key
	// holds nothing.
	ReadAt func(key string, position int64, size int64) (value []byte, found bool, err error)

	// Delete removes key. Removing a key that holds nothing is not an
	// error: absent before and absent after is the same outcome.
	Delete func(key string) error

	// Lock takes an advisory lease on key for seconds seconds. locked is
	// false when someone else already holds a lease that has not expired,
	// which is not an error. Keep never calls it itself — a database is
	// written by one writer at a time — so a backend with no leases may
	// report locked == true and do nothing.
	Lock func(key string, seconds int) (locked bool, err error)

	// UnLock releases a lease taken by Lock. Releasing a lease nobody holds
	// is not an error.
	UnLock func(key string) error
}
