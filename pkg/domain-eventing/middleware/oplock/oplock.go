package oplock

import "context"

type oplockType struct{}

var oplockKey = oplockType{}

// ContextWithOplock returns a new context with the oplock flag set.
// This flag indicates an optimistic lock failure has occurred.
func ContextWithOplock(ctx context.Context) context.Context {
	if OplockFromContext(ctx) {
		return ctx
	}
	return context.WithValue(ctx, oplockKey, true)
}

// OplockFromContext returns the oplock flag from the context.
// This flag indicates an optimistic lock failure has occurred.
func OplockFromContext(ctx context.Context) bool {
	v := ctx.Value(oplockKey)
	if v == nil {
		return false
	}
	return v.(bool)
}
