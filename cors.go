package utils

import (
	"strings"
)

// CheckCors check cors
func CheckCors(origin string, domain string) bool {
	return strings.HasSuffix(origin, domain) || strings.HasSuffix(origin, "127.0.0.1:3000")
}

// Cors sets the Cors (Cross-Origin Resource Sharing) related HTTP header information
func Cors(set func(key string, value string), origin string, allow []string) {

	set("Access-Control-Allow-Origin", origin)

	set("Access-Control-Allow-Credentials", "true")

	if len(allow) > 0 {
		set("Access-Control-Allow-Methods", strings.Join(allow, ", "))
	}

	set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Attrs")
	set("Access-Control-Max-Age", "86400")

	set("Vary", "Accept-Encoding, Origin")
}
