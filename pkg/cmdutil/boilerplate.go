package cmdutil

import (
	"context"
	"strings"
	"sync"

	"github.com/sweetloveinyourheart/planning-poker/pkg/config"
	log "github.com/sweetloveinyourheart/planning-poker/pkg/logger"
)

type AppRun struct {
	serviceType string
	serviceKey  string
	runType     string
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

func BoilerplateRun(serviceType string, runType string) (*AppRun, error) {
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
		runType:     runType,
		wg:          sync.WaitGroup{},
		ctx:         ctx,
		cancel:      cancel,
		readyChan:   readyChan,
	}, nil
}
