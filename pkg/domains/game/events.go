package game

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// registerEvents registers the event types for the game domain
func registerEvents(subjFunc eventing.SubjectFunc, subjRootFunc eventing.SubjectFunc, subjTokenPos int, tokensFunc eventing.TokensFunc) {
	args := make([]eventing.EventRegistrationOption, 0)
	args = append(args, eventing.WithRegisterSubjectRootFunc(subjRootFunc))
	args = append(args, eventing.WithRegisterSubjectFunc(subjFunc))
	args = append(args, eventing.WithRegisterTokensFunc(tokensFunc))
	args = append(args, eventing.WithRegisterSubjectTokenPosition(subjTokenPos))

	eventing.RegisterEventData[GameCreated](EventTypeGameCreated, args...)
}

// EventTypeGameCreated is the event type for when a game is created
var EventTypeGameCreated = (&GameCreated{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeGameCreated,
}

type GameCreated struct {
	GameID    uuid.UUID   `json:"game_id"`
	PlayerIDs []uuid.UUID `json:"player_ids"`
}

func (p *GameCreated) EventType() common.EventType { return "GAME_CREATED" }

func (p *GameCreated) GetGameID() uuid.UUID { return p.GameID }

func (p *GameCreated) GetPlayerIDs() []uuid.UUID { return p.PlayerIDs }
