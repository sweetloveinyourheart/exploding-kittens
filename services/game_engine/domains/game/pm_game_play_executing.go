package game

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	card_effects "github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/card-effects"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	nats2 "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
	dataProviderGrpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/interfaces"
)

var (
	GamePlayExecutingHandlerType            = "game-play-executing"
	GamePlayExecutingConsumerAndDurableName = "game-play-executing-consumer-and-durable"
)

type GamePlayExecutor struct {
	ctx context.Context
	*game.GameProjector

	queue        chan lo.Tuple3[context.Context, common.Event, jetstream.Msg]
	dataProvider dataProviderGrpc.DataProviderClient

	gameDeskID      map[string]uuid.UUID
	gameCardsToDraw map[string]int
}

func NewGamePlayExecutor(ctx context.Context) (*GamePlayExecutor, error) {
	gpe := &GamePlayExecutor{
		ctx:          ctx,
		dataProvider: do.MustInvoke[dataProviderGrpc.DataProviderClient](nil),
		queue:        make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),

		gameCardsToDraw: make(map[string]int),
		gameDeskID:      make(map[string]uuid.UUID),
	}

	gpe.GameProjector = game.NewGameProjection(gpe)

	gameMatcher := eventing.NewMatchEventSubject(game.SubjectFactory, game.AggregateType,
		game.EventTypeGameCreated,
		game.EventTypeGameInitialized,
		game.EventTypeCardPlayed,
		game.EventTypeActionCreated,
		game.EventTypeActionExecuted,
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
		Name:        GamePlayExecutingConsumerAndDurableName,
		Durable:     GamePlayExecutingConsumerAndDurableName,
		Description: "Responsible for reading game events related to game play",
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

				gpe.queue <- lo.T3(eventCtx, event, msg)
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
			case tuple := <-gpe.queue:
				if err := gpe.HandleEvent(tuple.A, tuple.B); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to handle event", zap.Error(err))
				}
				if err := tuple.C.Ack(); err != nil {
					log.Global().ErrorContext(tuple.A, "failed to ack message", zap.Error(err))
				}
			}
		}
	}()

	log.Global().InfoContext(ctx, "initialized game play executing")

	return gpe, nil
}

func (w *GamePlayExecutor) HandleGameCreated(ctx context.Context, event common.Event, data *game.GameCreated) error {
	w.gameCardsToDraw[data.GameID.String()] = 1

	return nil
}

func (w *GamePlayExecutor) HandleGameInitialized(ctx context.Context, event common.Event, data *game.GameInitialized) error {
	w.gameDeskID[data.GameID.String()] = data.Desk

	return nil
}

func (w *GamePlayExecutor) HandleCardPlayed(ctx context.Context, event common.Event, data *game.CardPlayed) error {
	cards := data.GetCardIDs()
	if len(cards) == 0 {
		return errors.Errorf("failed to play: no card to play")
	}

	cardDataRes, err := w.dataProvider.GetMapCards(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		log.Global().Error("error retrieving cards map", zap.String("game_id", data.GameID.String()))
		return errors.Errorf("error retrieving cards map: %w", err)
	}

	cardsMap := cardDataRes.Msg.GetCards()

	cardInformation, ok := cardsMap[cards[0].String()]
	if !ok || cardInformation == nil {
		return errors.Errorf("cannot recognize card information")
	}

	var cardEffects []interfaces.CardEffect
	if len(cards) > 1 {
		// combo effects
		err = json.Unmarshal(cardInformation.ComboEffects, &cardEffects)
		if err != nil {
			return errors.Errorf("failed to unmarshal card combo effects: %w", err)
		}
	} else {
		// single effect
		var cardEffect interfaces.CardEffect
		err = json.Unmarshal(cardInformation.Effects, &cardEffect)
		if err != nil {
			return errors.Errorf("failed to unmarshal card effects: %w", err)
		}

		cardEffects = append(cardEffects, cardEffect)
	}

	if len(cardEffects) == 0 {
		return errors.Errorf("no card effects found")
	}

	for _, effect := range cardEffects {
		if err := domains.CommandBus.HandleCommand(ctx, &game.CreateAction{
			GameID:   data.GetGameID(),
			PlayerID: data.GetPlayerID(),
			Effect:   effect.Type,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to execute action", zap.Error(err))
			return err
		}
	}

	return nil
}

func (w *GamePlayExecutor) HandleActionCreated(ctx context.Context, event common.Event, data *game.ActionCreated) error {
	switch data.Effect {
	case card_effects.PeekCards:
		// peek cards
	case card_effects.StealCard:
		// steal card
	case card_effects.StealNamedCard:
		// steal named card
	case card_effects.StealRandomCard:
		// steal random card
	case card_effects.TradeAnyDiscard:
		// trade any discard
	default:
		if err := domains.CommandBus.HandleCommand(ctx, &game.ExecuteAction{
			GameID:   data.GetGameID(),
			TargetID: data.GetPlayerID(),
			Effect:   data.Effect,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to execute action", zap.Error(err))
			return err
		}
	}

	return nil
}

func (w *GamePlayExecutor) HandleActionExecuted(ctx context.Context, event common.Event, data *game.ActionExecuted) error {
	switch data.Effect {
	case card_effects.ShuffleDesk:
		if err := domains.CommandBus.HandleCommand(ctx, &desk.ShuffleDesk{
			DeskID:   w.gameDeskID[data.GameID.String()],
			GameID:   data.GetGameID(),
			PlayerID: data.GetTargetID(),
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to shuffle desk", zap.Error(err))
			return err
		}

	case card_effects.SkipTurn:
		if err := domains.CommandBus.HandleCommand(ctx, &game.FinishTurn{
			GameID:   data.GetGameID(),
			PlayerID: data.GetTargetID(),
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to finish current turn", zap.Error(err))
			return err
		}

	case card_effects.SkipTurnAndAttack:
		w.gameCardsToDraw[data.GameID.String()] = card_effects.AttackBonusCount + w.gameCardsToDraw[data.GameID.String()]

		if err := domains.CommandBus.HandleCommand(ctx, &game.FinishTurn{
			GameID:   data.GetGameID(),
			PlayerID: data.GetTargetID(),
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to finish current turn", zap.Error(err))
			return err
		}

	case card_effects.CancelAction:
		// cancel action
	case card_effects.PeekCards:
		// peek cards
	case card_effects.StealCard:
		// steal card
	case card_effects.StealNamedCard:
		// steal named card
	case card_effects.StealRandomCard:
		// steal random card
	case card_effects.TradeAnyDiscard:
		// trade any discard
	default:
		log.Global().ErrorContext(ctx, "unknown action effect", zap.String("effect", data.Effect))
		return errors.Errorf("unknown action effect: %s", data.Effect)
	}

	return nil
}
