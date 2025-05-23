package desk

import "github.com/cockroachdb/errors"

var (
	ErrDeskAlreadyCreated = errors.New("desk already created")
	ErrDeskNotAvailable   = errors.New("desk is not available")
	ErrNoCardsToDiscard   = errors.New("no cards to discard")
	ErrInvalidPeekCount   = errors.New("invalid peek count")
	ErrInvalidDrawCount   = errors.New("invalid draw count")
	ErrNoCardToDraw       = errors.New("no cards to draw")
	ErrInvalidInsertIndex = errors.New("invalid insert index")
)
