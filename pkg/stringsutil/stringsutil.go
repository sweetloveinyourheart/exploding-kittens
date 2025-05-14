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

// ConvertUUIDToStringPtr converts a given UUID to its string representation
// and returns a pointer to the resulting string.
//
// Parameters:
//   - uuid: The UUID to be converted.
//
// Returns:
//   - A pointer to the string representation of the provided UUID.
func ConvertUUIDToStringPtr(uuid uuid.UUID) *string {
	if uuid.IsNil() {
		return nil
	}

	str := uuid.String()
	return &str
}

// ConvertStringsToUUIDs converts a slice of strings to a slice of UUIDs,
// skipping any strings that are not valid UUIDs.
//
// Parameters:
//
//	strings []string - A slice of strings to be converted.
//
// Returns:
//
//	[]uuid.UUID - A slice containing the UUIDs converted from the valid strings.
func ConvertStringsToUUIDs(strings []string) []uuid.UUID {
	var uuids []uuid.UUID
	for _, s := range strings {
		u, err := uuid.FromString(s)
		if err == nil {
			uuids = append(uuids, u)
		}
	}
	return uuids
}

// ConvertStringToUUID converts a string to a UUID.
//
// Parameters:
//
//	s string - The string to be converted.
//
// Returns:
//
//	uuid.UUID - The UUID converted from the string.
func ConvertStringToUUID(s string) uuid.UUID {
	u, err := uuid.FromString(s)
	if err != nil {
		return uuid.Nil
	}
	return u
}
