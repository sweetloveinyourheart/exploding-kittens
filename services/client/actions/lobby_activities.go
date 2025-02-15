package actions

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/lobby"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	userProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

func (a *actions) CreateLobby(ctx context.Context, request *connect.Request[proto.CreateLobbyRequest]) (response *connect.Response[proto.CreateLobbyResponse], err error) {
	hostUserId, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	getUserRequest := userProto.GetUserRequest{UserId: hostUserId.String()}
	profile, err := a.userServerClient.GetUser(ctx, connect.NewRequest(&getUserRequest))
	if err != nil {
		return nil, grpc.NotFoundError(err)
	}

	lobbyId := uuid.Must(uuid.NewV7())
	lobbyName := request.Msg.GetLobbyName()
	lobbyCode, _ := stringsutil.GenerateRandomCode(6) // LobbyCode have 6 characters

	if err := domains.CommandBus.HandleCommand(ctx, &lobby.CreateLobby{
		LobbyID:    lobbyId,
		LobbyCode:  lobbyCode,
		LobbyName:  lobbyName,
		HostUserID: uuid.Must(uuid.FromString(profile.Msg.GetUser().GetUserId())),
	}); err != nil {
		if errors.Is(err, lobby.ErrLobbyAlreadyCreated) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "table_id", "table doesn't exists"))
		}

		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&proto.CreateLobbyResponse{
		LobbyId: lobbyId.String(),
	}), nil
}
