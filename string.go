package utils

import (
	"crypto/rand"
	"math/big"
)

// RandomString random string
func RandomString(len int) (string, error) {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RandomStringWithCharset(len, characters)
}

// RandomStringWithCharset random string
func RandomStringWithCharset(length int, characters string) (string, error) {

	max := big.NewInt(int64(len(characters)))

	b := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = characters[randomIndex.Int64()]
	}

	return string(b), nil
}
