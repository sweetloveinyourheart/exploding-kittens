package cmdutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	fields "github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
)

const rootCmdName = "app"
const defaultShortDescription = "Unified EXPLODING kittens service launcher"

var (
	cfgFile        string
	shortCircuit   bool
	ServiceRootCmd = &cobra.Command{
		Use:   rootCmdName,
		Short: defaultShortDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.CalledAs() != rootCmdName {
				viper.Set("service", cmd.CalledAs())
				if err := viper.ReadInConfig(); err != nil {
					if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
						panic(fmt.Errorf("fatal error config file: %w", err))
					}
				}

				fields.ResolveRequireFlags(cmd)

				shortCircuit = true
			}

			service := viper.GetString("service")
			viper.SetConfigName(service)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if shortCircuit {
				return nil
			}

			// This is some trickery to allow us to determine which service to run if supplied via flag or environment variable
			// instead of via command
			service := viper.GetString("service")
			for _, subCmd := range cmd.Commands() {
				if subCmd.Name() == service {
					shortCircuit = true

					// Inject the command into the OS args for cobra to pick up
					os.Args = append(os.Args[0:1], append([]string{service}, os.Args[1:]...)...)
					return subCmd.Execute()
				}
			}

			return fmt.Errorf("no valid command or service specified")
		},
	}
)

func InitializeService(command ...*cobra.Command) {
	cobra.OnInitialize(initConfig)
	ServiceRootCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{
		UnknownFlags: true,
	}

	// Common
	ServiceRootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.EXPLODING-poker/app.yaml)")
	ServiceRootCmd.PersistentFlags().StringP("service", "s", "", "which service to run")
	ServiceRootCmd.PersistentFlags().String("log-level", "info", "log level to use")

	fields.BindWithDefault(ServiceRootCmd.PersistentFlags().Lookup("service"), "service", "", "EXPLODING_KITTENS_SERVICE")
	fields.BindWithDefault(ServiceRootCmd.PersistentFlags().Lookup("log-level"), "log.level", "info", "LOG_LEVEL")

	// Health check
	ServiceRootCmd.PersistentFlags().Int64("healthcheck-port", HealthCheckPortGRPC, "Port to listen on for services that support a health check")
	ServiceRootCmd.PersistentFlags().Int64("healthcheck-web-port", HealthCheckPortHTTP, "Port to listen on for services that support a health check")
	ServiceRootCmd.PersistentFlags().String("healthcheck-host", "localhost", "Host to listen on for services that support a health check")

	fields.BindWithDefault(ServiceRootCmd.PersistentFlags().Lookup("healthcheck-port"), "healthcheck.port", HealthCheckPortGRPC, "EXPLODING_KITTENS_HEALTHCHECK_PORT")
	fields.BindWithDefault(ServiceRootCmd.PersistentFlags().Lookup("healthcheck-web-port"), "healthcheck.web.port", HealthCheckPortHTTP, "EXPLODING_KITTENS_HEALTHCHECK_WEB_PORT")
	fields.BindWithDefault(ServiceRootCmd.PersistentFlags().Lookup("healthcheck-host"), "healthcheck.host", "localhost", "EXPLODING_KITTENS_HEALTHCHECK_HOST")

	// Monitoring, Logging and Tracing
	ServiceRootCmd.PersistentFlags().String("otel-url", "localhost:30080", "URL to send OpenTelemetry data to")
	fields.BindWithDefault(ServiceRootCmd.PersistentFlags().Lookup("otel-url"), "otel.url", "localhost:30080", "EXPLODING_KITTENS_OTEL_URL")

	for _, c := range command {
		ServiceRootCmd.AddCommand(c)
	}

	if len(os.Args) > 1 && os.Args[1] == "generate" {
		generateDocs(ServiceRootCmd)
		generateSchema(ServiceRootCmd)
		return
	}

	if err := ServiceRootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("EXPLODING_KITTENS")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("kittens")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/exploding-kittens")
		viper.AddConfigPath("$HOME/.exploding-kittens")
		viper.AddConfigPath("./cmd/exploding-kittens")
		viper.AddConfigPath(".")
	}
}
