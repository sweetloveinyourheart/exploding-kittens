package hand

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	nats2 "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
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

	handCardIDs map[uuid.UUID][]uuid.UUID
}

func NewHandStateProcessor(ctx context.Context) (*HandStateProcessor, error) {
	dsp := &HandStateProcessor{
		ctx:         ctx,
		queue:       make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),
		bus:         do.MustInvokeNamed[*nats.Conn](nil, fmt.Sprintf("%s-conn", constants.Bus)),
		handCardIDs: make(map[uuid.UUID][]uuid.UUID),
	}

	dsp.HandProjector = hand.NewHandProjection(dsp)

	handMatcher := eventing.NewMatchEventSubject(hand.SubjectFactory, hand.AggregateType,
		hand.EventTypeHandCreated,
		hand.EventTypeCardsPlayed,
		hand.EventTypeCardsReceived,
		hand.EventTypeCardsGiven,
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

func (w *HandStateProcessor) HandleHandCreated(ctx context.Context, event common.Event, data *hand.HandCreated) error {
	log.Global().InfoContext(ctx, "hand created", zap.String("handID", data.GetHandID().String()))

	w.handCardIDs[data.GetHandID()] = data.GetCardIDs()

	return nil
}

func (w *HandStateProcessor) HandleCardsPlayed(ctx context.Context, event common.Event, data *hand.CardsPlayed) error {
	// Remove the played cards from the hand
	handCardIDs := slices.Clone(w.handCardIDs[data.GetHandID()])
	for _, cardID := range data.GetCardIDs() {
		index := slices.IndexFunc(handCardIDs, func(cID uuid.UUID) bool {
			return cID == cardID
		})
		if index != -1 {
			handCardIDs = slices.Delete(handCardIDs, index, index+1)
		}
	}

	w.handCardIDs[data.GetHandID()] = handCardIDs

	log.Global().InfoContext(ctx, "card played",
		zap.String("handID", data.GetHandID().String()),
		zap.Strings("cardIDs", lo.Map(data.GetCardIDs(), func(c uuid.UUID, _ int) string {
			return c.String()
		})))

	// Emit hand state update event
	err := w.emitHandStateUpdateEvent(data.GetHandID())
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to emit hand state update event", zap.Error(err))
		return err
	}

	return nil
}

func (w *HandStateProcessor) HandleCardsReceived(ctx context.Context, event common.Event, data *hand.CardsReceived) error {
	// Add the received cards to the hand
	w.handCardIDs[data.GetHandID()] = append(w.handCardIDs[data.GetHandID()], data.GetCardIDs()...)

	log.Global().InfoContext(ctx, "cards received",
		zap.String("handID", data.GetHandID().String()),
		zap.Strings("cardIDs", lo.Map(data.GetCardIDs(), func(c uuid.UUID, _ int) string {
			return c.String()
		})))

	// Emit hand state update event
	err := w.emitHandStateUpdateEvent(data.GetHandID())
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to emit hand state update event", zap.Error(err))
		return err
	}

	return nil
}

func (w *HandStateProcessor) HandleCardsGiven(ctx context.Context, event common.Event, data *hand.CardsGiven) error {
	handCardIDs := slices.Clone(w.handCardIDs[data.GetHandID()])

	cardIDsToGive := make([]uuid.UUID, 0)
	if len(data.GetCardIDs()) > 0 && len(handCardIDs) > 0 {
		cardIDsToGive = data.GetCardIDs()

		for _, cardID := range data.GetCardIDs() {
			index := slices.IndexFunc(handCardIDs, func(cID uuid.UUID) bool {
				return cID == cardID
			})

			if index != -1 {
				handCardIDs = slices.Delete(handCardIDs, index, index+1)
			}
		}
	}

	if len(data.GetCardIndexes()) > 0 && len(handCardIDs) > 0 {
		// Remove cards at the specified indexes, in descending order to avoid shifting issues
		indexes := data.GetCardIndexes()
		slices.SortFunc(indexes, func(a, b int) int { return b - a }) // sort descending

		for _, idx := range indexes {
			if idx >= 0 && idx < len(handCardIDs) {
				cardIDsToGive = append(cardIDsToGive, handCardIDs[idx])
				handCardIDs = slices.Delete(handCardIDs, idx, idx+1)
			}
		}
	}

	w.handCardIDs[data.GetHandID()] = handCardIDs

	if err := domains.CommandBus.HandleCommand(ctx, &hand.ReceiveCards{
		HandID:  data.GetToHandID(),
		CardIDs: cardIDsToGive,
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to receive cards", zap.Error(err))

		return err
	}

	log.Global().InfoContext(ctx, "cards given",
		zap.String("handID", data.GetHandID().String()),
		zap.String("toHandID", data.GetToHandID().String()),
		zap.Strings("cardIDs", lo.Map(cardIDsToGive, func(c uuid.UUID, _ int) string {
			return c.String()
		})))

	// Emit hand state update event
	err := w.emitHandStateUpdateEvent(data.GetHandID())
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to emit hand state update event", zap.Error(err))
		return err
	}

	return nil
}

func (w *HandStateProcessor) emitHandStateUpdateEvent(handID uuid.UUID) error {
	msg := &hand.Hand{
		HandID: handID,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = w.bus.Publish(fmt.Sprintf("%s.%s", constants.HandStream, msg.GetHandID().String()), msgBytes)
	if err != nil {
		return err
	}

	return nil
}
