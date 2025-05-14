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
	eventing.RegisterEventData[CardsReceived](EventTypeCardsReceived, args...)
	eventing.RegisterEventData[CardsGiven](EventTypeCardsGiven, args...)
	eventing.RegisterEventData[CardsPlayed](EventTypeCardsPlayed, args...)
}

// EventTypeHandCreated is the event type for when a hand is created
var EventTypeHandCreated = (&HandCreated{}).EventType()

// EventTypeHandShuffled is the event type for when a hand is shuffled
var EventTypeHandShuffled = (&HandShuffled{}).EventType()

// EventTypeCardsReceived is the event type for when cards are received by a hand
var EventTypeCardsReceived = (&CardsReceived{}).EventType()

// EventTypeCardsGiven is the event type for when cards are given from a hand
var EventTypeCardsGiven = (&CardsGiven{}).EventType()

// EventTypeCardsPlayed is the event type for when cards are played from a hand
var EventTypeCardsPlayed = (&CardsPlayed{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeHandCreated,
	EventTypeHandShuffled,
	EventTypeCardsReceived,
	EventTypeCardsGiven,
	EventTypeCardsPlayed,
}

type HandCreated struct {
	HandID  uuid.UUID   `json:"hand_id"`
	CardIDs []uuid.UUID `json:"cards"`
}

func (p *HandCreated) EventType() common.EventType { return "HAND_CREATED" }

func (p *HandCreated) GetHandID() uuid.UUID { return p.HandID }

func (p *HandCreated) GetCardIDs() []uuid.UUID { return p.CardIDs }

type HandShuffled struct {
	HandID uuid.UUID `json:"hand_id"`
}

func (p *HandShuffled) EventType() common.EventType { return "HAND_SHUFFLED" }

func (p *HandShuffled) GetHandID() uuid.UUID { return p.HandID }

type CardsReceived struct {
	HandID  uuid.UUID   `json:"hand_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (p *CardsReceived) EventType() common.EventType { return "HAND_CARDS_RECEIVED" }

func (p *CardsReceived) GetHandID() uuid.UUID { return p.HandID }

func (p *CardsReceived) GetCardIDs() []uuid.UUID { return p.CardIDs }

type CardsPlayed struct {
	HandID  uuid.UUID   `json:"hand_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (p *CardsPlayed) EventType() common.EventType { return "HAND_CARDS_PLAYED" }

func (p *CardsPlayed) GetHandID() uuid.UUID { return p.HandID }

func (p *CardsPlayed) GetCardIDs() []uuid.UUID { return p.CardIDs }

type CardsGiven struct {
	HandID   uuid.UUID   `json:"hand_id"`
	CardIDs  []uuid.UUID `json:"card_ids"`
	ToHandID uuid.UUID   `json:"to_hand_id"`
}

func (p *CardsGiven) EventType() common.EventType { return "HAND_CARDS_GIVEN" }

func (p *CardsGiven) GetHandID() uuid.UUID { return p.HandID }

func (p *CardsGiven) GetToHandID() uuid.UUID { return p.ToHandID }

func (p *CardsGiven) GetCardIDs() []uuid.UUID { return p.CardIDs }
