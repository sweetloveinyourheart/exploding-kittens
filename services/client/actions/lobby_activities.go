package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

func (a *actions) CreateLobby(ctx context.Context, request *connect.Request[proto.CreateLobbyRequest]) (response *connect.Response[proto.CreateLobbyResponse], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	lobbyId := uuid.Must(uuid.NewV7())
	lobbyName := request.Msg.GetLobbyName()
	lobbyCode, _ := stringsutil.GenerateRandomCode(6) // LobbyCode have 6 characters

	if err := domains.CommandBus.HandleCommand(ctx, &lobby.CreateLobby{
		LobbyID:    lobbyId,
		LobbyCode:  lobbyCode,
		LobbyName:  lobbyName,
		HostUserID: userID,
	}); err != nil {
		if errors.Is(err, lobby.ErrLobbyAlreadyCreated) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "lobby_id", "lobby already created"))
		}

		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.CreateLobbyResponse{
		LobbyId: lobbyId.String(),
	}), nil
}

func (a *actions) JoinLobby(ctx context.Context, request *connect.Request[proto.JoinLobbyRequest]) (response *connect.Response[proto.JoinLobbyResponse], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	lobbies, err := domains.LobbyRepo.FindAll(ctx)
	if err != nil {
		return nil, grpc.NotFoundError(err)
	}

	var lobbyState *lobby.Lobby
	for _, lobby := range lobbies {
		if lobby.LobbyCode == request.Msg.GetLobbyCode() {
			lobbyState = lobby
			break
		}
	}

	if lobbyState == nil {
		return nil, grpc.NotFoundError(errors.New("lobby not found"))
	}

	lobbyID := lobbyState.GetLobbyID()
	if err := domains.CommandBus.HandleCommand(ctx, &lobby.JoinLobby{
		LobbyID: lobbyID,
		UserID:  userID,
	}); err != nil {
		if errors.Is(err, lobby.ErrLobbyNotAvailable) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "lobby_id", "lobby is not available"))
		}

		return nil, grpc.InternalError(err)
	}

	err = a.emitLobbyUpdateEvent(lobbyID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.JoinLobbyResponse{
		LobbyId: lobbyID.String(),
	}), nil
}

func (a *actions) LeaveLobby(ctx context.Context, request *connect.Request[proto.LeaveLobbyRequest]) (response *connect.Response[proto.LeaveLobbyResponse], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	lobbyIDString := request.Msg.GetLobbyId()
	lobbyID, err := uuid.FromString(lobbyIDString)
	if err != nil {
		return nil, grpc.InvalidArgumentError(errors.New("lobby_id must be in the correct format"))
	}

	if err := domains.CommandBus.HandleCommand(ctx, &lobby.LeaveLobby{
		LobbyID: lobbyID,
		UserID:  userID,
	}); err != nil {
		if errors.Is(err, lobby.ErrLobbyNotAvailable) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "lobby_id", "lobby is not availale"))
		}

		return nil, grpc.InternalError(err)
	}

	err = a.emitLobbyUpdateEvent(lobbyID)
	if err != nil {
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.LeaveLobbyResponse{
		LobbyId: lobbyID.String(),
	}), nil
}

func (a *actions) emitLobbyUpdateEvent(lobbyID uuid.UUID) error {
	msg := &proto.Lobby{
		LobbyId: lobbyID.String(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = a.bus.Publish(fmt.Sprintf("%s.%s", constants.LobbyStream, msg.GetLobbyId()), msgBytes)
	if err != nil {
		return err
	}

	return nil
}

type GetLobbyRequestValidator proto.GetLobbyRequest

func (request *GetLobbyRequestValidator) Validate() error {
	var fieldErrors []*errdetails.BadRequest_FieldViolation
	_, err := uuid.FromString(strings.TrimSpace(request.LobbyId))
	if err != nil {
		fieldErrors = append(fieldErrors, grpc.FieldViolation("lobby_id", err))
	}

	if fieldErrors == nil {
		return nil
	}

	return grpc.InvalidArgumentErrorWithField(fieldErrors...)
}
