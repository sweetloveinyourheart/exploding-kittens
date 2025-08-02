package dataprovider

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/db"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/interceptors"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/dataprovider/go/grpcconnect"
	dataprovider "github.com/sweetloveinyourheart/exploding-kittens/services/data_provider"
	"github.com/sweetloveinyourheart/exploding-kittens/services/data_provider/actions"
	"github.com/sweetloveinyourheart/exploding-kittens/services/data_provider/repos"
)

const DEFAULT_DATAPROVIDER_GRPC_PORT = 50055

const serviceType = "dataprovider"
const dbTablePrefix = "kittens_dataprovider"
const defDBName = "kittens_dataprovider"
const envPrefix = "DATAPROVIDER"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var dataProviderCommand = &cobra.Command{
		Use:   fmt.Sprintf("%s [flags]", serviceType),
		Short: fmt.Sprintf("Run as %s service", serviceType),
		Run: func(cmd *cobra.Command, args []string) {
			app, err := cmdutil.BoilerplateRun(serviceType)
			if err != nil {
				log.GlobalSugared().Fatal(err)
			}

			app.Migrations(dataprovider.FS, dbTablePrefix)

			if err := setupDependencies(); err != nil {
				log.GlobalSugared().Fatal(err)
			}

			if err := dataprovider.InitializeRepos(app.Ctx()); err != nil {
				log.GlobalSugared().Fatal(err)
			}

			redisURL := config.Instance().GetString("dataprovider.redis.url")
			address, dbIndex, err := stringsutil.RedisOptionsFromURL(redisURL)
			if err != nil {
				log.Global().Fatal("Unable to parse redis url", zap.Error(err))
			}
			redisClient := redis.NewClient(
				&redis.Options{
					Addr: address,
					DB:   dbIndex,
				},
			)

			// Enable tracing instrumentation.
			if err := redisotel.InstrumentTracing(redisClient); err != nil {
				panic(err)
			}

			// Enable metrics instrumentation.
			if err := redisotel.InstrumentMetrics(redisClient); err != nil {
				panic(err)
			}

			signingKey := config.Instance().GetString("dataprovider.secrets.token_signing_key")
			actions := actions.NewActions(app.Ctx(), signingKey, redisClient)

			opt := connect.WithInterceptors(
				interceptors.CommonConnectInterceptors(
					serviceType,
					signingKey,
					interceptors.ConnectServerAuthHandler(signingKey),
				)...,
			)
			path, handler := grpcconnect.NewDataProviderHandler(
				actions,
				opt,
			)
			go grpc.ServeBuf(app.Ctx(), path, handler, config.Instance().GetUint64("dataprovider.grpc.port"), serviceType)

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
	config.Int64Default(dataProviderCommand, "dataprovider.grpc.port", "grpc-port", DEFAULT_DATAPROVIDER_GRPC_PORT, "GRPC Port to listen on", "DATAPROVIDER_GRPC_PORT")

	cmdutil.BoilerplateFlagsCore(dataProviderCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(dataProviderCommand, serviceType)
	cmdutil.BoilerplateFlagsDB(dataProviderCommand, serviceType, envPrefix)
	cmdutil.BoilerplateFlagsRedisEdge(dataProviderCommand, serviceType, envPrefix)

	return dataProviderCommand
}

func setupDependencies() error {
	dbConn, err := db.NewDbWithWait(config.Instance().GetString("dataprovider.db.url"), db.DBOptions{
		TimeoutSec:      config.Instance().GetInt("dataprovider.db.postgres.timeout"),
		MaxOpenConns:    config.Instance().GetInt("dataprovider.db.postgres.max_open_connections"),
		MaxIdleConns:    config.Instance().GetInt("dataprovider.db.postgres.max_idle_connections"),
		ConnMaxLifetime: config.Instance().GetInt("dataprovider.db.postgres.max_lifetime"),
		ConnMaxIdleTime: config.Instance().GetInt("dataprovider.db.postgres.max_idletime"),
		EnableTracing:   config.Instance().GetBool("dataprovider.db.tracing"),
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
