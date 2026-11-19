package models

import "testing"

func TestDefaultSettings(t *testing.T) {
	s := DefaultSettings()
	if s.SessionTimeout != 5 {
		t.Fatalf("got SessionTimeout %d, want 5", s.SessionTimeout)
	}
	if s.ClipboardTimeout != 30 {
		t.Fatalf("got ClipboardTimeout %d, want 30", s.ClipboardTimeout)
	}
	if s.DefaultCategory != "general" {
		t.Fatalf("got DefaultCategory %q, want %q", s.DefaultCategory, "general")
	}
	if s.PasswordLength != 20 {
		t.Fatalf("got PasswordLength %d, want 20", s.PasswordLength)
	}
	if !s.IncludeSymbols {
		t.Fatal("expected IncludeSymbols to be true by default")
	}
}
