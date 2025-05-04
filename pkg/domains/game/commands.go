package game

import (
	"slices"

	"github.com/gofrs/uuid"

	card_effects "github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/card-effects"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateGame, *CreateGame]()
	eventing.RegisterCommand[InitializeGame, *InitializeGame]()
	eventing.RegisterCommand[StartTurn, *StartTurn]()
	eventing.RegisterCommand[FinishTurn, *FinishTurn]()

	eventing.RegisterCommand[PlayCard, *PlayCard]()
	eventing.RegisterCommand[CreateAction, *CreateAction]()
	eventing.RegisterCommand[ExecuteAction, *ExecuteAction]()
}

const (
	CreateGameCommand     = common.CommandType("game:create")
	InitializeGameCommand = common.CommandType("game:init")
	StartTurnCommand      = common.CommandType("game:turn:start")
	FinishTurnCommand     = common.CommandType("game:turn:finish")

	PlayCardCommand      = common.CommandType("game:card:play")
	CreateActionCommand  = common.CommandType("game:action:create")
	ExecuteActionCommand = common.CommandType("game:action:execute")
)

var AllCommands = []common.CommandType{
	CreateGameCommand,
	InitializeGameCommand,
	StartTurnCommand,
	FinishTurnCommand,

	PlayCardCommand,
	CreateActionCommand,
	ExecuteActionCommand,
}

var _ = eventing.Command(&CreateGame{})
var _ = eventing.Command(&InitializeGame{})
var _ = eventing.Command(&StartTurn{})
var _ = eventing.Command(&FinishTurn{})
var _ = eventing.Command(&PlayCard{})
var _ = eventing.Command(&CreateAction{})
var _ = eventing.Command(&ExecuteAction{})

type CreateGame struct {
	GameID    uuid.UUID   `json:"game_id"`
	PlayerIDs []uuid.UUID `json:"player_ids"`
}

func (c *CreateGame) AggregateType() common.AggregateType { return AggregateType }

func (c *CreateGame) AggregateID() string { return c.GameID.String() }

func (c *CreateGame) CommandType() common.CommandType { return CreateGameCommand }

func (c *CreateGame) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if len(c.PlayerIDs) == 0 {
		return &common.CommandFieldError{Field: "player_ids", Details: "empty list"}
	}

	return nil
}

type InitializeGame struct {
	GameID      uuid.UUID               `json:"game_id"`
	Desk        uuid.UUID               `json:"desk_id"`
	PlayerHands map[uuid.UUID]uuid.UUID `json:"player_hands"`
}

func (c *InitializeGame) AggregateType() common.AggregateType { return AggregateType }

func (c *InitializeGame) AggregateID() string { return c.GameID.String() }

func (c *InitializeGame) CommandType() common.CommandType { return InitializeGameCommand }

func (c *InitializeGame) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if c.Desk == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	return nil
}

type StartTurn struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (c *StartTurn) AggregateType() common.AggregateType { return AggregateType }
func (c *StartTurn) AggregateID() string                 { return c.GameID.String() }
func (c *StartTurn) CommandType() common.CommandType     { return StartTurnCommand }
func (c *StartTurn) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if c.PlayerID == uuid.Nil {
		return &common.CommandFieldError{Field: "player_id", Details: "empty field"}
	}

	return nil
}

type FinishTurn struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
}

func (c *FinishTurn) AggregateType() common.AggregateType { return AggregateType }
func (c *FinishTurn) AggregateID() string                 { return c.GameID.String() }
func (c *FinishTurn) CommandType() common.CommandType     { return FinishTurnCommand }
func (c *FinishTurn) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if c.PlayerID == uuid.Nil {
		return &common.CommandFieldError{Field: "player_id", Details: "empty field"}
	}

	return nil
}

type PlayCard struct {
	GameID   uuid.UUID   `json:"game_id"`
	PlayerID uuid.UUID   `json:"player_id"`
	CardIDs  []uuid.UUID `json:"card_ids"`
}

func (c *PlayCard) AggregateType() common.AggregateType { return AggregateType }

func (c *PlayCard) AggregateID() string { return c.GameID.String() }

func (c *PlayCard) CommandType() common.CommandType { return PlayCardCommand }

func (c *PlayCard) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if c.PlayerID == uuid.Nil {
		return &common.CommandFieldError{Field: "player_id", Details: "empty field"}
	}

	if len(c.CardIDs) == 0 {
		return &common.CommandFieldError{Field: "card_ids", Details: "empty list"}
	}

	return nil
}

type CreateAction struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id"`
	Effect   string    `json:"effect"`
}

func (c *CreateAction) AggregateType() common.AggregateType { return AggregateType }

func (c *CreateAction) AggregateID() string { return c.GameID.String() }

func (c *CreateAction) CommandType() common.CommandType { return CreateActionCommand }

func (c *CreateAction) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if c.PlayerID == uuid.Nil {
		return &common.CommandFieldError{Field: "player_id", Details: "empty field"}
	}

	if c.Effect == "" {
		return &common.CommandFieldError{Field: "effect", Details: "empty field"}
	}

	// Check if the effect is valid
	if !slices.Contains(card_effects.AllCardEffects, c.Effect) {
		return &common.CommandFieldError{Field: "effect", Details: "invalid effect"}
	}

	return nil
}

type ExecuteAction struct {
	GameID   uuid.UUID `json:"game_id"`
	TargetID uuid.UUID `json:"target_id"`
	Effect   string    `json:"effect"`
}

func (c *ExecuteAction) AggregateType() common.AggregateType { return AggregateType }

func (c *ExecuteAction) AggregateID() string { return c.GameID.String() }

func (c *ExecuteAction) CommandType() common.CommandType { return ExecuteActionCommand }

func (c *ExecuteAction) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if c.TargetID == uuid.Nil {
		return &common.CommandFieldError{Field: "target_id", Details: "empty field"}
	}

	if c.Effect == "" {
		return &common.CommandFieldError{Field: "effect", Details: "empty field"}
	}

	// Check if the effect is valid
	if !slices.Contains(card_effects.AllCardEffects, c.Effect) {
		return &common.CommandFieldError{Field: "effect", Details: "invalid effect"}
	}

	return nil
}
