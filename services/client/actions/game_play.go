package actions

import (
	"context"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	card_effects "github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/card-effects"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/desk"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	gameProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/domains"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

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

func (a *actions) PeekCards(ctx context.Context, request *connect.Request[proto.PeekCardsRequest]) (response *connect.Response[proto.PeekCardsResponse], err error) {
	_, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	deskState, err := domains.DeskRepo.Find(ctx, request.Msg.GetDeskId())
	if err != nil {
		if errors.Is(err, desk.ErrDeskNotAvailable) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "desk_id", "no such desk"))
		}

		return nil, grpc.NotFoundError(err)
	}

	if err := domains.CommandBus.HandleCommand(ctx, &game.ExecuteAction{
		GameID: stringsutil.ConvertStringToUUID(request.Msg.GetGameId()),
		Effect: card_effects.PeekCards,
	}); err != nil {
		if errors.Is(err, game.ErrGameNotFound) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}

		if errors.Is(err, game.ErrPlayerNotInTheirTurn) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "player_id", "not in their turn"))
		}

		return nil, grpc.InternalError(err)
	}

	peekedCards := deskState.GetCardIDs()
	if len(peekedCards) > card_effects.PeekCardsCount {
		peekedCards = peekedCards[len(peekedCards)-card_effects.PeekCardsCount:]
	} else if len(peekedCards) > 0 {
		peekedCards = peekedCards[:]
	}

	return connect.NewResponse(&proto.PeekCardsResponse{
		CardIds: stringsutil.ConvertUUIDsToStrings(peekedCards),
	}), nil
}

func (a *actions) SelectAffectedPlayer(ctx context.Context, request *connect.Request[proto.SelectAffectedPlayerRequest]) (response *connect.Response[emptypb.Empty], err error) {
	_, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	if err := domains.CommandBus.HandleCommand(ctx, &game.SelectAffectedPlayer{
		GameID:   stringsutil.ConvertStringToUUID(request.Msg.GetGameId()),
		PlayerID: stringsutil.ConvertStringToUUID(request.Msg.GetPlayerId()),
	}); err != nil {
		if errors.Is(err, game.ErrGameNotFound) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}

		if errors.Is(err, game.ErrGameNotInActionPhase) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "game not in action phase"))
		}

		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *actions) StealCard(ctx context.Context, request *connect.Request[proto.StealCardRequest]) (response *connect.Response[emptypb.Empty], err error) {
	_, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	var cardEffect string
	var args game.ActionArguments
	if request.Msg.GetCardId() != "" {
		cardEffect = card_effects.StealNamedCard
		args.CardIDs = []uuid.UUID{stringsutil.ConvertStringToUUID(request.Msg.GetCardId())}
	} else {
		cardEffect = card_effects.StealRandomCard
		args.CardIndexes = []int{int(request.Msg.GetCardIndex())}
	}

	if err := domains.CommandBus.HandleCommand(ctx, &game.ExecuteAction{
		GameID: stringsutil.ConvertStringToUUID(request.Msg.GetGameId()),
		Effect: cardEffect,
		Args:   &args,
	}); err != nil {
		if errors.Is(err, game.ErrGameNotFound) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}

		if errors.Is(err, game.ErrPlayerNotInTheirTurn) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "player_id", "not in their turn"))
		}

		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *actions) GiveCard(ctx context.Context, request *connect.Request[proto.GiveCardRequest]) (response *connect.Response[emptypb.Empty], err error) {
	_, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	if err := domains.CommandBus.HandleCommand(ctx, &game.ExecuteAction{
		GameID: stringsutil.ConvertStringToUUID(request.Msg.GetGameId()),
		Effect: card_effects.StealCard,
		Args: &game.ActionArguments{
			CardIDs: []uuid.UUID{stringsutil.ConvertStringToUUID(request.Msg.GetCardId())},
		},
	}); err != nil {
		if errors.Is(err, game.ErrGameNotFound) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}

		if errors.Is(err, game.ErrGameNotInActionPhase) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "game not in action phase"))
		}

		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *actions) DrawCard(ctx context.Context, request *connect.Request[proto.DrawCardRequest]) (response *connect.Response[emptypb.Empty], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	if err := domains.CommandBus.HandleCommand(ctx, &game.DrawCard{
		GameID:   stringsutil.ConvertStringToUUID(request.Msg.GetGameId()),
		PlayerID: userID,
	}); err != nil {
		if errors.Is(err, game.ErrGameNotFound) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}
		if errors.Is(err, game.ErrPlayerNotInTheirTurn) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "player_id", "not in their turn"))
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *actions) DefuseExplodingKitten(ctx context.Context, request *connect.Request[proto.DefuseExplodingKittenRequest]) (response *connect.Response[emptypb.Empty], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	getUserRequest := gameProto.DefuseExplodingKittenRequest{
		GameId: request.Msg.GetGameId(),
		UserId: userID.String(),
		CardId: request.Msg.CardId,
	}

	_, err = a.gameEngineServerClient.DefuseExplodingKitten(ctx, connect.NewRequest(&getUserRequest))
	if err != nil {
		return nil, grpc.InvalidArgumentError(errors.Wrap(err, "failed to defuse exploding kitten"))
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *actions) PlantExplodingKitten(ctx context.Context, request *connect.Request[proto.PlantExplodingKittenRequest]) (response *connect.Response[emptypb.Empty], err error) {
	userID, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	if err := domains.CommandBus.HandleCommand(ctx, &game.PlantTheKitten{
		GameID:   stringsutil.ConvertStringToUUID(request.Msg.GetGameId()),
		PlayerID: userID,
		Index:    int(request.Msg.GetCardIndex()),
	}); err != nil {
		if errors.Is(err, game.ErrGameNotFound) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "no such game"))
		}

		if errors.Is(err, game.ErrGameNotInActionPhase) {
			return nil, grpc.PreconditionError(grpc.PreconditionFailure("state", "game_id", "game not in action phase"))
		}

		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
