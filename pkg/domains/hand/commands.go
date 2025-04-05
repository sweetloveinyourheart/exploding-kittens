package hand

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func init() {
	eventing.RegisterCommand[CreateHand, *CreateHand]()
}

const (
	CreateHandCommand = common.CommandType("hand:create")
)

var AllCommands = []common.CommandType{
	CreateHandCommand,
}

var _ = eventing.Command(&CreateHand{})

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
