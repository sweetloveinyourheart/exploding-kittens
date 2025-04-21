package game

import "github.com/cockroachdb/errors"

var (
	ErrGameIsInitializing     = errors.New("game is initializing")
	ErrNotEnoughUserToPlay    = errors.New("game is not enough players to run")
	ErrGameAlreadyCreated     = errors.New("game already created")
	ErrGameNotAvailable       = errors.New("game is not available")
	ErrGameAlreadyInitialized = errors.New("game has already been initialized")
)
