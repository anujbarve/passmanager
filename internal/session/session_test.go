package session

import (
	"bytes"
	"testing"
	"time"

	"passmanager/internal/crypto"
)

func TestSessionLifecycle(t *testing.T) {
	s := &Session{timeout: 5 * time.Minute}
	if s.IsAuthenticated() {
		t.Fatal("expected new session to be unauthenticated")
	}

	salt := []byte("0123456789abcdef")
	cryptoSvc := crypto.NewCryptoService("pw", salt)
	s.Login(nil, cryptoSvc, salt)

	if !s.IsAuthenticated() {
		t.Fatal("expected session to be authenticated after Login")
	}
	if !bytes.Equal(s.GetSalt(), salt) {
		t.Fatalf("got salt %v, want %v", s.GetSalt(), salt)
	}
	if s.GetTimeRemaining() <= 0 {
		t.Fatal("expected positive remaining session time")
	}
	if s.GetCrypto() == nil {
		t.Fatal("expected non-nil CryptoService while logged in")
	}

	s.Logout()
	if s.IsAuthenticated() {
		t.Fatal("expected session to be unauthenticated after Logout")
	}
	if s.GetCrypto() != nil {
		t.Fatal("expected nil CryptoService after Logout")
	}
	if s.GetSalt() != nil {
		t.Fatal("expected nil salt after Logout")
	}
}

func TestSessionTimeout(t *testing.T) {
	s := &Session{timeout: 20 * time.Millisecond}
	salt := []byte("0123456789abcdef")
	s.Login(nil, crypto.NewCryptoService("pw", salt), salt)

	time.Sleep(35 * time.Millisecond)
	if s.IsAuthenticated() {
		t.Fatal("expected session to expire after timeout")
	}
}
