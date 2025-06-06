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
	eventing.RegisterEventData[GameStarted](EventTypeGameStarted, args...)
	eventing.RegisterEventData[TurnStarted](EventTypeTurnStarted, args...)
	eventing.RegisterEventData[TurnFinished](EventTypeTurnFinished, args...)
	eventing.RegisterEventData[TurnReversed](EventTypeTurnReversed, args...)
	eventing.RegisterEventData[GameFinished](EventTypeGameFinished, args...)
	eventing.RegisterEventData[CardsPlayed](EventTypeCardsPlayed, args...)
	eventing.RegisterEventData[ActionCreated](EventTypeActionCreated, args...)
	eventing.RegisterEventData[ActionExecuted](EventTypeActionExecuted, args...)
	eventing.RegisterEventData[AffectedPlayerSelected](EventTypeAffectedPlayerSelected, args...)
	eventing.RegisterEventData[CardDrawn](EventTypeCardDrawn, args...)
	eventing.RegisterEventData[ExplodingDrawn](EventTypeExplodingDrawn, args...)
	eventing.RegisterEventData[ExplodingDefused](EventTypeExplodingDefused, args...)
	eventing.RegisterEventData[PlayerEliminated](EventTypePlayerEliminated, args...)
	eventing.RegisterEventData[KittenPlanted](EventTypeKittenPlanted, args...)
}

// EventTypeGameCreated is the event type for when a game is created
var EventTypeGameCreated = (&GameCreated{}).EventType()

// EventTypeGameInitialized is the event type for when a game is initialized
var EventTypeGameInitialized = (&GameInitialized{}).EventType()

// EventTypeGameStarted is the event type for when a game is started
var EventTypeGameStarted = (&GameStarted{}).EventType()

// EventTypeCardsPlayed is the event type for when a card is played
var EventTypeCardsPlayed = (&CardsPlayed{}).EventType()

// EventTypeTurnStarted is the event type for when a turn starts
var EventTypeTurnStarted = (&TurnStarted{}).EventType()

// EventTypeTurnFinished is the event type for when a turn finishes
var EventTypeTurnFinished = (&TurnFinished{}).EventType()

// EventTypeTurnReversed is the event type for when a turn is reversed
var EventTypeTurnReversed = (&TurnReversed{}).EventType()

// EventTypeGameFinished is the event type for when a game is finished
var EventTypeGameFinished = (&GameFinished{}).EventType()

// EventTypeActionCreated is the event type for when an action is created
var EventTypeActionCreated = (&ActionCreated{}).EventType()

// EventTypeActionExecuted is the event type for when an action is executed
var EventTypeActionExecuted = (&ActionExecuted{}).EventType()

// EventTypeAffectedPlayerSelected is the event type for when an affected player is selected
var EventTypeAffectedPlayerSelected = (&AffectedPlayerSelected{}).EventType()

// EventTypeCardDrawn is the event type for when cards are drawn
var EventTypeCardDrawn = (&CardDrawn{}).EventType()

// EventTypeExplodingDrawn is the event type for when an exploding kitten is drawn
var EventTypeExplodingDrawn = (&ExplodingDrawn{}).EventType()

// EventTypeExplodingDefused is the event type for when an exploding
var EventTypeExplodingDefused = (&ExplodingDefused{}).EventType()

// EventTypePlayerEliminated is the event type for when a player is eliminated
var EventTypePlayerEliminated = (&PlayerEliminated{}).EventType()

// EventTypeKittenPlanted is the event type for when a kitten is planted
var EventTypeKittenPlanted = (&KittenPlanted{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeGameCreated,
	EventTypeGameInitialized,
	EventTypeGameStarted,
	EventTypeTurnStarted,
	EventTypeTurnFinished,
	EventTypeTurnReversed,
	EventTypeGameFinished,
	EventTypeCardsPlayed,
	EventTypeActionCreated,
	EventTypeActionExecuted,
	EventTypeAffectedPlayerSelected,
	EventTypeCardDrawn,
	EventTypeExplodingDrawn,
	EventTypeExplodingDefused,
	EventTypePlayerEliminated,
	EventTypeKittenPlanted,
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
	DeskID      uuid.UUID               `json:"desk_id"`
	PlayerHands map[uuid.UUID]uuid.UUID `json:"player_hands"`
}

func (p *GameInitialized) EventType() common.EventType { return "GAME_INITIALIZED" }

func (p *GameInitialized) GetGameID() uuid.UUID { return p.GameID }

func (p *GameInitialized) GetDeskID() uuid.UUID { return p.DeskID }

func (p *GameInitialized) GetPlayerHands() map[uuid.UUID]uuid.UUID { return p.PlayerHands }

type GameStarted struct {
	GameID uuid.UUID `json:"game_id"`
}

func (p *GameStarted) EventType() common.EventType { return "GAME_STARTED" }

func (p *GameStarted) GetGameID() uuid.UUID { return p.GameID }

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
	GameID uuid.UUID `json:"game_id"`
	Effect string    `json:"effect"`
}

func (p *ActionCreated) EventType() common.EventType { return "GAME_ACTION_CREATED" }

func (p *ActionCreated) GetGameID() uuid.UUID { return p.GameID }

func (p *ActionCreated) GetEffect() string { return p.Effect }

type AffectedPlayerSelected struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (p *AffectedPlayerSelected) EventType() common.EventType { return "GAME_AFFECTED_PLAYER_SELECTED" }

func (p *AffectedPlayerSelected) GetGameID() uuid.UUID { return p.GameID }

func (p *AffectedPlayerSelected) GetPlayerID() uuid.UUID { return p.PlayerID }

type ActionExecuted struct {
	GameID uuid.UUID        `json:"game_id"`
	Effect string           `json:"effect"`
	Args   *ActionArguments `json:"args"`
}

func (p *ActionExecuted) EventType() common.EventType { return "GAME_ACTION_EXECUTED" }

func (p *ActionExecuted) GetGameID() uuid.UUID { return p.GameID }

func (p *ActionExecuted) GetEffect() string { return p.Effect }

func (p *ActionExecuted) GetArgs() *ActionArguments { return p.Args }

type ActionArguments struct {
	CardIDs     []uuid.UUID `json:"card_ids"`
	CardIndexes []int       `json:"card_indexes"`
}

func (p *ActionArguments) GetCardIDs() []uuid.UUID { return p.CardIDs }

func (p *ActionArguments) GetCardIndexes() []int { return p.CardIndexes }

type CardDrawn struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (p *CardDrawn) EventType() common.EventType { return "GAME_CARDS_DRAWN" }

func (p *CardDrawn) GetGameID() uuid.UUID { return p.GameID }

func (p *CardDrawn) GetPlayerID() uuid.UUID { return p.PlayerID }

type ExplodingDrawn struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
	CardID   uuid.UUID `json:"card_id"`
}

func (p *ExplodingDrawn) EventType() common.EventType { return "GAME_EXPLODING_DRAWN" }

func (p *ExplodingDrawn) GetGameID() uuid.UUID { return p.GameID }

func (p *ExplodingDrawn) GetPlayerID() uuid.UUID { return p.PlayerID }

func (p *ExplodingDrawn) GetCardID() uuid.UUID { return p.CardID }

type ExplodingDefused struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
	CardID   uuid.UUID `json:"card_id"`
}

func (p *ExplodingDefused) EventType() common.EventType { return "GAME_EXPLODING_DEFUSED" }

func (p *ExplodingDefused) GetGameID() uuid.UUID { return p.GameID }

func (p *ExplodingDefused) GetPlayerID() uuid.UUID { return p.PlayerID }

func (p *ExplodingDefused) GetCardID() uuid.UUID { return p.CardID }

type PlayerEliminated struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (p *PlayerEliminated) EventType() common.EventType { return "GAME_PLAYER_ELIMINATED" }

func (p *PlayerEliminated) GetGameID() uuid.UUID { return p.GameID }

func (p *PlayerEliminated) GetPlayerID() uuid.UUID { return p.PlayerID }

type KittenPlanted struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
	Index    int       `json:"index"`
}

func (p *KittenPlanted) EventType() common.EventType { return "GAME_KITTEN_PLANTED" }

func (p *KittenPlanted) GetGameID() uuid.UUID { return p.GameID }

func (p *KittenPlanted) GetPlayerID() uuid.UUID { return p.PlayerID }

func (p *KittenPlanted) GetIndex() int { return p.Index }

type GameFinished struct {
	GameID   uuid.UUID `json:"game_id"`
	WinnerID uuid.UUID `json:"winner_id"`
}

func (p *GameFinished) EventType() common.EventType { return "GAME_FINISHED" }

func (p *GameFinished) GetGameID() uuid.UUID { return p.GameID }

func (p *GameFinished) GetWinnerID() uuid.UUID { return p.WinnerID }
