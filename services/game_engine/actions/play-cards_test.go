package actions_test

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/gofrs/uuid"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/actions"
)

func (as *ActionsSuite) Test_Validate_NultiCardPlay_CardMustPlayAlone() {
	as.setupEnvironment()
	_, cardsMapByCode := as.prepareCards()

	playerID := uuid.Must(uuid.NewV7())
	ctx := context.WithValue(context.Background(), grpc.AuthToken, playerID)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	playCardsRequest := &proto.PlayCardsRequest{
		CardIds: []string{
			cardsMapByCode[cards.Skip].GetCardId(),
			cardsMapByCode[cards.SeeTheFuture].GetCardId(),
		},
	}

	actions := actions.NewActions(ctx, "test")
	res, err := actions.PlayCards(
		ctx,
		connect.NewRequest(playCardsRequest),
	)

	msg := err.Error()
	fmt.Println(msg)

	as.Nil(res)
	as.ErrorContains(err, fmt.Sprintf("card '%s' must be played alone", cards.Skip))
}

func (as *ActionsSuite) Test_Validate_ComboCardPlay_NotAllowedCard() {
	as.setupEnvironment()
	_, cardsMapByCode := as.prepareCards()

	playerID := uuid.Must(uuid.NewV7())
	ctx := context.WithValue(context.Background(), grpc.AuthToken, playerID)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	playCardsRequest := &proto.PlayCardsRequest{
		CardIds: []string{
			cardsMapByCode[cards.TacoCat].GetCardId(),
			cardsMapByCode[cards.Shuffle].GetCardId(),
		},
	}

	actions := actions.NewActions(ctx, "test")
	res, err := actions.PlayCards(
		ctx,
		connect.NewRequest(playCardsRequest),
	)

	msg := err.Error()
	fmt.Println(msg)

	as.Nil(res)
	as.ErrorContains(err, fmt.Sprintf("card '%s' must be played alone", cards.Shuffle))
}

func (as *ActionsSuite) Test_Validate_ComboCardPlay_DifferenceCards() {
	as.setupEnvironment()
	_, cardsMapByCode := as.prepareCards()

	playerID := uuid.Must(uuid.NewV7())
	ctx := context.WithValue(context.Background(), grpc.AuthToken, playerID)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	playCardsRequest := &proto.PlayCardsRequest{
		CardIds: []string{
			cardsMapByCode[cards.TacoCat].GetCardId(),
			cardsMapByCode[cards.Catermelon].GetCardId(),
		},
	}

	actions := actions.NewActions(ctx, "test")
	res, err := actions.PlayCards(
		ctx,
		connect.NewRequest(playCardsRequest),
	)

	msg := err.Error()
	fmt.Println(msg)

	as.Nil(res)
	as.ErrorContains(err, "2-card combo must be two of the same combo card")
}
