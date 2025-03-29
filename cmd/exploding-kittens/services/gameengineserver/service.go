package gameengineserver

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/db"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/gameengineserver/go/grpcconnect"
	gameengine "github.com/sweetloveinyourheart/exploding-kittens/services/game_engine"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/actions"
	"github.com/sweetloveinyourheart/exploding-kittens/services/game_engine/repos"
)

const DEFAULT_GAMEENGINESERVER_GRPC_PORT = 50054

const serviceType = "gameengineserver"
const dbTablePrefix = "kittens_gameengineserver"
const defDBName = "kittens_gameengineserver"
const envPrefix = "GAMEENGINESERVER"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var userServerCommand = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", serviceType),
		Short: fmt.Sprintf("Run as %s service", serviceType),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := cmdutil.BoilerplateRun(serviceType)
			if err != nil {
				log.GlobalSugared().Fatal(err)
			}

			app.Migrations(gameengine.FS, dbTablePrefix)

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
			config.AddDefaultDatabase(cmd, defDBName)
			return nil
		},
	}

	// config options
	config.Int64Default(userServerCommand, "gameengineserver.grpc.port", "grpc-port", DEFAULT_GAMEENGINESERVER_GRPC_PORT, "GRPC Port to listen on", "GAMEENGINESERVER_GRPC_PORT")

	cmdutil.BoilerplateFlagsCore(userServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(userServerCommand, serviceType)
	cmdutil.BoilerplateFlagsDB(userServerCommand, serviceType, envPrefix)

	return userServerCommand
}

func setupDependencies() error {
	dbConn, err := db.NewDbWithWait(config.Instance().GetString("gameengineserver.db.url"), db.DBOptions{
		TimeoutSec:      config.Instance().GetInt("gameengineserver.db.postgres.timeout"),
		MaxOpenConns:    config.Instance().GetInt("gameengineserver.db.postgres.max_open_connections"),
		MaxIdleConns:    config.Instance().GetInt("gameengineserver.db.postgres.max_idle_connections"),
		ConnMaxLifetime: config.Instance().GetInt("gameengineserver.db.postgres.max_lifetime"),
		ConnMaxIdleTime: config.Instance().GetInt("gameengineserver.db.postgres.max_idletime"),
		EnableTracing:   config.Instance().GetBool("gameengineserver.db.tracing"),
	})
	if err != nil {
		return err
	}

	do.Provide[*pgxpool.Pool](nil, func(i *do.Injector) (*pgxpool.Pool, error) {
		return dbConn, nil
	})

	cardRepo := repos.NewCardRepository(dbConn)
	do.Provide[repos.ICardRepository](nil, func(i *do.Injector) (repos.ICardRepository, error) {
		return cardRepo, nil
	})

	return nil
}
