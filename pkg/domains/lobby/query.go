package lobby

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/samber/do"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/constants"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	pkgJson "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/codec/json"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
)

func GetLobbyEventByID(ctx context.Context, lobbyID uuid.UUID) ([]common.Event, error) {
	codec := &pkgJson.EventCodec{}

	connPool, err := do.InvokeNamed[*pool.ConnPool](nil, string(constants.ConnectionPool))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conn, err := connPool.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	//nolint:errcheck
	defer connPool.Put(conn)

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	consumer, err := js.OrderedConsumer(cancelCtx, constants.LobbyStream, jetstream.OrderedConsumerConfig{
		FilterSubjects: []string{
			fmt.Sprintf("%s.%s.%s.%s", constants.LobbyStream, AggregateType, lobbyID, lobbyID.String()),
		},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		ReplayPolicy:      jetstream.ReplayInstantPolicy,
		InactiveThreshold: time.Second * 5,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return nil, errors.WithStack(eventing.ErrEntityNotFound)
	}

	res, err := consumer.Fetch(int(initialPending), jetstream.FetchMaxWait(time.Second*10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	events := make([]common.Event, 0, initialPending)
	read := uint64(0)
	for msg := range res.Messages() {
		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}

		seq := md.Sequence.Stream
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(uint64(0), seq),
		}
		header := msg.Headers()
		if header.Get(nats.EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		if err != nil {
			return nil, err
		}

		events = append(events, event)

		read++
		if read == initialPending {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

func GetLobbyByID(ctx context.Context, lobbyID uuid.UUID) (*Lobby, error) {
	var result *Lobby
	codec := &pkgJson.EventCodec{}
	projector := NewProjector()

	connPool, err := do.InvokeNamed[*pool.ConnPool](nil, string(constants.ConnectionPool))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conn, err := connPool.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	//nolint:errcheck
	defer connPool.Put(conn)

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	consumer, err := js.OrderedConsumer(cancelCtx, constants.LobbyStream, jetstream.OrderedConsumerConfig{
		FilterSubjects: []string{
			fmt.Sprintf("%s.%s.%s.%s", constants.LobbyStream, AggregateType, lobbyID, lobbyID.String()),
		},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		ReplayPolicy:      jetstream.ReplayInstantPolicy,
		InactiveThreshold: time.Second * 5,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return nil, errors.WithStack(eventing.ErrEntityNotFound)
	}

	res, err := consumer.Fetch(int(initialPending), jetstream.FetchMaxWait(time.Second*10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	current := new(Lobby)
	read := uint64(0)
	for msg := range res.Messages() {
		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}

		seq := md.Sequence.Stream
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(uint64(0), seq),
		}
		header := msg.Headers()
		if header.Get(nats.EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		if err != nil {
			return nil, err
		}

		result, err = projector.Project(ctx, event, current)
		if err != nil {
			return nil, err
		}

		read++
		if read == initialPending {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return result, nil
}
