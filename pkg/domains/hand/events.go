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
	eventing.RegisterEventData[HandShuffled](EventTypeHandShuffled, args...)
	eventing.RegisterEventData[CardsAdded](EventTypeCardsAdded, args...)
	eventing.RegisterEventData[CardsRemoved](EventTypeCardsRemoved, args...)
	eventing.RegisterEventData[CardStolen](EventTypeCardStolen, args...)
}

// EventTypeHandCreated is the event type for when a hand is created
var EventTypeHandCreated = (&HandCreated{}).EventType()

// EventTypeHandShuffled is the event type for when a hand is shuffled
var EventTypeHandShuffled = (&HandShuffled{}).EventType()

// EventTypeCardsAdded is the event type for when cards are added to a hand
var EventTypeCardsAdded = (&CardsAdded{}).EventType()

// EventTypeCardsRemoved is the event type for when cards are removed from a hand
var EventTypeCardsRemoved = (&CardsRemoved{}).EventType()

// EventTypeCardStolen is the event type for when a card is stolen from a hand
var EventTypeCardStolen = (&CardStolen{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeHandCreated,
	EventTypeHandShuffled,
	EventTypeCardsAdded,
	EventTypeCardsRemoved,
	EventTypeCardStolen,
}

type HandCreated struct {
	HandID uuid.UUID   `json:"hand_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (p *HandCreated) EventType() common.EventType { return "HAND_CREATED" }

func (p *HandCreated) GetHandID() uuid.UUID { return p.HandID }

func (p *HandCreated) GetCards() []uuid.UUID { return p.Cards }

type HandShuffled struct {
	HandID uuid.UUID `json:"hand_id"`
}

func (p *HandShuffled) EventType() common.EventType { return "HAND_SHUFFLED" }

func (p *HandShuffled) GetHandID() uuid.UUID { return p.HandID }

type CardsAdded struct {
	HandID uuid.UUID   `json:"hand_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (p *CardsAdded) EventType() common.EventType { return "CARDS_ADDED" }

func (p *CardsAdded) GetHandID() uuid.UUID { return p.HandID }

func (p *CardsAdded) GetCards() []uuid.UUID { return p.Cards }

type CardsRemoved struct {
	HandID uuid.UUID   `json:"hand_id"`
	Cards  []uuid.UUID `json:"cards"`
}

func (p *CardsRemoved) EventType() common.EventType { return "CARDS_REMOVED" }

func (p *CardsRemoved) GetHandID() uuid.UUID { return p.HandID }

func (p *CardsRemoved) GetCards() []uuid.UUID { return p.Cards }

type CardStolen struct {
	HandID   uuid.UUID `json:"hand_id"`
	ToHandID uuid.UUID `json:"to_hand_id"`
	CardID   uuid.UUID `json:"card_id"`
}

func (p *CardStolen) EventType() common.EventType { return "CARD_STOLEN" }

func (p *CardStolen) GetHandID() uuid.UUID { return p.HandID }

func (p *CardStolen) GetToHandID() uuid.UUID { return p.ToHandID }

func (p *CardStolen) GetCardID() uuid.UUID { return p.CardID }
