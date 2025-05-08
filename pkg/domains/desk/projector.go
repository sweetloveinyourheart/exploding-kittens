package desk

import (
	"context"

	"github.com/samber/lo/mutable"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

type Projector struct {
	Projector *DeskProjector
}

func NewProjector() *Projector {
	p := &Projector{}
	p.Projector = NewDeskProjection(p)
	return p
}

var _ AllEventsProjector = (*Projector)(nil)
var _ AfterEntityHandler = (*Projector)(nil)

func (p *Projector) ProjectorType() common.ProjectorType {
	return common.ProjectorType(AggregateType.String())
}

func (p *Projector) AfterHandleEvent(ctx context.Context, event common.Event, data any, entity *Desk) (*Desk, error) {
	entity.UpdatedAt = timeutil.NowRoundedForGranularity()
	return entity, nil
}

func (p *Projector) Project(ctx context.Context, event common.Event, entity *Desk) (*Desk, error) {
	return p.Projector.Project(ctx, event, entity)
}

func (p *Projector) HandleDeskCreated(ctx context.Context, event common.Event, data *DeskCreated, entity *Desk) (*Desk, error) {
	entity.DeskID = data.GetDeskID()
	entity.CardIDs = data.GetCardIDs()

	return entity, nil
}

func (p *Projector) HandleDeskShuffled(ctx context.Context, event common.Event, data *DeskShuffled, entity *Desk) (*Desk, error) {
	cardIDs := entity.GetCardIDs()
	mutable.Shuffle(cardIDs)

	entity.DeskID = data.GetDeskID()
	entity.CardIDs = cardIDs
	entity.ShuffledAt = timeutil.NowRoundedForGranularity()

	return entity, nil
}
