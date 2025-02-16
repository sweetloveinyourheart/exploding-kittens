package match

import (
	"github.com/gofrs/uuid"
)

type Matcher func(subject string) bool

func MatchLobbyID(lobbyID uuid.UUID) Matcher {
	return func(subject string) bool {
		return lobbyID.String() == subject
	}
}

type LobbyIDer interface {
	GetLobbyId() uuid.UUID
}
