package desk

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/samber/lo"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/aggregate"
	aggregateCommandHandler "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/aggregate"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/command_handler/bus"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	natsEventBus "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_handler/projector"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_store/natsjs"
	consumeroptions "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/consumer_options"
	contexthook "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/context_hook"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/ephemeral"
	natsRepo "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/repo/natsjs_eventsourced"
)

func CreateNATSRepoDesk(ctx context.Context, appID string, mw ...eventing.EventHandlerMiddleware) (eventing.ReadRepo[Desk, *Desk], error) {
	connPool, err := do.InvokeNamed[*pool.ConnPool](nil, string(constants.ConnectionPool))
	if err != nil {
		return nil, err
	}
	rng := lo.RandomString(8, lo.LettersCharset)
	appID = fmt.Sprintf("%s_desk_repo_%s", appID, rng)

	neb, err := natsEventBus.NewEventBus(connPool, appID, natsEventBus.WithStreamName(constants.DeskStream))
	if err != nil {
		return nil, err
	}

	natsEventBus.BusErrors(ctx, neb)

	entityProjector := NewProjector()

	// Create repo for projector
	var domainRepo eventing.ReadWriteRepo[Desk, *Desk]
	domainRepo, err = natsRepo.NewRepo[Desk, *Desk](ctx,
		constants.DeskStream,
		SubjectFactory,
		entityProjector,
		natsRepo.WithEventBus(neb))
	if err != nil {
		return nil, err
	}

	// Create projector
	var domainProjector eventing.EventHandler
	domainProjector = projector.NewEventHandler[Desk, *Desk](entityProjector, domainRepo)

	domainProjector = ephemeral.NewMiddleware()(domainProjector)
	domainProjector = consumeroptions.NewDeliveryPolicyMiddleware(jetstream.DeliverNewPolicy, 0)(domainProjector)
	for _, m := range mw {
		domainProjector = m(domainProjector)
	}

	err = neb.AddHandler(context.Background(), eventing.NewMatchEventSubject(SubjectFactory, AggregateType), domainProjector)
	if err != nil {
		return nil, err
	}

	return domainRepo, nil
}

func AddNATSDeskCommandHandlers(ctx context.Context, appID string, commandBus *bus.CommandHandler, mw ...eventing.CommandHandlerMiddleware) error {
	connPool, err := do.InvokeNamed[*pool.ConnPool](nil, string(constants.ConnectionPool))
	if err != nil {
		return err
	}

	natsBus, err := natsEventBus.NewEventBus(connPool, fmt.Sprintf("%s-desk-command-read", appID), natsEventBus.WithStreamName(constants.DeskStream), natsEventBus.WithCodec(customCodec{}))
	if err != nil {
		return err
	}

	// Create in memory store for aggregate
	hookHandler := contexthook.NewMiddleware()(natsBus)
	natsEventStore, err := natsjs.NewEventStore(ctx, constants.DeskStream, SubjectFactory, natsjs.WithEventHandler(hookHandler), natsjs.WithEventBus(natsBus))
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = natsBus.Close()
				_ = natsEventStore.Close()
				return
			case err, ok := <-natsBus.Errors():
				natsEventBus.HandleError(ctx, err)
				if !ok {
					_ = natsEventStore.Close()
					return
				}
			}
		}
	}()

	wrappedEventStore := consumeroptions.NewMiddleware(jetstream.ConsumerConfig{
		DeliverPolicy:     jetstream.DeliverNewPolicy,
		AckPolicy:         jetstream.AckExplicitPolicy,
		AckWait:           10 * time.Second,
		MaxDeliver:        10,
		InactiveThreshold: 60 * time.Minute,
	})(natsEventStore)
	err = natsBus.AddHandler(context.Background(), eventing.NewMatchEventSubject(SubjectFactory, AggregateType), wrappedEventStore)
	if err != nil {
		return err
	}

	aggregateStore, err := aggregate.NewAggregateStore(natsEventStore, aggregate.WithSequencedStore())
	if err != nil {
		return err
	}

	var domainCommandHandler eventing.CommandHandler
	domainCommandHandler, err = aggregateCommandHandler.NewCommandHandler(AggregateType, aggregateStore, aggregateCommandHandler.WithDeadline(5*time.Second))
	if err != nil {
		return err
	}
	for _, m := range mw {
		domainCommandHandler = m(domainCommandHandler)
	}

	domainCommands := AllCommands
	for _, command := range domainCommands {
		err = commandBus.SetHandler(domainCommandHandler, command)
		if err != nil {
			return err
		}
	}

	return nil
}

type customCodec struct {
}

func (c customCodec) MarshalEvent(ctx context.Context, event common.Event) ([]byte, error) {
	return natsEventBus.DefaultEventCodec.MarshalEvent(ctx, event)
}

func (c customCodec) UnmarshalEvent(
	ctx context.Context,
	bytes []byte,
	option ...eventing.EventOption,
) (common.Event, context.Context, error) {
	event, ctx, err := natsEventBus.DefaultEventCodec.UnmarshalEvent(ctx, bytes, option...)
	if err != nil {
		return nil, ctx, err
	}
	event = eventing.ReplaceSubject(event, SubjectFactory)
	return event, ctx, nil
}

var _ eventing.EventCodec = customCodec{}
