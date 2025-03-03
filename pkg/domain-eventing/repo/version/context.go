package version

import (
	"context"
	"time"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
)

// DefaultMinVersionDeadline is the deadline to use when creating a min version
// context that waits.
var DefaultMinVersionDeadline = 10 * time.Second

func init() {
	// Register the version context.
	eventing.RegisterContextMarshaler(func(ctx context.Context, vals map[string]interface{}) {
		if v, ok := ctx.Value(minVersionKey).(uint64); ok {
			vals[minVersionKeyStr] = v
		}
	})

	eventing.RegisterContextUnmarshaler(func(ctx context.Context, vals map[string]interface{}) context.Context {
		if v, ok := vals[minVersionKeyStr].(uint64); ok {
			return NewContextWithMinVersion(ctx, v)
		}

		if v, ok := vals[minVersionKeyStr].(int); ok {
			return NewContextWithMinVersion(ctx, uint64(v))
		}

		// Support JSON-like marshaling of ints as floats.
		if v, ok := vals[minVersionKeyStr].(float64); ok {
			return NewContextWithMinVersion(ctx, uint64(v))
		}

		if v, ok := vals[minVersionKeyStr].(int64); ok {
			return NewContextWithMinVersion(ctx, uint64(v))
		}

		return ctx
	})
}

type contextKey int

const (
	minVersionKey contextKey = iota
	matchVersionKey
)

// Strings used to marshal context values.
const (
	minVersionKeyStr   = "eventing_minversion"
	matchVersionKeyStr = "eventing_matchversion"
)

// MinVersionFromContext returns the min version from the context.
func MinVersionFromContext(ctx context.Context) (uint64, bool) {
	minVersion, ok := ctx.Value(minVersionKey).(uint64)

	return minVersion, ok
}

// MatchVersionFromContext returns the version matcher from the context.
func MatchVersionFromContext[T any, PT eventing.GenericEntity[T]](ctx context.Context) (func(entity *T) bool, bool) {
	minVersion, ok := ctx.Value(matchVersionKey).(func(entity *T) bool)

	return minVersion, ok
}

// NewContextWithMinVersion returns the context with min version set.
func NewContextWithMinVersion(ctx context.Context, minVersion uint64) context.Context {
	return context.WithValue(ctx, minVersionKey, minVersion)
}

// NewContextWithMinVersionWait returns the context with min version and a
// default deadline set.
func NewContextWithMinVersionWait(ctx context.Context, minVersion uint64) (c context.Context, cancel func()) {
	ctx = context.WithValue(ctx, minVersionKey, minVersion)

	return context.WithTimeout(ctx, DefaultMinVersionDeadline)
}

// NewContextWithWaitBy returns the context with matcher and a
// default deadline set.
func NewContextWithWaitBy[T any, PT eventing.GenericEntity[T]](ctx context.Context, f func(event *T) bool) (c context.Context, cancel func()) {
	ctx = context.WithValue(ctx, matchVersionKey, f)

	return context.WithTimeout(ctx, DefaultMinVersionDeadline)
}
