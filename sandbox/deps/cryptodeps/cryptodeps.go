package cryptodeps

// This package is the sandbox's *copy* of the api a cipher library exposes —
// the same mechanic as hashdeps, stddeps and stringsdeps, for the same
// reason: the sandbox may import nothing but the sandbox, so `crypto/aes`,
// `crypto/cipher` and `crypto/rand` may not appear inside it. The contract is
// restated here, and the adapter — which lives outside the sandbox — is what
// fills it.
//
// Keep needs it for one thing: a database whose Props carries an EncryptKey
// passes every value it writes through Encrypt and every value it reads
// through Decrypt, so the backend only ever holds ciphertext. Replacing the
// two fields before sandbox.New swaps the cipher for every such database.

// Contract is the cipher library injected whole as the Deps.CryptoDeps field.
type Contract struct {
	// Encrypt returns content sealed under key. It draws a fresh random
	// nonce on every call, so the same content never encrypts to the same
	// bytes twice, and the result carries whatever Decrypt needs besides the
	// key.
	Encrypt func(content []byte, key string) ([]byte, error)

	// Decrypt returns the content Encrypt sealed under key. It fails when
	// content was sealed under another key, was never sealed at all, or was
	// altered since.
	Decrypt func(content []byte, key string) ([]byte, error)
}
