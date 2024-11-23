package cmdutil

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/sweetloveinyourheart/planning-poker/pkg/config"
	log "github.com/sweetloveinyourheart/planning-poker/pkg/logger"
)

type AppRun struct {
	serviceType string
	serviceKey  string
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	readyChan   chan bool
}

func st(s string) (serviceType string, serviceKey string) {
	serviceType = strings.ToLower(s)
	serviceKey = strings.ReplaceAll(serviceType, "_", "-")
	serviceKey = strings.ReplaceAll(serviceKey, "-", ".")
	return
}

func BoilerplateRun(serviceType string) (*AppRun, error) {
	serviceType, serviceKey := st(serviceType)

	ctx, cancel := context.WithCancel(context.Background())

	// Logger init.
	logLevel := config.Instance().GetString("log.level")
	logger := log.New().SetStringLevel(logLevel)

	serviceName := config.Instance().GetString("service")
	logger.Infow("starting service",
		"type", serviceType,
		"service", serviceName,
	)

	readyChan := StartHealthServices(ctx, serviceName, config.Instance().GetInt("healthcheck.port"), config.Instance().GetInt("healthcheck.web.port"))

	return &AppRun{
		serviceType: serviceType,
		serviceKey:  serviceKey,
		wg:          sync.WaitGroup{},
		ctx:         ctx,
		cancel:      cancel,
		readyChan:   readyChan,
	}, nil
}

func BoilerplateFlagsCore(command *cobra.Command, serviceType string, envPrefix string) {
	_, serviceKey := st(serviceType)
	envPrefix = strings.ToUpper(envPrefix)

	config.String(command, fmt.Sprintf("%s.id", serviceKey), "id", "Unique identifier for this services", fmt.Sprintf("%s_ID", envPrefix))
	config.String(command, fmt.Sprintf("%s.secrets.token_signing_key", serviceKey), "token-signing-key", "Signing key used for service to service tokens", fmt.Sprintf("%s_SECRETS_TOKEN_SIGNING_KEY", envPrefix))

	_ = command.MarkPersistentFlagRequired("id")
	_ = command.MarkPersistentFlagRequired("token-signing-key")
}

func BoilerplateMetaConfig(serviceType string) {
	_, serviceKey := st(serviceType)

	config.Instance().Set(config.ServerId, config.Instance().GetString(fmt.Sprintf("%s.id", serviceKey)))
	config.Instance().Set(config.ServerReplicaCount, config.Instance().GetInt64(fmt.Sprintf("%s.replicas", serviceKey)))
	config.Instance().Set(config.ServerReplicaNumber, config.Instance().GetInt64(fmt.Sprintf("%s.replica_num", serviceKey)))
}

func BoilerplateSecureFlags(command *cobra.Command, serviceType string) {
	_, serviceKey := st(serviceType)

	config.SecureFields(fmt.Sprintf("%s.db.url", serviceKey),
		fmt.Sprintf("%s.db.read.url", serviceKey),
		fmt.Sprintf("%s.db.migrations.url", serviceKey),
		fmt.Sprintf("%s.secrets.token_signing_key", serviceKey),
		fmt.Sprintf("%s.oci.plugins.registry.password", serviceKey),
	)
}

func (a *AppRun) Ctx() context.Context {
	return a.ctx
}

func (a *AppRun) Cancel() {
	a.cancel()
}

func (a *AppRun) Ready() {
	a.readyChan <- true
}

func (a *AppRun) Run() {
	signalMonitor := make(chan os.Signal, 1)
	signal.Notify(signalMonitor, os.Interrupt, syscall.SIGTERM)

	a.wg.Add(1)
	go func() {
		<-signalMonitor
		a.cancel()
		a.wg.Done()
	}()

	a.readyChan <- true

	// Wait for a signal or other termination event
	a.wg.Wait()
}
