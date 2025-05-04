package actions_test

import (
	goTesting "testing"

	"connectrpc.com/connect"
	goMock "github.com/stretchr/testify/mock"

	"github.com/samber/do"

	dataProviderProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
	dataProviderGrpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"

	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/mock"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/testing"
)

type ActionsSuite struct {
	*testing.Suite
	deferred []func()

	mockDataProviderClient *mock.MockDataProviderClient
}

func (as *ActionsSuite) SetupTest() {
	as.mockDataProviderClient = new(mock.MockDataProviderClient)
}

func (as *ActionsSuite) TearDownTest() {
	if len(as.deferred) > 0 {
		for _, f := range as.deferred {
			f()
		}
		as.deferred = nil
	}
}

func TestActionsSuite(t *goTesting.T) {
	as := &ActionsSuite{
		Suite: testing.MakeSuite(t),
	}

	suite.Run(t, as)
}

func (as *ActionsSuite) setupEnvironment() {
	do.Override[dataProviderGrpc.DataProviderClient](nil, func(i *do.Injector) (dataProviderGrpc.DataProviderClient, error) {
		return as.mockDataProviderClient, nil
	})
}

func (as *ActionsSuite) prepareCards() (map[string]*dataProviderProto.Card, map[string]*dataProviderProto.Card) {
	cardsMap := map[string]*dataProviderProto.Card{
		"123e4567-e89b-12d3-a456-426655440001": {CardId: "123e4567-e89b-12d3-a456-426655440001", Code: cards.ExplodingKitten, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440002": {CardId: "123e4567-e89b-12d3-a456-426655440002", Code: cards.Defuse, Quantity: 6},
		"123e4567-e89b-12d3-a456-426655440003": {CardId: "123e4567-e89b-12d3-a456-426655440003", Code: cards.Attack, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440004": {CardId: "123e4567-e89b-12d3-a456-426655440004", Code: cards.Nope, Quantity: 5},
		"123e4567-e89b-12d3-a456-426655440005": {CardId: "123e4567-e89b-12d3-a456-426655440005", Code: cards.SeeTheFuture, Quantity: 5},
		"123e4567-e89b-12d3-a456-426655440006": {CardId: "123e4567-e89b-12d3-a456-426655440006", Code: cards.Shuffle, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440007": {CardId: "123e4567-e89b-12d3-a456-426655440007", Code: cards.Skip, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440008": {CardId: "123e4567-e89b-12d3-a456-426655440008", Code: cards.Favor, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440009": {CardId: "123e4567-e89b-12d3-a456-426655440009", Code: cards.BeardCat, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440010": {CardId: "123e4567-e89b-12d3-a456-426655440010", Code: cards.Catermelon, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440011": {CardId: "123e4567-e89b-12d3-a456-426655440011", Code: cards.HairyPotatoCat, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440012": {CardId: "123e4567-e89b-12d3-a456-426655440012", Code: cards.TacoCat, Quantity: 4},
		"123e4567-e89b-12d3-a456-426655440013": {CardId: "123e4567-e89b-12d3-a456-426655440013", Code: cards.RainbowRalphingCat, Quantity: 4},
	}

	cardsMapByCode := map[string]*dataProviderProto.Card{
		cards.ExplodingKitten:    {CardId: "123e4567-e89b-12d3-a456-426655440001", Code: cards.ExplodingKitten, Quantity: 4, Effects: []byte(`{"type": "explode"}`)},
		cards.Defuse:             {CardId: "123e4567-e89b-12d3-a456-426655440002", Code: cards.Defuse, Quantity: 6, Effects: []byte(`{"type": "prevent_explode"}`)},
		cards.Attack:             {CardId: "123e4567-e89b-12d3-a456-426655440003", Code: cards.Attack, Quantity: 4, Effects: []byte(`{"type": "skip_turn_and_attack"}`)},
		cards.Nope:               {CardId: "123e4567-e89b-12d3-a456-426655440004", Code: cards.Nope, Quantity: 5, Effects: []byte(`{"type": "cancel_action"}`)},
		cards.SeeTheFuture:       {CardId: "123e4567-e89b-12d3-a456-426655440005", Code: cards.SeeTheFuture, Quantity: 5, Effects: []byte(`{"type": "peek_cards"}`)},
		cards.Shuffle:            {CardId: "123e4567-e89b-12d3-a456-426655440006", Code: cards.Shuffle, Quantity: 4, Effects: []byte(`{"type": "shuffle_deck"}`)},
		cards.Skip:               {CardId: "123e4567-e89b-12d3-a456-426655440007", Code: cards.Skip, Quantity: 4, Effects: []byte(`{"type": "skip_turn"}`)},
		cards.Favor:              {CardId: "123e4567-e89b-12d3-a456-426655440008", Code: cards.Favor, Quantity: 4, Effects: []byte(`{"type": "steal_card"}`)},
		cards.BeardCat:           {CardId: "123e4567-e89b-12d3-a456-426655440009", Code: cards.BeardCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card"}, {"type": "steal_named_card"}, {"type": "trade_any_discard"}]`)},
		cards.Catermelon:         {CardId: "123e4567-e89b-12d3-a456-426655440010", Code: cards.Catermelon, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card"}, {"type": "steal_named_card"}, {"type": "trade_any_discard"}]`)},
		cards.HairyPotatoCat:     {CardId: "123e4567-e89b-12d3-a456-426655440011", Code: cards.HairyPotatoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card"}, {"type": "steal_named_card"}, {"type": "trade_any_discard"}]`)},
		cards.TacoCat:            {CardId: "123e4567-e89b-12d3-a456-426655440012", Code: cards.TacoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card"}, {"type": "steal_named_card"}, {"type": "trade_any_discard"}]`)},
		cards.RainbowRalphingCat: {CardId: "123e4567-e89b-12d3-a456-426655440013", Code: cards.RainbowRalphingCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card"}, {"type": "steal_named_card"}, {"type": "trade_any_discard"}]`)},
	}

	as.mockDataProviderClient.On("GetMapCards", goMock.Anything, goMock.Anything).Return(connect.NewResponse(&dataProviderProto.GetMapCardsResponse{
		Cards: cardsMap,
	}), nil)

	return cardsMap, cardsMapByCode
}
