package actions

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/zmwangx/debounce"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/hand"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
	gameClientDomain "github.com/sweetloveinyourheart/exploding-kittens/services/client/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains/match"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

func (a *actions) StreamGame(ctx context.Context, request *connect.Request[proto.StreamGameRequest], stream *connect.ServerStream[proto.StreamGameReply]) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var streamError error

	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	if err := (*StreamGameRequestValidator)(request.Msg).Validate(); err != nil {
		return err
	}

	gameState, err := domains.GameRepo.Find(ctx, request.Msg.GetGameId())
	if err != nil {
		if errors.Is(err, eventing.ErrEntityNotFound) {
			return grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}

		return grpc.NotFoundError(err)
	}

	isAuthorized := false
	for _, player := range gameState.GetPlayers() {
		if player.GetPlayerID() == userID {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return grpc.NotFoundError(errors.Errorf("Game not found"))
	}

	gameChan := make(chan *nats.Msg, constants.NatsChannelBufferSize)
	gameUpdateStream, err := a.bus.ChanSubscribe(fmt.Sprintf("%s.%s", constants.GameStream, request.Msg.GetGameId()), gameChan)
	if err != nil {
		log.Global().ErrorContext(ctx, "Error subscribing to game update stream", zap.Error(err), zap.String("game_id", request.Msg.GetGameId()))
		return grpc.InternalError(err)
	}
	defer func() {
		err := gameUpdateStream.Unsubscribe()
		if err != nil {
			log.Global().ErrorContext(ctx, "Error unsubscribing from game update stream", zap.Error(err), zap.String("game_id", request.Msg.GetGameId()))
		}
	}()

	deskChan := make(chan *nats.Msg, constants.NatsChannelBufferSize)
	if gameState.GetDeskID() != uuid.Nil {
		// Subscribe to desk updates if the game has a desk
		deskUpdateStream, err := a.bus.ChanSubscribe(fmt.Sprintf("%s.%s", constants.DeskStream, gameState.GetDeskID()), deskChan)
		if err != nil {
			log.Global().ErrorContext(ctx, "Error subscribing to desk update stream", zap.Error(err), zap.String("desk_id", gameState.GetDeskID().String()))
			return grpc.InternalError(err)
		}
		defer func() {
			err := deskUpdateStream.Unsubscribe()
			if err != nil {
				log.Global().ErrorContext(ctx, "Error unsubscribing from desk update stream", zap.Error(err), zap.String("desk_id", gameState.GetDeskID().String()))
			}
		}()
	}

	handChan := make(chan *nats.Msg, constants.NatsChannelBufferSize)
	if _, ok := gameState.GetPlayerHands()[userID]; ok {
		// Subscribe to hand updates if the user has a hand
		handUpdateStream, err := a.bus.ChanSubscribe(fmt.Sprintf("%s.%s", constants.HandStream, gameState.GetPlayerHands()[userID]), handChan)
		if err != nil {
			log.Global().ErrorContext(ctx, "Error subscribing to hand update stream", zap.Error(err), zap.String("hand_id", gameState.GetPlayerHands()[userID].String()))
			return grpc.InternalError(err)
		}
		defer func() {
			err := handUpdateStream.Unsubscribe()
			if err != nil {
				log.Global().ErrorContext(ctx, "Error unsubscribing from hand update stream", zap.Error(err), zap.String("hand_id", gameState.GetPlayerHands()[userID].String()))
			}
		}()
	}

	builder := gameClientDomain.NewGameResponseBuilder(userID)

	mux := &sync.Mutex{}
	sendData := func() {
		mux.Lock()
		defer mux.Unlock()

		if cancelCtx.Err() != nil {
			return
		}

		gameState, err := domains.GameRepo.Find(ctx, request.Msg.GetGameId())
		if err != nil {
			streamError = err
			cancel()
			return
		}

		deskState, err := domains.DeskRepo.Find(ctx, gameState.DeskID.String())
		if err != nil {
			streamError = err
			cancel()
			return
		}

		handStates := make(map[string]*hand.Hand)
		for _, player := range gameState.GetPlayers() {
			playerID := player.GetPlayerID()

			handID := hand.NewPlayerHandID(gameState.GetGameID(), playerID)
			handState, err := domains.HandRepo.Find(ctx, handID.String())
			if err != nil {
				streamError = err
				cancel()
				return
			}

			handStates[playerID.String()] = handState
		}

		game, err := builder.Build(gameState, deskState, handStates)
		if err != nil {
			log.Global().ErrorContext(ctx, "Error fetching game data for stream", zap.Error(err))
			streamError = err
			cancel()
			return
		}

		response := &proto.StreamGameReply{
			GameState: game,
		}

		err = stream.Send(response)
		if err != nil {
			log.Global().ErrorContext(ctx, "Error sending game data", zap.Error(err), zap.String("user_id", userID.String()))
			streamError = err
			cancel()
			return
		}
	}

	sendData()

	debounced, control := debounce.Throttle(sendData, 500*time.Millisecond)
	keepAlive := time.NewTimer(KeepAliveTimeout)
	defer func() {
		if !keepAlive.Stop() {
			// drain the timer chan
			select {
			case <-keepAlive.C:
			default:
			}
		}
		control.Cancel()
	}()

	if streamError != nil {
		return streamError
	}

	unsubscribeGame := domains.GameSubscriber.SubscribeMatch(match.MatchGameID(gameState.GetGameID()), func(_ string, _ any) {
		debounced()
		keepAlive.Reset(KeepAliveTimeout)
	})
	defer unsubscribeGame()

	unsubscribeDesk := domains.DeskSubscriber.SubscribeMatch(match.MatchDeskID(gameState.GetDeskID()), func(_ string, _ any) {
		debounced()
		keepAlive.Reset(KeepAliveTimeout)
	})
	defer unsubscribeDesk()

	unsubscribeHand := domains.HandSubscriber.SubscribeMatch(match.MatchHandID(gameState.GetPlayerHands()[userID]), func(_ string, _ any) {
		debounced()
		keepAlive.Reset(KeepAliveTimeout)
	})
	defer unsubscribeHand()

	for {
		select {
		case <-ctx.Done():
			log.Global().InfoContext(ctx, "stream context done, closing stream", zap.String("user_id", userID.String()))
			return ctx.Err()
		case <-cancelCtx.Done():
			log.Global().InfoContext(ctx, "stream context cancelled, closing stream", zap.String("user_id", userID.String()))
			return streamError
		case <-keepAlive.C:
			debounced()
			keepAlive.Reset(KeepAliveTimeout)
		case <-gameChan:
			debounced()
			keepAlive.Reset(KeepAliveTimeout)
		case <-deskChan:
			debounced()
			keepAlive.Reset(KeepAliveTimeout)
		case <-handChan:
			debounced()
			keepAlive.Reset(KeepAliveTimeout)
		}
	}
}

type StreamGameRequestValidator proto.StreamGameRequest

func (request *StreamGameRequestValidator) Validate() error {
	var fieldErrors []*errdetails.BadRequest_FieldViolation
	_, err := uuid.FromString(strings.TrimSpace(request.GameId))
	if err != nil {
		fieldErrors = append(fieldErrors, grpc.FieldViolation("game_id", err))
	}

	if fieldErrors == nil {
		return nil
	}

	return grpc.InvalidArgumentErrorWithField(fieldErrors...)
}
