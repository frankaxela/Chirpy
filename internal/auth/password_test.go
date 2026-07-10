package auth

import (
	"strings"
	"testing"
)

func TestHashPasswordReturnsArgon2idHash(t *testing.T) {
	password := "correct-horse-battery-staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("expected hash to start with $argon2id$, got %q", hash)
	}
	if hash == password {
		t.Error("hash must not equal the plain password")
	}
}

func TestHashPasswordProducesUniqueHashes(t *testing.T) {
	password := "same-password"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("first HashPassword returned error: %v", err)
	}
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("second HashPassword returned error: %v", err)
	}
	if hash1 == hash2 {
		t.Error("hashing the same password twice should produce different hashes (random salt)")
	}
}

func TestCheckPasswordHashCorrectPassword(t *testing.T) {
	password := "s3cret!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if !match {
		t.Error("expected correct password to match its hash")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	hash, err := HashPassword("right-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	match, err := CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned error: %v", err)
	}
	if match {
		t.Error("expected wrong password not to match the hash")
	}
}
