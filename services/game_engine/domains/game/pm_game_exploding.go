package game

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
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
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
	GameExplodingProcessManagerType                   = "game-exploding-process-manager"
	GameExplodingProcessManagerConsumerAndDurableName = "game-exploding-process-manager-consumer-and-durable"
)

type GameExplodingProcessManager struct {
	ctx context.Context
	*game.GameProjector

	dataProvider dataProviderGrpc.DataProviderClient

	queue chan lo.Tuple3[context.Context, common.Event, jetstream.Msg]
	bus   *nats.Conn

	gameActiveExplodingCardID map[string]uuid.UUID
	gameDeskID                map[string]uuid.UUID
}

func NewGameExplodingProcessManager(ctx context.Context) (*GameExplodingProcessManager, error) {
	gip := &GameExplodingProcessManager{
		ctx:                       ctx,
		dataProvider:              do.MustInvoke[dataProviderGrpc.DataProviderClient](nil),
		queue:                     make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),
		bus:                       do.MustInvokeNamed[*nats.Conn](nil, fmt.Sprintf("%s-conn", constants.Bus)),
		gameActiveExplodingCardID: make(map[string]uuid.UUID),
		gameDeskID:                make(map[string]uuid.UUID),
	}

	gip.GameProjector = game.NewGameProjection(gip)

	gameMatcher := eventing.NewMatchEventSubject(game.SubjectFactory, game.AggregateType,
		game.EventTypeGameInitialized,
		game.EventTypeExplodingDrawn,
		game.EventTypeExplodingDefused,
		game.EventTypeKittenPlanted,
		game.EventTypePlayerEliminated,
	)

	gameSubject := nats2.CreateConsumerSubject(constants.GameStream, gameMatcher)

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

	gameConsumer, err := js.CreateOrUpdateConsumer(ctx, constants.GameStream, jetstream.ConsumerConfig{
		Name:        GameExplodingProcessManagerConsumerAndDurableName,
		Durable:     GameExplodingProcessManagerConsumerAndDurableName,
		Description: "Responsible for reading game events related to game state",
		FilterSubjects: []string{
			gameSubject,
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
			messages, err := gameConsumer.Fetch(BatchSize, jetstream.FetchMaxWait(1*time.Second))
			if err != nil {
				if errors.Is(err, jetstream.ErrNoMessages) ||
					errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch game events", zap.Error(err))
			}

			if messages.Error() != nil {
				if errors.Is(messages.Error(), jetstream.ErrNoMessages) ||
					errors.Is(messages.Error(), context.Canceled) ||
					errors.Is(messages.Error(), context.DeadlineExceeded) {
					continue
				}
				log.Global().FatalContext(ctx, "failed to fetch game events, messages error", zap.Error(messages.Error()))
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

				if !gameMatcher.Match(event) {
					if err := msg.Ack(); err != nil {
						log.Global().ErrorContext(ctx, "failed to ack message", zap.Error(err))
					}
					continue
				}

				gip.queue <- lo.T3(eventCtx, event, msg)
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
			case tuple := <-gip.queue:
				if err := gip.HandleEvent(tuple.A, tuple.B); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to handle event", zap.Error(err))
				}
				if err := tuple.C.Ack(); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to ack message", zap.Error(err))
				}
			}
		}
	}()

	log.Global().InfoContext(ctx, "initialized game exploding process manager")

	return gip, nil
}

func (w *GameExplodingProcessManager) HandlerType() common.EventHandlerType {
	return common.EventHandlerType(GameExplodingProcessManagerType)
}

func (w *GameExplodingProcessManager) HandleEvent(ctx context.Context, event common.Event) (err error) {
	if event == nil {
		return errors.WithStack(errors.New("event is nil"))
	}

	if event.AggregateType() == game.AggregateType {
		return w.GameProjector.HandleEvent(ctx, event)
	}

	return errors.WithStack(fmt.Errorf("unknown aggregate type %s", event.AggregateType()))
}

func (w *GameExplodingProcessManager) HandleGameInitialized(ctx context.Context, event common.Event, data *game.GameInitialized) error {
	w.gameDeskID[data.GameID.String()] = data.GetDeskID()

	return nil
}

func (w *GameExplodingProcessManager) HandleExplodingDrawn(ctx context.Context, event common.Event, data *game.ExplodingDrawn) error {
	w.gameActiveExplodingCardID[data.GetGameID().String()] = data.GetCardID()

	log.Global().InfoContext(ctx, "Exploding drawn", zap.String("gameID", data.GetGameID().String()), zap.String("playerID", data.GetPlayerID().String()))

	return nil
}

func (w *GameExplodingProcessManager) HandleExplodingDefused(ctx context.Context, event common.Event, data *game.ExplodingDefused) error {
	handID := hand.NewPlayerHandID(data.GetGameID(), data.GetPlayerID())
	if err := domains.CommandBus.HandleCommand(ctx, &hand.PlayCards{
		HandID:  handID,
		CardIDs: []uuid.UUID{data.GetCardID()},
	}); err != nil {
		log.Global().ErrorContext(ctx, "play hand cards error", zap.Error(err))
		return err
	}

	if err := domains.CommandBus.HandleCommand(ctx, &desk.DiscardCards{
		DeskID:  w.gameDeskID[data.GameID.String()],
		CardIDs: []uuid.UUID{data.GetCardID()},
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to discard cards", zap.Error(err))
		return err
	}

	log.Global().InfoContext(ctx, "Exploding defused", zap.String("gameID", data.GetGameID().String()), zap.String("playerID", data.GetPlayerID().String()))

	return nil
}

func (w *GameExplodingProcessManager) HandleKittenPlanted(ctx context.Context, event common.Event, data *game.KittenPlanted) error {
	cardID, ok := w.gameActiveExplodingCardID[data.GetGameID().String()]
	if !ok || cardID == uuid.Nil {
		return errors.WithStack(fmt.Errorf("game %s does not have an active exploding card", data.GetGameID().String()))
	}

	if err := domains.CommandBus.HandleCommand(ctx, &desk.InsertCard{
		DeskID: w.gameDeskID[data.GameID.String()],
		CardID: cardID,
		Index:  data.GetIndex(),
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to insert cards", zap.Error(err))
		return err
	}

	if err := domains.CommandBus.HandleCommand(ctx, &game.FinishTurn{
		GameID:   data.GetGameID(),
		PlayerID: data.GetPlayerID(),
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to finish current turn", zap.Error(err))
		return err
	}

	delete(w.gameActiveExplodingCardID, data.GetGameID().String())

	log.Global().InfoContext(ctx, "Kitten planted", zap.String("gameID", data.GetGameID().String()), zap.String("playerID", data.GetPlayerID().String()))

	return nil
}

func (w *GameExplodingProcessManager) HandlePlayerEliminated(ctx context.Context, event common.Event, data *game.PlayerEliminated) error {
	cardID, ok := w.gameActiveExplodingCardID[data.GetGameID().String()]
	if !ok || cardID == uuid.Nil {
		return errors.WithStack(fmt.Errorf("game %s does not have an active exploding card", data.GetGameID().String()))
	}

	if err := domains.CommandBus.HandleCommand(ctx, &desk.DiscardCards{
		DeskID:  w.gameDeskID[data.GameID.String()],
		CardIDs: []uuid.UUID{cardID},
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to discard cards", zap.Error(err))
		return err
	}

	delete(w.gameActiveExplodingCardID, data.GetGameID().String())

	return nil
}
