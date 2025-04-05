package desk

import "github.com/cockroachdb/errors"

var (
	ErrDeskAlreadyCreated = errors.New("desk already created")
	ErrDeskNotAvailable   = errors.New("desk is not available")
)
