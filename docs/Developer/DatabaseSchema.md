# Dense Record Pattern for Key-Value Stores

> This is the internal design document behind Keep's storage layer (`pkg/database`). You don't need it to *use* the library — start at [Library Initialization](../Consumer/UseCases/LibInitialization.md). Read this if you want to understand how the data is laid out or verify an implementation.

A storage pattern for managing collections of records on top of a plain key-value store, using only single-item reads and writes. The pattern never relies on key listing, prefix scans, or range queries, which makes it portable to any KV backend and keeps every operation bounded to a constant number of key accesses.

## Design Goals

1. **No listing.** Every operation is expressed exclusively as reads, writes, and deletes of individual keys. The pattern must work on backends that offer nothing beyond `get`, `set`, and `delete`.
2. **Constant-time operations.** Insertion, deletion, and lookup by unique field each touch a fixed, small number of keys, regardless of how many records the collection holds.
3. **Stable identity.** A record's id never changes and is never reused, even after deletion.

## Key Layout

Each collection (e.g. `users`) is represented by four families of keys.

**Metadata keys** hold the collection's global state:

- `{collection}-size` — the current number of live records.
- `{collection}-last-id` — the highest id ever assigned. It only grows; it is never decremented, and ids are never recycled. This guarantees that an id observed by an external system always refers to the same logical record.

**List keys** form a dense positional array mapping positions to ids:

- `{collection}-list-{position}` — the id of the record currently occupying that position.

Positions run contiguously from `1` to `size`, with no holes. This density is the property that makes iteration possible without listing: a consumer reads position 1, 2, 3, and so on, up to `size`. Note that the list is a *dense set*, not an ordered sequence — deletion may reorder it (see the deletion procedure). If insertion order matters to a consumer, it must be stored as a field on the record itself.

**Unique index keys** map a hashed field value to a record id:

- `{collection}-keys-{field}-{sha(value)}` — the id of the record holding that value.

Hashing the value serves two purposes: it bounds the key length regardless of the value's size, and it neutralizes any characters in the value that could collide with the key layout's separator.

**Record keys** hold the record's own data:

- `{collection}-{id}-position` — the record's current position in the list. This is the back-pointer that makes constant-time deletion possible.
- `{collection}-{id}-values-{field}` — one key per field of the record.

## Normalization Rules

Consistency between the index keys and the stored values is mandatory, otherwise lookups silently fail.

- Field names must be spelled identically everywhere they appear (index keys and value keys). Pick one canonical casing per field and never deviate.
- Values of indexed fields must be normalized before hashing, and the same normalization must be applied at both write time and lookup time. For case-insensitive fields such as email addresses, lowercase the value before hashing.
- Collection names and field names must not contain the separator character used in the key layout. Either forbid the separator in names outright, or choose a separator that cannot occur in them.

## Insertion

Given a new record with its field values:

1. **Check uniqueness.** For each unique field, normalize the incoming value, compute its hash, and read `{collection}-keys-{field}-{hash}`. If any of these keys already exists, abort the insertion with a uniqueness violation. Nothing has been written yet, so no cleanup is needed.
2. **Allocate the id.** Read `{collection}-last-id`, increment it, and write it back. The incremented value is the new record's id.
3. **Determine the position.** Read `{collection}-size`. Its current value plus one is the new record's position, since positions are one-based and dense.
4. **Write the record's data.** Write one `{collection}-{id}-values-{field}` key per field, and write `{collection}-{id}-position` with the position determined in step 3.
5. **Write the index entries.** For each unique field, write `{collection}-keys-{field}-{hash(value)}` pointing to the new id.
6. **Publish the record.** Write `{collection}-list-{position}` with the new id, then write `{collection}-size` incremented by one.

The size key is written **last**, deliberately. Until size is incremented, the new position lies outside the valid range `[1, size]`, so a reader iterating the list never observes a half-written record. If the writer crashes mid-insertion, everything written so far is unreachable garbage rather than corrupt data, and can be reconciled later (see Recovery).

### Insertion Example

Inserting the following record into `users`:

```
email:    user2@gmail.com
username: User2
password: 12345
```

**Database before insertion:**

```
users-size: 1
users-last-id: 1
users-list-1: 1

users-keys-email-sha(user1@gmail.com): 1
users-keys-username-sha(user1): 1

users-1-position: 1
users-1-values-username: User1
users-1-values-password: 12345
users-1-values-email: user1@gmail.com
```

Walking through the steps:

- *Step 1:* read `users-keys-email-sha(user2@gmail.com)` and `users-keys-username-sha(user2)` — note the value is lowercased before hashing. Both are missing, so the record is unique.
- *Step 2:* read `users-last-id` (1), write it back as 2. The new id is **2**.
- *Step 3:* read `users-size` (1). The new position is **2**.
- *Steps 4–6:* write the value keys, position, index entries, then the list slot, and finally size.

**Database after insertion:**

```
users-size: 2
users-last-id: 2
users-list-1: 1
users-list-2: 2

users-keys-email-sha(user1@gmail.com): 1
users-keys-email-sha(user2@gmail.com): 2
users-keys-username-sha(user1): 1
users-keys-username-sha(user2): 2

users-1-position: 1
users-1-values-username: User1
users-1-values-password: 12345
users-1-values-email: user1@gmail.com

users-2-position: 2
users-2-values-username: User2
users-2-values-password: 12345
users-2-values-email: user2@gmail.com
```

The stored value keeps its original casing (`User2`); only the index entry uses the normalized form (`sha(user2)`).

## Deletion (Swap-With-Last)

Naively removing a record from the middle of the list would leave a hole, and closing that hole by shifting would cost one write per remaining record. Instead, the pattern always fills the hole with the **last** record of the list, keeping deletion at a constant cost. Given the id to delete:

1. **Read the victim's position.** Read `{collection}-{id}-position`; call it `p`. If the key does not exist, the record is already gone and the operation is a no-op.
2. **Locate the last record.** Read `{collection}-size`; the last occupied position is `size`. Read `{collection}-list-{size}` to obtain the id of the record living there; call it `lastId`.
3. **Move the last record into the hole** (skip this step if the victim *is* the last record, i.e. `p == size`). Write `{collection}-list-{p}` with `lastId`, and write `{collection}-{lastId}-position` with `p`.
4. **Shrink the list.** Delete `{collection}-list-{size}` and write `{collection}-size` decremented by one.
5. **Remove the index entries.** For each unique field, read the victim's stored value from `{collection}-{id}-values-{field}`, normalize and hash it, and delete `{collection}-keys-{field}-{hash}`.
6. **Remove the record's data.** Delete every `{collection}-{id}-values-{field}` key and `{collection}-{id}-position`.

Note that `{collection}-last-id` is untouched: ids are permanent and the deleted id will never be assigned again.

### Deletion Example

To make the swap visible, this example uses **three** records and deletes the one in the middle — deleting the first or last record is just a simpler case of the same procedure.

Deleting `id: 2` from the following state:

**Database before deletion:**

```
users-size: 3
users-last-id: 3
users-list-1: 1
users-list-2: 2
users-list-3: 3

users-keys-email-sha(user1@gmail.com): 1
users-keys-email-sha(user2@gmail.com): 2
users-keys-email-sha(user3@gmail.com): 3
users-keys-username-sha(user1): 1
users-keys-username-sha(user2): 2
users-keys-username-sha(user3): 3

users-1-position: 1
users-1-values-username: User1
users-1-values-password: 12345
users-1-values-email: user1@gmail.com

users-2-position: 2
users-2-values-username: User2
users-2-values-password: 12345
users-2-values-email: user2@gmail.com

users-3-position: 3
users-3-values-username: User3
users-3-values-password: 12345
users-3-values-email: user3@gmail.com
```

Walking through the steps:

- *Step 1:* read `users-2-position` → `p = 2`.
- *Step 2:* read `users-size` (3), so the last position is 3. Read `users-list-3` → `lastId = 3`.
- *Step 3:* the victim is not the last record, so record 3 moves into the hole: write `users-list-2 = 3` and `users-3-position = 2`.
- *Step 4:* delete `users-list-3`, write `users-size = 2`.
- *Step 5:* read the victim's values, hash them, delete `users-keys-email-sha(user2@gmail.com)` and `users-keys-username-sha(user2)`.
- *Step 6:* delete all `users-2-*` keys.

**Database after deletion:**

```
users-size: 2
users-last-id: 3
users-list-1: 1
users-list-2: 3

users-keys-email-sha(user1@gmail.com): 1
users-keys-email-sha(user3@gmail.com): 3
users-keys-username-sha(user1): 1
users-keys-username-sha(user3): 3

users-1-position: 1
users-1-values-username: User1
users-1-values-password: 12345
users-1-values-email: user1@gmail.com

users-3-position: 2
users-3-values-username: User3
users-3-values-password: 12345
users-3-values-email: user3@gmail.com
```

Two things worth noticing: `users-last-id` remains 3 even though id 2 is gone, and record 3 now occupies position 2 — the list order changed, which is the documented price of constant-time deletion. Whether the collection has 3 records or 3 million, this deletion touched the same number of keys.

The consequence of swap-with-last is that list order is **not stable** — deleting a record moves an unrelated record to a new position. This is the price of constant-time deletion and must be documented as a contract of the pattern.

## Updating an Indexed Field

Updates to non-indexed fields are a single key write. Updates to a **unique indexed field** are the subtle case, because the index must move without leaving orphans:

1. Read the current value from `{collection}-{id}-values-{field}` — the old value is needed to locate the old index entry.
2. Normalize and hash the new value, and read `{collection}-keys-{field}-{newHash}`. If it exists and points to a different id, abort with a uniqueness violation.
3. Write the new index entry `{collection}-keys-{field}-{newHash}` pointing to the id.
4. Write the new value to `{collection}-{id}-values-{field}`.
5. Delete the old index entry `{collection}-keys-{field}-{oldHash}`.

Writing the new entry before deleting the old one means a crash mid-update leaves the record reachable through at least one of the two values, never through neither.

### Update Example

Changing the email of `id: 1` from `user1@gmail.com` to `newmail@gmail.com`.

**Database before update (relevant keys only):**

```
users-keys-email-sha(user1@gmail.com): 1

users-1-values-email: user1@gmail.com
```

Walking through the steps:

- *Step 1:* read `users-1-values-email` → `user1@gmail.com`. This is what allows the old index entry to be located; without reading it first, the entry `users-keys-email-sha(user1@gmail.com)` would be left orphaned forever.
- *Step 2:* read `users-keys-email-sha(newmail@gmail.com)` — missing, so the new value is free.
- *Step 3:* write `users-keys-email-sha(newmail@gmail.com) = 1`.
- *Step 4:* write `users-1-values-email = newmail@gmail.com`.
- *Step 5:* delete `users-keys-email-sha(user1@gmail.com)`.

**Database after update (relevant keys only):**

```
users-keys-email-sha(newmail@gmail.com): 1

users-1-values-email: newmail@gmail.com
```

## Lookup

- **By id:** read the `{collection}-{id}-values-{field}` keys directly.
- **By unique field:** normalize the value, hash it, read `{collection}-keys-{field}-{hash}` to obtain the id, then proceed as a lookup by id.
- **Full iteration:** read `{collection}-size`, then read `{collection}-list-1` through `{collection}-list-{size}`, resolving each id to its record.

### Lookup Example

Finding the user with email `User3@Gmail.com` (note the mixed casing as typed by a caller):

- Normalize: `user3@gmail.com`. Hash it.
- Read `users-keys-email-sha(user3@gmail.com)` → `3`.
- Read `users-3-values-username`, `users-3-values-email`, and any other fields needed.

Three reads total, regardless of collection size — and the normalization step is what makes the lookup succeed even though the caller typed the email with different casing than it was stored with.

## Concurrency and Atomicity

The pattern's write sequences are safe against crashes (via the write ordering described above) but not, by themselves, against concurrent writers: two simultaneous insertions could both pass the uniqueness check in step 1, or both read the same `last-id`.

- If the backend offers transactions or atomic batches, wrap each operation in one. The write orderings above then matter only as documentation of intent.
- If it does not, the pattern assumes a **single writer**. Multiple readers are always safe, subject to the visibility guarantee provided by writing `size` last on insertion.

This assumption must be stated explicitly by any implementation of the pattern.

## Recovery

Because `size` acts as the commit point for insertion, recovery after a crash is mechanical: any record whose stored position is greater than or equal to the current `size`, or whose position in the list does not point back to it, is an incomplete insertion and can be safely garbage-collected by deleting its value keys and any index entries pointing to its id. Orphaned index entries — index keys whose target id no longer has a position key — are the corresponding artifact of an interrupted deletion or update, and are likewise safe to delete.

### Recovery Example

Suppose a crash happened during an insertion, after step 5 but before step 6. The database looks like this:

```
users-size: 1
users-last-id: 2
users-list-1: 1

users-keys-email-sha(user2@gmail.com): 2
users-keys-username-sha(user2): 2

users-2-position: 2
users-2-values-username: User2
users-2-values-email: user2@gmail.com
```

Record 2 claims position 2, but `users-size` is 1, so the valid range is only `[1, 1]` — position 2 was never published. A reader iterating the list never sees record 2, and a recovery pass identifies it as garbage (its position > size) and deletes its value keys and index entries. Note that `users-last-id` stays at 2: even a failed insertion consumes its id, preserving the no-reuse guarantee.

## Invariants

An implementation is correct if, at every quiescent point, all of the following hold:

1. `{collection}-list-{p}` exists exactly for `p` in `[1, size]`.
2. For every live record, `list[position(id)] == id` — the list and the back-pointers agree.
3. Every unique index entry points to a live record whose current normalized field value hashes back to that entry.
4. `last-id` is greater than or equal to every id that has ever existed in the collection.