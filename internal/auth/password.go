package auth

import (
	"github.com/alexedwards/argon2id"
)

// HashPassword hashes the plain password using argon2id.CreateHash.
func HashPassword(password string) (string, error) {
	// Use default parameters from the library
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

// CheckPasswordHash compares a plain password with a stored hash using argon2id.ComparePasswordAndHash.
func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
