package hand

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	nats2 "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
	"go.uber.org/zap"
)

var (
	ProcessingTimeout               = 10 * time.Second
	BatchSize                       = 1024
	HandStateHandlerType            = "hand-state-processing"
	HandStateConsumerAndDurableName = "hand-state-processing-consumer-and-durable"
)

type HandStateProcessor struct {
	ctx context.Context
	*hand.HandProjector

	queue chan lo.Tuple3[context.Context, common.Event, jetstream.Msg]
	bus   *nats.Conn
}

func NewHandStateProcessor(ctx context.Context) (*HandStateProcessor, error) {
	dsp := &HandStateProcessor{
		ctx:   ctx,
		queue: make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),
		bus:   do.MustInvokeNamed[*nats.Conn](nil, fmt.Sprintf("%s-conn", constants.Bus)),
	}

	dsp.HandProjector = hand.NewHandProjection(dsp)

	handMatcher := eventing.NewMatchEventSubject(hand.SubjectFactory, hand.AggregateType,
		hand.EventTypeCardStolen,
	)

	handSubject := nats2.CreateConsumerSubject(constants.HandStream, handMatcher)

	connPool, err := do.InvokeNamed[*pool.ConnPool](nil, string(constants.ConnectionPool))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conn, err := connPool.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	replicas := 1
	if reps := config.Instance().GetInt(config.NatsConsumerReplicas); reps > 0 {
		replicas = reps
	}
	memory := false
	if ss := config.Instance().GetString(config.NatsConsumerStorage); strings.EqualFold(ss, "memory") {
		memory = true
	}

	handConsumer, err := js.CreateOrUpdateConsumer(ctx, constants.HandStream, jetstream.ConsumerConfig{
		Name:        HandStateConsumerAndDurableName,
		Durable:     HandStateConsumerAndDurableName,
		Description: "Responsible for reading hand events related to hand state",
		FilterSubjects: []string{
			handSubject,
		},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		ReplayPolicy:      jetstream.ReplayInstantPolicy,
		AckWait:           ProcessingTimeout,
		AckPolicy:         jetstream.AckExplicitPolicy,
		InactiveThreshold: time.Hour * 24 * 7,
		MaxDeliver:        10,
		MaxWaiting:        1,
		Replicas:          replicas,
		MemoryStorage:     memory,
		MaxAckPending:     BatchSize,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			messages, err := handConsumer.Fetch(BatchSize, jetstream.FetchMaxWait(1*time.Second))
			if err != nil {
				if errors.Is(err, jetstream.ErrNoMessages) ||
					errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch hand events", zap.Error(err))
			}

			if messages.Error() != nil {
				if errors.Is(messages.Error(), jetstream.ErrNoMessages) ||
					errors.Is(messages.Error(), context.Canceled) ||
					errors.Is(messages.Error(), context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch hand events, messages error", zap.Error(messages.Error()))
			}

			for msg := range messages.Messages() {
				ctx := context.WithoutCancel(ctx)

				event, eventCtx, err := nats2.JSMsgToEvent(ctx, msg)
				if err != nil {
					log.Global().ErrorContext(ctx, "failed to convert message to event", zap.Error(err))
					if err := msg.Ack(); err != nil {
						log.Global().ErrorContext(ctx, "failed to ack message", zap.Error(err))
					}
					continue
				}

				if !handMatcher.Match(event) {
					if err := msg.Ack(); err != nil {
						log.Global().ErrorContext(ctx, "failed to ack message", zap.Error(err))
					}
					continue
				}

				dsp.queue <- lo.T3(eventCtx, event, msg)
			}
		}
	}()

	go func() {
		timer := timeutil.Clock.Timer(100 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				if _, err := connPool.Put(conn); err != nil {
					log.Global().WarnContext(ctx, "failed to return connection to pool", zap.Error(err))
				}
				timer.Stop()
				return
			case tuple := <-dsp.queue:
				if err := dsp.HandleEvent(tuple.A, tuple.B); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to handle event", zap.Error(err))
				}
				if err := tuple.C.Ack(); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to ack message", zap.Error(err))
				}
			}
		}
	}()

	log.Global().InfoContext(ctx, "initialized hand state processing")

	return dsp, nil
}

func (w *HandStateProcessor) HandleCardStolen(ctx context.Context, event common.Event, data *hand.CardStolen) error {
	if err := domains.CommandBus.HandleCommand(ctx, &hand.AddCards{
		HandID: data.ToHandID,
		Cards:  []uuid.UUID{data.CardID},
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to add cards", zap.Error(err))
		return err
	}

	return nil
}
