package game_test

import (
	"context"
	goTesting "testing"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/testing"

	gameengine "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine"

	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/repos"
	mockCard "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/repos/mock"
)

type GameSuite struct {
	*testing.Suite
	deferred []func()

	mockCardRepository *mockCard.MockCardRepository
}

func (as *GameSuite) SetupTest() {
	as.mockCardRepository = new(mockCard.MockCardRepository)
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
		Suite:              testing.MakeSuite(t),
		mockCardRepository: new(mockCard.MockCardRepository),
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

	do.Override[repos.ICardRepository](nil, func(i *do.Injector) (repos.ICardRepository, error) {
		return gs.mockCardRepository, nil
	})

	do.OverrideNamed[*pool.ConnPool](nil, string(constants.ConnectionPool),
		func(i *do.Injector) (*pool.ConnPool, error) {
			return connPool, nil
		})

	config.Instance().Set("gameengineserver.id", uuid.Must(uuid.NewV7()).String())

	ctx, cancel := context.WithCancel(context.Background())
	gs.deferred = append(gs.deferred, cancel)
	err = gameengine.InitializeRepos(ctx)
	gs.NoError(err)
}
