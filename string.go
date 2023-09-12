package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomString random string
func RandomString(len int) (string, error) {

	b := make([]byte, len/2)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
