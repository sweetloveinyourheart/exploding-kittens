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

func MatchGameID(gameID uuid.UUID) Matcher {
	return func(subject string) bool {
		return gameID.String() == subject
	}
}

type GameIDer interface {
	GetGameId() uuid.UUID
}

func MatchDeskID(deskID uuid.UUID) Matcher {
	return func(subject string) bool {
		return deskID.String() == subject
	}
}

type DeskIDer interface {
	GetDeskId() uuid.UUID
}

func MatchHandID(handID uuid.UUID) Matcher {
	return func(subject string) bool {
		return handID.String() == subject
	}
}

type HandIDer interface {
	GetHandId() uuid.UUID
}
