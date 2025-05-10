package consumeroptions

import (
	"github.com/nats-io/nats.go/jetstream"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
)

// ConsumerOptionsProvider is used to check for an observer middleware in a chain.
type ConsumerOptionsProvider interface {
	ConsumerOptions(defaultConfig *jetstream.ConsumerConfig) *jetstream.ConsumerConfig
}

type DeliveryPolicyProvider interface {
	DeliveryPolicy() jetstream.DeliverPolicy
	StartSequence() uint64
}

type eventHandler struct {
	eventing.EventHandler
	config jetstream.ConsumerConfig
}

type deliveryEventHandler struct {
	eventing.EventHandler
	policy        jetstream.DeliverPolicy
	startSequence uint64
}

// NewMiddleware creates a new middleware that can be examined for options.
func NewMiddleware(config jetstream.ConsumerConfig) func(eventing.EventHandler) eventing.EventHandler {
	return func(h eventing.EventHandler) eventing.EventHandler {
		return &eventHandler{
			EventHandler: h,
			config:       config,
		}
	}
}

// NewDeliveryPolicyMiddleware creates a new middleware that can be examined for options.
func NewDeliveryPolicyMiddleware(policy jetstream.DeliverPolicy, startSequence uint64) func(eventing.EventHandler) eventing.EventHandler {
	return func(h eventing.EventHandler) eventing.EventHandler {
		return &deliveryEventHandler{
			EventHandler:  h,
			policy:        policy,
			startSequence: startSequence,
		}
	}
}

// ConsumerOptions returns true if the handler values
func (h *eventHandler) ConsumerOptions(defaultOptions *jetstream.ConsumerConfig) *jetstream.ConsumerConfig {
	return &h.config
}

// DeliveryPolicy returns true if the handler values
func (h *deliveryEventHandler) DeliveryPolicy() jetstream.DeliverPolicy {
	return h.policy
}

// StartSequence returns true if the handler values
func (h *deliveryEventHandler) StartSequence() uint64 {
	return h.startSequence
}

// InnerHandler implements MiddlewareChain
func (h *eventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}

// InnerHandler implements MiddlewareChain
func (h *deliveryEventHandler) InnerHandler() eventing.EventHandler {
	return h.EventHandler
}
