package eventing

// CommandHandlerMiddleware is a function that middlewares can implement to be
// able to chain.
type CommandHandlerMiddleware func(CommandHandler) CommandHandler

// UseCommandHandlerMiddleware wraps a CommandHandler in one or more middleware.
func UseCommandHandlerMiddleware(h CommandHandler, middleware ...CommandHandlerMiddleware) CommandHandler {
	// Apply in reverse order.
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		h = m(h)
	}

	return h
}

// EventHandlerMiddleware is a function that middlewares can implement to be
// able to chain.
type EventHandlerMiddleware func(EventHandler) EventHandler

// EventHandlerChain declares InnerHandler that returns the inner handler of a event handler middleware.
// This enables an endpoint or other middlewares to traverse the chain of handlers
// in order to find a specific middleware that can be interacted with.
//
// For handlers who's intrinsic properties requires them to be the last responder of a chain, or
// can't produce an InnerHandler, a nil response can be implemented thereby hindering any
// further attempt to traverse the chain.
type EventHandlerChain interface {
	InnerHandler() EventHandler
}

// UseEventHandlerMiddleware wraps a EventHandler in one or more middleware.
func UseEventHandlerMiddleware(h EventHandler, middleware ...EventHandlerMiddleware) EventHandler {
	// Apply in reverse order.
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		h = m(h)
	}

	return h
}
