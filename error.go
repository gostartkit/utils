package utils

import "errors"

var (
	// ErrNotVerified data not verified
	ErrNotVerified = errors.New("data not verified")
	// ErrEmailNotVerified email not verified
	ErrEmailNotVerified = errors.New("email not verified")
	// ErrEmailAddressUsed email address used
	ErrEmailAddressUsed = errors.New("email address used")
	// ErrPhoneNotVerified phone not verified
	ErrPhoneNotVerified = errors.New("phone not verified")
	// ErrPhoneNumberUsed phone number used
	ErrPhoneNumberUsed = errors.New("phone number used")
	// ErrInvalid data invalid
	ErrInvalid = errors.New("data invalid")
	// ErrVersionInvalid version invalid
	ErrVersionInvalid = errors.New("version invalid")
	// ErrLengthInvalid invalid length
	ErrLengthInvalid = errors.New("length invalid")
	// ErrFilterInvalid invalid filter
	ErrFilterInvalid = errors.New("filter invalid")
	// ErrOrderByInvalid invalid orderby
	ErrOrderByInvalid = errors.New("orderby invalid")
	// ErrEcPublicKeyInvalid ec public key invalid
	ErrEcPublicKeyInvalid = errors.New("ec public key invalid")
	// ErrPemBlockInvalid pem block invalid
	ErrPemBlockInvalid = errors.New("pem block invalid")
)
