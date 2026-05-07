// Package encryptor provides AES-GCM encryption and decryption for .env values.
package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidCiphertext is returned when decryption fails due to malformed input.
var ErrInvalidCiphertext = errors.New("encryptor: invalid ciphertext")

// Encryptor encrypts and decrypts string values using AES-256-GCM.
type Encryptor struct {
	key []byte
}

// New creates an Encryptor from a passphrase. The passphrase is hashed with
// SHA-256 to produce a 32-byte AES key.
func New(passphrase string) *Encryptor {
	hash := sha256.Sum256([]byte(passphrase))
	return &Encryptor{key: hash[:]}
}

// Encrypt encrypts plaintext and returns a base64-encoded ciphertext string.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded ciphertext string produced by Encrypt.
func (e *Encryptor) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", ErrInvalidCiphertext
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plaintext), nil
}

// EncryptMap encrypts every value in the provided map, returning a new map.
func (e *Encryptor) EncryptMap(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		enc, err := e.Encrypt(v)
		if err != nil {
			return nil, err
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts every value in the provided map, returning a new map.
func (e *Encryptor) DecryptMap(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		dec, err := e.Decrypt(v)
		if err != nil {
			return nil, err
		}
		out[k] = dec
	}
	return out, nil
}
