package nats

import codecJson "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/codec/json"

var DefaultEventCodec = codecJson.EventCodec{}
