package kittens_userserver

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/db"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/userserver/go/grpcconnect"

	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/services/user"
	"github.com/sweetloveinyourheart/exploding-kittens/services/user/actions"
	"github.com/sweetloveinyourheart/exploding-kittens/services/user/repos"
)

const DEFAULT_USERSERVER_GRPC_PORT = 50052

const serviceType = "userserver"
const dbTablePrefix = "kittens_userserver"
const defDBName = "kittens_userserver"
const envPrefix = "USERSERVER"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var userServerCommand = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", serviceType),
		Short: fmt.Sprintf("Run as %s service", serviceType),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := cmdutil.BoilerplateRun(serviceType)
			if err != nil {
				log.GlobalSugared().Fatal(err)
			}

			app.Migrations(user.FS, dbTablePrefix)

			if err := setupDependencies(); err != nil {
				log.GlobalSugared().Fatal(err)
			}

			user.InitializeRepos(app.Ctx())

			signingKey := config.Instance().GetString("userserver.secrets.token_signing_key")
			actions := actions.NewActions(app.Ctx(), signingKey)

			path, handler := grpcconnect.NewUserServerHandler(
				actions,
			)
			go grpc.ServeBuf(app.Ctx(), path, handler, config.Instance().GetUint64("userserver.grpc.port"), serviceType)

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
	config.Int64Default(userServerCommand, "userserver.grpc.port", "grpc-port", DEFAULT_USERSERVER_GRPC_PORT, "GRPC Port to listen on", "USERSERVER_GRPC_PORT")

	cmdutil.BoilerplateFlagsCore(userServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(userServerCommand, serviceType)
	cmdutil.BoilerplateFlagsDB(userServerCommand, serviceType, envPrefix)

	return userServerCommand
}

func setupDependencies() error {
	dbConn, err := db.NewDbWithWait(config.Instance().GetString("userserver.db.url"), db.DBOptions{
		TimeoutSec:      config.Instance().GetInt("userserver.db.postgres.timeout"),
		MaxOpenConns:    config.Instance().GetInt("userserver.db.postgres.max_open_connections"),
		MaxIdleConns:    config.Instance().GetInt("userserver.db.postgres.max_idle_connections"),
		ConnMaxLifetime: config.Instance().GetInt("userserver.db.postgres.max_lifetime"),
		ConnMaxIdleTime: config.Instance().GetInt("userserver.db.postgres.max_idletime"),
		EnableTracing:   config.Instance().GetBool("userserver.db.tracing"),
	})
	if err != nil {
		return err
	}

	do.Provide[*pgxpool.Pool](nil, func(i *do.Injector) (*pgxpool.Pool, error) {
		return dbConn, nil
	})

	userRepo := repos.NewUserRepository(dbConn)
	do.Provide[repos.IUserRepository](nil, func(i *do.Injector) (repos.IUserRepository, error) {
		return userRepo, nil
	})

	userCredentialRepo := repos.NewUserCredentialRepository(dbConn)
	do.Provide[repos.IUserCredentialRepository](nil, func(i *do.Injector) (repos.IUserCredentialRepository, error) {
		return userCredentialRepo, nil
	})

	userSessionRepo := repos.NewUserSessionRepository(dbConn)
	do.Provide[repos.IUserSessionRepository](nil, func(i *do.Injector) (repos.IUserSessionRepository, error) {
		return userSessionRepo, nil
	})

	return nil
}
