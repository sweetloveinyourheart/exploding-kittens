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
	eventing.RegisterEventData[GameArgsInitialized](EventTypeGameArgsInitialized, args...)
}

// EventTypeGameCreated is the event type for when a game is created
var EventTypeGameCreated = (&GameCreated{}).EventType()

// EventTypeGameCreated is the event type for when a game is created
var EventTypeGameArgsInitialized = (&GameArgsInitialized{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeGameCreated,
	EventTypeGameArgsInitialized,
}

type GameCreated struct {
	GameID    uuid.UUID   `json:"game_id"`
	PlayerIDs []uuid.UUID `json:"player_ids"`
}

func (p *GameCreated) EventType() common.EventType { return "GAME_CREATED" }

func (p *GameCreated) GetGameID() uuid.UUID { return p.GameID }

func (p *GameCreated) GetPlayerIDs() []uuid.UUID { return p.PlayerIDs }

type GameArgsInitialized struct {
	GameID      uuid.UUID               `json:"game_id"`
	Desk        uuid.UUID               `json:"desk"`
	PlayerHands map[uuid.UUID]uuid.UUID `json:"player_hands"`
	PlayerTurn  uuid.UUID               `json:"player_turn"`
}

func (p *GameArgsInitialized) EventType() common.EventType { return "GAME_ARGS_INITIALIZED" }

func (p *GameArgsInitialized) GetGameID() uuid.UUID { return p.GameID }

func (p *GameArgsInitialized) GetDesk() uuid.UUID { return p.Desk }

func (p *GameArgsInitialized) GetPlayerHands() map[uuid.UUID]uuid.UUID { return p.PlayerHands }

func (p *GameArgsInitialized) GetPlayerTurn() uuid.UUID { return p.PlayerTurn }
