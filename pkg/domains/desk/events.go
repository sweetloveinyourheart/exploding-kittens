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
	eventing.RegisterEventData[DeskShuffled](EventTypeDeskShuffled, args...)
	eventing.RegisterEventData[CardsDiscarded](EventTypeCardsDiscarded, args...)
	eventing.RegisterEventData[CardsPeeked](EventTypeCardsPeeked, args...)
	eventing.RegisterEventData[CardDrawn](EventTypeCardDrawn, args...)
}

// EventTypeDeskCreated is the event type for when a desk is created
var EventTypeDeskCreated = (&DeskCreated{}).EventType()

// EventTypeDeskShuffled is the event type for when a desk is shuffled
var EventTypeDeskShuffled = (&DeskShuffled{}).EventType()

// EventTypeCardsDiscarded is the event type for when cards are discarded
var EventTypeCardsDiscarded = (&CardsDiscarded{}).EventType()

// EventTypeCardsPeeked is the event type for when cards are peeked
var EventTypeCardsPeeked = (&CardsPeeked{}).EventType()

// EventTypeCardDrawn is the event type for when a card is drawn
var EventTypeCardDrawn = (&CardDrawn{}).EventType()

var AllEventTypes = []common.EventType{
	EventTypeDeskCreated,
	EventTypeDeskShuffled,
	EventTypeCardsDiscarded,
	EventTypeCardsPeeked,
	EventTypeCardDrawn,
}

type DeskCreated struct {
	DeskID  uuid.UUID   `json:"desk_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (p *DeskCreated) EventType() common.EventType { return "DESK_CREATED" }

func (p *DeskCreated) GetDeskID() uuid.UUID { return p.DeskID }

func (p *DeskCreated) GetCardIDs() []uuid.UUID { return p.CardIDs }

type DeskShuffled struct {
	DeskID uuid.UUID `json:"desk_id"`
}

func (p *DeskShuffled) EventType() common.EventType { return "DESK_SHUFFLED" }

func (p *DeskShuffled) GetDeskID() uuid.UUID { return p.DeskID }

type CardsDiscarded struct {
	DeskID  uuid.UUID   `json:"desk_id"`
	CardIDs []uuid.UUID `json:"card_ids"`
}

func (p *CardsDiscarded) EventType() common.EventType { return "DESK_CARDS_DISCARDED" }

func (p *CardsDiscarded) GetDeskID() uuid.UUID { return p.DeskID }

func (p *CardsDiscarded) GetCardIDs() []uuid.UUID { return p.CardIDs }

type CardsPeeked struct {
	DeskID uuid.UUID `json:"desk_id"`
	Count  int       `json:"count"`
}

func (p *CardsPeeked) EventType() common.EventType { return "DESK_CARDS_PEEKED" }

func (p *CardsPeeked) GetDeskID() uuid.UUID { return p.DeskID }

func (p *CardsPeeked) GetCount() int { return p.Count }

type CardDrawn struct {
	DeskID        uuid.UUID `json:"desk_id"`
	GameID        uuid.UUID `json:"game_id"`
	PlayerID      uuid.UUID `json:"player_id"`
	CanFinishTurn bool      `json:"can_finish_turn"` // Indicates if the player can finish their turn after drawing a card
}

func (p *CardDrawn) EventType() common.EventType { return "DESK_CARD_DRAWN" }

func (p *CardDrawn) GetDeskID() uuid.UUID { return p.DeskID }

func (p *CardDrawn) GetGameID() uuid.UUID { return p.GameID }

func (p *CardDrawn) GetPlayerID() uuid.UUID { return p.PlayerID }

func (p *CardDrawn) GetCanFinishTurn() bool { return p.CanFinishTurn }
