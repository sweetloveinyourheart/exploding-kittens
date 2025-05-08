package desk

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateDesk, *CreateDesk]()
	eventing.RegisterCommand[ShuffleDesk, *ShuffleDesk]()
}

const (
	CreateDeskCommand  = common.CommandType("desk:create")
	ShuffleDeskCommand = common.CommandType("desk:shuffle")
)

var AllCommands = []common.CommandType{
	CreateDeskCommand,
	ShuffleDeskCommand,
}

var _ = eventing.Command(&CreateDesk{})
var _ = eventing.Command(&ShuffleDesk{})

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
