package game

import "github.com/cockroachdb/errors"

var (
	ErrNotEnoughUsersToPlay       = errors.New("game is not enough players to run")
	ErrTooManyUsersToPlay         = errors.New("game has too many players to run")
	ErrGameAlreadyCreated         = errors.New("game already created")
	ErrGameNotAvailable           = errors.New("game is not available")
	ErrGameNotFound               = errors.New("game is not found")
	ErrGameAlreadyInitialized     = errors.New("game has already been initialized")
	ErrPlayerNotInTheirTurn       = errors.New("player is not in their turn")
	ErrPlayerIsAlreadyInTheirTurn = errors.New("player is already in their turn")
	ErrActionEffectNotProvided    = errors.New("action effect is not provided")
	ErrGameNotInActionPhase       = errors.New("game is not in action phase")
	ErrGameNotInPlayPhase         = errors.New("game is not in play phase")
	ErrInvalidIndexToPlant        = errors.New("invalid index to plant")
)
