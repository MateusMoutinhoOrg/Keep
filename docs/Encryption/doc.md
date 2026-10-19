# Encryption

A database is encrypted by one field of its `Props`. Leave it empty and nothing changes:
every value is stored as given, byte for byte as before.

```go
var Props = api.Props{
	Path:       "test-dir/database/",
	EncryptKey: os.Getenv("KEEP_KEY"), // empty: no encryption
	Schemas:    []api.Schema{ /* … */ },
}
```

Reads and writes take and hand back plain values; only the backend sees ciphertext.
[encrypted-database](../../examples/lib/encrypted-database/example.go) runs it end to end.

## What is encrypted

| Stored | With an `EncryptKey` |
|---|---|
| every value: field values, `size`, `last-id`, `list/<n>`, `position`, index entries | ciphertext, `Deps.CryptoDeps.Encrypt(value, EncryptKey)` |
| an index entry's name, `keys/<field>/<hash>` | `sha256(EncryptKey ‖ 0x00 ‖ folded value)`, so it cannot be checked against a guessed value |
| a key's segments: `Path`, schema and field names, ids, positions | in the clear: they are the layout, see [DenseRecordPattern](../DenseRecordPattern/doc.md) |

The number of `list/<n>` keys still shows how many records a collection holds.

## How

| Piece | Where | Does |
|---|---|---|
| `Props.EncryptKey` | `sandbox/api/databases.go` | turns encryption on for one database |
| `encryption.Sandbox` | `sandbox/internal/encryption/` | builds the copy of the sandbox the database runs over: `StorageDeps` encrypting every write and decrypting every read, `HashDeps.Sha256Hex` keyed |
| `cryptodeps.Contract` | `sandbox/deps/cryptodeps/` | `Encrypt(content, key)` and `Decrypt(content, key)` |
| `aesgcmcrypto` | `adapters/impls/aesgcmcrypto/` | AES-256-GCM, AES key = `sha256(key)`, stored as `nonce(12) ‖ ciphertext ‖ tag(16)`: 28 bytes over each value |

Every key is unchanged and every write is still one write, so the dense, record and
collection packages run over the encrypted storage untouched, and every ordering of
[DenseRecordPattern](../DenseRecordPattern/doc.md) holds as it does in the clear.

## Another cipher

`CryptoDeps` is a dep like any other: replace its two fields before `sandbox.New`.

```go
deps := standard.New()
deps.CryptoDeps.Encrypt = func(content []byte, key string) ([]byte, error) { /* … */ }
deps.CryptoDeps.Decrypt = func(content []byte, key string) ([]byte, error) { /* … */ }
lib := sandbox.New(&deps)
```

`Decrypt` has to fail on content it did not seal under that key: that failure is what
reports a wrong key.

## Opening with the wrong key

| Opened with | Gets |
|---|---|
| the key it was written with | the plain values |
| another key | [`Internal`](../Errors/doc.md), `cannot decrypt`, on the first value read |
| no key | a lookup finds nothing; an insert, a listing or a repair fails with `Internal` before it writes |
| a key, over a database written without one | `Internal` on the first value read |

## Limits

| Limit | Because |
|---|---|
| no migration in place: an existing database stays in the clear | copy it record by record into a database with an `EncryptKey` |
| a weak key can be brute-forced offline | `aesgcmcrypto` derives with one SHA-256: use a long random secret, not a password |
| a value is not bound to its key | someone who can write to the backend can move one ciphertext under another key of the same database |
| `Database.Props` hands the key back | it is a copy of the `Props` the caller gave |
| roughly +30% CPU per operation over `memstorage`, ~1% over `filestorage` | one AES-GCM seal or open per value read or written |
