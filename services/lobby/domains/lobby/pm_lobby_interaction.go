package lobby

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/samber/lo"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"

	"go.uber.org/zap"

	nats2 "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

var (
	ProcessingTimeout                      = 10 * time.Second
	BatchSize                              = 1024
	LobbyInteractionHandlerType            = "lobby-interaction-processing"
	LobbyInteractionConsumerandDurableName = "lobby-cancelled-void-choices-consumer"
)

type LobbyInteractionProcessor struct {
	ctx context.Context
	*lobby.LobbyProjector

	queue chan lo.Tuple3[context.Context, common.Event, jetstream.Msg]
}

func NewLobbyInteractionProcessor(ctx context.Context) (*LobbyInteractionProcessor, error) {
	lip := &LobbyInteractionProcessor{
		ctx:   ctx,
		queue: make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),
	}

	lip.LobbyProjector = lobby.NewLobbyProjection(lip)

	lobbyMatcher := eventing.NewMatchEventSubject(lobby.SubjectFactory, lobby.AggregateType,
		lobby.EventTypeLobbyLeft,
	)

	lobbySubject := nats2.CreateConsumerSubject(constants.LobbyStream, lobbyMatcher)

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

	lobbyConsumer, err := js.CreateOrUpdateConsumer(ctx, constants.LobbyStream, jetstream.ConsumerConfig{
		Name:        LobbyInteractionConsumerandDurableName,
		Durable:     LobbyInteractionConsumerandDurableName,
		Description: "Responsible for reading lobby events related to lobby leaves",
		FilterSubjects: []string{
			lobbySubject,
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
			messages, err := lobbyConsumer.Fetch(BatchSize, jetstream.FetchMaxWait(1*time.Second))
			if err != nil {
				if errors.Is(err, jetstream.ErrNoMessages) ||
					errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch lobby events", zap.Error(err))
			}

			if messages.Error() != nil {
				if errors.Is(messages.Error(), jetstream.ErrNoMessages) ||
					errors.Is(messages.Error(), context.Canceled) ||
					errors.Is(messages.Error(), context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch lobby events, messages error", zap.Error(messages.Error()))
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

				if !lobbyMatcher.Match(event) {
					if err := msg.Ack(); err != nil {
						log.Global().ErrorContext(ctx, "failed to ack message", zap.Error(err))
					}
					continue
				}

				lip.queue <- lo.T3(eventCtx, event, msg)
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
			case tuple := <-lip.queue:
				if err := lip.HandleEvent(tuple.A, tuple.B); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to handle event", zap.Error(err))
				}
				if err := tuple.C.Ack(); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to ack message", zap.Error(err))
				}
			}
		}
	}()

	log.Global().InfoContext(ctx, "initialized lobby interaction processing")

	return lip, nil
}

func (w *LobbyInteractionProcessor) HandlerType() common.EventHandlerType {
	return common.EventHandlerType(LobbyInteractionHandlerType)
}

func (w *LobbyInteractionProcessor) HandleEvent(ctx context.Context, event common.Event) (err error) {
	if event == nil {
		return errors.WithStack(errors.New("event is nil"))
	}

	if event.AggregateType() == lobby.AggregateType {
		return w.LobbyProjector.HandleEvent(ctx, event)
	}

	return errors.WithStack(fmt.Errorf("unknown aggregate type %s", event.AggregateType()))
}

func (w *LobbyInteractionProcessor) HandleLobbyLeft(ctx context.Context, event common.Event, data *lobby.LobbyLeft) error {
	lobby, err := lobby.GetLobbyByID(ctx, data.GetLobbyID())
	if err != nil {
		log.Global().Error("error getting lobby with aggregate", zap.Error(err))
		return err
	}

	log.Global().InfoContext(ctx, "lobby details",
		zap.String("id", lobby.GetLobbyID().String()),
		zap.String("code", lobby.GetLobbyCode()),
		zap.String("name", lobby.GetLobbyName()),
		zap.String("owner", lobby.GetHostUserID().String()),
		zap.Int("player_count", len(lobby.GetParticipants())),
		zap.Time("created_at", lobby.CreatedAt),
		zap.Time("updated_at", lobby.UpdatedAt),
	)

	return nil
}
