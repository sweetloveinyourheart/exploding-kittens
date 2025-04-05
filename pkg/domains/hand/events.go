package hand

import (
	"github.com/gofrs/uuid"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

// registerEvents registers the event types for the hand domain
func registerEvents(subjFunc eventing.SubjectFunc, subjRootFunc eventing.SubjectFunc, subjTokenPos int, tokensFunc eventing.TokensFunc) {
	args := make([]eventing.EventRegistrationOption, 0)
	args = append(args, eventing.WithRegisterSubjectRootFunc(subjRootFunc))
	args = append(args, eventing.WithRegisterSubjectFunc(subjFunc))
	args = append(args, eventing.WithRegisterTokensFunc(tokensFunc))
	args = append(args, eventing.WithRegisterSubjectTokenPosition(subjTokenPos))

	eventing.RegisterEventData[HandCreated](EventTypeHandCreated, args...)
}

// EventTypeHandCreated is the event type for when a hand is created
var EventTypeHandCreated = (&HandCreated{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeHandCreated,
}

type HandCreated struct {
	HandID uuid.UUID   `json:"hand_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (p *HandCreated) EventType() common.EventType { return "HAND_CREATED" }

func (p *HandCreated) GetHandID() uuid.UUID { return p.HandID }

func (p *HandCreated) GetCards() []uuid.UUID { return p.Cards }
