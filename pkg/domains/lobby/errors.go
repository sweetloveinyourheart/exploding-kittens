package lobby

import "github.com/cockroachdb/errors"

var (
	ErrLobbyAlreadyCreated   = errors.New("lobby already created")
	ErrLobbyInWaitingMode    = errors.New("lobby is in wating mode")
	ErrLobbyNotAvailable     = errors.New("lobby is not available")
	ErrHostUserNotRecognized = errors.New("host user is not recognized")
	ErrGameIsAlreadyStarted  = errors.New("game is already stared")
	ErrGameIsNotEnoughPlayer = errors.New("game is not enough player to start")
)
