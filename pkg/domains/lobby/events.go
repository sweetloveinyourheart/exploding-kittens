package lobby

import (
	"time"

	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// registerEvents registers the event types for the lobby domain
func registerEvents(subjFunc eventing.SubjectFunc, subjRootFunc eventing.SubjectFunc, subjTokenPos int, tokensFunc eventing.TokensFunc) {
	args := make([]eventing.EventRegistrationOption, 0)
	args = append(args, eventing.WithRegisterSubjectRootFunc(subjRootFunc))
	args = append(args, eventing.WithRegisterSubjectFunc(subjFunc))
	args = append(args, eventing.WithRegisterTokensFunc(tokensFunc))
	args = append(args, eventing.WithRegisterSubjectTokenPosition(subjTokenPos))

	eventing.RegisterEventData[LobbyCreated](EventTypeLobbyCreated, args...)
}

// EventTypeLobbyCreated is the event type for when a lobby is created
var EventTypeLobbyCreated = (&LobbyCreated{}).EventType()

type LobbyCreated struct {
	LobbyID      uuid.UUID   `json:"lobby_id"`
	LobbyCode    string      `json:"lobby_code"`
	LobbyName    string      `json:"lobby_name"`
	HostUserID   uuid.UUID   `json:"host_user_id"`
	Participants []uuid.UUID `json:"participants"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (p *LobbyCreated) EventType() common.EventType { return "LOBBY_CREATED" }

func (p *LobbyCreated) GetLobbyID() uuid.UUID { return p.LobbyID }

func (p *LobbyCreated) GetLobbyCode() string { return p.LobbyCode }

func (p *LobbyCreated) GetLobbyName() string { return p.LobbyName }

func (p *LobbyCreated) GetHostUserID() uuid.UUID { return p.HostUserID }

func (p *LobbyCreated) GetParticipants() []uuid.UUID { return p.Participants }

func (p *LobbyCreated) GetCreatedAt() time.Time { return p.CreatedAt }

func (p *LobbyCreated) GetUpdatedAt() time.Time { return p.UpdatedAt }
