package actions

import (
	"context"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/helpers"
)

func (a *actions) RetrieveCardsData(ctx context.Context, request *connect.Request[emptypb.Empty]) (response *connect.Response[proto.RetrieveCardsDataResponse], err error) {
	_, ok := ctx.Value(grpc.AuthToken).(uuid.UUID)
	if !ok {
		// This should never happen as this endpoint should be authenticated
		return nil, grpc.UnauthenticatedError(helpers.ErrInvalidSession)
	}

	res, err := a.dataProviderClient.GetCards(ctx, connect.NewRequest(&emptypb.Empty{}))
	if err != nil {
		return nil, grpc.NotFoundError(err)
	}

	cardsData := make([]*proto.Card, 0)
	for _, card := range res.Msg.GetCards() {
		cardData := &proto.Card{
			CardId:      card.GetCardId(),
			Name:        card.GetName(),
			Code:        card.GetCode(),
			Description: card.GetDescription(),
		}

		cardsData = append(cardsData, cardData)
	}

	return connect.NewResponse(&proto.RetrieveCardsDataResponse{
		Cards: cardsData,
	}), nil
}
