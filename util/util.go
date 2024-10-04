package util

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gorilla/sessions"
)

func GenerateRandomKey(length int) (string, error) {
	// Create a slice to hold the random bytes
	key := make([]byte, length)

	// Fill the slice with random bytes
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to a hexadecimal string
	return hex.EncodeToString(key), nil
}

func InitializeStore() *sessions.CookieStore {
	// Generate a 32-byte random key
	key, err := GenerateRandomKey(32)
	if err != nil {
		panic(err)
	}

	// Use the key for session encryption
	store := sessions.NewCookieStore([]byte(key))
	return store
}

func ValidatePassword(password string) bool {
	passLen := len(password)
	var containsDigits bool
	var containsCaps bool

	if passLen < 6 {
		return false
	}

	for i := 0; i < passLen; i++ {
		if password[i] >= '0' && password[i] <= '9' {
			containsDigits = true
		}

		if password[i] >= 'A' && password[i] <= 'Z' {
			containsCaps = true
		}
	}

	if !containsDigits || !containsCaps {
		return false
	}

	return true
}

func ValidateName(name string) bool {
	nameLen := len(name)

	if nameLen < 2 {
		return false
	}

	for i := 0; i < nameLen; i++ {
		if !(name[i] >= 'a' && name[i] <= 'z') && !(name[i] >= 'A' && name[i] <= 'Z') {
			return false
		}

	}

	return true
}
