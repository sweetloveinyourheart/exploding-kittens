package pocker_gateway

import (
	"fmt"
	"net/http"

	"github.com/samber/do"
	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/planning-pocker/pkg/cmdutil"
	"github.com/sweetloveinyourheart/planning-pocker/pkg/config"
	log "github.com/sweetloveinyourheart/planning-pocker/pkg/logger"
	"github.com/sweetloveinyourheart/planning-pocker/proto/code/userserver/go/grpcconnect"
)

const DEFAULT_GATEWAY_HTTP_PORT = 9000

const serviceType = "gateway"
const envPrefix = "API_GATEWAY"

func Command(rootCmd *cobra.Command) *cobra.Command {
	var gatewayCmd = &cobra.Command{
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
	config.Int64Default(gatewayCmd, "gateway.http.port", "http-port", DEFAULT_GATEWAY_HTTP_PORT, "HTTP Port to listen on", "API_GATEWAY_HTTP_PORT")
	config.StringDefault(gatewayCmd, "gateway.userserver.url", "userserver-url", "http://userserver:50051", "User Server connection URL", "GATEWAY_USERSERVER_URL")

	cmdutil.BoilerplateFlagsCore(gatewayCmd, serviceType, envPrefix)

	return gatewayCmd
}

func setupDependencies() error {
	userServerClient := grpcconnect.NewUserServerClient(
		http.DefaultClient,
		config.Instance().GetString("gateway.userserver.url"),
	)
	do.Provide[grpcconnect.UserServerClient](nil, func(i *do.Injector) (grpcconnect.UserServerClient, error) {
		return userServerClient, nil
	})

	return nil
}
