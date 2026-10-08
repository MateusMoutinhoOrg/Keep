# DenseRecordPattern

The key layout that lets a schema database run over a backend with nothing but `get`, `set`
and `delete`. Not needed to use Keep — [Schemas](../Schemas/doc.md) is — but it is what
`sandbox/internal/dense` and `sandbox/internal/record` implement, and what a change to
either has to preserve.

Three properties, and everything below follows from them:

1. **No listing.** Every operation is reads, writes and deletes of individual keys. Nothing
   scans a prefix, lists keys or queries a range.
2. **Fixed cost.** Insert, look-up-by-key and the removal of a record with no nested
   collection each touch a number of keys that does not grow with the size of the collection.
3. **Permanent ids.** A record's id never changes and is never reused, not even after the
   record is removed.

## Key layout

A key is a **list of segments**, never one joined string — that is what
[storagedeps](../StorageContract/doc.md) takes, and the adapter is what flattens it, by
convention with a slash. Written `a/b` below, a key is really `["a", "b"]`, so no character
of a field name or of a value can be read back as a boundary between two segments.

`{c}` is the collection prefix — the slash-separated segments of `Props.Path` followed by
the schema name for a top-level collection, `{c}/{id}/{field}` for one nested inside a
record. The two are the same thing: every key family below applies unchanged to a nested
collection.

| Key | Holds | Family |
|---|---|---|
| `{c}/size` | the number of positions of the list — the highest occupied one | metadata |
| `{c}/last-id` | the highest id ever allocated. Only grows | metadata |
| `{c}/list/{position}` | the id living at that position | position list |
| `{c}/keys/{field}/{sha256(fold(value))}` | the id holding that value for that `Key` field | unique index |
| `{c}/{id}/position` | the record's position — its back-pointer | record |
| `{c}/{id}/values/{field}` | one field value of one record | record |
| `{c}/{id}/{field}` | the prefix of a nested collection owned by the record | record |

`position` and `values` are therefore names no `Nested` field may take; `Databases.New`
refuses them.

Positions run from `1` to `size`. That density is what makes iteration possible without
listing: read `size`, then read positions 1 through `size`. It is a dense set, not an
ordered sequence — a removal reorders it. A consumer that needs insertion order stores it as
a field.

Index values are case-folded before they are hashed — Unicode's canonical caseless match,
`NFD(casefold(NFD(v)))`, which lower-cases ASCII — so a key lookup ignores case in every
script, and hashed at all so that key length stays bounded whatever a value holds.

## What is live

A write that stops part-way — a crash, or a backend failing between two writes — leaves keys
behind that look like data. Three readings tell the data from the debris, and every
operation that trusts a key reads it through one of them (`sandbox/internal/dense/live.go`):

| Thing | Live when | Reads |
|---|---|---|
| a record `id` | `{c}/{id}/position` → `p`, `1 ≤ p ≤ size`, and `{c}/list/{p}` → `id` | 3 |
| a list slot `p` | `{c}/list/{p}` → `id`, and `{c}/{id}/position` → `p` | 2 |
| an index entry | it names a live record whose current value of that field hashes back to the entry | 5 |

A key that is absent or holds something that is not a decimal id reads as "nothing here" —
never as id 0, and never as a failure. An index entry that is not valid is **free**: an
insert or an update claims it by writing over it, and a lookup through it finds nothing.

## Insert

| Step | Does | Fails with |
|---|---|---|
| 0 | validate the fields against the schema | `InvalidField`, `MissingField` |
| — | take the collection's write lock; when nested, check the owner is live | `Removed` |
| 1 | read `{c}/keys/{field}/{hash}` for every `Key` field; a valid entry is a conflict | `KeyConflict` |
| 2 | read `{c}/last-id`, write it back +1. That is the new id | |
| 3 | read `{c}/size`. The new position is size+1 | |
| 4 | write every `{c}/{id}/values/{field}`, then `{c}/{id}/position` | |
| 5 | write every `{c}/keys/{field}/{hash}` pointing at the new id | |
| 6 | write `{c}/list/{position}`, then `{c}/size` | |

Steps 0 and 1 decide every failure a caller can see but `Internal`, and they write nothing —
a refused insert leaves the database exactly as it was.

`{c}/size` is written **last, deliberately**. Until it grows, the new position lies outside
`[1, size]`, so the record is not live: no reader iterating the list observes it, and
neither `FindByID` nor `FindByKey` resolves it. A step from 3 on that fails deletes what the
insert had written, as far as the backend lets it; a crash leaves it for
[Repair](#recovery). The id it consumed is not given back: a failed insert still spends its
id, which is what keeps ids unique.

## Remove — swap with last

Closing the hole by shifting would cost one write per remaining record. The last record
fills it instead.

| Step | Does |
|---|---|
| 1 | read `{c}/{id}/position` → `p`. Absent means already gone: a no-op |
| 2 | read `{c}/size`; drop every dead slot at the end of the list, so the last slot is live |
| 3 | when slot `p` holds `id`, or is dead: unless `p` is the last, write `{c}/list/{p}` = `lastID`, then `{c}/{lastID}/position` = `p` |
| 4 | … then delete `{c}/list/{size}`, write `{c}/size` = size-1 |
| 5 | for every `Key` field: read the value, delete `{c}/keys/{field}/{hash}` when it still names `id` |
| 6 | delete every value key, clear every nested collection, delete `{c}/{id}/position` last |

`{c}/last-id` is untouched. **List order is not stable across removals** — an unrelated
record moves — and that is the documented price of step 3.

Every order here is what makes a removal that stopped part-way finish when it is run again,
instead of starting another one:

- The record stops being live at the first write of step 3: its slot no longer holds it.
- The moved record is live at every point. It lives at the end until slot `p` is written,
  and at `p` from the moment its back-pointer is.
- A slot left behind at the end is dead, so the next removal's step 2 drops it.
- A slot `p` written but whose moved record still points at the end is dead too, and step 3
  of the retry refills it.
- The back-pointer goes last, so the retry still finds `p`.
- An index entry is deleted only while it names the record, because one it no longer holds
  may already belong to another record.

Step 6 clears a nested collection from the last position backwards, so no swap is ever
needed. Each record is emptied **while it is still in the list** — its data, then its
back-pointer, which kills its slot, then the slot and the size. A clear that stops part-way
and runs again therefore finds what is left at the end of the list. It then deletes the
collection's `size` and `last-id`, and recurses to any depth.

## Update

A field that carries no index is one write, after the lock and a check that the record is
live (`Removed` otherwise). A `nil` value deletes it — the value first, then, for a `Key`,
the index entry it owned.

A `Key` field moves its index entry:

| Step | Does |
|---|---|
| 1 | read `{c}/{id}/values/{field}` — the old value locates the old index entry |
| 2 | read `{c}/keys/{field}/{newHash}`. A valid entry naming another id is a `KeyConflict`, and nothing is written |
| 3 | write `{c}/keys/{field}/{newHash}` = id |
| 4 | write `{c}/{id}/values/{field}` = the new value |
| 5 | delete `{c}/keys/{field}/{oldHash}`, when it differs from the new one and still names id |

New entry before old delete: a crash part-way leaves the record reachable through one of
the two values, never through neither. Whichever entry the record does not hold the value
of is not valid, so it resolves to nothing and is free for the next writer — it can neither
return the record under a value it no longer has nor block another record from taking it.

## Read

| Operation | Costs |
|---|---|
| by id | 3 reads to check the record is live, then one per field read |
| by unique key | 5 reads — the entry, the 3 of liveness, the value it indexes — then one per field |
| one page | one read for `size`, then 2 per position in the page |
| whole collection | one read for `size`, then 2 per record |

A field is read on demand — nothing loads a whole record. A listing skips a dead slot and
never hands out the same id twice, so it never returns a record that is not there.

## Concurrency

Every write to a top-level collection — and to every collection nested under one of its
records — holds that collection's **write lock** from its first read to its last write
(`sandbox/internal/writelock`):

1. **An in-process lock**, one per top-level collection per sandbox. Goroutines of one
   program hand the collection over without polling.
2. **A storage lease**, `StorageDeps.Lock` on the collection's prefix, for 60 seconds. A
   writer in another process sharing the backend waits, polling from 1 ms up to 100 ms
   between attempts, and gives up with `Internal` after 120 s. A lease left by a crashed
   writer expires on its own.

A backend with no leases reports `locked == true` at once, and then only writers within one
process are serialized. The lease names no holder — `Unlock` takes none — so an operation
that outlived its 60 s could release a lease another writer has since taken. One operation
can come near that: removing a record whose nested collections hold thousands of records,
over a backend where a write costs milliseconds — `filestorage` with `fsync`, for one.

**Readers take no lock.** What makes them safe beside a writer:

- each key is written atomically by the backend;
- every reading above checks a key against the one it points at;
- a listing skips what does not check out.

A listing is not a snapshot: a record removed while it runs may be left out, and so may the
record that removal moved. The same holds across pages of `List`.

## Invariants

At every quiescent point after operations that all succeeded:

1. `{c}/list/{p}` exists exactly for `p` in `[1, size]`, and every slot is live.
2. For every live record, `list[position(id)] == id`.
3. Every index entry points at a live record whose current value hashes back to that entry.
4. `last-id` ≥ every id the collection has ever held.
5. No id that is not live owns a key.

A failure part-way can break 1, 3 and 5 — never 2 or 4 — and every reading above tolerates
what it leaves. `Collection.Repair` restores all five.

## Recovery

`Collection.Repair` runs under the write lock and reads only what the layout names:
positions 1 to `size` and ids 1 to `last-id`.

| Step | Does |
|---|---|
| 1 | refills every dead slot from the end of the list and drops dead slots at the end, so the list is dense again |
| 2 | raises `last-id` to the highest id in the list, should it be below |
| 3 | for every id up to `last-id` that is not live and still owns a key: deletes its index entries that name it, its values, its nested collections and its back-pointer |
| 4 | for every live record, writes every missing index entry, then repairs every nested collection of it |

Step 4 is also how a field turned from `String` into `Key` gets its index. Two live records
holding the same value for a `Key` are reported as one `KeyConflict` once everything else
is repaired.

Two kinds of debris are out of reach, because nothing the layout names leads to them. Both
are harmless:

- An index entry under a hash that no live record holds. It is not valid, so it resolves to
  nothing and the next insert of that value claims it.
- Keys a nested insert left behind when both it and its own clean-up failed, under an owner
  that was removed before any `Repair` ran. Nothing reaches them once the owner is gone.

`last-id` is never rewound during recovery: the ids those records consumed stay spent.
