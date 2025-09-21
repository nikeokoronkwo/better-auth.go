package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

// Hash the given password using scrypt
//
// Based on memory we use scrypt instead of argon2id
func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}


	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", base64.RawStdEncoding.EncodeToString(salt), 
	base64.RawStdEncoding.EncodeToString(hash)), nil
}

// Verify the given password against the given hash using scrypt
func VerifyPassword(password, hash string) (bool, error) {
	var saltB64, hashB64 string
	n, _ := fmt.Sscanf(hash, "%[^:]:%s", &saltB64, &hashB64)
	if n != 2 {
		return false, fmt.Errorf("invalid encoded hash format")
	}

	salt, _ := base64.RawStdEncoding.DecodeString(saltB64)
	expectedHash, _ := base64.RawStdEncoding.DecodeString(hashB64)

	// Same parameters as hashing
	reHash, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(reHash, expectedHash) == 1, nil
}