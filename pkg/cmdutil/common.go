package cmdutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	log "github.com/sweetloveinyourheart/planning-poker/pkg/logger"
)

const HealthCheckPortGRPC = 5051
const HealthCheckPortHTTP = 5052

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
