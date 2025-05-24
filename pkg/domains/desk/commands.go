package desk

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateDesk, *CreateDesk]()
	eventing.RegisterCommand[ShuffleDesk, *ShuffleDesk]()
	eventing.RegisterCommand[DiscardCards, *DiscardCards]()
	eventing.RegisterCommand[PeekCards, *PeekCards]()
	eventing.RegisterCommand[DrawCard, *DrawCard]()
	eventing.RegisterCommand[InsertCard, *InsertCard]()
}

const (
	CreateDeskCommand   = common.CommandType("desk:create")
	ShuffleDeskCommand  = common.CommandType("desk:shuffle")
	DiscardCardsCommand = common.CommandType("desk:cards:discard")
	PeekCardsCommand    = common.CommandType("desk:cards:peek")
	DrawCardCommand     = common.CommandType("desk:card:draw")
	InsertCardCommand   = common.CommandType("desk:card:insert")
)

var AllCommands = []common.CommandType{
	CreateDeskCommand,
	ShuffleDeskCommand,
	DiscardCardsCommand,
	PeekCardsCommand,
	DrawCardCommand,
	InsertCardCommand,
}

var _ = eventing.Command(&CreateDesk{})
var _ = eventing.Command(&ShuffleDesk{})
var _ = eventing.Command(&DiscardCards{})
var _ = eventing.Command(&PeekCards{})
var _ = eventing.Command(&DrawCard{})
var _ = eventing.Command(&InsertCard{})

type CreateDesk struct {
	DeskID  uuid.UUID   `json:"desk_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (c *CreateDesk) AggregateType() common.AggregateType { return AggregateType }

func (c *CreateDesk) AggregateID() string { return c.DeskID.String() }

func (c *CreateDesk) CommandType() common.CommandType { return CreateDeskCommand }

func (c *CreateDesk) Validate() error {
	if c.DeskID == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	if len(c.CardIDs) == 0 {
		return &common.CommandFieldError{Field: "card_ids", Details: "empty list"}
	}

	return nil
}

type ShuffleDesk struct {
	DeskID uuid.UUID `json:"desk_id"`
}

func (c *ShuffleDesk) AggregateType() common.AggregateType { return AggregateType }

func (c *ShuffleDesk) AggregateID() string { return c.DeskID.String() }

func (c *ShuffleDesk) CommandType() common.CommandType { return ShuffleDeskCommand }

func (c *ShuffleDesk) Validate() error {
	if c.DeskID == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	return nil
}

type DiscardCards struct {
	DeskID  uuid.UUID   `json:"desk_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (c *DiscardCards) AggregateType() common.AggregateType { return AggregateType }

func (c *DiscardCards) AggregateID() string { return c.DeskID.String() }

func (c *DiscardCards) CommandType() common.CommandType { return DiscardCardsCommand }

func (c *DiscardCards) Validate() error {
	if c.DeskID == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	if len(c.CardIDs) == 0 {
		return &common.CommandFieldError{Field: "card_ids", Details: "empty list"}
	}

	return nil
}

type PeekCards struct {
	DeskID uuid.UUID `json:"desk_id"`
	Count  int       `json:"count"`
}

func (c *PeekCards) AggregateType() common.AggregateType { return AggregateType }

func (c *PeekCards) AggregateID() string { return c.DeskID.String() }

func (c *PeekCards) CommandType() common.CommandType { return PeekCardsCommand }

func (c *PeekCards) Validate() error {
	if c.DeskID == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	if c.Count <= 0 {
		return &common.CommandFieldError{Field: "count", Details: "invalid count"}
	}

	return nil
}

type DrawCard struct {
	DeskID        uuid.UUID `json:"desk_id"`
	GameID        uuid.UUID `json:"game_id"`
	PlayerID      uuid.UUID `json:"player_id"`
	CanFinishTurn bool      `json:"can_finish_turn"` // Indicates if the player can finish their turn after drawing a card
}

func (c *DrawCard) AggregateType() common.AggregateType { return AggregateType }

func (c *DrawCard) AggregateID() string { return c.DeskID.String() }

func (c *DrawCard) CommandType() common.CommandType { return DrawCardCommand }

func (c *DrawCard) Validate() error {
	if c.DeskID == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	if c.GameID == uuid.Nil {
		return &common.CommandFieldError{Field: "game_id", Details: "empty field"}
	}

	if c.PlayerID == uuid.Nil {
		return &common.CommandFieldError{Field: "player_id", Details: "empty field"}
	}

	return nil
}

type InsertCard struct {
	DeskID uuid.UUID `json:"desk_id"`
	CardID uuid.UUID `json:"card_id"`
	Index  int       `json:"index"`
}

func (c *InsertCard) AggregateType() common.AggregateType { return AggregateType }

func (c *InsertCard) AggregateID() string { return c.DeskID.String() }

func (c *InsertCard) CommandType() common.CommandType { return InsertCardCommand }

func (c *InsertCard) Validate() error {
	if c.DeskID == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	if c.CardID == uuid.Nil {
		return &common.CommandFieldError{Field: "card_id", Details: "empty field"}
	}

	if c.Index < 0 {
		return &common.CommandFieldError{Field: "index", Details: "invalid index"}
	}

	return nil
}
