package natsjs

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	codecJson "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/codec/json"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	suppressedloader "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats/suppressed_loader"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/oplock"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/ttlcache"
)

// EventStore is an eventing.EventStore where all events are stored in
// memory and not persisted. Useful for testing and experimenting.
type EventStore struct {
	db *ttlcache.Cache[string, *aggregateRecord]

	dbMu         sync.RWMutex
	eventHandler eventing.EventHandler

	group          *singleflight.Group
	eventBus       eventing.EventBus
	jetstream      jetstream.JetStream
	connectionPool *pool.ConnPool
	streamName     string
	subject        common.EventSubject
	codec          eventing.EventCodec

	noDefaultFallback bool
}

// NewEventStore creates a new EventStore using memory as storage.
func NewEventStore(ctx context.Context, streamName string, subject common.EventSubject, options ...Option) (*EventStore, error) {
	cache := ttlcache.New[string, *aggregateRecord](ttlcache.WithTTL[string, *aggregateRecord](time.Minute * 60))
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Global().Error("panic in eventstore", zap.Any("recover", r))
			}
		}()
		timer := timeutil.Clock.Ticker(time.Minute * 3)
		for {
			select {
			case <-timer.C:
				cache.DeleteExpired()
			case <-ctx.Done():
				timer.Stop()
				return
			}
		}
	}()
	s := &EventStore{
		db:         cache,
		streamName: streamName,
		subject:    subject,
		group:      new(singleflight.Group),
		codec:      &codecJson.EventCodec{},
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.WithStack(fmt.Errorf("error while applying option: %v", err))
		}
	}

	if s.eventBus == nil && s.jetstream == nil && s.connectionPool == nil {
		return nil, errors.WithStack(fmt.Errorf("no event bus or jetstream connection"))
	}

	return s, nil
}

var _ eventing.EventStore = (*EventStore)(nil)
var _ eventing.EventHandler = (*EventStore)(nil)

// Option is an option setter used to configure creation.
type Option func(*EventStore) error

// WithEventHandler adds an event handler that will be called when saving events.
// An example would be to add an event bus to publish events.
func WithEventHandler(h eventing.EventHandler) Option {
	return func(s *EventStore) error {
		s.eventHandler = h

		return nil
	}
}

func WithConnectionPool(pool *pool.ConnPool) Option {
	return func(s *EventStore) error {
		s.connectionPool = pool

		return nil
	}
}

func WithJetStream(js jetstream.JetStream) Option {
	return func(s *EventStore) error {
		s.jetstream = js

		return nil
	}
}

func WithEventBus(bus *nats.EventBus) Option {
	return func(s *EventStore) error {
		s.eventBus = bus

		return nil
	}
}

func WithNoDefaultFallback() Option {
	return func(s *EventStore) error {
		s.noDefaultFallback = true

		return nil
	}
}

func (s *EventStore) HandlerType() common.EventHandlerType {
	return "natsjs"
}

func (s *EventStore) HandleEvent(ctx context.Context, event common.Event) error {
	s.dbMu.Lock()
	defer s.dbMu.Unlock()

	id := event.AggregateID()

	var cacheMiss bool
	if event.Version() == 0 || event.Version() == 1 {
		aggregate := &aggregateRecord{
			AggregateID: id,
			Version:     1,
			Events: []common.Event{
				event,
			},
		}

		s.db.Set(id, aggregate, ttlcache.DefaultTTL)
	} else {
		subject := fmt.Sprintf("%s.%s", s.streamName, event.Subject(ctx).Subject())

		loader := suppressedloader.NewSuppressedLoader[string, *aggregateRecord](ttlcache.LoaderFunc[string, *aggregateRecord](
			func(cache *ttlcache.Cache[string, *aggregateRecord], key string) *ttlcache.Item[string, *aggregateRecord] {
				ctx := context.WithoutCancel(ctx)

				cacheMiss = true
				var ar *aggregateRecord
				if s.jetstream != nil {
					events, err := nats.LoadJetStream(ctx, s.jetstream, s.streamName, subject, s.codec)
					if err != nil || len(events) == 0 {
						return nil
					}

					ar = &aggregateRecord{
						AggregateID: id,
						Version:     events[len(events)-1].Version(),
						Events:      events,
					}
				} else if s.connectionPool != nil {
					events, err := nats.Load(ctx, s.connectionPool, s.streamName, subject, s.codec)
					if err != nil || len(events) == 0 {
						return nil
					}

					ar = &aggregateRecord{
						AggregateID: id,
						Version:     events[len(events)-1].Version(),
						Events:      events,
					}
				} else if s.eventBus != nil {
					events, err := nats.LoadBus(ctx, s.eventBus, s.streamName, subject, s.codec)
					if err != nil || len(events) == 0 {
						return nil
					}

					ar = &aggregateRecord{
						AggregateID: id,
						Version:     events[len(events)-1].Version(),
						Events:      events,
					}
				} else {
					return nil
				}

				ret := cache.Set(key, ar, ttlcache.DefaultTTL)
				return ret
			},
		), s.group)
		entry := s.db.Get(id, ttlcache.WithLoader[string, *aggregateRecord](loader))
		var aggregate *aggregateRecord
		if entry != nil {
			aggregate = entry.Value()
		}

		if aggregate != nil && aggregate.Version+1 != event.Version() {
			if aggregate.Version >= event.Version() {
				if cacheMiss {
					if aggregate.Version == event.Version() &&
						aggregate.Events[len(aggregate.Events)-1].Version() == event.Version() &&
						aggregate.Events[len(aggregate.Events)-1].EventType() == event.EventType() {
						return nil
					}
				}
				log.Global().WarnContext(ctx, "event store event out of sequence", zap.String("aggregate_id", id), zap.String("aggregate_type", string(event.AggregateType())), zap.String("subject", subject), zap.Uint64("version", event.Version()), zap.Bool("cache_miss", cacheMiss), zap.Uint64("aggregate_version", aggregate.Version), zap.String("event_type", string(event.EventType())), zap.Any("event", event.Data()))
				return nil
			} else if !cacheMiss {
				log.Global().WarnContext(ctx, "event store event out of sequence", zap.String("aggregate_id", id), zap.String("aggregate_type", string(event.AggregateType())), zap.String("subject", subject), zap.Uint64("version", event.Version()), zap.Bool("cache_miss", cacheMiss), zap.Uint64("aggregate_version", aggregate.Version), zap.String("event_type", string(event.EventType())), zap.Any("event", event.Data()))

				// Invalidate the cache and reload the aggregate
				s.db.Delete(id)
				entry := s.db.Get(id, ttlcache.WithLoader[string, *aggregateRecord](loader))
				if entry != nil {
					aggregate = entry.Value()
				}
			} else {
				log.Global().ErrorContext(ctx, "event store event out of sequence", zap.String("aggregate_id", id), zap.String("aggregate_type", string(event.AggregateType())), zap.String("subject", subject), zap.Uint64("version", event.Version()), zap.Bool("cache_miss", cacheMiss), zap.Uint64("aggregate_version", aggregate.Version), zap.String("event_type", string(event.EventType())), zap.Any("event", event.Data()))
				return errors.WithStack(fmt.Errorf("event store event out of sequence"))
			}
		}

		if aggregate != nil {
			if cacheMiss {
				if aggregate.Version == event.Version() &&
					aggregate.Events[len(aggregate.Events)-1].Version() == event.Version() &&
					aggregate.Events[len(aggregate.Events)-1].EventType() == event.EventType() {
					return nil
				}
			}

			if aggregate.Version+1 != event.Version() {
				if aggregate.Version >= event.Version() {
					log.Global().WarnContext(ctx, "event store event out of sequence, after invalidation", zap.String("aggregate_id", id), zap.String("aggregate_type", string(event.AggregateType())), zap.String("subject", subject), zap.Uint64("version", event.Version()), zap.Bool("cache_miss", cacheMiss), zap.Uint64("aggregate_version", aggregate.Version), zap.String("event_type", string(event.EventType())), zap.Any("event", event.Data()))
					return nil
				} else {
					log.Global().ErrorContext(ctx, "event store event out of sequence, after invalidation", zap.String("aggregate_id", id), zap.String("aggregate_type", string(event.AggregateType())), zap.String("subject", subject), zap.Uint64("version", event.Version()), zap.Bool("cache_miss", cacheMiss), zap.Uint64("aggregate_version", aggregate.Version), zap.String("event_type", string(event.EventType())), zap.Any("event", event.Data()))
					return errors.WithStack(fmt.Errorf("event store event out of sequence, after invalidation"))
				}
			}
			aggregate.Version = event.Version()
			aggregate.Events = append(aggregate.Events, event)

			s.db.Set(id, aggregate, ttlcache.DefaultTTL)
		} else {
			log.Global().ErrorContext(ctx, "could not load aggregate, but version is greater than 1", zap.String("aggregate_id", id), zap.String("aggregate_type", string(event.AggregateType())), zap.String("subject", subject), zap.Uint64("version", event.Version()), zap.String("event_type", string(event.EventType())), zap.Any("event", event.Data()))
			return errors.WithStack(fmt.Errorf("could not load aggregate, but version is greater than 1"))
		}
	}

	return nil
}

func (s *EventStore) Sequenced() bool {
	return true
}

// Save implements the Save method of the eventing.EventStore interface.
func (s *EventStore) Save(ctx context.Context, events []common.Event, originalVersion uint64) error {
	if len(events) == 0 {
		return &eventing.EventStoreError{
			Err: eventing.ErrMissingEvents,
			Op:  eventing.EventStoreOpSave,
		}
	}

	id := events[0].AggregateID()
	at := events[0].AggregateType()
	for _, event := range events {
		// Only accept events belonging to the same aggregate.
		if event.AggregateID() != id {
			return &eventing.EventStoreError{
				Err:              eventing.ErrMismatchedEventAggregateIDs,
				Op:               eventing.EventStoreOpSave,
				AggregateType:    at,
				AggregateID:      id,
				AggregateVersion: originalVersion,
				Events:           events,
			}
		}

		if event.AggregateType() != at {
			return &eventing.EventStoreError{
				Err:              eventing.ErrMismatchedEventAggregateTypes,
				Op:               eventing.EventStoreOpSave,
				AggregateType:    at,
				AggregateID:      id,
				AggregateVersion: originalVersion,
				Events:           events,
			}
		}
	}

	// Let the optional event handler handle the events. Aborts the transaction
	// in case of error.
	if s.eventHandler != nil {
		priorVersion := uint64(0)
		for _, e := range events {
			if priorVersion > 0 {
				e = eventing.WithSequence(e, priorVersion, priorVersion+1)
				priorVersion += 1
			}
			if err := s.eventHandler.HandleEvent(ctx, e); err != nil {
				return &eventing.EventStoreError{
					Err: err,
					Op:  eventing.EventStoreOpSave,
				}
			}
			if originalVersion == 0 && len(events) > 1 {
				subject := fmt.Sprintf("%s.%s", s.streamName, e.Subject(ctx).Subject())

				if s.jetstream != nil {
					_, version, err := nats.LoadLatestJetStream(ctx, s.jetstream, s.streamName, subject, s.codec)
					if err != nil {
						return err
					}
					priorVersion = version
				} else if s.connectionPool != nil {
					_, version, err := nats.LoadLatest(ctx, s.connectionPool, s.streamName, subject, s.codec)
					if err != nil {
						return err
					}
					priorVersion = version
				} else if s.eventBus != nil {
					_, version, err := nats.LoadLatestBus(ctx, s.eventBus, s.streamName, subject, s.codec)
					if err != nil {
						return err
					}
					priorVersion = version
				} else {
					return errors.WithStack(fmt.Errorf("no event bus or jetstream connection"))
				}
			}
		}
	}

	return nil
}

// Load implements the Load method of the eventing.EventStore interface.
func (s *EventStore) Load(ctx context.Context, id string) ([]common.Event, error) {
	return s.LoadFrom(ctx, id, 1)
}

// LoadFrom loads all events from version for the aggregate id from the store.
func (s *EventStore) LoadFrom(ctx context.Context, id string, version uint64) ([]common.Event, error) {
	s.dbMu.RLock()
	defer s.dbMu.RUnlock()

	var subject string
	tokens := s.subject.Tokens()
	tokenAggregateIDPos := s.subject.SubjectTokenPosition()

	subjectTokens := make([]string, len(tokens))
	for i, token := range tokens {
		if token.Key() == "aggregate_type" {
			subjectTokens[i] = fmt.Sprint(token.Value())
			continue
		}
		if token.Position() == tokenAggregateIDPos {
			subjectTokens[i] = id
			continue
		}
		subjectTokens[i] = "*"
	}
	subject = fmt.Sprintf("%s.%s", s.streamName, strings.Join(subjectTokens, "."))

	var loadOptions []ttlcache.Option[string, *aggregateRecord]

	loader := suppressedloader.NewSuppressedLoader[string, *aggregateRecord](ttlcache.LoaderFunc[string, *aggregateRecord](
		func(cache *ttlcache.Cache[string, *aggregateRecord], key string) *ttlcache.Item[string, *aggregateRecord] {
			ctx := context.WithoutCancel(ctx)

			var ar *aggregateRecord
			if s.jetstream != nil {
				events, err := nats.LoadJetStreamFrom(ctx, s.jetstream, s.streamName, subject, version, s.codec)
				if err != nil || len(events) == 0 {
					return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
				}

				ar = &aggregateRecord{
					AggregateID: id,
					Version:     events[len(events)-1].Version(),
					Events:      events,
				}
			} else if s.connectionPool != nil {
				events, err := nats.LoadFrom(ctx, s.connectionPool, s.streamName, subject, version, s.codec)
				if err != nil || len(events) == 0 {
					return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
				}

				ar = &aggregateRecord{
					AggregateID: id,
					Version:     events[len(events)-1].Version(),
					Events:      events,
				}
			} else if s.eventBus != nil {
				events, err := nats.LoadBusFrom(ctx, s.eventBus, s.streamName, subject, version, s.codec)
				if err != nil || len(events) == 0 {
					return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
				}

				ar = &aggregateRecord{
					AggregateID: id,
					Version:     events[len(events)-1].Version(),
					Events:      events,
				}
			} else {
				return nil
			}

			return cache.CompareAndSet(key, ar, int64(ar.Version), ttlcache.DefaultTTL)
		},
	), s.group)

	if !s.noDefaultFallback || oplock.OplockFromContext(ctx) {
		loadOptions = append(loadOptions, ttlcache.WithLoader[string, *aggregateRecord](loader))
	}

	retried := false
retryOnce:
	entry := s.db.Get(id, loadOptions...)
	var aggregate *aggregateRecord
	if entry != nil {
		aggregate = entry.Value()
	}
	if aggregate == nil {
		if entry != nil && oplock.OplockFromContext(ctx) && !retried {
			s.db.Delete(id)
			retried = true
			goto retryOnce
		}

		return nil, &eventing.EventStoreError{
			Err:         eventing.ErrAggregateNotFound,
			Op:          eventing.EventStoreOpLoad,
			AggregateID: id,
			Subject:     subject,
		}
	}

	events := make([]common.Event, len(aggregate.Events))

	for i, event := range aggregate.Events {
		if event.Version() < version {
			continue
		}

		e, err := event.Clone(ctx)
		if err != nil {
			return nil, &eventing.EventStoreError{
				Err:              fmt.Errorf("could not copy event: %w", err),
				Op:               eventing.EventStoreOpLoad,
				AggregateType:    event.AggregateType(),
				AggregateID:      id,
				AggregateVersion: event.Version(),
				Events:           events,
				Subject:          subject,
			}
		}

		events[i] = e
	}

	return events, nil
}

type aggregateRecord struct {
	AggregateID string
	Version     uint64
	Events      []common.Event
}

// Close implements the Close method of the eventing.EventStore interface.
func (s *EventStore) Close() error {
	return nil
}
