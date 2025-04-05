package desk

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateDesk, *CreateDesk]()
}

const (
	CreateDeskCommand = common.CommandType("desk:create")
)

var AllCommands = []common.CommandType{
	CreateDeskCommand,
}

var _ = eventing.Command(&CreateDesk{})

type CreateDesk struct {
	DeskID uuid.UUID   `json:"desk_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (c *CreateDesk) AggregateType() common.AggregateType { return AggregateType }

func (c *CreateDesk) AggregateID() string { return c.DeskID.String() }

func (c *CreateDesk) CommandType() common.CommandType { return CreateDeskCommand }

func (c *CreateDesk) Validate() error {
	if c.DeskID == uuid.Nil {
		return &common.CommandFieldError{Field: "desk_id", Details: "empty field"}
	}

	if len(c.Cards) == 0 {
		return &common.CommandFieldError{Field: "cards", Details: "empty list"}
	}

	return nil
}
