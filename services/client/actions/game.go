package actions

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	gameProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

func (a *actions) GetGameMetaData(ctx context.Context, request *connect.Request[proto.GetGameMetaDataRequest]) (response *connect.Response[proto.GetGameMetaDataResponse], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	gameState, err := domains.GameRepo.Find(ctx, request.Msg.GetGameId())
	if err != nil {
		if errors.Is(err, eventing.ErrEntityNotFound) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}

		return nil, grpc.NotFoundError(err)
	}

	isAuthorized := false
	players := make([]string, 0)
	for _, player := range gameState.GetPlayers() {
		if player.GetPlayerID() == userID {
			isAuthorized = true
		}

		players = append(players, player.GetPlayerID().String())
	}

	if !isAuthorized {
		return nil, grpc.NotFoundError(errors.Errorf("Game not found"))
	}

	return connect.NewResponse(&proto.GetGameMetaDataResponse{
		Meta: &proto.GameMetaData{
			GameId:  gameState.GetGameID().String(),
			Players: players,
		},
	}), nil
}

func (a *actions) PlayCards(ctx context.Context, request *connect.Request[proto.PlayCardsRequest]) (response *connect.Response[emptypb.Empty], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	getUserRequest := gameProto.PlayCardsRequest{
		GameId:  request.Msg.GetGameId(),
		UserId:  userID.String(),
		CardIds: request.Msg.GetCardIds(),
	}

	_, err = a.gameEngineServerClient.PlayCards(ctx, connect.NewRequest(&getUserRequest))
	if err != nil {
		return nil, grpc.InvalidArgumentError(errors.Wrap(err, "failed to play cards"))
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *actions) ExecuteAction(ctx context.Context, request *connect.Request[proto.ExecuteActionRequest]) (response *connect.Response[emptypb.Empty], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	getUserRequest := gameProto.ExecuteActionRequest{
		GameId:     request.Msg.GetGameId(),
		UserId:     userID.String(),
		Effect:     request.Msg.GetEffect(),
		TargetUser: request.Msg.TargetUser,
		TargetCard: request.Msg.TargetCard,
	}

	_, err = a.gameEngineServerClient.ExecuteAction(ctx, connect.NewRequest(&getUserRequest))
	if err != nil {
		return nil, grpc.InvalidArgumentError(errors.Wrap(err, "failed to execute action"))
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
