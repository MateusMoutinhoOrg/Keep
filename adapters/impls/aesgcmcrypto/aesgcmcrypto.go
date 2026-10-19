package aesgcmcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"sync"

	cryptodeps "github.com/MateusMoutinhoOrg/Keep/sandbox/deps/cryptodeps"

	"github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)

// ciphers caches the AEAD of every key seen, so a key is turned into an
// AES-256 cipher once rather than on every value.
var ciphers sync.Map

// aeadOf returns the AES-256-GCM AEAD of key, whose 32-byte AES key is the
// SHA-256 digest of key.
func aeadOf(key string) (cipher.AEAD, error) {
	if cached, ok := ciphers.Load(key); ok {
		return cached.(cipher.AEAD), nil
	}
	sum := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphers.Store(key, aead)
	return aead, nil
}

// encrypt fills cryptodeps.Contract.Encrypt: a random nonce followed by the
// GCM ciphertext and its tag.
func encrypt(content []byte, key string) ([]byte, error) {
	aead, err := aeadOf(key)
	if err != nil {
		return nil, err
	}
	sealed := make([]byte, aead.NonceSize(), aead.NonceSize()+len(content)+aead.Overhead())
	if _, err := rand.Read(sealed); err != nil {
		return nil, err
	}
	return aead.Seal(sealed, sealed, content, nil), nil
}

// decrypt fills cryptodeps.Contract.Decrypt, reading the nonce back off the
// front of content.
func decrypt(content []byte, key string) ([]byte, error) {
	aead, err := aeadOf(key)
	if err != nil {
		return nil, err
	}
	if len(content) < aead.NonceSize()+aead.Overhead() {
		return nil, errors.New("aesgcmcrypto: content is too short to have been encrypted")
	}
	nonce, sealed := content[:aead.NonceSize()], content[aead.NonceSize():]
	return aead.Open(nil, nonce, sealed, nil)
}

// Bind fills deps.Deps.CryptoDeps with AES-256-GCM over the standard
// library's crypto/aes, crypto/cipher and crypto/rand.
func Bind(deps *deps.Deps) {
	deps.CryptoDeps = cryptodeps.Contract{
		Encrypt: func(content []byte, key string) ([]byte, error) {
			return encrypt(content, key)
		},
		Decrypt: func(content []byte, key string) ([]byte, error) {
			return decrypt(content, key)
		},
	}
}
