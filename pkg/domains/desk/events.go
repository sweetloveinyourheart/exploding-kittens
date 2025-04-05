package desk

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// registerEvents registers the event types for the desk domain
func registerEvents(subjFunc eventing.SubjectFunc, subjRootFunc eventing.SubjectFunc, subjTokenPos int, tokensFunc eventing.TokensFunc) {
	args := make([]eventing.EventRegistrationOption, 0)
	args = append(args, eventing.WithRegisterSubjectRootFunc(subjRootFunc))
	args = append(args, eventing.WithRegisterSubjectFunc(subjFunc))
	args = append(args, eventing.WithRegisterTokensFunc(tokensFunc))
	args = append(args, eventing.WithRegisterSubjectTokenPosition(subjTokenPos))

	eventing.RegisterEventData[DeskCreated](EventTypeDeskCreated, args...)
}

// EventTypeDeskCreated is the event type for when a desk is created
var EventTypeDeskCreated = (&DeskCreated{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeDeskCreated,
}

type DeskCreated struct {
	DeskID uuid.UUID   `json:"desk_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (p *DeskCreated) EventType() common.EventType { return "DESK_CREATED" }

func (p *DeskCreated) GetDeskID() uuid.UUID { return p.DeskID }

func (p *DeskCreated) GetCards() []uuid.UUID { return p.Cards }
