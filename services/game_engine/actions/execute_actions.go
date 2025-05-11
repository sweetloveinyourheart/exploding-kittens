package actions

import (
	"context"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
)

func (a *actions) ExecuteAction(ctx context.Context, request *connect.Request[proto.ExecuteActionRequest]) (response *connect.Response[emptypb.Empty], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	gameID := stringsutil.ConvertStringToUUID(request.Msg.GetGameId())
	if err := domains.CommandBus.HandleCommand(ctx, &game.ExecuteAction{
		GameID:         gameID,
		PlayerID:       userID,
		Effect:         request.Msg.GetEffect(),
		TargetPlayerID: stringsutil.ConvertStringToUUID(request.Msg.GetTargetUser()),
		TargetCardID:   stringsutil.ConvertStringToUUID(request.Msg.GetTargetCard()),
	}); err != nil {
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
