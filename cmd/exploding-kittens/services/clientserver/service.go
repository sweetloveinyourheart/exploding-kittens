package kittens_clientserver

import (
	"fmt"
	"net/http"

	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/grpc"
	"github.com/sweetloveinyourheart/exploding-kittens/proto/code/clientserver/go/grpcconnect"

	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
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

			client.InitializeRepos(app.Ctx())

			signingKey := config.Instance().GetString("clientserver.secrets.token_signing_key")
			actions := actions.NewActions(app.Ctx(), signingKey)

			path, handler := grpcconnect.NewClientServerHandler(
				actions,
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

	cmdutil.BoilerplateFlagsCore(clientServerCommand, serviceType, envPrefix)
	cmdutil.BoilerplateSecureFlags(clientServerCommand, serviceType)

	return clientServerCommand
}

func setupDependencies() error {
	userServerClient := userServerConnect.NewUserServerClient(
		http.DefaultClient,
		config.Instance().GetString("clientserver.userserver.url"),
	)
	do.Provide[userServerConnect.UserServerClient](nil, func(i *do.Injector) (userServerConnect.UserServerClient, error) {
		return userServerClient, nil
	})

	return nil
}
