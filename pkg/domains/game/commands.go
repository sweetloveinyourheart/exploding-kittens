package game

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateGame, *CreateGame]()
	eventing.RegisterCommand[InitGameArgs, *InitGameArgs]()
}

const (
	CreateGameCommand   = common.CommandType("game:create")
	InitGameArgsCommand = common.CommandType("game:init_args")
)

var AllCommands = []common.CommandType{
	CreateGameCommand,
	InitGameArgsCommand,
}

var _ = eventing.Command(&CreateGame{})

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

type InitGameArgs struct {
	GameID      uuid.UUID               `json:"game_id"`
	Desk        uuid.UUID               `json:"desk"`
	PlayerHands map[uuid.UUID]uuid.UUID `json:"player_hands"`
	PlayerTurn  uuid.UUID               `json:"player_turn"`
}

func (c *InitGameArgs) AggregateType() common.AggregateType { return AggregateType }

func (c *InitGameArgs) AggregateID() string { return c.GameID.String() }

func (c *InitGameArgs) CommandType() common.CommandType { return InitGameArgsCommand }

func (c *InitGameArgs) Validate() error {
	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	return nil
}
