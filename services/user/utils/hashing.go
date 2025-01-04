package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/pbkdf2"
)

// This file contains the hashing and random string genration functions, this should be
// analyzed for security.
const (
	RAND_NUMBERS    = "0123456789"
	RAND_UPPER_CASE = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	RAND_LOWER_CASE = "abcdefghijklmnopqrstuvwxyz"

	SYSTEM_SALT              = "cWNZhnAiHeO8Y6UFW2HKor5WE6vbq8SC"
	PASSWORD_HASH_ITERATIONS = 1000
	SESSION_HASH_ITERATIONS  = 10
)

// Generates a random string of inputed length
func GenerateRandomString(length int) (string, error) {
	characters := RAND_NUMBERS + RAND_UPPER_CASE + RAND_LOWER_CASE
	charactersLength := int64(len(characters))

	randomString := ""

	for i := 0; i < length; i++ {
		bigIndex, err := rand.Int(rand.Reader, big.NewInt(charactersLength))
		if err != nil {
			return "", err
		}

		// Convert to int
		index := int(bigIndex.Int64())

		randomString += string(characters[index])
	}

	return randomString, nil
}

func SessionHash(userId uuid.UUID) (string, error) {
	// Review this for security, while it seems unsafe to store the system salt in the password hash function
	// as losing it would result in users no longer being able to log in without resetting their password (disaster)
	// losing the salt for the session hash would be less of a disaster as it would simply result in all current users
	// being logged out and having to log back in.  Considering the session hashes are sent out and are techinically
	// visible, it might be worthwhile

	// Start with a random string
	randomStr, err := GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	hashBase := make([]byte, 0, len(randomStr)+len(SYSTEM_SALT)+len(userId))
	hashBase = append(hashBase, []byte(randomStr)...)

	// Add system salt
	hashBase = append(hashBase, []byte(SYSTEM_SALT)...)

	// Add user id
	hashBase = append(hashBase, userId.Bytes()...)

	// Run pbkdf2 some amount of times
	hashBytes := pbkdf2.Key(hashBase, []byte(SYSTEM_SALT), PASSWORD_HASH_ITERATIONS, 32, sha256.New)

	// Convert to string
	encodedHash := base64.URLEncoding.EncodeToString(hashBytes)

	// Truncate string to 32 characters
	return encodedHash[:32], nil
}
