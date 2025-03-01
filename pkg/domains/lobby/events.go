package lobby

import (
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
	eventing.RegisterEventData[LobbyJoined](EventTypeLobbyJoined, args...)
	eventing.RegisterEventData[LobbyLeft](EventTypeLobbyLeft, args...)
}

// EventTypeLobbyCreated is the event type for when a lobby is created
var EventTypeLobbyCreated = (&LobbyCreated{}).EventType()

// EventTypeLobbyJoined is the event type for when a user joins a lobby
var EventTypeLobbyJoined = (&LobbyJoined{}).EventType()

// EventTypeLobbyLeft is the event type for when a user leaves a lobby
var EventTypeLobbyLeft = (&LobbyLeft{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeLobbyCreated,
	EventTypeLobbyJoined,
	EventTypeLobbyLeft,
}

type LobbyCreated struct {
	LobbyID      uuid.UUID   `json:"lobby_id"`
	LobbyCode    string      `json:"lobby_code"`
	LobbyName    string      `json:"lobby_name"`
	HostUserID   uuid.UUID   `json:"host_user_id"`
	Participants []uuid.UUID `json:"participants"`
}

func (p *LobbyCreated) EventType() common.EventType { return "LOBBY_CREATED" }

func (p *LobbyCreated) GetLobbyID() uuid.UUID { return p.LobbyID }

func (p *LobbyCreated) GetLobbyCode() string { return p.LobbyCode }

func (p *LobbyCreated) GetLobbyName() string { return p.LobbyName }

func (p *LobbyCreated) GetHostUserID() uuid.UUID { return p.HostUserID }

func (p *LobbyCreated) GetParticipants() []uuid.UUID { return p.Participants }

type LobbyJoined struct {
	LobbyID uuid.UUID `json:"lobby_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (p *LobbyJoined) EventType() common.EventType { return "LOBBY_JOINED" }

func (p *LobbyJoined) GetLobbyID() uuid.UUID { return p.LobbyID }

func (p *LobbyJoined) GetUserID() uuid.UUID { return p.UserID }

type LobbyLeft struct {
	LobbyID uuid.UUID `json:"lobby_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (p *LobbyLeft) EventType() common.EventType { return "LOBBY_LEAVED" }

func (p *LobbyLeft) GetLobbyID() uuid.UUID { return p.LobbyID }

func (p *LobbyLeft) GetUserID() uuid.UUID { return p.UserID }
