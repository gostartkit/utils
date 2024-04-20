package utils

import (
	"fmt"
	"regexp"
)

const (
	strAlphabet = "^([a-zA-Z]+)$"
	strDomain   = `^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`
)

var (
	regexAlphabet = regexp.MustCompile(strAlphabet)
	regexDomain   = regexp.MustCompile(strDomain)

	regexEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	regexPhone = regexp.MustCompile(`^[1-9][0-9]{10,11}$`)
)

// Int validate name
func Int(name string, val int) error {

	if val == 0 {
		return fmt.Errorf("%s invalid", name)
	}

	return nil
}

// Uint64 validate name
func Uint64(name string, val uint64) error {

	if val == 0 {
		return fmt.Errorf("%s invalid", name)
	}

	return nil
}

// String validate name
func String(name string, val string) error {

	if val == "" {
		return fmt.Errorf("%s required", name)
	}

	return nil
}

// RString validate name
func RString(name string, val *string) error {

	if val == nil {
		return fmt.Errorf("%s invalid", name)
	}

	if *val == "" {
		return fmt.Errorf("%s required", name)
	}

	return nil
}

// ValidateName validate name
func ValidateName(name ...string) error {

	for i := 0; i < len(name); i++ {

		if !MatchAlphabet(name[i]) {
			return fmt.Errorf("'%s' invalid, must match regex /%s/", name[i], strAlphabet)
		}
	}

	return nil
}

// ValidateDomain validate domain
func ValidateDomain(domain ...string) error {

	for i := 0; i < len(domain); i++ {

		if !MatchDomain(domain[i]) {
			return fmt.Errorf("'%s' invalid, must match regex /%s/", domain[i], strDomain)
		}
	}

	return nil
}

// ValidateEmail validate email
func ValidateEmail(email ...string) error {

	for i := 0; i < len(email); i++ {

		if !MatchEmail(email[i]) {
			return fmt.Errorf("email '%s' invalid", email[i])
		}
	}

	return nil
}

// ValidatePhone validate phone number
func ValidatePhone(phone ...string) error {

	for i := 0; i < len(phone); i++ {

		if !MatchPhone(phone[i]) {
			return fmt.Errorf("phone '%s' invalid", phone[i])
		}
	}

	return nil
}

// MatchAlphabet check alphabet
func MatchAlphabet(val string) bool {
	return regexAlphabet.MatchString(val)
}

// MatchEmail check email
func MatchEmail(val string) bool {
	return regexEmail.MatchString(val)
}

// MatchPhone check phone
func MatchPhone(val string) bool {
	return regexPhone.MatchString(val)
}

// MatchDomain check domain
func MatchDomain(val string) bool {
	return regexDomain.MatchString(val)
}
