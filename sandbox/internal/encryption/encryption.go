package encryption

import (
	api "github.com/MateusMoutinhoOrg/Keep/sandbox/api"
	storagedeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/storagedeps"
)

// The sandbox an encrypted database runs over. A database whose Props
// carries an EncryptKey is built over a copy of the sandbox whose deps differ
// in two fields, and in nothing else:
//
//   - StorageDeps encrypts every value it writes and decrypts every value it
//     reads with Deps.CryptoDeps, so the backend only ever holds ciphertext;
//   - HashDeps.Sha256Hex hashes the key ahead of its content, which keys the
//     digest dense.HashIndexValue names an index entry by — its one caller —
//     so an entry no longer tells whoever reads the backend which value it
//     indexes.
//
// Every key stays exactly what it was, and every write is still one write,
// so the dense, record and collection packages run over the copy unchanged
// and every invariant of the dense record pattern holds as it does in the
// clear. The copy is taken when the database is built: a field of the sandbox
// replaced after that is not seen by it.

// Sandbox returns the sandbox a database encrypted under key runs over. An
// empty key hands sandbox back unchanged: the database stores every value as
// given.
func Sandbox(sandbox *api.Sandbox, key string) *api.Sandbox {
	if key == "" {
		return sandbox
	}
	deps := *sandbox.Deps
	deps.StorageDeps = Storage(sandbox, key)
	deps.HashDeps.Sha256Hex = keyedDigest(sandbox, key)

	encrypted := *sandbox
	encrypted.Deps = &deps
	return &encrypted
}

// keyedDigest returns the digest an encrypted database names its index
// entries by: the SHA-256 of the key, a zero byte and the content, in that
// order.
func keyedDigest(sandbox *api.Sandbox, key string) func(content []byte) string {
	digest := sandbox.Deps.HashDeps.Sha256Hex
	return func(content []byte) string {
		keyed := make([]byte, 0, len(key)+1+len(content))
		keyed = append(keyed, key...)
		keyed = append(keyed, 0)
		return digest(append(keyed, content...))
	}
}

// Storage wraps the sandbox's storage so every value crosses it encrypted
// under key. Keys, existence, deletion and leases pass through untouched;
// the fields that work on part of a value decrypt the whole of it first.
func Storage(sandbox *api.Sandbox, key string) storagedeps.Contract {
	storage := sandbox.Deps.StorageDeps
	crypto := sandbox.Deps.CryptoDeps

	read := func(at []string) ([]byte, bool, error) {
		sealed, found, err := storage.Read(at)
		if err != nil || !found {
			return nil, found, err
		}
		value, err := crypto.Decrypt(sealed, key)
		if err != nil {
			return nil, false, decryptError(sandbox, at, err)
		}
		return value, true, nil
	}
	write := func(at []string, value []byte) error {
		sealed, err := crypto.Encrypt(value, key)
		if err != nil {
			return err
		}
		return storage.Write(at, sealed)
	}

	return storagedeps.Contract{
		Write: write,
		WriteIfAbsent: func(at []string, value []byte) (bool, error) {
			sealed, err := crypto.Encrypt(value, key)
			if err != nil {
				return false, err
			}
			return storage.WriteIfAbsent(at, sealed)
		},
		// The comparison is on the decrypted value, and the write is
		// conditional on the ciphertext it was decrypted from, so it stays as
		// atomic as the backend makes it.
		WriteIfValueEquals: func(at []string, value []byte, oldValue []byte) (bool, error) {
			current, found, err := storage.Read(at)
			if err != nil || !found {
				return false, err
			}
			plain, err := crypto.Decrypt(current, key)
			if err != nil {
				return false, decryptError(sandbox, at, err)
			}
			if string(plain) != string(oldValue) {
				return false, nil
			}
			sealed, err := crypto.Encrypt(value, key)
			if err != nil {
				return false, err
			}
			return storage.WriteIfValueEquals(at, sealed, current)
		},
		Append: func(at []string, value []byte) error {
			current, _, err := read(at)
			if err != nil {
				return err
			}
			return write(at, append(current, value...))
		},
		InsertAt: func(at []string, position int64, value []byte) error {
			current, _, err := read(at)
			if err != nil {
				return err
			}
			if position < 0 || position > int64(len(current)) {
				return sandbox.Deps.StdDeps.Errorf("keep: position %d is outside a value of %d bytes", position, len(current))
			}
			spliced := make([]byte, 0, len(current)+len(value))
			spliced = append(spliced, current[:position]...)
			spliced = append(spliced, value...)
			return write(at, append(spliced, current[position:]...))
		},
		Exists: storage.Exists,
		Read:   read,
		ReadAt: func(at []string, position int64, size int64) ([]byte, bool, error) {
			value, found, err := read(at)
			if err != nil || !found {
				return nil, found, err
			}
			if position < 0 || position > int64(len(value)) {
				return []byte{}, true, nil
			}
			end := int64(len(value))
			if size >= 0 && size < end-position {
				end = position + size
			}
			return value[position:end], true, nil
		},
		Delete: storage.Delete,
		Lock:   storage.Lock,
		Unlock: storage.Unlock,
	}
}

// decryptError is the failure of a value that does not decrypt under the
// database's key: it reaches the caller as Internal.
func decryptError(sandbox *api.Sandbox, at []string, err error) error {
	return sandbox.Deps.StdDeps.Errorf("keep: cannot decrypt %q: %v — the database was opened with another EncryptKey, or the value was written without one",
		sandbox.Deps.StringsDeps.Join(at, "/"), err)
}
