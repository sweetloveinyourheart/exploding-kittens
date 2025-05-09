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
	eventing.RegisterEventData[GameInitialized](EventTypeGameInitialized, args...)
	eventing.RegisterEventData[TurnStarted](EventTypeTurnStarted, args...)
	eventing.RegisterEventData[TurnFinished](EventTypeTurnFinished, args...)
	eventing.RegisterEventData[TurnReversed](EventTypeTurnReversed, args...)
	eventing.RegisterEventData[CardsPlayed](EventTypeCardsPlayed, args...)
	eventing.RegisterEventData[ActionCreated](EventTypeActionCreated, args...)
	eventing.RegisterEventData[ActionExecuted](EventTypeActionExecuted, args...)
}

// EventTypeGameCreated is the event type for when a game is created
var EventTypeGameCreated = (&GameCreated{}).EventType()

// EventTypeGameInitialized is the event type for when a game is initialized
var EventTypeGameInitialized = (&GameInitialized{}).EventType()

// EventTypeCardsPlayed is the event type for when a card is played
var EventTypeCardsPlayed = (&CardsPlayed{}).EventType()

// EventTypeTurnStarted is the event type for when a turn starts
var EventTypeTurnStarted = (&TurnStarted{}).EventType()

// EventTypeTurnFinished is the event type for when a turn finishes
var EventTypeTurnFinished = (&TurnFinished{}).EventType()

// EventTypeTurnReversed is the event type for when a turn is reversed
var EventTypeTurnReversed = (&TurnReversed{}).EventType()

// EventTypeActionCreated is the event type for when an action is created
var EventTypeActionCreated = (&ActionCreated{}).EventType()

// EventTypeActionExecuted is the event type for when an action is executed
var EventTypeActionExecuted = (&ActionExecuted{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeGameCreated,
	EventTypeGameInitialized,
	EventTypeTurnStarted,
	EventTypeTurnFinished,
	EventTypeTurnReversed,
	EventTypeCardsPlayed,
	EventTypeActionCreated,
	EventTypeActionExecuted,
}

type GameCreated struct {
	GameID    uuid.UUID   `json:"game_id"`
	PlayerIDs []uuid.UUID `json:"player_ids"`
}

func (p *GameCreated) EventType() common.EventType { return "GAME_CREATED" }

func (p *GameCreated) GetGameID() uuid.UUID { return p.GameID }

func (p *GameCreated) GetPlayerIDs() []uuid.UUID { return p.PlayerIDs }

type GameInitialized struct {
	GameID      uuid.UUID               `json:"game_id"`
	Desk        uuid.UUID               `json:"desk"`
	PlayerHands map[uuid.UUID]uuid.UUID `json:"player_hands"`
}

func (p *GameInitialized) EventType() common.EventType { return "GAME_INITIALIZED" }

func (p *GameInitialized) GetGameID() uuid.UUID { return p.GameID }

func (p *GameInitialized) GetDesk() uuid.UUID { return p.Desk }

func (p *GameInitialized) GetPlayerHands() map[uuid.UUID]uuid.UUID { return p.PlayerHands }

type TurnStarted struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (p *TurnStarted) EventType() common.EventType { return "GAME_TURN_STARTED" }

func (p *TurnStarted) GetGameID() uuid.UUID { return p.GameID }

func (p *TurnStarted) GetPlayerID() uuid.UUID { return p.PlayerID }

type TurnFinished struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (p *TurnFinished) EventType() common.EventType { return "GAME_TURN_FINISHED" }

func (p *TurnFinished) GetGameID() uuid.UUID { return p.GameID }

func (p *TurnFinished) GetPlayerID() uuid.UUID { return p.PlayerID }

type TurnReversed struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (p *TurnReversed) EventType() common.EventType { return "GAME_TURN_REVERSED" }

func (p *TurnReversed) GetGameID() uuid.UUID { return p.GameID }

func (p *TurnReversed) GetPlayerID() uuid.UUID { return p.PlayerID }

type CardsPlayed struct {
	GameID   uuid.UUID   `json:"game_id"`
	PlayerID uuid.UUID   `json:"player_id"`
	CardIDs  []uuid.UUID `json:"card_ids"`
}

func (p *CardsPlayed) EventType() common.EventType { return "GAME_CARDS_PLAYED" }

func (p *CardsPlayed) GetGameID() uuid.UUID { return p.GameID }

func (p *CardsPlayed) GetPlayerID() uuid.UUID { return p.PlayerID }

func (p *CardsPlayed) GetCardIDs() []uuid.UUID { return p.CardIDs }

type ActionCreated struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
	Effect   string    `json:"effect"`
}

func (p *ActionCreated) EventType() common.EventType { return "GAME_ACTION_CREATED" }

func (p *ActionCreated) GetGameID() uuid.UUID { return p.GameID }

func (p *ActionCreated) GetPlayerID() uuid.UUID { return p.PlayerID }

func (p *ActionCreated) GetEffect() string { return p.Effect }

type ActionExecuted struct {
	GameID         uuid.UUID `json:"game_id"`
	PlayerID       uuid.UUID `json:"player_id"`
	Effect         string    `json:"effect"`
	TargetPlayerID uuid.UUID `json:"target_player_id"`
	TargetCardID   uuid.UUID `json:"target_card_id"`
}

func (p *ActionExecuted) EventType() common.EventType { return "GAME_ACTION_EXECUTED" }

func (p *ActionExecuted) GetGameID() uuid.UUID { return p.GameID }

func (p *ActionExecuted) GetPlayerID() uuid.UUID { return p.PlayerID }

func (p *ActionExecuted) GetEffect() string { return p.Effect }

func (p *ActionExecuted) GetTargetPlayerID() uuid.UUID { return p.TargetPlayerID }

func (p *ActionExecuted) GetTargetCardID() uuid.UUID { return p.TargetCardID }
