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
