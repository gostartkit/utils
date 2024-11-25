package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

// CreatePrivateKeyPEM create private key pem
func CreatePrivateKeyPEM() ([]byte, error) {

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)

	if err != nil {
		return nil, err
	}

	privateKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyPEM := pem.EncodeToMemory(privateKeyPEMBlock)

	return privateKeyPEM, nil
}

// CreatePublicKeyPEM create public key pem
func CreatePublicKeyPEM(privateKey *ecdsa.PrivateKey) ([]byte, error) {

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return nil, err
	}

	publicKeyPEMBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyPEM := pem.EncodeToMemory(publicKeyPEMBlock)

	return publicKeyPEM, nil
}

// CreatePrivateKey return *ecdsa.PrivateKey
func CreatePrivateKey(privateKeyPEM []byte) (*ecdsa.PrivateKey, error) {

	block, _ := pem.Decode(privateKeyPEM)

	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, ErrPemBlockInvalid
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// CreatePublicKey return *ecdsa.PublicKey
func CreatePublicKey(publicKeyPEM []byte) (*ecdsa.PublicKey, error) {

	block, _ := pem.Decode(publicKeyPEM)

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, ErrPemBlockInvalid
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		return nil, ErrEcPublicKeyInvalid
	}

	return ecdsaPublicKey, nil
}
