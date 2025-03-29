package kittens_lobbyserver

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	"github.com/sweetloveinyourheart/exploding-kittens/services/lobby"

	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
)

const DEFAULT_LOBBYSERVER_GRPC_PORT = 50053

const serviceType = "lobbyserver"
const envPrefix = "LOBBYSERVER"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var lobbyServerCommand = &cobra.Command{
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

			if err := lobby.InitializeRepos(app.Ctx()); err != nil {
				log.GlobalSugared().Fatal(err)
			}

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
	config.Int64Default(lobbyServerCommand, "lobbyserver.grpc.port", "grpc-port", DEFAULT_LOBBYSERVER_GRPC_PORT, "GRPC Port to listen on", "LOBBYSERVER_GRPC_PORT")

	cmdutil.BoilerplateFlagsCore(lobbyServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsNats(lobbyServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(lobbyServerCommand, serviceType)

	return lobbyServerCommand
}

func setupDependencies() error {
	timeout := 2 * time.Second

	busConnection, err := nats.Connect(config.Instance().GetString("lobbyserver.nats.url"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.Name("kittens/lobbyserver/1.0/single"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.String("type", "nats"), zap.Error(err))
		}))
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "failed to connect to  nats"))
	}

	if err := cmdutil.WaitForNatsConnection(timeout, busConnection); err != nil {
		return errors.WithStack(errors.Wrap(err, "failed to connect to  nats"))
	}

	connPool := pool.New(100, config.Instance().GetString("lobbyserver.nats.url"),
		nats.NoEcho(),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.Name("kittens/lobbyserver/1.0"),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Global().Error("nats error", zap.String("type", "nats"), zap.Error(err))
		}),
	)

	do.ProvideNamed[*pool.ConnPool](nil, string(constants.ConnectionPool),
		func(i *do.Injector) (*pool.ConnPool, error) {
			return connPool, nil
		})

	return nil
}
