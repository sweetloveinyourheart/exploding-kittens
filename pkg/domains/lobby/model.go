package lobby

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

type Lobby struct {
	LobbyID      uuid.UUID   `json:"lobby_id"`
	LobbyCode    string      `json:"lobby_code"`
	LobbyName    string      `json:"lobby_name"`
	HostUserID   uuid.UUID   `json:"host_user_id"`
	Participants []uuid.UUID `json:"participants"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

var _ = common.Entity(&Lobby{})

func (t *Lobby) EntityID() string {
	return t.LobbyID.String()
}

func (t *Lobby) GetLobbyID() uuid.UUID {
	return t.LobbyID
}

func (t *Lobby) GetLobbyCode() string {
	return t.LobbyCode
}

func (t *Lobby) GetLobbyName() string {
	return t.LobbyName
}

func (t *Lobby) GetHostUserID() uuid.UUID {
	return t.HostUserID
}

func (t *Lobby) GetParticipants() []uuid.UUID {
	return t.Participants
}

func (t *Lobby) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Lobby) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}
