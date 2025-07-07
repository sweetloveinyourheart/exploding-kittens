package tracing

import (
	"context"

	"go.opentelemetry.io/otel"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
)

const TracerName = "com.sweetloveinyourheart.kittens.eventing"

// RegisterContext registers the tracing span to be marshaled/unmarshaled on the
// context. This enables propagation of the tracing spans for backends that
// supports it (like Jaeger).
//
// For usage with Elastic APM which doesn't support submitting of child spans
// for the same parent span multiple times outside of a single transaction don't
// register the context. This will provide a new context upon handling in the
// event bus or outbox, which currently is the best Elastic APM can support.
//
// See: https://github.com/elastic/apm/issues/122
func init() {
	eventing.RegisterContextMarshaler(func(ctx context.Context, vals map[string]any) {
		carrier := NewMapMessageCarrier(vals)
		propagator := otel.GetTextMapPropagator()
		propagator.Inject(ctx, carrier)
	})
	eventing.RegisterContextUnmarshaler(func(_ context.Context, vals map[string]any) context.Context {
		carrier := NewMapMessageCarrier(vals)
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(context.Background(), carrier)
		return ctx
	})
}
