package actions

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	cardsConst "github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
)

func (a *actions) DefuseExplodingKitten(ctx context.Context, request *connect.Request[proto.DefuseExplodingKittenRequest]) (response *connect.Response[emptypb.Empty], err error) {
	gameID := stringsutil.ConvertStringToUUID(request.Msg.GetGameId())
	playerId := stringsutil.ConvertStringToUUID(request.Msg.GetUserId())

	if request.Msg.GetCardId() == "" {
		// If no card ID is provided, we assume the player is eliminated because of an exploding kitten
		if err := domains.CommandBus.HandleCommand(ctx, &game.EliminatePlayer{
			GameID:   gameID,
			PlayerID: playerId,
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

	// Validate the defuse card
	err = a.validateDefuse(ctx, request.Msg.GetCardId())
	if err != nil {
		return nil, grpc.InvalidArgumentError(err)
	}

	// If a card ID is provided, we assume the player is defusing the exploding kitten
	if err := domains.CommandBus.HandleCommand(ctx, &game.DefuseExplodingKitten{
		GameID:   gameID,
		PlayerID: playerId,
		CardID:   stringsutil.ConvertStringToUUID(request.Msg.GetCardId()),
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

func (a *actions) validateDefuse(ctx context.Context, cardID string) error {
	cardsDataRes, err := a.dataProvider.GetMapCards(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		return err
	}

	cardMap := cardsDataRes.Msg.GetCards()

	if card, ok := cardMap[cardID]; ok {
		if card.GetCode() == cardsConst.Defuse {
			return nil
		}
		return errors.New("card is not a defuse card")
	}

	return errors.New("card is not recognized")
}
