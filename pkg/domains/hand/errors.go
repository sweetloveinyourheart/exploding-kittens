package hand

import "github.com/cockroachdb/errors"

var (
	ErrHandAlreadyCreated = errors.New("hand already created")
	ErrHandNotAvailable   = errors.New("hand is not available")
)
