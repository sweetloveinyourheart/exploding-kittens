package kittens_clientserver

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go/grpcconnect"

	auth_interceptors "github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors/auth"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	dataProviderConnect "github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
	gameEngineServerConnect "github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go/grpcconnect"
	userServerConnect "github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go/grpcconnect"

	"github.com/sweetloveinyourheart/exploding-kittens/services/client"
	"github.com/sweetloveinyourheart/exploding-kittens/services/client/actions"
)

const DEFAULT_CLIENTSERVER_GRPC_PORT = 50051

const serviceType = "clientserver"
const envPrefix = "CLIENTSERVER"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var clientServerCommand = &cobra.Command{
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

			if err := client.InitializeRepos(app.Ctx()); err != nil {
				log.GlobalSugared().Fatal(err)
			}

			signingKey := config.Instance().GetString("clientserver.secrets.token_signing_key")
			actions := actions.NewActions(app.Ctx(), signingKey)

			opt := connect.WithInterceptors(
				interceptors.CommonConnectInterceptors(
					serviceType,
					signingKey,
					interceptors.ConnectAuthHandler(signingKey),
					auth_interceptors.WithOverride(actions),
				)...,
			)
			path, handler := grpcconnect.NewClientServerHandler(
				actions,
				opt,
			)
			go grpc.ServeBuf(app.Ctx(), path, handler, config.Instance().GetUint64("clientserver.grpc.port"), serviceType)

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

			config.AddServicePort(cmd, cmd.PersistentFlags().Lookup("grpc-port"), config.ServicePort{
				Name:         serviceType,
				WireProtocol: "tcp",
				Protocol:     "grpc",
				Public:       true,
			})

			config.AddDefaultServicePorts(cmd, rootCmd)
			return nil
		},
	}

	// config options
	config.Int64Default(clientServerCommand, "clientserver.grpc.port", "grpc-port", DEFAULT_CLIENTSERVER_GRPC_PORT, "GRPC Port to listen on", "CLIENTSERVER_GRPC_PORT")
	config.StringDefault(clientServerCommand, "clientserver.userserver.url", "userserver-url", "http://userserver:50052", "Userserver connection URL", "CLIENTSERVER_USERSERVER_URL")
	config.StringDefault(clientServerCommand, "clientserver.gameengineserver.url", "gameengineserver-url", "http://gameengineserver:50054", "Game Engine Server connection URL", "CLIENTSERVER_GAMEENGINESERVER_URL")
	config.StringDefault(clientServerCommand, "clientserver.dataprovider.url", "dataprovider-url", "http://dataprovider:50055", "Dataprovider connection URL", "CLIENTSERVER_DATAPROVIDER_URL")

	cmdutil.BoilerplateFlagsCore(clientServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsNats(clientServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(clientServerCommand, serviceType)

	return clientServerCommand
}

func setupDependencies() error {
	signingKey := config.Instance().GetString("clientserver.secrets.token_signing_key")

	connPool := pool.New(100, config.Instance().GetString("clientserver.nats.url"),
		nats.NoEcho(),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.Name("kittens/clientserver/1.0"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.String("type", "nats"), zap.Error(err))
		}),
	)

	busConnection, err := nats.Connect(config.Instance().GetString("clientserver.nats.url"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.Name("kittens/clientserver/1.0/single"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.String("type", "nats"), zap.Error(err))
		}))
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "failed to connect to nats"))
	}

	userServerClient := userServerConnect.NewUserServerClient(
		http.DefaultClient,
		config.Instance().GetString("clientserver.userserver.url"),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			signingKey,
		)...),
	)

	dataProviderClient := dataProviderConnect.NewDataProviderClient(
		http.DefaultClient,
		config.Instance().GetString("clientserver.dataprovider.url"),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			signingKey,
		)...),
	)

	gameEngineServerClient := gameEngineServerConnect.NewGameEngineServerClient(
		http.DefaultClient,
		config.Instance().GetString("clientserver.gameengineserver.url"),
		connect.WithInterceptors(interceptors.CommonConnectClientInterceptors(
			serviceType,
			signingKey,
		)...),
	)

	do.Provide[userServerConnect.UserServerClient](nil, func(i *do.Injector) (userServerConnect.UserServerClient, error) {
		return userServerClient, nil
	})

	do.Provide[dataProviderConnect.DataProviderClient](nil, func(i *do.Injector) (dataProviderConnect.DataProviderClient, error) {
		return dataProviderClient, nil
	})

	do.Provide[gameEngineServerConnect.GameEngineServerClient](nil, func(i *do.Injector) (gameEngineServerConnect.GameEngineServerClient, error) {
		return gameEngineServerClient, nil
	})

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
