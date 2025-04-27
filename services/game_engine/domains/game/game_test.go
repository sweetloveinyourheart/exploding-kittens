package game_test

import (
	"context"
	"fmt"
	goTesting "testing"

	"connectrpc.com/connect"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	goMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants/cards"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/testing"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/mock"
	gameengine "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine"

	"go.uber.org/zap"

	dataProviderProto "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go"
	dataProviderGrpc "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
)

type GameSuite struct {
	*testing.Suite
	deferred []func()

	mockDataProviderClient *mock.MockDataProviderClient
}

func (as *GameSuite) SetupTest() {
	as.mockDataProviderClient = new(mock.MockDataProviderClient)
}

func (as *GameSuite) TearDownTest() {
	if len(as.deferred) > 0 {
		for _, f := range as.deferred {
			f()
		}
		as.deferred = nil
	}
}

func TestGameSuite(t *goTesting.T) {
	gs := &GameSuite{
		Suite:                  testing.MakeSuite(t),
		mockDataProviderClient: new(mock.MockDataProviderClient),
	}
	suite.Run(t, gs)
}

func (gs *GameSuite) setupEnvironment() {
	bus, shutdown := testing.StartLocalNATSServer(gs.T())
	gs.deferred = append(gs.deferred, shutdown)

	busAddress := "nats://" + bus.Addr().String() // use server address
	busConnection, err := nats.Connect(busAddress)
	gs.NoError(err)

	testing.NATSWaitConnected(gs.T(), busConnection) // wait connection if not connected yet

	gs.NoError(err)

	connPool := pool.New(100, busAddress,
		nats.NoEcho(),
		nats.Name("kittens/gameengineserver/1.0"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.Error(err))
		}),
	)

	do.Override[dataProviderGrpc.DataProviderClient](nil, func(i *do.Injector) (dataProviderGrpc.DataProviderClient, error) {
		return gs.mockDataProviderClient, nil
	})

	do.OverrideNamed[*pool.ConnPool](nil, string(constants.ConnectionPool),
		func(i *do.Injector) (*pool.ConnPool, error) {
			return connPool, nil
		})

	do.ProvideNamed[*nats.Conn](nil, fmt.Sprintf("%s-conn", string(constants.Bus)),
		func(i *do.Injector) (*nats.Conn, error) {
			return busConnection, nil
		})

	config.Instance().Set("gameengineserver.id", uuid.Must(uuid.NewV7()).String())

	ctx, cancel := context.WithCancel(context.Background())
	gs.deferred = append(gs.deferred, cancel)
	err = gameengine.InitializeRepos(ctx)
	gs.NoError(err)
}

func (gs *GameSuite) prepareCards() []*dataProviderProto.Card {
	cards := []*dataProviderProto.Card{
		{CardId: "123e4567-e89b-12d3-a456-426655440001", Code: cards.ExplodingKitten, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440002", Code: cards.Defuse, Quantity: 6},
		{CardId: "123e4567-e89b-12d3-a456-426655440003", Code: cards.Attack, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440004", Code: cards.Nope, Quantity: 5},
		{CardId: "123e4567-e89b-12d3-a456-426655440005", Code: cards.SeeTheFuture, Quantity: 5},
		{CardId: "123e4567-e89b-12d3-a456-426655440006", Code: cards.Shuffle, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440007", Code: cards.Skip, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440008", Code: cards.Favor, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440009", Code: cards.BeardCat, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440010", Code: cards.Catermelon, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440011", Code: cards.HairyPotatoCat, Quantity: 4},
		{CardId: "123e4567-e89b-12d3-a456-426655440012", Code: cards.RainbowRalphingCat, Quantity: 4},
	}

	gs.mockDataProviderClient.On("GetCards", goMock.Anything, goMock.Anything).Return(connect.NewResponse(&dataProviderProto.GetCardsResponse{
		Cards: cards,
	}), nil)

	return cards
}
