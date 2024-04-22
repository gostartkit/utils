package utils

import "errors"

var (
	// ErrEmailOrPasswordInvalid email or password invalid
	ErrEmailOrPasswordInvalid = errors.New("email or password invalid")
	// ErrEmailNotVerified email not verified
	ErrEmailNotVerified = errors.New("email not verified")
	// ErrEmailAddressIsAlreadyInUse eemail address is already in use
	ErrEmailAddressIsAlreadyInUse = errors.New("email address is already in use")
	// ErrPhoneOrPasswordInvalid phone or password invalid
	ErrPhoneOrPasswordInvalid = errors.New("phone or password invalid")
	// ErrPhoneNotVerified phone not verified
	ErrPhoneNotVerified = errors.New("phone not verified")
	// ErrPhoneNumberIsAlreadyInUse phone number is already in use
	ErrPhoneNumberIsAlreadyInUse = errors.New("phone number is already in use")
	// ErrDecodingPemBlock decoding PEM block error
	ErrDecodingPemBlock = errors.New("decoding PEM block error")
	// ErrConvertingEcPublicKey converting EC public key error
	ErrConvertingEcPublicKey = errors.New("converting EC public key error")
	// ErrVersionHeaderInvalid version header invalid
	ErrVersionHeaderInvalid = errors.New("version header invalid")
)
