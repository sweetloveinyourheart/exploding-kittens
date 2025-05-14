package hand

import "github.com/cockroachdb/errors"

var (
	ErrHandAlreadyCreated = errors.New("hand already created")
	ErrHandNotAvailable   = errors.New("hand is not available")
	ErrNoCardsAvailable   = errors.New("no cards available")
	ErrCardsNotFound      = errors.New("cards not found")
)
