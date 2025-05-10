package hand

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateHand, *CreateHand]()
	eventing.RegisterCommand[ShuffleHand, *ShuffleHand]()
	eventing.RegisterCommand[AddCards, *AddCards]()
	eventing.RegisterCommand[RemoveCards, *RemoveCards]()
}

const (
	CreateHandCommand  = common.CommandType("hand:create")
	ShuffleHandCommand = common.CommandType("hand:shuffle")
	AddCardsCommand    = common.CommandType("hand:cards:add")
	RemoveCardsCommand = common.CommandType("hand:cards:remove")
	StealCardCommand   = common.CommandType("hand:cards:steal")
)

var AllCommands = []common.CommandType{
	CreateHandCommand,
	ShuffleHandCommand,
	AddCardsCommand,
	RemoveCardsCommand,
	StealCardCommand,
}

var _ = eventing.Command(&CreateHand{})
var _ = eventing.Command(&ShuffleHand{})
var _ = eventing.Command(&AddCards{})
var _ = eventing.Command(&RemoveCards{})
var _ = eventing.Command(&StealCard{})

type CreateHand struct {
	HandID uuid.UUID   `json:"hand_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (c *CreateHand) AggregateType() common.AggregateType { return AggregateType }

func (c *CreateHand) AggregateID() string { return c.HandID.String() }

func (c *CreateHand) CommandType() common.CommandType { return CreateHandCommand }

func (c *CreateHand) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if len(c.Cards) == 0 {
		return &common.CommandFieldError{Field: "cards", Details: "empty list"}
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

type AddCards struct {
	HandID uuid.UUID   `json:"hand_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (c *AddCards) AggregateType() common.AggregateType { return AggregateType }

func (c *AddCards) AggregateID() string { return c.HandID.String() }

func (c *AddCards) CommandType() common.CommandType { return AddCardsCommand }

func (c *AddCards) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if len(c.Cards) == 0 {
		return &common.CommandFieldError{Field: "cards", Details: "empty list"}
	}

	return nil
}

type RemoveCards struct {
	HandID uuid.UUID   `json:"hand_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (c *RemoveCards) AggregateType() common.AggregateType { return AggregateType }

func (c *RemoveCards) AggregateID() string { return c.HandID.String() }

func (c *RemoveCards) CommandType() common.CommandType { return RemoveCardsCommand }

func (c *RemoveCards) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if len(c.Cards) == 0 {
		return &common.CommandFieldError{Field: "cards", Details: "empty list"}
	}

	return nil
}

type StealCard struct {
	HandID   uuid.UUID `json:"hand_id"`
	CardID   uuid.UUID `json:"card_id"`
	ToHandID uuid.UUID `json:"to_hand_id"`
}

func (c *StealCard) AggregateType() common.AggregateType { return AggregateType }

func (c *StealCard) AggregateID() string { return c.HandID.String() }

func (c *StealCard) CommandType() common.CommandType { return StealCardCommand }

func (c *StealCard) Validate() error {
	if c.HandID == uuid.Nil {
		return &common.CommandFieldError{Field: "hand_id", Details: "empty field"}
	}

	if c.ToHandID == uuid.Nil {
		return &common.CommandFieldError{Field: "to_hand_id", Details: "empty field"}
	}

	if c.CardID == uuid.Nil {
		return &common.CommandFieldError{Field: "card_id", Details: "empty field"}
	}

	return nil
}
