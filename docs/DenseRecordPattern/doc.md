# DenseRecordPattern

The key layout that lets a schema database run over a backend with nothing but `get`, `set`
and `delete`. Not needed to use Keep — [Schemas](../Schemas/doc.md) is — but it is what
`sandbox/internal/dense` and `sandbox/internal/schemaitem` implement, and what a change to
either has to preserve.

Three properties, and everything below follows from them:

1. **No listing.** Every operation is reads, writes and deletes of individual keys. Nothing
   scans a prefix, lists keys or queries a range.
2. **Fixed cost.** Insert, remove and look-up-by-key each touch a number of keys that does
   not grow with the size of the collection.
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
| `{c}/size` | the number of live records — the highest occupied position | metadata |
| `{c}/last-id` | the highest id ever allocated. Only grows | metadata |
| `{c}/list/{position}` | the id living at that position | position list |
| `{c}/keys/{field}/{sha256(lower(value))}` | the id holding that value for that `Key` field | unique index |
| `{c}/{id}/position` | the record's position. Its presence is what marks the record live | record |
| `{c}/{id}/values/{field}` | one field value of one record | record |
| `{c}/{id}/{field}` | the prefix of a nested collection owned by the record | record |

Positions run from `1` to `size` with **no gap**. That density is what makes iteration
possible without listing: read `size`, then read positions 1 through `size`. It is a dense
set, not an ordered sequence — a removal reorders it. A consumer that needs insertion order
stores it as a field.

Index values are lower-cased before they are hashed, which is what makes a key lookup
case-insensitive, and hashed at all so that key length stays bounded whatever a value
holds.

## Insert

| Step | Does | Fails with |
|---|---|---|
| 0 | validate the fields against the schema | `InvalidField`, `MissingField` |
| 1 | read `{c}/keys/{field}/{hash}` for every `Key` field | `KeyConflict` |
| 2 | read `{c}/last-id`, write it back +1. That is the new id | |
| 3 | read `{c}/size`. The new position is size+1 | |
| 4 | write every `{c}/{id}/values/{field}`, then `{c}/{id}/position` | |
| 5 | write every `{c}/keys/{field}/{hash}` pointing at the new id | |
| 6 | write `{c}/list/{position}`, then `{c}/size` | |

Steps 0 and 1 decide every failure a caller can see, and they write nothing — a refused
insert leaves the database exactly as it was.

`{c}/size` is written **last, deliberately**. Until it grows, the new position lies outside
`[1, size]` and no reader iterating the list can observe the record. A crash before that
leaves unreachable keys, never a half-visible record. The id it consumed is not given back:
a failed insert still spends its id, which is what keeps ids unique.

## Remove — swap with last

Closing the hole by shifting would cost one write per remaining record. The last record
fills it instead.

| Step | Does |
|---|---|
| 1 | read `{c}/{id}/position` → `p`. Absent means already gone: a no-op |
| 2 | read `{c}/size`, then `{c}/list/{size}` → `lastId` |
| 3 | when `p != size`: write `{c}/list/{p}` = `lastId`, `{c}/{lastId}/position` = `p` |
| 4 | delete `{c}/list/{size}`, write `{c}/size` = size-1 |
| 5 | for every `Key` field: read the value, delete `{c}/keys/{field}/{hash}` |
| 6 | delete every value key, clear every nested collection, delete `{c}/{id}/position` |

`{c}/last-id` is untouched. **List order is not stable across removals** — an unrelated
record moves — and that is the documented price of step 3.

Step 6 clears a nested collection by removing its records from the last position backwards,
so no swap is ever needed, then deletes its `size` and `last-id` keys. It recurses to any
depth.

## Update

A field that carries no index is one write.

A `Key` field moves its index entry:

| Step | Does |
|---|---|
| 1 | read `{c}/{id}/values/{field}` — the old value locates the old index entry |
| 2 | read `{c}/keys/{field}/{newHash}`. Another id there is a `KeyConflict`, and nothing is written |
| 3 | write `{c}/keys/{field}/{newHash}` = id |
| 4 | write `{c}/{id}/values/{field}` = the new value |
| 5 | delete `{c}/keys/{field}/{oldHash}`, when it differs from the new one |

New entry before old delete: a crash part-way leaves the record reachable through one of
the two values, never through neither. Step 1 is what stops the old entry being orphaned
forever.

## Read

| Operation | Costs |
|---|---|
| by id | one read for the position back-pointer, then one per field read |
| by unique key | one read for the index entry, one for the back-pointer, then one per field |
| one page | one read for `size`, then one per position in the page |
| whole collection | one read for `size`, then one per record |

A field is read on demand — nothing loads a whole record.

## Concurrency

The orderings above are safe against a crash, not against two writers. Two concurrent
inserts can both pass step 1, or both read the same `last-id`.

**Keep assumes a single writer.** Readers are always safe, given the visibility `size`
guarantees. A backend with transactions or atomic batches should wrap each operation in
one, and the orderings then only document intent. `Deps.Storagedeps` carries `Lock` and
`UnLock` for a backend that offers advisory leases; Keep never calls them itself — see
[StorageContract](../StorageContract/doc.md).

## Invariants

At every quiescent point:

1. `{c}/list/{p}` exists exactly for `p` in `[1, size]`.
2. For every live record, `list[position(id)] == id`.
3. Every index entry points at a live record whose current value hashes back to that entry.
4. `last-id` ≥ every id the collection has ever held.

## Recovery

`size` is the commit point, so a crash leaves only two kinds of debris, both safe to delete:

| Debris | Recognised by | Left by |
|---|---|---|
| an unpublished record | `position(id) > size`, or `list[position(id)] != id` | an insert that stopped before step 6 |
| an orphaned index entry | its target id has no `{c}/{id}/position` key | a removal or an update that stopped part-way |

`last-id` is never rewound during recovery: the ids those records consumed stay spent.
