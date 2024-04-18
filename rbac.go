package utils

import (
	"strings"
)

// CheckKeys check right
func CheckKeys(userRight int64, val func(c string) int64, keys ...string) bool {

	if len(keys) == 0 {
		return true
	}

	return checkKeys(userRight, val, keys...)
}

// CheckVal check right
func CheckVals(userRight int64, vals ...int64) bool {

	if len(vals) == 0 {
		return true
	}

	return checkVals(userRight, vals...)
}

// checkKeys check right by keys
func checkKeys(userRight int64, val func(c string) int64, keys ...string) bool {

	for _, key := range keys {
		if !checkKey(userRight, val, key) {
			return false
		}
	}

	return true
}

// checkKey check right by key
func checkKey(userRight int64, val func(c string) int64, key string) bool {

	orKeys := strings.Split(key, "|")

	for _, orKey := range orKeys {

		if orKey == "" {
			return userRight > 0
		}

		right := val(orKey)

		if right > 0 {
			return right&userRight > 0
		}
	}

	return false
}

// checkVals check right
func checkVals(userRight int64, vals ...int64) bool {

	for _, val := range vals {
		if !checkVal(userRight, val) {
			return false
		}
	}

	return true
}

// checkVal check right
func checkVal(userRight int64, val int64) bool {

	if val > 0 {
		return val&userRight > 0
	}

	return false
}
