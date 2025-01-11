package helpers

import "github.com/cockroachdb/errors"

var (
	ErrInvalidLogin   = errors.New("invalid login")
	ErrInvalidSession = errors.New("invalid session")
)
