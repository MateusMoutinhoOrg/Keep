# `deps.CryptoDeps`

`sandbox/deps/cryptodeps`

## `Contract`

Contract is the cipher library injected whole as the Deps.CryptoDeps field.

| Field | Type | Description |
| --- | --- | --- |
| `Encrypt` | `func(content []byte, key string) ([]byte, error)` | Encrypt returns content sealed under key. It draws a fresh random nonce on every call, so the same content never encrypts to the same bytes twice, and the result carries whatever Decrypt needs besides the key. |
| `Decrypt` | `func(content []byte, key string) ([]byte, error)` | Decrypt returns the content Encrypt sealed under key. It fails when content was sealed under another key, was never sealed at all, or was altered since. |

[every contract](doc.md)
