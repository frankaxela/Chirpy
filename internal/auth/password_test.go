package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func TestMakeJWTAndValidateJWTRoundTrip(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}
	if token == "" {
		t.Fatal("MakeJWT returned an empty token")
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}
	if gotID != userID {
		t.Errorf("expected user ID %s, got %s", userID, gotID)
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-time.Hour)),
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-time.Minute)),
	})
	token, err := expiredToken.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("signing expired token returned error: %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected an error validating an expired token, got nil")
	}
	if gotID != uuid.Nil {
		t.Errorf("expected uuid.Nil for expired token, got %s", gotID)
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	userID := uuid.New()

	token, err := MakeJWT(userID, "correct-secret")
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	gotID, err := ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected an error validating a token with the wrong secret, got nil")
	}
	if gotID != uuid.Nil {
		t.Errorf("expected uuid.Nil for wrong-secret token, got %s", gotID)
	}
}

func TestValidateJWTRejectsMalformedToken(t *testing.T) {
	gotID, err := ValidateJWT("not-a-jwt", "test-secret")
	if err == nil {
		t.Fatal("expected an error validating a malformed token, got nil")
	}
	if gotID != uuid.Nil {
		t.Errorf("expected uuid.Nil for malformed token, got %s", gotID)
	}
}

func TestGetBearerTokenValidHeader(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer my-token-123")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}
	if token != "my-token-123" {
		t.Errorf("expected token %q, got %q", "my-token-123", token)
	}
}

func TestGetBearerTokenPreservesTokenWithSpaces(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer part1 part2")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}
	if token != "part1 part2" {
		t.Errorf("expected token %q, got %q", "part1 part2", token)
	}
}

func TestGetBearerTokenEmptyTokenAfterBearer(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}
	if token != "" {
		t.Errorf("expected empty token, got %q", token)
	}
}

func TestGetBearerTokenMissingHeader(t *testing.T) {
	headers := http.Header{}

	token, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error for a missing Authorization header, got nil")
	}
	if token != "" {
		t.Errorf("expected empty token for missing header, got %q", token)
	}
}

func TestGetBearerTokenEmptyHeaderValue(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "")

	token, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error for an empty Authorization header, got nil")
	}
	if token != "" {
		t.Errorf("expected empty token for empty header, got %q", token)
	}
}

func TestGetBearerTokenWrongScheme(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Basic dXNlcjpwYXNz")

	token, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error for a non-Bearer Authorization header, got nil")
	}
	if token != "" {
		t.Errorf("expected empty token for wrong scheme, got %q", token)
	}
}

func TestGetBearerTokenMissingSpaceAfterBearer(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearermy-token-123")

	token, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error when Bearer is not followed by a space, got nil")
	}
	if token != "" {
		t.Errorf("expected empty token, got %q", token)
	}
}

func TestGetBearerTokenLowercaseSchemeRejected(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "bearer my-token-123")

	token, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("expected an error for lowercase bearer scheme, got nil")
	}
	if token != "" {
		t.Errorf("expected empty token for lowercase scheme, got %q", token)
	}
}
