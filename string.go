package utils

import (
	"crypto/rand"
	"math/big"
)

// RandomString random string
func RandomString(l int) (string, error) {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RandomStringWithCharset(l, characters)
}

// RandomStringWithCharset random string
func RandomStringWithCharset(l int, characters string) (string, error) {

	max := big.NewInt(int64(len(characters)))

	b := make([]byte, l)

	for i := 0; i < l; i++ {
		randomIndex, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = characters[randomIndex.Int64()]
	}

	return string(b), nil
}
