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
	entity.CardIDs = data.GetCardIDs()

	return entity, nil
}

func (p *Projector) HandleCardsPlayed(ctx context.Context, event common.Event, data *CardsPlayed, entity *Hand) (*Hand, error) {
	cardIDs := entity.GetCardIDs()

	for _, cardID := range data.GetCardIDs() {
		index := slices.IndexFunc(cardIDs, func(cID uuid.UUID) bool {
			return cID == cardID
		})
		if index != -1 {
			cardIDs = slices.Delete(cardIDs, index, index+1)
		}
	}

	entity.CardIDs = cardIDs

	return entity, nil
}

func (p *Projector) HandleHandShuffled(ctx context.Context, event common.Event, data *HandShuffled, entity *Hand) (*Hand, error) {
	cardIDs := entity.GetCardIDs()
	mutable.Shuffle(cardIDs)

	entity.HandID = data.GetHandID()
	entity.CardIDs = cardIDs
	entity.ShuffledAt = timeutil.NowRoundedForGranularity()

	return entity, nil
}

func (p *Projector) HandleCardsReceived(ctx context.Context, event common.Event, data *CardsReceived, entity *Hand) (*Hand, error) {
	cardIDs := entity.GetCardIDs()
	cardIDs = append(cardIDs, data.GetCardIDs()...)

	entity.HandID = data.GetHandID()
	entity.CardIDs = cardIDs

	return entity, nil
}

func (p *Projector) HandleCardsGiven(ctx context.Context, event common.Event, data *CardsGiven, entity *Hand) (*Hand, error) {
	cardIDs := entity.GetCardIDs()

	for _, cardID := range data.GetCardIDs() {
		index := slices.IndexFunc(cardIDs, func(cID uuid.UUID) bool {
			return cID == cardID
		})
		if index != -1 {
			cardIDs = slices.Delete(cardIDs, index, index+1)
		}
	}

	entity.CardIDs = cardIDs

	return entity, nil
}
