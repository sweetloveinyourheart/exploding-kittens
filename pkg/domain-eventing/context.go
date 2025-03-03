package eventing

import (
	"context"
	"sync"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

type contextKey int

const (
	aggregateIDKey contextKey = iota
	aggregateTypeKey
	commandTypeKey
	namespaceDetchedTypeKey
	eventSubjectKey
)

// Private context marshaling funcs.
var (
	contextMarshalFuncs   = []ContextMarshalFunc{}
	contextMarshalFuncsMu = &sync.RWMutex{}

	contextUnmarshalFuncs   = []ContextUnmarshalFunc{}
	contextUnmarshalFuncsMu = &sync.RWMutex{}
)

// ContextMarshalFunc is a function that marshalls any context values to a map,
// used for sending context on the wire.
type ContextMarshalFunc func(context.Context, map[string]interface{})

// RegisterContextMarshaler registers a marshaler function used by MarshalContext.
func RegisterContextMarshaler(f ContextMarshalFunc) {
	contextMarshalFuncsMu.Lock()
	defer contextMarshalFuncsMu.Unlock()

	contextMarshalFuncs = append(contextMarshalFuncs, f)
}

// MarshalContext marshals a context into a map.
func MarshalContext(ctx context.Context) map[string]interface{} {
	contextMarshalFuncsMu.RLock()
	defer contextMarshalFuncsMu.RUnlock()

	allVals := map[string]interface{}{}

	for _, f := range contextMarshalFuncs {
		vals := map[string]interface{}{}
		f(ctx, vals)

		for key, val := range vals {
			if _, ok := allVals[key]; ok {
				panic("duplicate context entry for: " + key)
			}

			allVals[key] = val
		}
	}

	return allVals
}

// ContextUnmarshalFunc is a function that marshals any context values to a map,
// used for sending context on the wire.
type ContextUnmarshalFunc func(context.Context, map[string]interface{}) context.Context

// RegisterContextUnmarshaler registers a marshaler function used by UnmarshalContext.
func RegisterContextUnmarshaler(f ContextUnmarshalFunc) {
	contextUnmarshalFuncsMu.Lock()
	defer contextUnmarshalFuncsMu.Unlock()

	contextUnmarshalFuncs = append(contextUnmarshalFuncs, f)
}

// UnmarshalContext unmarshals a context from a map.
func UnmarshalContext(ctx context.Context, vals map[string]interface{}) context.Context {
	contextUnmarshalFuncsMu.RLock()
	defer contextUnmarshalFuncsMu.RUnlock()

	if vals == nil {
		return ctx
	}

	for _, f := range contextUnmarshalFuncs {
		ctx = f(ctx, vals)
	}

	return ctx
}

// DetachedNamespaceFromContext return the command type from the context.
func DetachedNamespaceFromContext(ctx context.Context) bool {
	detached, ok := ctx.Value(namespaceDetchedTypeKey).(bool)

	return detached && ok
}

// NewContextWithAggregateID adds a aggregate ID on the context.
func NewContextWithAggregateID(ctx context.Context, aggregateID string) context.Context {
	return context.WithValue(ctx, aggregateIDKey, aggregateID)
}

// NewContextWithAggregateType adds a aggregate type on the context.
func NewContextWithAggregateType(ctx context.Context, aggregateType common.AggregateType) context.Context {
	return context.WithValue(ctx, aggregateTypeKey, aggregateType)
}

// NewContextWithDetachedNamespace allows for multiple subjects to form a single aggregate
func NewContextWithDetachedNamespace(ctx context.Context) context.Context {
	return context.WithValue(ctx, namespaceDetchedTypeKey, true)
}
