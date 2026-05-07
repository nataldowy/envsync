// Package encryptor provides symmetric encryption helpers for protecting
// sensitive .env values at rest.
//
// Values are encrypted with AES-256-GCM, a standard authenticated encryption
// scheme that guarantees both confidentiality and integrity. A random 96-bit
// nonce is generated for every Encrypt call, so repeated encryption of the
// same plaintext yields different ciphertexts.
//
// Usage:
//
//	enc := encryptor.New("my-secret-passphrase")
//
//	cipher, err := enc.Encrypt("s3cr3t")
//	plain,  err := enc.Decrypt(cipher)
//
//	encryptedEnv, err := enc.EncryptMap(envMap)
//	decryptedEnv, err := enc.DecryptMap(encryptedEnv)
package encryptor
