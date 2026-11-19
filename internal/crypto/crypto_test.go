package crypto

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	s1, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() error = %v", err)
	}
	if len(s1) != saltLen {
		t.Fatalf("expected salt length %d, got %d", saltLen, len(s1))
	}

	s2, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt() second call error = %v", err)
	}
	if bytes.Equal(s1, s2) {
		t.Fatal("expected two generated salts to be different")
	}
}

func TestDeriveKey(t *testing.T) {
	salt := []byte("0123456789abcdef")
	k1 := DeriveKey("master-pass", salt)
	k2 := DeriveKey("master-pass", salt)
	k3 := DeriveKey("other-pass", salt)

	if len(k1) != argonKeyLen {
		t.Fatalf("expected key length %d, got %d", argonKeyLen, len(k1))
	}
	if !bytes.Equal(k1, k2) {
		t.Fatal("expected DeriveKey to be deterministic for same inputs")
	}
	if bytes.Equal(k1, k3) {
		t.Fatal("expected DeriveKey to differ for different passwords")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	salt := []byte("0123456789abcdef")
	svc := NewCryptoService("super-secret", salt)

	cases := []string{
		"",
		"hello world",
		"p@ssw0rd!#%&*()_+",
		"unicode-🔐-пароль",
	}

	for _, tc := range cases {
		enc, err := svc.Encrypt(tc)
		if err != nil {
			t.Fatalf("Encrypt(%q) error = %v", tc, err)
		}
		dec, err := svc.Decrypt(enc)
		if err != nil {
			t.Fatalf("Decrypt(%q) error = %v", tc, err)
		}
		if dec != tc {
			t.Fatalf("round-trip mismatch: got %q, want %q", dec, tc)
		}
	}

	enc1, err := svc.Encrypt("same-plaintext")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	enc2, err := svc.Encrypt("same-plaintext")
	if err != nil {
		t.Fatalf("Encrypt() second call error = %v", err)
	}
	if enc1 == enc2 {
		t.Fatal("expected distinct ciphertexts due to random nonce")
	}
}

func TestDecryptErrors(t *testing.T) {
	salt := []byte("0123456789abcdef")
	svc := NewCryptoService("correct-password", salt)
	wrongSvc := NewCryptoService("wrong-password", salt)

	enc, err := svc.Encrypt("secret")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if _, err := wrongSvc.Decrypt(enc); err == nil {
		t.Fatal("expected error decrypting with wrong key")
	}
	if _, err := svc.Decrypt("not-valid-base64!!!"); err == nil {
		t.Fatal("expected error for invalid base64")
	}
	if _, err := svc.Decrypt("AAAA"); err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}

func TestHashMasterPassword(t *testing.T) {
	salt := []byte("0123456789abcdef")
	h1 := HashMasterPassword("pw1", salt)
	h2 := HashMasterPassword("pw1", salt)
	h3 := HashMasterPassword("pw2", salt)

	if h1 == "" {
		t.Fatal("expected non-empty password hash")
	}
	if h1 != h2 {
		t.Fatal("expected HashMasterPassword to be deterministic")
	}
	if h1 == h3 {
		t.Fatal("expected different passwords to produce different hashes")
	}
}

func TestGeneratePassword(t *testing.T) {
	const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const symbols = "!@#$%^&*()_+-=[]{}|;:,.<>?"

	pwNoSym, err := GeneratePassword(24, false)
	if err != nil {
		t.Fatalf("GeneratePassword(24, false) error = %v", err)
	}
	if len(pwNoSym) != 24 {
		t.Fatalf("expected length 24, got %d", len(pwNoSym))
	}
	for _, r := range pwNoSym {
		if !strings.ContainsRune(alnum, r) {
			t.Fatalf("unexpected non-alphanumeric rune %q in %q", r, pwNoSym)
		}
	}

	pwWithSym, err := GeneratePassword(32, true)
	if err != nil {
		t.Fatalf("GeneratePassword(32, true) error = %v", err)
	}
	if len(pwWithSym) != 32 {
		t.Fatalf("expected length 32, got %d", len(pwWithSym))
	}
	allowed := alnum + symbols
	for _, r := range pwWithSym {
		if !strings.ContainsRune(allowed, r) {
			t.Fatalf("unexpected rune %q in %q", r, pwWithSym)
		}
	}
}

func TestSecureClear(t *testing.T) {
	svc := NewCryptoService("secret", []byte("0123456789abcdef"))
	svc.SecureClear()
	for i, b := range svc.masterKey {
		if b != 0 {
			t.Fatalf("expected masterKey[%d] to be zeroed, got %d", i, b)
		}
	}
}
