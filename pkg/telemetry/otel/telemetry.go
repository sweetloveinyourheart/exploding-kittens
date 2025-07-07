package otel

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/version"
)

func otelTracerProcessor(url string) (tracesdk.SpanProcessor, error) {
	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	timeout := time.Second * 15
	if strings.Contains(url, "localhost") {
		timeout = time.Second * 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, url,
		// Insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to create gRPC connection to collector: %w", err))
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to create trace exporter: %w", err))
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := tracesdk.NewBatchSpanProcessor(traceExporter)
	return bsp, nil
}

func StartTracer(ctx context.Context, serviceName, otelURL string) error {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = fmt.Sprintf("unknown-%s", uuid.Must(uuid.NewV7()))
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(serviceName),
			semconv.ServiceNamespace(config.Instance().GetString(config.ServerNamespace)),
			semconv.ServiceInstanceID(fmt.Sprintf("%s-%s", config.Instance().GetString(config.ServerId), hostName)),
			semconv.ServiceVersion(version.GetVersion()),
		),
		resource.WithContainerID(),
		resource.WithFromEnv(),   // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
	)
	if err != nil {
		return errors.WithStack(fmt.Errorf("failed to create resource: %w", err))
	}

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(TraceContext{})

	if stringsutil.IsBlank(otelURL) {
		log.Global().WarnContext(ctx, "unable to start telemetry", zap.Error(errors.WithStack(fmt.Errorf("no telemetry url provided"))))
		return nil
	}

	go func() {
		var otelExp tracesdk.SpanProcessor

		if !stringsutil.IsBlank(otelURL) {
			otelExp, err = otelTracerProcessor(otelURL)
			if err != nil {
				log.Global().ErrorContext(ctx, "failed to create otel exporter", zap.Error(err))
				otelExp = nil
			}
		}

		if otelExp == nil {
			return
		}

		opts := []tracesdk.TracerProviderOption{
			tracesdk.WithSampler(tracesdk.AlwaysSample()),
			tracesdk.WithResource(res),
		}

		opts = append(opts, tracesdk.WithSpanProcessor(otelExp))

		tracerProvider := tracesdk.NewTracerProvider(
			opts...,
		)

		go func() {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
				log.Global().ErrorContext(ctx, "failed to shutdown tracer provider", zap.Error(err))
			}
		}()
		otel.SetTracerProvider(tracerProvider)
		log.Global().InfoContext(ctx, "telemetry started")
	}()
	return nil
}
