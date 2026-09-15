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
| `memstorage` | `native` | a map | no |

## Keys

A key is an opaque, slash-separated string the sandbox builds — the layout is in
[DenseRecordPattern](../DenseRecordPattern/doc.md). An adapter may escape it however its
backend requires, on one condition: **the escaping is injective**. Two different keys must
never resolve to the same place. `filestorage` escapes each slash-separated segment on its
own, which also keeps a key from reaching outside its base directory.

A value is an arbitrary byte slice. Keep stores short ASCII in practice — decimal integers
and field values — but nothing in the contract bounds either.

## No sentinel errors

**No field reports an expected condition as an error.** That is the one rule the whole
contract turns on, and why the package imports nothing at all:

| Condition | Reported as | Never as |
|---|---|---|
| the key holds nothing | `found == false` | an error |
| a conditional write did not apply | `written == false` | an error |
| someone else holds the lease | `locked == false` | an error |
| deleting a key that holds nothing | `nil` | an error |
| the backend actually failed | a non-nil `error` | |

An `error` any of these returns reaches a caller as
[`api.Internal`](../Errors/doc.md) carrying its `Error()` text, so returning one for an
absent key turns an ordinary lookup into a reported failure.

## Field by field

| Field | Must |
|---|---|
| `Write` | overwrite an existing value, create an absent key, and create whatever the backend needs around it (a parent directory, a bucket) |
| `WriteIfKeyNotExists` | be atomic against a concurrent writer where the backend can be, and report `written == false` for a key that already exists |
| `WriteIfValueEquals` | compare the whole current value byte for byte; report `written == false` both for an absent key and for a different value |
| `Append` | create the key when it is absent, so appending to nothing yields the value |
| `InsertAt` | count `position` in bytes from the start; refuse a position past the current length — it would leave a hole |
| `Exists` | answer without reading the value, where the backend allows it |
| `Read` | return the whole value, and `nil` with `found == false` for a key that holds nothing |
| `ReadAt` | count `position` in bytes; truncate a range reaching past the end rather than refusing it |
| `Delete` | succeed on a key that holds nothing |
| `Lock` | take an advisory lease expiring after `seconds`; a backend with no leases may report `locked == true` and do nothing |
| `UnLock` | succeed when nobody holds the lease |

Keep itself calls `Write`, `Read`, `Exists` and `Delete`. The other seven are part of the
contract so that a backend is described once and completely, and so a program holding the
same `Deps` can use them directly; a database is written by one writer at a time, so Keep
never takes a lease of its own.

## What Keep does not require

| Not required | Because |
|---|---|
| listing, prefix scan, range query | the position list stands in for all three |
| transactions, atomic batches | the write orderings of the dense record pattern make a partial write recoverable |
| ordering between keys | every operation addresses one key |
| a value size limit, a key length limit | keys are hashed where a value could make them grow |
| durability on return | a crash mid-sequence leaves debris that recovery deletes, never corrupt data |

A backend that does offer transactions should wrap each Keep operation in one. The orderings
then only document intent.

## Writing one

Two files and a line of yaml, the same as any adapter — the recipe is in
[Workflow](../Workflow/doc.md#add-a-dependency):

```
adapters/libs/<name>/<name>.go       func Bind(deps *deps.Deps) { deps.Storagedeps = … }
adapters/libs/<name>/adapter.yaml    dep: storagedeps
```

Then point an available at it:

```bash
agnos add-available <name>
agnos set-adapter storagedeps <name> --available <name>
agnos list-adapters
```

`adapters/libs/memstorage/memstorage.go` is the shortest complete implementation to copy
from: it is 130 lines over a map, and it fills every field.
