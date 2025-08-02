package actions

import (
	"context"
	"encoding/json"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/cockroachdb/errors"

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

	cardList := make([]*proto.Card, 0)

	hKey := "dataprovider:cards"
	cachedCards, err := a.redisClient.HGet(ctx, hKey, GlobalCacheKey).Bytes()
	if err != nil {
		cardList, err = a.retrieveCards(ctx, hKey)
		if err != nil {
			log.Global().Error("error retrieving cards from repository", zap.Error(err))
			return nil, grpc.InternalError(errors.WithStack(err))
		}
	} else {
		err = json.Unmarshal(cachedCards, &cardList)
		if err != nil {
			return nil, grpc.InternalError(errors.WithStack(err))
		}
	}

	response = connect.NewResponse(&proto.GetCardsResponse{
		Cards: cardList,
	})
	return response, nil
}

func (a *actions) retrieveCards(ctx context.Context, hKey string) ([]*proto.Card, error) {
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

	bBytes, err := json.Marshal(cardList)
	if err != nil {
		return nil, grpc.InternalError(errors.WithStack(err))
	}

	err = a.redisClient.HSet(ctx, hKey, GlobalCacheKey, bBytes).Err()
	if err != nil {
		return nil, grpc.InternalError(errors.WithStack(err))
	}

	return cardList, nil
}

func (a *actions) GetMapCards(ctx context.Context, request *connect.Request[emptypb.Empty]) (response *connect.Response[proto.GetMapCardsResponse], err error) {
	opName := "dataprovider.GetMapCards()"
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindInternal),
	}
	ctx, span := a.tracer.Start(ctx, opName, opts...)
	defer span.End()

	cardMap := make(map[string]*proto.Card, 0)

	hKey := "dataprovider:map-cards"
	cachedCards, err := a.redisClient.HGet(ctx, hKey, GlobalCacheKey).Bytes()
	if err != nil {
		cardMap, err = a.retrieveMapCards(ctx, hKey)
		if err != nil {
			log.Global().Error("error retrieving cards from repository", zap.Error(err))
			return nil, grpc.InternalError(errors.WithStack(err))
		}
	} else {
		err = json.Unmarshal(cachedCards, &cardMap)
		if err != nil {
			return nil, grpc.InternalError(errors.WithStack(err))
		}
	}

	response = connect.NewResponse(&proto.GetMapCardsResponse{
		Cards: cardMap,
	})
	return response, nil
}

func (a *actions) retrieveMapCards(ctx context.Context, hKey string) (map[string]*proto.Card, error) {
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

	bBytes, err := json.Marshal(cardMap)
	if err != nil {
		return nil, grpc.InternalError(errors.WithStack(err))
	}

	err = a.redisClient.HSet(ctx, hKey, GlobalCacheKey, bBytes).Err()
	if err != nil {
		return nil, grpc.InternalError(errors.WithStack(err))
	}

	return cardMap, nil
}
