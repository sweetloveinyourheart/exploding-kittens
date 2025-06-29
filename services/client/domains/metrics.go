package domains

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/do"
	"go.opentelemetry.io/otel/trace"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
)

type ClientMetrics struct {
	TotalLobbies *prometheus.CounterVec
}

var _ cmdutil.Initializer = (*ClientMetrics)(nil)

func (m *ClientMetrics) Initialize() {
	m.TotalLobbies = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_lobbies_started",
			Help: "Total number of lobbies",
		},
		[]string{"lobby_id"},
	)

	prometheus.MustRegister(m.TotalLobbies)

	do.Provide[*ClientMetrics](nil, func(i *do.Injector) (*ClientMetrics, error) {
		return m, nil
	})
}

func (m *ClientMetrics) TotalLobbiesInc(ctx context.Context, LobbyID uuid.UUID, amount int) {
	if m == nil {
		return
	}

	traceID := ""
	span := trace.SpanFromContext(ctx)
	if span != nil && span.SpanContext().IsValid() && span.SpanContext().HasTraceID() {
		traceID = span.SpanContext().TraceID().String()
	}

	counter := m.TotalLobbies.WithLabelValues(LobbyID.String())
	if adder, ok := counter.(prometheus.ExemplarAdder); amount >= 0 && ok && traceID != "" {
		adder.AddWithExemplar(float64(amount), prometheus.Labels{"traceID": traceID})
	} else {
		counter.Add(float64(amount))
	}
}
