package lobby

import "github.com/cockroachdb/errors"

var (
	ErrLobbyAlreadyCreated = errors.New("lobby already created")
	ErrLobbyInWaitingMode  = errors.New("lobby is in wating mode")
	ErrLobbyNotAvailable   = errors.New("lobby is not available")
)
