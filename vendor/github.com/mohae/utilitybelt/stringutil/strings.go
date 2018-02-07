// Package stringutil provieds some functions for working with strings.
package stringutil

import (
	"strconv"
	"strings"
)

// SetPrefix checks the string to see if it already has the prefix. If it
// doesn't have the prefix, the prefix is prefixed to the string.
//
// Either way, the returned string will not repeat the prefix, unless the
// source string already repeats the prefix.
func SetPrefix(s string, prefix string) string {
	if prefix == "" {
		return s
	}

	if strings.HasPrefix(s, prefix) {
		return s
	}

	return prefix + s
}

// SetSuffix checks the string to see if it already has the suffix. If it
// doesn't have the suffix, the suffix is appended to the string.
//
// Either way, the returned string will not repeat the suffix, unless the
// source string already repeats the suffix.
func SetSuffix(s string, suffix string) string {
	if suffix == "" {
		return s
	}

	if strings.HasSuffix(s, suffix) {
		return s
	}

	return s + suffix
}

// ParseBool wraps strconv.ParseBool. It treats an error condition as a false.
func ParseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return b
}
