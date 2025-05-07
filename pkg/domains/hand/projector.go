package hand

import (
	"context"
	"slices"

	"github.com/gofrs/uuid"
	"github.com/samber/lo/mutable"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

type Projector struct {
	Projector *HandProjector
}

func NewProjector() *Projector {
	p := &Projector{}
	p.Projector = NewHandProjection(p)
	return p
}

var _ AllEventsProjector = (*Projector)(nil)
var _ AfterEntityHandler = (*Projector)(nil)

func (p *Projector) ProjectorType() common.ProjectorType {
	return common.ProjectorType(AggregateType.String())
}

func (p *Projector) AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Hand) (*Hand, error) {
	entity.UpdatedAt = timeutil.NowRoundedForGranularity()
	return entity, nil
}

func (p *Projector) Project(ctx context.Context, event common.Event, entity *Hand) (*Hand, error) {
	return p.Projector.Project(ctx, event, entity)
}

func (p *Projector) HandleHandCreated(ctx context.Context, event common.Event, data *HandCreated, entity *Hand) (*Hand, error) {
	entity.HandID = data.GetHandID()
	entity.Cards = data.GetCards()

	return entity, nil
}

func (p *Projector) HandleHandShuffled(ctx context.Context, event common.Event, data *HandShuffled, entity *Hand) (*Hand, error) {
	cards := entity.GetCards()
	mutable.Shuffle(cards)

	entity.HandID = data.GetHandID()
	entity.Cards = cards
	entity.ShuffledAt = timeutil.NowRoundedForGranularity()

	return entity, nil
}

func (p *Projector) HandleCardsAdded(ctx context.Context, event common.Event, data *CardsAdded, entity *Hand) (*Hand, error) {
	cards := entity.GetCards()
	cards = append(cards, data.GetCards()...)

	entity.Cards = cards

	return entity, nil
}

func (p *Projector) HandleCardsRemoved(ctx context.Context, event common.Event, data *CardsRemoved, entity *Hand) (*Hand, error) {
	cards := entity.GetCards()
	removedCards := data.GetCards()

	for _, cardID := range removedCards {
		cards = slices.DeleteFunc(cards, func(cID uuid.UUID) bool {
			return cID == cardID
		})
	}

	entity.Cards = cards

	return entity, nil
}

func (p *Projector) HandleCardStolen(ctx context.Context, event common.Event, data *CardStolen, entity *Hand) (*Hand, error) {
	cards := entity.GetCards()
	cards = slices.DeleteFunc(cards, func(cID uuid.UUID) bool {
		return cID == data.GetCardID()
	})

	entity.Cards = cards

	return entity, nil
}
