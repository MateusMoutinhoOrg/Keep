# StorageContract

`sandbox/deps/storagedeps` is the whole of what Keep asks of a backend: eleven functions,
each addressing exactly one key. Nothing lists keys, scans a prefix or queries a range,
which is what lets the same library run over a directory of files, a bucket, a cache, a row
store or a remote key-value service. Signatures are in
[PublicApi](../PublicApi/doc.md#depsstoragedeps); this page is what an implementation has to
guarantee beyond them.

Two adapters ship with the repository, and [Adapters](../Adapters/doc.md) is how a program
picks between them:

| Adapter | Available | Backed by | Survives the process |
|---|---|---|---|
| `filestorage` | `standard` | one file per key, under the working directory | yes |
| `memstorage` | `memory` | a map | no |

## Keys

A key is a **list of opaque segments** the sandbox builds — `[]string`, not one joined
string — and the layout is in [DenseRecordPattern](../DenseRecordPattern/doc.md). The
sandbox hands the segments apart precisely so that no separator has to survive inside them:
a field name or a path component holding a slash can no longer be read back as a boundary
and collide with another key.

An adapter flattens the list the way its backend wants, conventionally joining with a
slash — `["a", "b.txt"]` becomes `a/b.txt` — on one condition: **the flattening is
injective**. Two different segment lists must never resolve to the same place — not even on
a backend that ignores case, or that reads `.`, `..` or an empty segment as a path of its
own. Both shipped adapters encode each segment *before* joining, so `["a/b"]` and
`["a", "b"]` stay two different keys.

A value is an arbitrary byte slice. Keep stores short ASCII in practice — decimal integers
and field values — but nothing in the contract bounds either.

## No sentinel errors

**No field reports an expected condition as an error.** That is the one rule the whole
contract turns on, and why the package imports nothing at all:

| Condition | Reported as | Never as |
|---|---|---|
| the key holds nothing | `found == false` | an error |
| a conditional write did not apply | `written == false` | an error |
| someone holds the lease | `locked == false` | an error |
| deleting a key that holds nothing | `nil` | an error |
| the backend actually failed | a non-nil `error` | |

An `error` any of these returns reaches a caller as
[`api.Internal`](../Errors/doc.md) carrying its `Error()` text, so returning one for an
absent key turns an ordinary lookup into a reported failure.

## Field by field

| Field | Must |
|---|---|
| `Write` | overwrite an existing value, create an absent key, and create whatever the backend needs around it (a parent directory, a bucket). **Atomic**: a reader, and a process that crashed part-way, sees the old value or the new one — never an empty or a torn one |
| `WriteIfAbsent` | be atomic against a concurrent writer where the backend can be, and report `written == false` for a key that already exists |
| `WriteIfValueEquals` | compare the whole current value byte for byte; report `written == false` both for an absent key and for a different value |
| `Append` | create the key when it is absent, so appending to nothing yields the value |
| `InsertAt` | count `position` in bytes from the start; refuse a position past the current length — it would leave a hole |
| `Exists` | answer without reading the value, where the backend allows it |
| `Read` | return the whole value, and `nil` with `found == false` for a key that holds nothing |
| `ReadAt` | count `position` in bytes; refuse a negative size; truncate a range reaching past the end rather than refusing it |
| `Delete` | succeed on a key that holds nothing, and leave behind nothing the backend created only for that key (a directory, a bucket prefix) |
| `Lock` | take an advisory lease expiring after `seconds`, refused (`locked == false`) while anyone holds one that has not expired — this process included. A lease lives beside the keys: locking `["a"]` never reads, writes or shadows the key `["a"]`. A backend with no leases may report `locked == true` and do nothing |
| `Unlock` | release the lease whoever took it, and succeed when nobody holds it |

Keep itself calls `Write`, `Read`, `Exists`, `Delete`, `Lock` and `Unlock`. It takes a lease
on the prefix of a top-level collection around every write to it — see
[Concurrency](../DenseRecordPattern/doc.md#concurrency) — so on a backend whose leases work
across processes, several processes can write one database. On a backend with no leases,
only one process may. The other five fields are part of the contract so that a backend is
described once and completely, and so a program holding the same `Deps` can use them
directly.

## `filestorage`

One file per key, one directory per segment, under the base directory — `.` for the
`standard` binding, so where a database lands depends on the working directory of the
process. `filestorage.New(base)` roots it anywhere else.

| Concern | What it does |
|---|---|
| encoding | `a-z`, `0-9`, `-` and `_` stay as they are; every other byte becomes `%XX` with upper-case hex digits; `""` becomes `%`; an encoding longer than 200 bytes becomes `%%` + the SHA-256 of the segment. Injective on a filesystem that ignores case, and no segment can be `.` or `..` or escape the base |
| writes | to a temporary file `.keep-tmp-*` beside the key, flushed, then renamed over it. `NewWithOptions(base, Options{NoSync: true})` skips the flush of the file and of its directory: still atomic, no longer durable across a power loss |
| leases | a `<key>.keeplock` file beside the key's own path, created whole by a hard link. An expired one is taken over by one process at a time, behind a `.takeover` guard file |
| permissions | files `0600`, directories `0700` |
| empty directories | removed as soon as the last key under them is deleted, so a removed record leaves nothing behind on disk |
| prefixes | a name is a file or a directory, never both: `["a"]` cannot hold a value while `["a", "b"]` does. `Read` and `Exists` report such a prefix as absent, `Write` to it fails. Keep's own layout never needs both |
| limits | the operating system's path length bounds how deep nested collections go: each level adds `/{id}/{field}` to every path under it, and macOS caps a path at 1024 bytes. `memstorage` has no such bound |

What it costs on an SSD under macOS, measured on 3000 records:

| Operation | `memstorage` | `filestorage`, `NoSync` | `filestorage` |
|---|---|---|---|
| insert | 14 µs | 3.4 ms | 68 ms |
| remove | 14 µs | 3.8 ms | 32 ms |
| `FindByKey` | 7 µs | 0.26 ms | 0.16 ms |
| `ListAll` + one `Get`, per record | 3.5 µs | 0.10 ms | 0.08 ms |

A flush under macOS is `F_FULLFSYNC`, which empties the drive's cache. That is what a write
costs there when it has to survive a power loss, and it dominates the durable column; an
fsync under Linux costs a fraction of it.

## What Keep does not require

| Not required | Because |
|---|---|
| listing, prefix scan, range query | the position list stands in for all three |
| transactions, atomic batches | the write orderings of the dense record pattern make a partial write recoverable |
| ordering between keys | every operation addresses one key |
| a value size limit, a key length limit | keys are hashed where a value could make them grow |

A backend that does offer transactions should wrap each Keep operation in one. The orderings
then only document intent.

## Writing one

Two files and a line of yaml, the same as any adapter — the recipe is in
[Workflow](../Workflow/doc.md#add-a-dependency):

```
adapters/impls/<name>/<name>.go      func Bind(deps *deps.Deps) { deps.StorageDeps = … }
adapters/impls/<name>/adapter.yaml   dep: storagedeps
```

Then point a binding at it:

```bash
agnos add-binding <name>
agnos set-adapter storagedeps <name> --binding <name>
agnos list-adapters
```

`adapters/impls/memstorage/memstorage.go` is the shortest complete implementation to copy
from: it is roughly 160 lines over a map, and it fills every field.
