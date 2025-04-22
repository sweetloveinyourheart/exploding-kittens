package game

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	dataProviderProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
	dataProviderGrpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"

	"go.uber.org/zap"

	cardConstants "github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	nats2 "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
)

const (
	BASE_CARDS_PER_PLAYER = 7
)

var (
	ProcessingTimeout                     = 10 * time.Second
	BatchSize                             = 1024
	GameInteractionHandlerType            = "game-interaction-processing"
	GameInteractionConsumerandDurableName = "game-cancelled-void-choices-consumer"
)

type GameInteractionProcessor struct {
	ctx context.Context
	*game.GameProjector

	dataProvider dataProviderGrpc.DataProviderClient

	queue chan lo.Tuple3[context.Context, common.Event, jetstream.Msg]
}

func NewGameInteractionProcessor(ctx context.Context) (*GameInteractionProcessor, error) {
	lip := &GameInteractionProcessor{
		ctx:          ctx,
		dataProvider: do.MustInvoke[dataProviderGrpc.DataProviderClient](nil),
		queue:        make(chan lo.Tuple3[context.Context, common.Event, jetstream.Msg], BatchSize*2),
	}

	lip.GameProjector = game.NewGameProjection(lip)

	gameMatcher := eventing.NewMatchEventSubject(game.SubjectFactory, game.AggregateType,
		game.EventTypeGameCreated,
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
		Name:        GameInteractionConsumerandDurableName,
		Durable:     GameInteractionConsumerandDurableName,
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

	log.Global().InfoContext(ctx, "initialized game interaction processing")

	return lip, nil
}

func (w *GameInteractionProcessor) HandlerType() common.EventHandlerType {
	return common.EventHandlerType(GameInteractionHandlerType)
}

func (w *GameInteractionProcessor) HandleEvent(ctx context.Context, event common.Event) (err error) {
	if event == nil {
		return errors.WithStack(errors.New("event is nil"))
	}

	if event.AggregateType() == game.AggregateType {
		return w.GameProjector.HandleEvent(ctx, event)
	}

	return errors.WithStack(fmt.Errorf("unknown aggregate type %s", event.AggregateType()))
}

func (w *GameInteractionProcessor) HandleGameCreated(ctx context.Context, event common.Event, data *game.GameCreated) error {
	playerCount := len(data.GetPlayerIDs())

	// Pre setup cards
	cards, err := w.setupCards(ctx, playerCount)
	if err != nil {
		return fmt.Errorf("error setting up cards")
	}

	// Ensure we have enough cards for all players
	playerIDs := data.GetPlayerIDs()
	if len(cards.StandardCards) < len(playerIDs)*BASE_CARDS_PER_PLAYER || len(cards.DefuseCards) < len(playerIDs) {
		return errors.New("not enough cards to deal to all players")
	}

	// Shuffle standard cards
	mutable.Shuffle(cards.StandardCards)

	playerHands := make(map[uuid.UUID]uuid.UUID)
	for _, playerID := range playerIDs {
		standardCards := slices.Clone(cards.StandardCards[:BASE_CARDS_PER_PLAYER]) // deal 7 standard cards
		handCards := append(standardCards, cards.DefuseCards[0])                   // deal 1 defuse card

		// Trim used cards
		cards.StandardCards = slices.Delete(cards.StandardCards, 0, BASE_CARDS_PER_PLAYER)
		cards.DefuseCards = slices.Delete(cards.DefuseCards, 0, 1)

		// Shuffle player's cards
		mutable.Shuffle(handCards)

		handID := hand.NewPlayerHandID(data.GetGameID(), playerID)
		if err := domains.CommandBus.HandleCommand(ctx, &hand.CreateHand{
			HandID: handID,
			Cards:  handCards,
		}); err != nil {
			return fmt.Errorf("failed to create hand for player %s: %w", playerID, err)
		}
		playerHands[playerID] = handID
	}

	// Create a desk by calculating the total length needed for deskCards
	deskID := uuid.Must(uuid.NewV7())
	totalLength := len(cards.StandardCards) + len(cards.ExplodingKittenCards) + len(cards.DefuseCards)
	deskCards := make([]uuid.UUID, 0, totalLength)
	deskCards = append(deskCards, cards.StandardCards...)
	deskCards = append(deskCards, cards.ExplodingKittenCards...)
	deskCards = append(deskCards, cards.DefuseCards...)

	// Shuffle desk's cards
	mutable.Shuffle(deskCards)

	if err := domains.CommandBus.HandleCommand(ctx, &desk.CreateDesk{
		DeskID: deskID,
		Cards:  deskCards,
	}); err != nil {
		return err
	}

	// Init game args
	if err := domains.CommandBus.HandleCommand(ctx, &game.InitGameArgs{
		GameID:      data.GetGameID(),
		Desk:        deskID,
		PlayerHands: playerHands,
		PlayerTurn:  data.GetPlayerIDs()[rand.Intn(playerCount)],
	}); err != nil {
		return err
	}

	return nil
}

type CardSetup struct {
	StandardCards        []uuid.UUID
	ExplodingKittenCards []uuid.UUID
	DefuseCards          []uuid.UUID
}

func (w *GameInteractionProcessor) setupCards(ctx context.Context, playerNum int) (*CardSetup, error) {
	// Get cards registry
	response, err := w.dataProvider.GetCards(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		return nil, err
	}

	// Num of Exploding Kitten cards must be playerNum - 1
	explodingToAdd := playerNum - 1

	explodingKittenCards := make([]*dataProviderProto.Card, 0, explodingToAdd)
	defuseCards := make([]*dataProviderProto.Card, 0)
	standardCards := make([]*dataProviderProto.Card, 0)

	for _, card := range response.Msg.GetCards() {
		switch card.Name {
		case cardConstants.ExplodingKitten:
			for range explodingToAdd {
				explodingKittenCards = append(explodingKittenCards, card)
			}
		case cardConstants.Defuse:
			for range card.Quantity {
				defuseCards = append(defuseCards, card)
			}
		default:
			for range card.Quantity {
				standardCards = append(standardCards, card)
			}
		}
	}

	var standardCardIDs []uuid.UUID
	for _, card := range standardCards {
		standardCardIDs = append(standardCardIDs, uuid.FromStringOrNil(card.CardId))
	}

	var explodingKittenCardIDs []uuid.UUID
	for _, card := range explodingKittenCards {
		explodingKittenCardIDs = append(explodingKittenCardIDs, uuid.FromStringOrNil(card.CardId))
	}

	var defuseCardIDs []uuid.UUID
	for _, card := range defuseCards {
		defuseCardIDs = append(defuseCardIDs, uuid.FromStringOrNil(card.CardId))
	}

	return &CardSetup{
		StandardCards:        standardCardIDs,
		ExplodingKittenCards: explodingKittenCardIDs,
		DefuseCards:          defuseCardIDs,
	}, nil
}
