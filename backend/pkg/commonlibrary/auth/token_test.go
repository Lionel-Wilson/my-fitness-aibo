package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHashing(t *testing.T) {
	hash, err := HashPassword("s3cret-pass")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !CheckPassword(hash, "s3cret-pass") {
		t.Error("CheckPassword should succeed for correct password")
	}
	if CheckPassword(hash, "wrong") {
		t.Error("CheckPassword should fail for wrong password")
	}
}

func TestTokenRoundTrip(t *testing.T) {
	m := NewTokenManager("test-secret", time.Hour)
	id := uuid.New()

	tok, err := m.Issue(id)
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}

	got, err := m.Verify(tok)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if got != id {
		t.Errorf("Verify returned %v, want %v", got, id)
	}
}

func TestVerifyRejectsExpired(t *testing.T) {
	m := NewTokenManager("test-secret", -time.Minute)
	tok, _ := m.Issue(uuid.New())
	if _, err := m.Verify(tok); err == nil {
		t.Error("Verify should reject an expired token")
	}
}

func TestVerifyRejectsWrongSecret(t *testing.T) {
	issuer := NewTokenManager("secret-a", time.Hour)
	verifier := NewTokenManager("secret-b", time.Hour)
	tok, _ := issuer.Issue(uuid.New())
	if _, err := verifier.Verify(tok); err == nil {
		t.Error("Verify should reject a token signed with a different secret")
	}
}
