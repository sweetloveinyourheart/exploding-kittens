package desk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	nats2 "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
	dataProviderGrpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
)

var (
	ProcessingTimeout               = 10 * time.Second
	BatchSize                       = 1024
	DeskStateHandlerType            = "desk-state-processing"
	DeskStateConsumerAndDurableName = "desk-state-processing-consumer-and-durable"
)

type DeskStateProcessor struct {
	ctx context.Context
	*desk.DeskProjector

	dataProvider dataProviderGrpc.DataProviderClient
	queue        chan lo.Tuple3[context.Context, common.Event, jetstream.Msg]
	bus          *nats.Conn

	deskCardIDs map[uuid.UUID][]uuid.UUID
}

func NewDeskStateProcessor(ctx context.Context) (*DeskStateProcessor, error) {
	dsp := &DeskStateProcessor{
		ctx:          ctx,
		dataProvider: do.MustInvoke[dataProviderGrpc.DataProviderClient](nil),
		queue:        make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),
		bus:          do.MustInvokeNamed[*nats.Conn](nil, fmt.Sprintf("%s-conn", constants.Bus)),

		deskCardIDs: make(map[uuid.UUID][]uuid.UUID),
	}

	dsp.DeskProjector = desk.NewDeskProjection(dsp)

	deskMatcher := eventing.NewMatchEventSubject(desk.SubjectFactory, desk.AggregateType,
		desk.EventTypeDeskCreated,
		desk.EventTypeDeskShuffled,
		desk.EventTypeCardsPeeked,
		desk.EventTypeCardsDrawn,
	)

	deskSubject := nats2.CreateConsumerSubject(constants.DeskStream, deskMatcher)

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

	deskConsumer, err := js.CreateOrUpdateConsumer(ctx, constants.DeskStream, jetstream.ConsumerConfig{
		Name:        DeskStateConsumerAndDurableName,
		Durable:     DeskStateConsumerAndDurableName,
		Description: "Responsible for reading desk events related to desk state",
		FilterSubjects: []string{
			deskSubject,
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
			messages, err := deskConsumer.Fetch(BatchSize, jetstream.FetchMaxWait(1*time.Second))
			if err != nil {
				if errors.Is(err, jetstream.ErrNoMessages) ||
					errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch desk events", zap.Error(err))
			}

			if messages.Error() != nil {
				if errors.Is(messages.Error(), jetstream.ErrNoMessages) ||
					errors.Is(messages.Error(), context.Canceled) ||
					errors.Is(messages.Error(), context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch desk events, messages error", zap.Error(messages.Error()))
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

				if !deskMatcher.Match(event) {
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

	log.Global().InfoContext(ctx, "initialized desk state processing")

	return dsp, nil
}

func (w *DeskStateProcessor) HandleDeskCreated(ctx context.Context, event common.Event, data *desk.DeskCreated) error {
	w.deskCardIDs[data.GetDeskID()] = data.GetCardIDs()

	log.Global().InfoContext(ctx, "desk created", zap.String("desk_id", data.GetDeskID().String()))

	// Emit desk state update event
	err := w.emitDeskStateUpdateEvent(data.GetDeskID())
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to emit desk state update event", zap.Error(err))
		return err
	}

	return nil
}

func (w *DeskStateProcessor) HandleDeskShuffled(ctx context.Context, event common.Event, data *desk.DeskShuffled) error {
	log.Global().InfoContext(ctx, "desk shuffled", zap.String("desk_id", data.GetDeskID().String()))

	// Emit desk state update event
	err := w.emitDeskStateUpdateEvent(data.GetDeskID())
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to emit desk state update event", zap.Error(err))
		return err
	}

	return nil
}

func (w *DeskStateProcessor) HandleCardsPeeked(ctx context.Context, event common.Event, data *desk.CardsPeeked) error {
	log.Global().InfoContext(ctx, "cards peeked", zap.String("desk_id", data.GetDeskID().String()), zap.Int("count", data.GetCount()))

	// Emit desk state update event
	err := w.emitDeskStateUpdateEvent(data.GetDeskID())
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to emit desk state update event", zap.Error(err))
		return err
	}

	return nil
}

func (w *DeskStateProcessor) HandleCardsDrawn(ctx context.Context, event common.Event, data *desk.CardsDrawn) error {
	cardIDs := w.deskCardIDs[data.GetDeskID()]

	var drawnCardIDs []uuid.UUID
	if data.GetCount() > 0 && data.GetCount() <= len(cardIDs) {
		drawnCardIDs = cardIDs[len(cardIDs)-data.GetCount():]
	} else {
		drawnCardIDs = cardIDs[:]
	}

	cardDataRes, err := w.dataProvider.GetMapCards(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		log.Global().Error("error retrieving cards map", zap.String("game_id", data.GameID.String()))
		return errors.Errorf("error retrieving cards map: %w", err)
	}

	cardsMap := cardDataRes.Msg.GetCards()

	cardIDsToSend := make([]uuid.UUID, 0)
	isExplodingKitten := false
	for _, cardID := range drawnCardIDs {
		cardInformation, ok := cardsMap[cardID.String()]
		if !ok || cardInformation == nil {
			return errors.Errorf("cannot recognize card information")
		}

		if cardInformation.GetCode() == cards.ExplodingKitten {
			isExplodingKitten = true
		} else {
			cardIDsToSend = append(cardIDsToSend, cardID)
		}
	}

	handID := hand.NewPlayerHandID(data.GetGameID(), data.GetPlayerID())
	if err := domains.CommandBus.HandleCommand(ctx, &hand.ReceiveCards{
		HandID:  handID,
		CardIDs: cardIDsToSend,
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to receive cards", zap.Error(err))
		return err
	}

	if isExplodingKitten {
		log.Global().InfoContext(ctx, "exploding kitten drawn", zap.String("desk_id", data.GetDeskID().String()), zap.String("player_id", data.GetPlayerID().String()))

		// TODO: handle exploding kitten
	} else {
		if err := domains.CommandBus.HandleCommand(ctx, &game.FinishTurn{
			GameID:   data.GetGameID(),
			PlayerID: data.GetPlayerID(),
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to finish turn", zap.Error(err))
			return err
		}
	}

	log.Global().InfoContext(ctx, "cards drawn", zap.String("desk_id", data.GetDeskID().String()), zap.Int("count", data.GetCount()))

	if data.GetCount() > 0 && data.GetCount() <= len(cardIDs) {
		w.deskCardIDs[data.GetDeskID()] = cardIDs[:len(cardIDs)-data.GetCount()]
	} else {
		w.deskCardIDs[data.GetDeskID()] = []uuid.UUID{}
	}

	// Emit desk state update event
	err = w.emitDeskStateUpdateEvent(data.GetDeskID())
	if err != nil {
		log.Global().ErrorContext(ctx, "failed to emit desk state update event", zap.Error(err))
		return err
	}

	return nil
}

func (w *DeskStateProcessor) emitDeskStateUpdateEvent(deskID uuid.UUID) error {
	msg := &desk.Desk{
		DeskID: deskID,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = w.bus.Publish(fmt.Sprintf("%s.%s", constants.DeskStream, msg.GetDeskID().String()), msgBytes)
	if err != nil {
		return err
	}

	return nil
}
