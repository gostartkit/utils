package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString random string with [a-zA-Z0-9]
func RandomString(l int) (string, error) {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RandomStringWithCharset(l, characters)
}

// RandStringWithCharset rand string
func RandomStringWithCharset(l int, characters string) (string, error) {

	if l <= 0 {
		return "", ErrLengthInvalid
	}

	max := len(characters)

	if max == 0 {
		return "", ErrLengthInvalid
	}

	b := make([]byte, l)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	for i := 0; i < l; i++ {
		idx := int(b[i]) % max
		b[i] = characters[idx]
	}

	return string(b), nil
}

// RandString random string
func RandString(l int) (string, error) {

	b := make([]byte, l/2)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
