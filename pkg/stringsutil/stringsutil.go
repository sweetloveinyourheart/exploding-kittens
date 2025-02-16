package stringsutil

import (
	"math/rand"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// IsBlank returns true if a string is empty or contains only whitespace.
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// GenerateRandomCode generates a random code of the given length.
func GenerateRandomCode(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("randome code length must be greater than 0")
	}

	src := rand.NewSource(time.Now().UnixNano()) // Avoid global rand source for thread safety.
	r := rand.New(src)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b), nil
}

// ConvertUUIDsToStrings converts a slice of UUIDs to a slice of their string representations.
//
// Parameters:
//
//	uuids []uuid.UUID - A slice of UUIDs to be converted.
//
// Returns:
//
//	[]string - A slice containing the string representations of the provided UUIDs.
func ConvertUUIDsToStrings(uuids []uuid.UUID) []string {
	strings := make([]string, len(uuids))
	for i, u := range uuids {
		strings[i] = u.String()
	}
	return strings
}
