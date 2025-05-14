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
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
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

	gameDeskID            map[string]uuid.UUID
	gameCardsToDraw       map[string]int
	gamePlayerHands       map[string]map[uuid.UUID]uuid.UUID
	gameAffectingPlayerID map[string]uuid.UUID
	gamePlayerTurnID      map[string]uuid.UUID
}

func NewGamePlayExecutor(ctx context.Context) (*GamePlayExecutor, error) {
	gpe := &GamePlayExecutor{
		ctx:          ctx,
		dataProvider: do.MustInvoke[dataProviderGrpc.DataProviderClient](nil),
		queue:        make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),

		gameDeskID:            make(map[string]uuid.UUID),
		gameCardsToDraw:       make(map[string]int),
		gamePlayerHands:       make(map[string]map[uuid.UUID]uuid.UUID),
		gameAffectingPlayerID: make(map[string]uuid.UUID),
		gamePlayerTurnID:      make(map[string]uuid.UUID),
	}

	gpe.GameProjector = game.NewGameProjection(gpe)

	gameMatcher := eventing.NewMatchEventSubject(game.SubjectFactory, game.AggregateType,
		game.EventTypeGameInitialized,
		game.EventTypeTurnStarted,
		game.EventTypeCardsPlayed,
		game.EventTypeActionCreated,
		game.EventTypeAffectedPlayerSelected,
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

func (w *GamePlayExecutor) HandleGameInitialized(ctx context.Context, event common.Event, data *game.GameInitialized) error {
	w.gameCardsToDraw[data.GameID.String()] = 1
	w.gameDeskID[data.GameID.String()] = data.GetDesk()
	w.gamePlayerHands[data.GameID.String()] = data.GetPlayerHands()
	w.gameAffectingPlayerID[data.GameID.String()] = uuid.Nil
	w.gamePlayerTurnID[data.GameID.String()] = uuid.Nil

	return nil
}

func (w *GamePlayExecutor) HandleTurnStarted(ctx context.Context, event common.Event, data *game.TurnStarted) error {
	w.gamePlayerTurnID[data.GameID.String()] = data.GetPlayerID()

	return nil
}

func (w *GamePlayExecutor) HandleCardsPlayed(ctx context.Context, event common.Event, data *game.CardsPlayed) error {
	if err := domains.CommandBus.HandleCommand(ctx, &hand.PlayCards{
		HandID:  w.gamePlayerHands[data.GameID.String()][data.PlayerID],
		CardIDs: data.GetCardIDs(),
	}); err != nil {
		log.Global().ErrorContext(ctx, "play hand cards error", zap.Error(err))
		return err
	}

	if err := domains.CommandBus.HandleCommand(ctx, &desk.DiscardCards{
		DeskID:  w.gameDeskID[data.GameID.String()],
		CardIDs: data.GetCardIDs(),
	}); err != nil {
		log.Global().ErrorContext(ctx, "failed to discard cards", zap.Error(err))
		return err
	}

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

	var effects []string
	if len(cards) > 1 {
		// combo effects
		var comboEffect []interfaces.CardComboEffect
		err = json.Unmarshal(cardInformation.ComboEffects, &comboEffect)
		if err != nil {
			return errors.Errorf("failed to unmarshal card combo effects: %w", err)
		}

		for _, effect := range comboEffect {
			if effect.RequiredCards == len(cards) {
				effects = append(effects, effect.Type)
			}
		}
	} else {
		// single effect
		var cardEffect interfaces.CardEffect
		err = json.Unmarshal(cardInformation.Effects, &cardEffect)
		if err != nil {
			return errors.Errorf("failed to unmarshal card effects: %w", err)
		}

		effects = append(effects, cardEffect.Type)
	}

	if len(effects) == 0 {
		return errors.Errorf("no card effects found")
	}

	for _, effect := range effects {
		if err := domains.CommandBus.HandleCommand(ctx, &game.CreateAction{
			GameID: data.GetGameID(),
			Effect: effect,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to execute action", zap.Error(err))
			return err
		}
	}

	return nil
}

func (w *GamePlayExecutor) HandleActionCreated(ctx context.Context, event common.Event, data *game.ActionCreated) error {
	switch data.Effect {
	// Manual effects
	case card_effects.PeekCards:
	case card_effects.StealCard:
	case card_effects.StealNamedCard:
	case card_effects.StealRandomCard:

	// Auto effects
	default:
		if err := domains.CommandBus.HandleCommand(ctx, &game.ExecuteAction{
			GameID: data.GetGameID(),
			Effect: data.Effect,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to execute action", zap.Error(err))
			return err
		}
	}

	return nil
}

func (w *GamePlayExecutor) HandleAffectedPlayerSelected(ctx context.Context, event common.Event, data *game.AffectedPlayerSelected) error {
	w.gameAffectingPlayerID[data.GameID.String()] = data.GetPlayerID()

	return nil
}

func (w *GamePlayExecutor) HandleActionExecuted(ctx context.Context, event common.Event, data *game.ActionExecuted) error {
	switch data.Effect {
	case card_effects.ShuffleDesk: // Shuffle the desk
		if err := domains.CommandBus.HandleCommand(ctx, &desk.ShuffleDesk{
			DeskID: w.gameDeskID[data.GameID.String()],
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to shuffle desk", zap.Error(err))
			return err
		}

	case card_effects.SkipTurn: // Skip the current turn
		currentPlayerID, ok := w.gamePlayerTurnID[data.GameID.String()]
		if !ok {
			return errors.Errorf("no player turn found for game ID: %s", data.GameID.String())
		}
		if currentPlayerID == uuid.Nil {
			return errors.Errorf("no target player found for game ID: %s", data.GameID.String())
		}

		if err := domains.CommandBus.HandleCommand(ctx, &game.FinishTurn{
			GameID:   data.GetGameID(),
			PlayerID: currentPlayerID,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to finish current turn", zap.Error(err))
			return err
		}

	case card_effects.SkipTurnAndAttack: // Skip the current turn and attack
		currentPlayerID, ok := w.gamePlayerTurnID[data.GameID.String()]
		if !ok {
			return errors.Errorf("no player turn found for game ID: %s", data.GameID.String())
		}
		if currentPlayerID == uuid.Nil {
			return errors.Errorf("no target player found for game ID: %s", data.GameID.String())
		}

		if err := domains.CommandBus.HandleCommand(ctx, &game.FinishTurn{
			GameID:   data.GetGameID(),
			PlayerID: currentPlayerID,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to finish current turn", zap.Error(err))
			return err
		}
		w.gameCardsToDraw[data.GameID.String()] = card_effects.AttackBonusCount + w.gameCardsToDraw[data.GameID.String()]

	case card_effects.PeekCards:
		// No action needed, just a notification

	case card_effects.StealCard: // Target player must give a card to the current player
		currentPlayerID, ok := w.gamePlayerTurnID[data.GameID.String()]
		if !ok {
			return errors.Errorf("no player turn found for game ID: %s", data.GameID.String())
		}
		if currentPlayerID == uuid.Nil {
			return errors.Errorf("no target player found for game ID: %s", data.GameID.String())
		}

		targetPlayerID := w.gameAffectingPlayerID[data.GameID.String()]
		if targetPlayerID == uuid.Nil {
			return errors.Errorf("no target player found for game ID: %s", data.GameID.String())
		}

		handID, ok := w.gamePlayerHands[data.GameID.String()][targetPlayerID]
		if !ok {
			return errors.Errorf("no hand found for player ID: %s", targetPlayerID.String())
		}

		toHandID, ok := w.gamePlayerHands[data.GameID.String()][currentPlayerID]
		if !ok {
			return errors.Errorf("no hand found for player ID: %s", currentPlayerID.String())
		}

		cardIDs := make([]uuid.UUID, 0)
		if data.GetCardID() != uuid.Nil {
			cardIDs = append(cardIDs, data.GetCardID())
		} else {
			return errors.Errorf("no target card ID found for game ID: %s", data.GameID.String())
		}

		if err := domains.CommandBus.HandleCommand(ctx, &hand.GiveCards{
			HandID:   handID,
			ToHandID: toHandID,
			CardIDs:  cardIDs,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to steal card", zap.Error(err))
			return err
		}
		w.gameAffectingPlayerID[data.GameID.String()] = uuid.Nil

	case card_effects.StealNamedCard, card_effects.StealRandomCard: // Steal a card from the target player
		currentPlayerID := w.gamePlayerTurnID[data.GameID.String()]
		if currentPlayerID == uuid.Nil {
			return errors.Errorf("no current player found for game ID: %s", data.GameID.String())
		}

		targetPlayerID, ok := w.gameAffectingPlayerID[data.GameID.String()]
		if !ok {
			return errors.Errorf("no affecting player found for game ID: %s", data.GameID.String())
		}
		if targetPlayerID == uuid.Nil {
			return errors.Errorf("no target player found for game ID: %s", data.GameID.String())
		}

		handID, ok := w.gamePlayerHands[data.GameID.String()][targetPlayerID]
		if !ok {
			return errors.Errorf("no hand found for player ID: %s", targetPlayerID.String())
		}

		toHandID, ok := w.gamePlayerHands[data.GameID.String()][currentPlayerID]
		if !ok {
			return errors.Errorf("no hand found for player ID: %s", currentPlayerID.String())
		}

		cardIDs := make([]uuid.UUID, 0)
		if data.GetCardID() != uuid.Nil {
			cardIDs = append(cardIDs, data.GetCardID())
		}

		if err := domains.CommandBus.HandleCommand(ctx, &hand.GiveCards{
			HandID:   handID,
			ToHandID: toHandID,
			CardIDs:  cardIDs,
		}); err != nil {
			log.Global().ErrorContext(ctx, "failed to steal card", zap.Error(err))
			return err
		}
		w.gameAffectingPlayerID[data.GameID.String()] = uuid.Nil

	case card_effects.CancelAction: // Cancel the action
		currentPlayerID := w.gamePlayerTurnID[data.GameID.String()]
		if currentPlayerID == uuid.Nil {
			return errors.Errorf("no current player found for game ID: %s", data.GameID.String())
		}

		if w.gameAffectingPlayerID[data.GameID.String()] == uuid.Nil {
			if err := domains.CommandBus.HandleCommand(ctx, &game.ReverseTurn{
				GameID:   data.GetGameID(),
				PlayerID: currentPlayerID,
			}); err != nil {
				log.Global().ErrorContext(ctx, "failed to reverse current turn", zap.Error(err))
				return err
			}
		}

		w.gameCardsToDraw[data.GameID.String()] = card_effects.AttackBonusCount
		w.gameAffectingPlayerID[data.GameID.String()] = uuid.Nil

	default:
		log.Global().ErrorContext(ctx, "unknown action effect", zap.String("effect", data.Effect))
		return errors.Errorf("unknown action effect: %s", data.Effect)
	}

	log.Global().InfoContext(ctx, "Action executed", zap.String("gameID", data.GetGameID().String()), zap.String("effect", data.GetEffect()))

	return nil
}
