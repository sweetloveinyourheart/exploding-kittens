package domains

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/samber/do"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/cmdutil"
)

type ClientMetrics struct {
	lobbyStreamCounter metric.Int64Counter
	lobbyStreamGauge   metric.Int64UpDownCounter
	gameStreamCounter  metric.Int64Counter
	gameStreamGauge    metric.Int64UpDownCounter
}

var _ cmdutil.Initializer = (*ClientMetrics)(nil)

func (m *ClientMetrics) Initialize() {
	meter := otel.GetMeterProvider().Meter("ClientMetrics")

	var err error
	m.lobbyStreamCounter, err = meter.Int64Counter(
		"lobby_streams_total",
		metric.WithDescription("How many lobby streams have been opened for this client"),
	)
	if err != nil {
		panic(err)
	}

	m.lobbyStreamGauge, err = meter.Int64UpDownCounter(
		"lobby_streams_gauge",
		metric.WithDescription("The current number of lobby streams open for this client"),
	)
	if err != nil {
		panic(err)
	}

	m.gameStreamCounter, err = meter.Int64Counter(
		"game_streams_total",
		metric.WithDescription("How many game streams have been opened for this client"),
	)
	if err != nil {
		panic(err)
	}

	m.gameStreamGauge, err = meter.Int64UpDownCounter(
		"game_streams_gauge",
		metric.WithDescription("The current number of game streams open for this client"),
	)
	if err != nil {
		panic(err)
	}

	do.Provide[*ClientMetrics](nil, func(i *do.Injector) (*ClientMetrics, error) {
		return m, nil
	})
}

func (m *ClientMetrics) LobbyStreamCounterInc(ctx context.Context, lobbyID uuid.UUID) {
	if m == nil {
		return
	}
	m.lobbyStreamCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("lobby_id", lobbyID.String())))
}

func (m *ClientMetrics) LobbyStreamGaugeAdd(ctx context.Context, lobbyID uuid.UUID, delta int64) {
	if m == nil {
		return
	}
	m.lobbyStreamGauge.Add(ctx, delta, metric.WithAttributes(attribute.String("lobby_id", lobbyID.String())))
}

func (m *ClientMetrics) GameStreamCounterInc(ctx context.Context, gameID uuid.UUID) {
	if m == nil {
		return
	}
	m.gameStreamCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("game_id", gameID.String())))
}

func (m *ClientMetrics) GameStreamGaugeAdd(ctx context.Context, gameID uuid.UUID, delta int64) {
	if m == nil {
		return
	}
	m.gameStreamGauge.Add(ctx, delta, metric.WithAttributes(attribute.String("game_id", gameID.String())))
}
