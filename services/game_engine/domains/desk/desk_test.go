package desk_test

import (
	"context"
	"fmt"
	goTesting "testing"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/testing"
	dataProviderProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
	gameengine "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine"
	deskDomain "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/domains/desk"

	"go.uber.org/zap"
)

type DeskSuite struct {
	*testing.Suite
	deferred []func()
}

func (as *DeskSuite) SetupTest() {}

func (as *DeskSuite) TearDownTest() {
	if len(as.deferred) > 0 {
		for _, f := range as.deferred {
			f()
		}
		as.deferred = nil
	}
}

func TestDeskSuite(t *goTesting.T) {
	hs := &DeskSuite{
		Suite: testing.MakeSuite(t),
	}
	suite.Run(t, hs)
}

func (hs *DeskSuite) setupEnvironment() {
	bus, shutdown := testing.StartLocalNATSServer(hs.T())
	hs.deferred = append(hs.deferred, shutdown)

	busAddress := "nats://" + bus.Addr().String() // use server address
	busConnection, err := nats.Connect(busAddress)
	hs.NoError(err)

	testing.NATSWaitConnected(hs.T(), busConnection) // wait connection if not connected yet

	jetStream, err := busConnection.JetStream()
	hs.NoError(err)

	connPool := pool.New(100, busAddress,
		nats.NoEcho(),
		nats.Name("kittens/gameengineserver/1.0"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.Error(err))
		}),
	)

	do.OverrideNamed[*pool.ConnPool](nil, string(constants.ConnectionPool),
		func(i *do.Injector) (*pool.ConnPool, error) {
			return connPool, nil
		})

	do.OverrideNamed[*nats.Conn](nil, fmt.Sprintf("%s-conn", string(constants.Bus)),
		func(i *do.Injector) (*nats.Conn, error) {
			return busConnection, nil
		})

	do.OverrideNamed[nats.JetStreamContext](nil, string(constants.Bus),
		func(i *do.Injector) (nats.JetStreamContext, error) {
			return jetStream, nil
		})

	ctx, cancel := context.WithCancel(context.Background())
	hs.deferred = append(hs.deferred, cancel)

	err = gameengine.InitializeCoreRepos(uuid.Must(uuid.NewV7()).String(), ctx)
	hs.NoError(err)

	_, err = deskDomain.NewDeskStateProcessor(ctx)

	hs.NoError(err)
}

func (gs *DeskSuite) prepareCards() ([]*dataProviderProto.Card, map[string]*dataProviderProto.Card, map[string]*dataProviderProto.Card) {
	cardsData := []*dataProviderProto.Card{
		{CardId: "123e4567-e89b-12d3-a456-426655440001", Code: cards.ExplodingKitten, Quantity: 4, Effects: []byte(`{"type": "explode"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440002", Code: cards.Defuse, Quantity: 6, Effects: []byte(`{"type": "prevent_explode"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440003", Code: cards.Attack, Quantity: 4, Effects: []byte(`{"type": "skip_turn_and_attack"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440004", Code: cards.Nope, Quantity: 5, Effects: []byte(`{"type": "cancel_action"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440005", Code: cards.SeeTheFuture, Quantity: 5, Effects: []byte(`{"type": "peek_cards"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440006", Code: cards.Shuffle, Quantity: 4, Effects: []byte(`{"type": "shuffle_deck"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440007", Code: cards.Skip, Quantity: 4, Effects: []byte(`{"type": "skip_turn"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440008", Code: cards.Favor, Quantity: 4, Effects: []byte(`{"type": "steal_card"}`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440009", Code: cards.BeardCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440010", Code: cards.Catermelon, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440011", Code: cards.HairyPotatoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440011", Code: cards.TacoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		{CardId: "123e4567-e89b-12d3-a456-426655440012", Code: cards.RainbowRalphingCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
	}

	cardsMap := map[string]*dataProviderProto.Card{
		"123e4567-e89b-12d3-a456-426655440001": {CardId: "123e4567-e89b-12d3-a456-426655440001", Code: cards.ExplodingKitten, Quantity: 4, Effects: []byte(`{"type": "explode"}`)},
		"123e4567-e89b-12d3-a456-426655440002": {CardId: "123e4567-e89b-12d3-a456-426655440002", Code: cards.Defuse, Quantity: 6, Effects: []byte(`{"type": "prevent_explode"}`)},
		"123e4567-e89b-12d3-a456-426655440003": {CardId: "123e4567-e89b-12d3-a456-426655440003", Code: cards.Attack, Quantity: 4, Effects: []byte(`{"type": "skip_turn_and_attack"}`)},
		"123e4567-e89b-12d3-a456-426655440004": {CardId: "123e4567-e89b-12d3-a456-426655440004", Code: cards.Nope, Quantity: 5, Effects: []byte(`{"type": "cancel_action"}`)},
		"123e4567-e89b-12d3-a456-426655440005": {CardId: "123e4567-e89b-12d3-a456-426655440005", Code: cards.SeeTheFuture, Quantity: 5, Effects: []byte(`{"type": "peek_cards"}`)},
		"123e4567-e89b-12d3-a456-426655440006": {CardId: "123e4567-e89b-12d3-a456-426655440006", Code: cards.Shuffle, Quantity: 4, Effects: []byte(`{"type": "shuffle_deck"}`)},
		"123e4567-e89b-12d3-a456-426655440007": {CardId: "123e4567-e89b-12d3-a456-426655440007", Code: cards.Skip, Quantity: 4, Effects: []byte(`{"type": "skip_turn"}`)},
		"123e4567-e89b-12d3-a456-426655440008": {CardId: "123e4567-e89b-12d3-a456-426655440008", Code: cards.Favor, Quantity: 4, Effects: []byte(`{"type": "steal_card"}`)},
		"123e4567-e89b-12d3-a456-426655440009": {CardId: "123e4567-e89b-12d3-a456-426655440009", Code: cards.BeardCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		"123e4567-e89b-12d3-a456-426655440010": {CardId: "123e4567-e89b-12d3-a456-426655440010", Code: cards.Catermelon, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		"123e4567-e89b-12d3-a456-426655440011": {CardId: "123e4567-e89b-12d3-a456-426655440011", Code: cards.HairyPotatoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		"123e4567-e89b-12d3-a456-426655440012": {CardId: "123e4567-e89b-12d3-a456-426655440012", Code: cards.TacoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		"123e4567-e89b-12d3-a456-426655440013": {CardId: "123e4567-e89b-12d3-a456-426655440013", Code: cards.RainbowRalphingCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
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
		cards.BeardCat:           {CardId: "123e4567-e89b-12d3-a456-426655440009", Code: cards.BeardCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		cards.Catermelon:         {CardId: "123e4567-e89b-12d3-a456-426655440010", Code: cards.Catermelon, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		cards.HairyPotatoCat:     {CardId: "123e4567-e89b-12d3-a456-426655440011", Code: cards.HairyPotatoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		cards.TacoCat:            {CardId: "123e4567-e89b-12d3-a456-426655440012", Code: cards.TacoCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
		cards.RainbowRalphingCat: {CardId: "123e4567-e89b-12d3-a456-426655440013", Code: cards.RainbowRalphingCat, Quantity: 4, ComboEffects: []byte(`[{"type": "steal_random_card", "required_cards": 2}, {"type": "steal_named_card", "required_cards": 3}]`)},
	}

	return cardsData, cardsMap, cardsMapByCode
}
