package util

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"

	"github.com/gorilla/sessions"
)

// GenerateRandomKey generates a random key based on length of input
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

// InitializeStore generates a random key and creates a new session
// to store user data in
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

// ValidatePassword checks that a password contains at least 6 characters
// and contains at least one digit and one uppercase letter
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

// ValidateName checks whether the given contains at least two letters
// and doesn't contain any non-letter characters
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

// ValidateImage checks whether the filename contains a valid file extension
func ValidateImage(filename string) bool {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

// ParseImageString takes a ';' separated list of image file paths and returns them as an array of strings
func ParseImageString(imagesString string) []string {
	var images []string
	start := 0
	for i := 0; i < len(imagesString); i++ {
		if imagesString[i] == ';' {
			images = append(images, imagesString[start:i])
			start = i + 1
		}
	}

	if images == nil {
		return []string{imagesString}
	}

	return images
}

func GetFirstImageFromString(imagesString string) string {
	endChar := 0
	for i := 0; i < len(imagesString); i++ {
		if imagesString[i] == ';' {
			endChar = i
			break
		}
	}

	if endChar == 0 {
		return imagesString
	}

	return imagesString[0:endChar]
}

func RemoveQueryBrackets(queryString string) string {
	length := len(queryString)

	if queryString[0] == '%' && queryString[length-1] == '%' {
		return queryString[1 : length-1]
	}

	return queryString
}
