package utils

import "errors"

var (
	// ErrDecodingPemBlock decoding PEM block error
	ErrDecodingPemBlock = errors.New("decoding PEM block error")
	// ErrConvertingEcPublicKey converting EC public key error
	ErrConvertingEcPublicKey = errors.New("converting EC public key error")
)
