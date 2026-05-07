package encryptor_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/encryptor"
)

func newEnc() *encryptor.Encryptor {
	return encryptor.New("test-passphrase")
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	enc := newEnc()
	plaintext := "super-secret-value"
	cipher, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if cipher == plaintext {
		t.Fatal("ciphertext should differ from plaintext")
	}
	got, err := enc.Decrypt(cipher)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if got != plaintext {
		t.Fatalf("want %q, got %q", plaintext, got)
	}
}

func TestEncrypt_NonDeterministic(t *testing.T) {
	enc := newEnc()
	c1, _ := enc.Encrypt("value")
	c2, _ := enc.Encrypt("value")
	if c1 == c2 {
		t.Fatal("expected different ciphertexts for same plaintext")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	enc := newEnc()
	_, err := enc.Decrypt("!!!not-base64!!!")
	if err != encryptor.ErrInvalidCiphertext {
		t.Fatalf("want ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	enc := newEnc()
	cipher, _ := enc.Encrypt("original")
	// flip last character to corrupt the MAC
	tampered := cipher[:len(cipher)-1] + strings.Map(func(r rune) rune {
		if r == 'A' {
			return 'B'
		}
		return 'A'
	}, string(cipher[len(cipher)-1]))
	_, err := enc.Decrypt(tampered)
	if err != encryptor.ErrInvalidCiphertext {
		t.Fatalf("want ErrInvalidCiphertext on tampered input, got %v", err)
	}
}

func TestEncryptMap_DecryptMap_Roundtrip(t *testing.T) {
	enc := newEnc()
	original := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
		"PORT":        "5432",
	}
	encrypted, err := enc.EncryptMap(original)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	for k, v := range original {
		if encrypted[k] == v {
			t.Errorf("key %s: expected encrypted value to differ", k)
		}
	}
	decrypted, err := enc.DecryptMap(encrypted)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	for k, want := range original {
		if got := decrypted[k]; got != want {
			t.Errorf("key %s: want %q, got %q", k, want, got)
		}
	}
}

func TestDecryptMap_InvalidEntry(t *testing.T) {
	enc := newEnc()
	bad := map[string]string{"KEY": "not-valid-cipher"}
	_, err := enc.DecryptMap(bad)
	if err == nil {
		t.Fatal("expected error decrypting invalid map entry")
	}
}
