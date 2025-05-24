package actions

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/cockroachdb/errors"

	cardsConst "github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domains/game"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	proto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains"
)

func (a *actions) PlayCards(ctx context.Context, request *connect.Request[proto.PlayCardsRequest]) (response *connect.Response[emptypb.Empty], err error) {
	err = a.validatePlayable(ctx, request.Msg.GetCardIds())
	if err != nil {
		return nil, grpc.InvalidArgumentError(err)
	}

	gameID := stringsutil.ConvertStringToUUID(request.Msg.GetGameId())
	playerId := stringsutil.ConvertStringToUUID(request.Msg.GetUserId())
	if err := domains.CommandBus.HandleCommand(ctx, &game.PlayCards{
		GameID:   gameID,
		PlayerID: playerId,
		CardIDs:  stringsutil.ConvertStringsToUUIDs(request.Msg.GetCardIds()),
	}); err != nil {
		return nil, grpc.InternalError(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (a *actions) validatePlayable(ctx context.Context, cardIds []string) error {
	cardsDataRes, err := a.dataProvider.GetMapCards(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		return err
	}

	cardMap := cardsDataRes.Msg.GetCards()

	cards := make([]string, 0, len(cardIds))
	for _, cardId := range cardIds {
		if card, ok := cardMap[cardId]; ok {
			cards = append(cards, card.GetCode())
		} else {
			return errors.Errorf("card '%s' is not recognized", cardId)
		}
	}

	if len(cards) == 0 {
		return errors.New("no cards played")
	}

	// Exploding kitten and defuse cannot be played manually
	for _, card := range cards {
		if card == cardsConst.ExplodingKitten {
			return errors.New("exploding kitten cannot be played manually")
		}
		if card == cardsConst.Defuse {
			return errors.New("defuse card cannot be played on your turn")
		}
	}

	// If playing one card, allow it if it's a known card
	if len(cards) == 1 {
		card := cards[0]
		if cardsConst.MustPlayAlone[card] {
			return nil
		}

		return errors.Errorf("card '%s' cannot be played alone", card)
	}

	// If multiple cards
	// Rule 1: no MustPlayAlone card can be combined
	for _, card := range cards {
		if cardsConst.MustPlayAlone[card] {
			return errors.Errorf("card '%s' must be played alone", card)
		}
	}

	// Rule 2: handle combo validation
	counts := make(map[string]int)
	unique := make(map[string]bool)
	for _, card := range cards {
		if !cardsConst.ComboCards[card] {
			return errors.Errorf("card '%s' cannot be used in combos", card)
		}
		counts[card]++
		unique[card] = true
	}

	switch len(cards) {
	case 2:
		// must be 2 of the same combo card
		for _, count := range counts {
			if count != 2 {
				return errors.New("2-card combo must be two of the same combo card")
			}
		}
	case 3:
		// must be 3 of the same combo card
		for _, count := range counts {
			if count != 3 {
				return errors.New("3-card combo must be three of the same combo card")
			}
		}
	default:
		return errors.New("invalid number of cards for a combo play")
	}

	return nil
}
