package game

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateGame, *CreateGame]()
	eventing.RegisterCommand[InitializeGame, *InitializeGame]()

	eventing.RegisterCommand[PlayCard, *PlayCard]()
}

const (
	CreateGameCommand     = common.CommandType("game:create")
	InitializeGameCommand = common.CommandType("game:init")

	PlayCardCommand = common.CommandType("game:card:play")
)

var AllCommands = []common.CommandType{
	CreateGameCommand,
	InitializeGameCommand,

	PlayCardCommand,
}

var _ = eventing.Command(&CreateGame{})
var _ = eventing.Command(&InitializeGame{})
var _ = eventing.Command(&PlayCard{})

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
