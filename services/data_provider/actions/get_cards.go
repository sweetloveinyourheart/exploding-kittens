package actions

import (
	"context"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"

	"google.golang.org/protobuf/types/known/emptypb"

	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
)

func (a *actions) GetCards(ctx context.Context, request *connect.Request[emptypb.Empty]) (response *connect.Response[proto.GetCardsResponse], err error) {
	opName := "dataprovider.GetCards()"
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindInternal),
	}
	ctx, span := a.tracer.Start(ctx, opName, opts...)
	defer span.End()

	cards, err := a.cardRepo.GetCards(ctx)
	if err != nil {
		log.Global().Error("error getting card list", zap.Error(err))
		return nil, grpc.NotFoundError(err)
	}

	cardList := make([]*proto.Card, 0)
	for _, card := range cards {
		cardList = append(cardList, &proto.Card{
			CardId:       card.CardID.String(),
			Code:         card.Code,
			Name:         card.Name,
			Description:  card.Description,
			Quantity:     int32(card.Quantity),
			Effects:      card.Effects,
			ComboEffects: card.ComboEffects,
		})
	}

	response = connect.NewResponse(&proto.GetCardsResponse{
		Cards: cardList,
	})
	return response, nil
}

func (a *actions) GetMapCards(ctx context.Context, request *connect.Request[emptypb.Empty]) (response *connect.Response[proto.GetMapCardsResponse], err error) {
	cards, err := a.cardRepo.GetCards(ctx)
	if err != nil {
		log.Global().Error("error getting card map", zap.Error(err))
		return nil, grpc.NotFoundError(err)
	}

	cardMap := make(map[string]*proto.Card, 0)
	for _, card := range cards {
		cardMap[card.CardID.String()] = &proto.Card{
			CardId:       card.CardID.String(),
			Code:         card.Code,
			Name:         card.Name,
			Description:  card.Description,
			Quantity:     int32(card.Quantity),
			Effects:      card.Effects,
			ComboEffects: card.ComboEffects,
		}
	}

	response = connect.NewResponse(&proto.GetMapCardsResponse{
		Cards: cardMap,
	})
	return response, nil
}
