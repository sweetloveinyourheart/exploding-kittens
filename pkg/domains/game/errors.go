package game

import "github.com/cockroachdb/errors"

var (
	ErrNotEnoughUsersToPlay   = errors.New("game is not enough players to run")
	ErrTooManyUsersToPlay     = errors.New("game has too many players to run")
	ErrGameAlreadyCreated     = errors.New("game already created")
	ErrGameNotAvailable       = errors.New("game is not available")
	ErrGameAlreadyInitialized = errors.New("game has already been initialized")
)
