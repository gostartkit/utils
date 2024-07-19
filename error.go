package utils

import "errors"

var (
	// ErrEmailOrPasswordInvalid email or password invalid
	ErrEmailOrPasswordInvalid = errors.New("email or password invalid")
	// ErrEmailNotVerified email not verified
	ErrEmailNotVerified = errors.New("email not verified")
	// ErrEmailAddressUsed email address used
	ErrEmailAddressUsed = errors.New("email address used")
	// ErrPhoneOrPasswordInvalid phone or password invalid
	ErrPhoneOrPasswordInvalid = errors.New("phone or password invalid")
	// ErrPhoneNotVerified phone not verified
	ErrPhoneNotVerified = errors.New("phone not verified")
	// ErrPhoneNumberUsed phone number used
	ErrPhoneNumberUsed = errors.New("phone number used")
	// ErrDecodingPemBlock decoding PEM block error
	ErrDecodingPemBlock = errors.New("decoding PEM block error")
	// ErrConvertingEcPublicKey converting EC public key error
	ErrConvertingEcPublicKey = errors.New("converting EC public key error")
	// ErrVersionInvalid version invalid
	ErrVersionInvalid = errors.New("invalid version")
	// ErrLengthInvalid invalid length
	ErrLengthInvalid = errors.New("invalid length")
	// ErrFilterInvalid invalid filter
	ErrFilterInvalid = errors.New("invalid filter")
	// ErrOrderByInvalid invalid orderby
	ErrOrderByInvalid = errors.New("invalid orderby")
)
