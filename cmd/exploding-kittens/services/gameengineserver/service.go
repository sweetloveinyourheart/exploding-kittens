package gameengineserver

import (
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go"
	pool "github.com/octu0/nats-pool"
	"go.uber.org/zap"

	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	dataProviderConnect "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go/grpcconnect"
	gameengine "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/actions"
)

const DEFAULT_GAMEENGINESERVER_GRPC_PORT = 50054

const serviceType = "gameengineserver"
const envPrefix = "GAMEENGINESERVER"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var gameEngineServerCommand = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", serviceType),
		Short: fmt.Sprintf("Run as %s service", serviceType),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := cmdutil.BoilerplateRun(serviceType)
			if err != nil {
				log.GlobalSugared().Fatal(err)
			}

			if err := setupDependencies(); err != nil {
				log.GlobalSugared().Fatal(err)
			}

			if err := gameengine.InitializeRepos(app.Ctx()); err != nil {
				log.GlobalSugared().Fatal(err)
			}

			signingKey := config.Instance().GetString("gameengineserver.secrets.token_signing_key")
			actions := actions.NewActions(app.Ctx(), signingKey)

			opt := connect.WithInterceptors(
				interceptors.CommonConnectInterceptors(
					serviceType,
					signingKey,
					interceptors.ConnectServerAuthHandler(signingKey),
				)...,
			)
			path, handler := grpcconnect.NewGameEngineServerHandler(
				actions,
				opt,
			)
			go grpc.ServeBuf(app.Ctx(), path, handler, config.Instance().GetUint64("gameengineserver.grpc.port"), serviceType)

			app.Run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.Instance().Set("service_prefix", serviceType)

			cmdutil.BoilerplateMetaConfig(serviceType)

			config.RegisterService(cmd, config.Service{
				Command: cmd,
			})
			config.AddDefaultServicePorts(cmd, rootCmd)
			return nil
		},
	}

	// config options
	config.Int64Default(gameEngineServerCommand, "gameengineserver.grpc.port", "grpc-port", DEFAULT_GAMEENGINESERVER_GRPC_PORT, "GRPC Port to listen on", "GAMEENGINESERVER_GRPC_PORT")
	config.StringDefault(gameEngineServerCommand, "gameengineserver.dataprovider.url", "dataprovider-url", "http://dataprovider:50055", "Data provider connection URL", "GAMEENGINESERVER_DATAPROVIDER_URL")

	cmdutil.BoilerplateFlagsCore(gameEngineServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsNats(gameEngineServerCommand, serviceType, envPrefix)

	return gameEngineServerCommand
}

func setupDependencies() error {
	timeout := 2 * time.Second

	signingKey := config.Instance().GetString("gameengineserver.secrets.token_signing_key")

	dataProviderClient := dataProviderConnect.NewDataProviderClient(
		http.DefaultClient,
		config.Instance().GetString("gameengineserver.dataprovider.url"),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			signingKey,
		)...),
	)

	busConnection, err := nats.Connect(config.Instance().GetString("gameengineserver.nats.url"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.Name("kittens/gameengineserver/1.0/single"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.String("type", "nats"), zap.Error(err))
		}))

	if err != nil {
		return errors.WithStack(errors.Wrap(err, "failed to connect to nats"))
	}

	if err := cmdutil.WaitForNatsConnection(timeout, busConnection); err != nil {
		return errors.WithStack(errors.Wrap(err, "failed to connect to nats"))
	}

	do.Provide[dataProviderConnect.DataProviderClient](nil, func(i *do.Injector) (dataProviderConnect.DataProviderClient, error) {
		return dataProviderClient, nil
	})

	connPool := pool.New(100, config.Instance().GetString("gameengineserver.nats.url"),
		nats.NoEcho(),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.Name("kittens/gameengineserver/1.0"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.String("type", "nats"), zap.Error(err))
		}),
	)

	do.ProvideNamed[*pool.ConnPool](nil, string(constants.ConnectionPool),
		func(i *do.Injector) (*pool.ConnPool, error) {
			return connPool, nil
		})

	do.ProvideNamed[*nats.Conn](nil, fmt.Sprintf("%s-conn", string(constants.Bus)),
		func(i *do.Injector) (*nats.Conn, error) {
			return busConnection, nil
		})

	return nil
}
