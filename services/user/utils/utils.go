package utils

import (
	"math/rand"
)

func GenerateSessionID() int64 {
	return rand.Int63()
}
