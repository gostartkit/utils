package utils

import "strings"

// CheckCors check cors
func CheckCors(origin string, domain string) bool {
	return strings.HasSuffix(origin, domain) || strings.HasSuffix(origin, "127.0.0.1:3000")
}
