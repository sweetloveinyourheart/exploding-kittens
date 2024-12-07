package stringsutil

import "strings"

// IsBlank returns false if a string has a non-zero length and doesn't contain only spaces.
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
