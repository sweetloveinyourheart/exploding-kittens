package cmdutil

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	otelPrometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"

	_ "golang.org/x/tools/go/packages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/version"
)

const HealthCheckPortGRPC = 5051
const HealthCheckPortHTTP = 5052

func WaitForNatsConnection(wait time.Duration, connection *nats.Conn) error {
	timeout := time.Now().Add(wait)
	for time.Now().Before(timeout) {
		if connection.IsConnected() {
			return nil
		}
		time.Sleep(25 * time.Millisecond)
	}
	return errors.New("timeout waiting for nats connection")
}

func StartHealthServices(ctx context.Context, serviceName string, grpcPort int, webPort int) chan bool {
	readyHTTP := make(chan bool)
	readyGRPC := make(chan bool)
	ready := make(chan bool)
	startGRPCHealth(ctx, serviceName, grpcPort, readyHTTP)
	startHTTPHealth(ctx, serviceName, webPort, readyGRPC, ready)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case isReady := <-ready:
				readyGRPC <- isReady
				readyHTTP <- isReady
			}
		}
	}()

	return ready
}

func startGRPCHealth(ctx context.Context, serviceName string, grpcPort int, ready chan bool) {
	log.Global().InfoContext(ctx, "GRPCHealth: binding to port", zap.Int("port", grpcPort))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))
	if err != nil {
		log.Global().FatalContext(ctx, "failed to listen", zap.Error(err))
	}

	srv := grpc.NewServer()
	server := health.NewServer()
	reflection.Register(srv)
	grpc_health_v1.RegisterHealthServer(srv, server)
	server.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_UNKNOWN)

	go func() {
		log.Global().InfoContext(ctx, fmt.Sprintf("starting grpc health %s server", serviceName), zap.Int("port", grpcPort))
		if err := srv.Serve(listener); err != nil {
			log.Global().FatalContext(ctx, "failed to serve", zap.Error(err))
		}
	}()

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case isReady := <-ready:
				if isReady {
					server.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)
				} else {
					server.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
				}
			}
		}
	}()
}

func startHTTPHealth(ctx context.Context, serviceName string, webPort int, ready chan bool, readySet chan bool) {
	log.Global().InfoContext(ctx, "HTTPHealth: binding to port", zap.Int("port", webPort))

	srv := &healthServer{
		router:     mux.NewRouter(),
		healthy:    1,
		readyState: ready,
		readySet:   readySet,
	}

	srv.router.HandleFunc("/healthz", srv.healthzHandler).Methods("GET")
	srv.router.HandleFunc("/readyz", srv.readyzHandler).Methods("GET")
	srv.router.HandleFunc("/readyz/enable", srv.enableReadyHandler).Methods("POST")
	srv.router.HandleFunc("/readyz/disable", srv.disableReadyHandler).Methods("POST")

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%v", webPort),
		Handler:           srv.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Global().InfoContext(ctx, fmt.Sprintf("starting HTTP health %s server", serviceName), zap.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Global().FatalContext(ctx, "HTTP health server stopped", zap.Error(err))
		}
	}()

	go func() {
		<-ctx.Done()
		_ = httpServer.Shutdown(ctx)
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case isReady := <-ready:
				if isReady {
					atomic.StoreInt32(&srv.ready, 1)
				} else {
					atomic.StoreInt32(&srv.ready, 0)
				}
			}
		}
	}()
}

type healthServer struct {
	router     *mux.Router
	healthy    int32
	ready      int32
	readyState chan bool
	readySet   chan bool
}

// Healthz godoc
// @Summary Liveness check
// @Description used by Kubernetes liveness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /healthz [get]
// @Success 200 {string} string "OK"
func (s *healthServer) healthzHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&s.healthy) == 1 {
		s.JSONResponse(w, r, map[string]string{"status": "OK"})
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// Readyz godoc
// @Summary Readiness check
// @Description used by Kubernetes readiness probe
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz [get]
// @Success 200 {string} string "OK"
func (s *healthServer) readyzHandler(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&s.ready) == 1 {
		s.JSONResponse(w, r, map[string]string{"status": "OK"})
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

// EnableReady godoc
// @Summary Enable ready state
// @Description signals the Kubernetes LB that this instance is ready to receive traffic
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz/enable [post]
// @Success 202 {string} string "OK"
func (s *healthServer) enableReadyHandler(w http.ResponseWriter, r *http.Request) {
	s.readySet <- true
	w.WriteHeader(http.StatusAccepted)
}

// DisableReady godoc
// @Summary Disable ready state
// @Description signals the Kubernetes LB to stop sending requests to this instance
// @Tags Kubernetes
// @Accept json
// @Produce json
// @Router /readyz/disable [post]
// @Success 202 {string} string "OK"
func (s *healthServer) disableReadyHandler(w http.ResponseWriter, r *http.Request) {
	s.readySet <- false
	w.WriteHeader(http.StatusAccepted)
}

func (s *healthServer) JSONResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Global().Error("failed to marshal response", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(prettyJSON(body))
}

func prettyJSON(b []byte) []byte {
	var out bytes.Buffer
	_ = json.Indent(&out, b, "", "  ")
	return out.Bytes()
}

type Initializer interface {
	Initialize()
}

func StartMetricServer(ctx context.Context, serviceName string, serviceID string, port int, metricsProvider ...Initializer) {
	SetAppInfoMetrics(metricsProvider...)

	exporter, err := otelPrometheus.New()
	if err != nil {
		log.Global().FatalContext(ctx, "failed to create the Prometheus exporter", zap.Error(err))
	}
	hostName, err := os.Hostname()
	if err != nil {
		hostName = fmt.Sprintf("unknown-%s", uuid.Must(uuid.NewV7()))
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.HostName(hostName),
			semconv.ServiceName(serviceName),
			semconv.ServiceNamespace(config.Instance().GetString(config.ServerNamespace)),
			semconv.ServiceInstanceID(fmt.Sprintf("%s-%s", serviceID, hostName)),
			semconv.ServiceVersion(version.GetVersion()),
		),
		resource.WithContainerID(),
		resource.WithFromEnv(),   // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		resource.WithHost(),      // This option configures a set of Detectors that discover host information
	)
	if err != nil {
		log.Global().FatalContext(ctx, "failed to create the Prometheus exporter", zap.Error(err))
	}
	provider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	startHttpServer := func(port int) *http.Server {
		srv := &http.Server{
			Addr:              fmt.Sprintf(":%v", port),
			ReadHeaderTimeout: 5 * time.Second,
		}
		serveMux := http.NewServeMux()
		serveMux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}))
		srv.Handler = serveMux

		go func() {
			// always returns error. ErrServerClosed on graceful close
			log.Global().InfoContext(ctx, "starting HTTP metric server", zap.String("addr", fmt.Sprintf(":%v", port)))
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Global().FatalContext(ctx, "HTTP metric server stopped", zap.Error(err))
			}
		}()

		// returning reference so caller can call Shutdown()
		return srv
	}

	srv := startHttpServer(port)
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx) // gracefully shutdown the server, waiting max 30 seconds for current operations to complete
		if err != nil {
			log.Global().Error("HTTP metric server shutdown failed", zap.Error(err))
		}
	}()
}

func SetAppInfoMetrics(metricsProvider ...Initializer) {
	for _, provider := range metricsProvider {
		provider.Initialize()
	}
}
