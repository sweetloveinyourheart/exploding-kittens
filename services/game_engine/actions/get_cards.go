package actions

import (
	"context"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"

	"google.golang.org/protobuf/types/known/emptypb"

	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go"
)

func (a *actions) GetCards(ctx context.Context, request *connect.Request[emptypb.Empty]) (response *connect.Response[proto.GetCardsResponse], err error) {
	cards, err := a.cardRepo.GetCardsInformation(ctx)
	if err != nil {
		log.Global().Error("error getting card list", zap.Error(err))
		return nil, grpc.NotFoundError(err)
	}

	cardList := make([]*proto.Card, 0)
	for _, card := range cards {
		cardList = append(cardList, &proto.Card{
			CardId:      card.CardID.String(),
			Name:        card.Name,
			Description: card.Description,
		})
	}

	response = connect.NewResponse(&proto.GetCardsResponse{
		Cards: cardList,
	})
	return response, nil
}
