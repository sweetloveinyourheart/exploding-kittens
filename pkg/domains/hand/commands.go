package hand

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateHand, *CreateHand]()
	eventing.RegisterCommand[ShuffleHand, *ShuffleHand]()
	eventing.RegisterCommand[PlayCards, *PlayCards]()
	eventing.RegisterCommand[ReceiveCards, *ReceiveCards]()
	eventing.RegisterCommand[GiveCards, *GiveCards]()
}

const (
	CreateHandCommand   = common.CommandType("hand:create")
	ShuffleHandCommand  = common.CommandType("hand:shuffle")
	PlayCardsCommand    = common.CommandType("hand:cards:play")
	ReceiveCardsCommand = common.CommandType("hand:cards:receive")
	GiveCardsCommand    = common.CommandType("hand:cards:give")
)

var AllCommands = []common.CommandType{
	CreateHandCommand,
	ShuffleHandCommand,
	ReceiveCardsCommand,
	PlayCardsCommand,
	GiveCardsCommand,
}

var _ = eventing.Command(&CreateHand{})
var _ = eventing.Command(&ShuffleHand{})
var _ = eventing.Command(&ReceiveCards{})
var _ = eventing.Command(&PlayCards{})
var _ = eventing.Command(&GiveCards{})

type CreateHand struct {
	HandID  uuid.UUID   `json:"hand_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (c *CreateHand) AggregateType() common.AggregateType { return AggregateType }

func (c *CreateHand) AggregateID() string { return c.HandID.String() }

func (c *CreateHand) CommandType() common.CommandType { return CreateHandCommand }

func (c *CreateHand) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if len(c.CardIDs) == 0 {
		return &common.CommandFieldError{Field: "cards", Details: "empty list"}
	}

	return nil
}

type PlayCards struct {
	HandID  uuid.UUID   `json:"hand_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (c *PlayCards) AggregateType() common.AggregateType { return AggregateType }

func (c *PlayCards) AggregateID() string { return c.HandID.String() }

func (c *PlayCards) CommandType() common.CommandType { return PlayCardsCommand }

func (c *PlayCards) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if len(c.CardIDs) == 0 {
		return &common.CommandFieldError{Field: "card_ids", Details: "empty list"}
	}

	return nil
}

type ShuffleHand struct {
	HandID uuid.UUID `json:"hand_id"`
}

func (c *ShuffleHand) AggregateType() common.AggregateType { return AggregateType }

func (c *ShuffleHand) AggregateID() string { return c.HandID.String() }

func (c *ShuffleHand) CommandType() common.CommandType { return ShuffleHandCommand }

func (c *ShuffleHand) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	return nil
}

type ReceiveCards struct {
	HandID  uuid.UUID   `json:"hand_id"`
	CardIDs []uuid.UUID `json:"cards"`
}

func (c *ReceiveCards) AggregateType() common.AggregateType { return AggregateType }

func (c *ReceiveCards) AggregateID() string { return c.HandID.String() }

func (c *ReceiveCards) CommandType() common.CommandType { return ReceiveCardsCommand }

func (c *ReceiveCards) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if len(c.CardIDs) == 0 {
		return &common.CommandFieldError{Field: "cards", Details: "empty list"}
	}

	return nil
}

type GiveCards struct {
	HandID   uuid.UUID   `json:"hand_id"`
	CardIDs  []uuid.UUID `json:"card_ids"`
	ToHandID uuid.UUID   `json:"to_hand_id"`
}

func (c *GiveCards) AggregateType() common.AggregateType { return AggregateType }

func (c *GiveCards) AggregateID() string { return c.HandID.String() }

func (c *GiveCards) CommandType() common.CommandType { return GiveCardsCommand }

func (c *GiveCards) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if c.ToHandID == uuid.Nil {
		return &common.CommandFieldError{Field: "to_hand_id", Details: "empty field"}
	}

	if len(c.CardIDs) == 0 {
		return &common.CommandFieldError{Field: "card_ids", Details: "empty list"}
	}

	return nil
}
