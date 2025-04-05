package game

import "github.com/cockroachdb/errors"

var (
	ErrGameIsInitializing = errors.New("game is initializing")
	ErrGameAlreadyCreated = errors.New("game already created")
	ErrGameNotAvailable   = errors.New("game is not available")
)
