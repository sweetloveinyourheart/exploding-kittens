package actions_test

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/actions"
)

func (as *ActionsSuite) Test_Validate_NultiCardPlay_CardMustPlayAlone() {
	as.setupEnvironment()
	_, cardsMapByCode := as.prepareCards()

	ctx, cancel := context.WithCancel(context.Background())
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

	ctx, cancel := context.WithCancel(context.Background())
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

	ctx, cancel := context.WithCancel(context.Background())
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

func (as *ActionsSuite) Test_Validate_ComboCardPlay_MustPlayInCombo() {
	as.setupEnvironment()
	_, cardsMapByCode := as.prepareCards()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	playCardsRequest := &proto.PlayCardsRequest{
		CardIds: []string{
			cardsMapByCode[cards.TacoCat].GetCardId(),
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
	as.ErrorContains(err, fmt.Sprintf("card '%s' cannot be played alone", cards.TacoCat))
}
