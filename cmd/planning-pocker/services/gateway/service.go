package pocker_gateway

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/planning-pocker/pkg/cmdutil"
	"github.com/sweetloveinyourheart/planning-pocker/pkg/config"
	log "github.com/sweetloveinyourheart/planning-pocker/pkg/logger"
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
	config.Int64Default(gatewayCmd, "userserver.http.port", "http-port", DEFAULT_GATEWAY_HTTP_PORT, "HTTP Port to listen on", "API_GATEWAY_HTTP_PORT")

	cmdutil.BoilerplateFlagsCore(gatewayCmd, serviceType, envPrefix)

	return gatewayCmd
}
