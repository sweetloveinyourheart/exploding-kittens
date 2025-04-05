package hand

import (
	"context"

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
