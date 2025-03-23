package nats

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go/jetstream"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	codecJson "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/codec/json"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var DefaultEventCodec = codecJson.EventCodec{}

func JSMsgToEvent(ctx context.Context, msg jetstream.Msg) (common.Event, context.Context, error) {
	md, err := msg.Metadata()
	if err != nil {
		return nil, ctx, errors.WithStack(fmt.Errorf("unpack: failed to get metadata: %s", err))
	}
	seq := md.Sequence.Stream

	headers := msg.Headers()
	eventOpts := []eventing.EventOption{
		eventing.ForSequence(0, seq),
	}
	if headers.Get(EventUnregisteredHdr) == "true" {
		eventOpts = append(eventOpts, eventing.AsUnregistered())
	}

	event, ctx, err := DefaultEventCodec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
	if err != nil {
		return nil, ctx, errors.WithStack(err)
	}

	return event, ctx, nil
}
