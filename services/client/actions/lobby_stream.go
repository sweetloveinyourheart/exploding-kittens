package actions

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"github.com/zmwangx/debounce"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains/match"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

const KeepAliveTimeout = 150 * time.Second

func (a *actions) StreamLobby(ctx context.Context, request *connect.Request[proto.GetLobbyRequest], stream *connect.ServerStream[proto.GetLobbyReply]) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var streamError error

	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	if err := (*GetLobbyRequestValidator)(request.Msg).Validate(); err != nil {
		return err
	}

	lobbyState, err := domains.LobbyRepo.Find(ctx, request.Msg.GetLobbyId())
	if err != nil {
		if errors.Is(err, eventing.ErrEntityNotFound) {
			return grpc.PreconditionError(grpc.PreconditionFailure("state", "lobby_id", "no such lobby"))
		}

		return grpc.NotFoundError(err)
	}

	isAuthorized := slices.Contains(lobbyState.GetParticipants(), userID)
	if !isAuthorized {
		return grpc.NotFoundError(errors.Errorf("Lobby not found"))
	}

	metrics := a.metrics
	metrics.LobbyStreamCounterInc(ctx, lobbyState.GetLobbyID())
	metrics.LobbyStreamGaugeAdd(ctx, lobbyState.GetLobbyID(), 1)
	defer func() {
		metrics.LobbyStreamGaugeAdd(ctx, lobbyState.GetLobbyID(), -1)
	}()

	lobbyChan := make(chan *nats.Msg, constants.NatsChannelBufferSize)
	lobbyUpdateStream, err := a.bus.ChanSubscribe(fmt.Sprintf("%s.%s", constants.LobbyStream, request.Msg.GetLobbyId()), lobbyChan)
	if err != nil {
		log.Global().ErrorContext(ctx, "Error subscribing to lobby update stream", zap.Error(err), zap.String("lobby_id", request.Msg.GetLobbyId()))
		return grpc.InternalError(err)
	}
	defer func() {
		err := lobbyUpdateStream.Unsubscribe()
		if err != nil {
			log.Global().ErrorContext(ctx, "Error unsubscribing from lobby update stream", zap.Error(err), zap.String("lobby_id", request.Msg.GetLobbyId()))
		}
	}()

	mux := &sync.Mutex{}
	sendData := func() {
		mux.Lock()
		defer mux.Unlock()

		lobbyState, err := domains.LobbyRepo.Find(ctx, request.Msg.GetLobbyId())
		if err != nil {
			if errors.Is(err, eventing.ErrEntityNotFound) {
				streamError = grpc.PreconditionError(grpc.PreconditionFailure("state", "lobby_id", "no such lobby"))
				cancel()
				return
			}
			streamError = err
			cancel()
			return
		}

		if ctx.Err() != nil {
			return
		}

		defer func() {
			if r := recover(); r != nil {
				// Handle and log the panic so it does not cause an entire system crash
				log.Global().ErrorContext(ctx, "recovered stream panic", zap.Any("panic", r), zap.String("user_id", userID.String()))
				streamError = errors.WithStack(errors.New("recovered stream panic"))
				cancel()
				return
			}
		}()

		// If the user is not in the room, close the stream
		if !slices.Contains(lobbyState.GetParticipants(), userID) {
			cancel()
			return
		}

		reply := &proto.GetLobbyReply{
			Lobby: &proto.Lobby{
				LobbyId:      lobbyState.GetLobbyID().String(),
				LobbyCode:    lobbyState.GetLobbyCode(),
				LobbyName:    lobbyState.GetLobbyName(),
				HostUserId:   lobbyState.GetHostUserID().String(),
				Participants: stringsutil.ConvertUUIDsToStrings(lobbyState.GetParticipants()),
				MatchId:      stringsutil.ConvertUUIDToStringPtr(lobbyState.GetMatchID()),
			},
		}

		err = stream.Send(reply)
		if err != nil {
			log.Global().ErrorContext(ctx, "Error sending lobby data", zap.Error(err), zap.String("user_id", userID.String()))
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

	unsubscribeLobby := domains.LobbySubscriber.SubscribeMatch(match.MatchLobbyID(lobbyState.GetLobbyID()), func(_ string, _ any) {
		debounced()
		keepAlive.Reset(KeepAliveTimeout)
	})
	defer unsubscribeLobby()

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
		case <-lobbyChan:
			debounced()
			keepAlive.Reset(KeepAliveTimeout)
		}
	}
}
