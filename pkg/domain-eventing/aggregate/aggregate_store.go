package aggregate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
)

// AggregateStore is an aggregate store using event sourcing. It
// uses an event store for loading and saving events used to build the aggregate
// and an event handler to handle resulting events.
type AggregateStore struct {
	store            eventing.EventStore
	streamName       string
	snapshotStore    eventing.SnapshotStore
	isSnapshotStore  bool
	snapshotStrategy eventing.SnapshotStrategy
	isSequencedStore bool
}

var (
	// ErrInvalidEventStore is when a dispatcher is created with a nil event store.
	ErrInvalidEventStore = errors.New("invalid event store")
	// ErrAggregateNotVersioned is when the aggregate does not implement the VersionedAggregate interface.
	ErrAggregateNotVersioned = errors.New("aggregate is not versioned")
	// ErrMismatchedEventType occurs when loaded events from ID does not match aggregate type.
	ErrMismatchedEventType = errors.New("mismatched event type and aggregate type")
)

var _ eventing.AggregateStore = (*AggregateStore)(nil)

// NewAggregateStore creates an aggregate store with an event store and an event
// handler that will handle resulting events (for example by publishing them
// on an event bus).
func NewAggregateStore(store eventing.EventStore, options ...Option) (*AggregateStore, error) {
	if store == nil {
		return nil, ErrInvalidEventStore
	}

	d := &AggregateStore{
		store: store,
	}

	d.snapshotStrategy = &NoSnapshotStrategy{}

	for _, option := range options {
		if err := option(d); err != nil {
			return nil, fmt.Errorf("error while applying option: %w", err)
		}
	}

	d.snapshotStore, d.isSnapshotStore = store.(eventing.SnapshotStore)

	return d, nil
}

// Option is an option setter used to configure creation.
type Option func(*AggregateStore) error

// WithSnapshotStrategy add the strategy to use when determining if a snapshot should be taken
func WithSnapshotStrategy(s eventing.SnapshotStrategy) Option {
	return func(as *AggregateStore) error {
		as.snapshotStrategy = s

		return nil
	}
}

func WithSequencedStore() Option {
	return func(as *AggregateStore) error {
		as.isSequencedStore = true

		return nil
	}
}

func WithStreamName(streamName string) Option {
	return func(as *AggregateStore) error {
		as.streamName = streamName

		return nil
	}
}

// Load implements the Load method of the eventing.AggregateStore interface.
// It loads an aggregate from the event store by creating a new aggregate of the
// type with the ID and then applies all events to it, thus making it the most
// current version of the aggregate.
func (r *AggregateStore) Load(ctx context.Context, aggregateType common.AggregateType, id string) (eventing.Aggregate, error) {
	agg, err := eventing.CreateAggregate(aggregateType, id)
	if err != nil {
		return nil, &eventing.AggregateStoreError{
			Err:           err,
			Op:            eventing.AggregateStoreOpLoad,
			AggregateType: aggregateType,
			AggregateID:   id,
		}
	}

	a, ok := agg.(VersionedAggregate)
	if !ok {
		return nil, &eventing.AggregateStoreError{
			Err:           ErrAggregateNotVersioned,
			Op:            eventing.AggregateStoreOpLoad,
			AggregateType: aggregateType,
			AggregateID:   id,
		}
	}

	fromVersion := uint64(1)

	if sa, ok := a.(eventing.Snapshotable); ok && r.isSnapshotStore {
		snapshot, err := r.snapshotStore.LoadSnapshot(ctx, id)
		if err != nil {
			return nil, &eventing.AggregateStoreError{
				Err:           err,
				Op:            eventing.AggregateStoreOpLoad,
				AggregateType: aggregateType,
				AggregateID:   id,
			}
		}

		if snapshot != nil {
			sa.ApplySnapshot(snapshot)
			fromVersion = snapshot.Version + 1
		}
	}

	ctx = eventing.NewContextWithAggregateType(ctx, aggregateType)
	ctx = eventing.NewContextWithAggregateID(ctx, id)
	events, err := r.store.LoadFrom(ctx, a.EntityID(), fromVersion)
	if err != nil && !errors.Is(err, eventing.ErrAggregateNotFound) {
		return nil, &eventing.AggregateStoreError{
			Err:           err,
			Op:            eventing.AggregateStoreOpLoad,
			AggregateType: aggregateType,
			AggregateID:   id,
		}
	}

	if err := r.applyEvents(ctx, a, events); err != nil {
		return nil, &eventing.AggregateStoreError{
			Err:           err,
			Op:            eventing.AggregateStoreOpLoad,
			AggregateType: aggregateType,
			AggregateID:   id,
		}
	}

	return a, nil
}

// Save implements the Save method of the eventing.AggregateStore interface.
// It saves all uncommitted events from an aggregate to the event store.
func (r *AggregateStore) Save(ctx context.Context, agg eventing.Aggregate) ([]common.Event, error) {
	a, ok := agg.(VersionedAggregate)
	if !ok {
		return nil, &eventing.AggregateStoreError{
			Err:           ErrAggregateNotVersioned,
			Op:            eventing.AggregateStoreOpSave,
			AggregateType: agg.AggregateType(),
			AggregateID:   agg.EntityID(),
		}
	}

	// Retrieve any new events to store.
	events := a.UncommittedEvents()
	if len(events) == 0 {
		return []common.Event{}, nil
	}

	if r.isSequencedStore {
		if err := r.sequenceEvents(ctx, agg, events); err != nil {
			return nil, err
		}
	}

	metadata, ok := ctx.Value(MetadataKey).(map[string]any)
	if ok {
		for i := range events {
			eventing.WithMetadata(metadata)(events[i])
		}
	}

	if err := r.store.Save(ctx, events, a.AggregateVersion()); err != nil {
		return nil, &eventing.AggregateStoreError{
			Err:           err,
			Op:            eventing.AggregateStoreOpSave,
			AggregateType: agg.AggregateType(),
			AggregateID:   agg.EntityID(),
		}
	}

	a.ClearUncommittedEvents()

	// Apply the events in case the aggregate needs to be further used
	// after this save. Currently it is not reused.
	if err := r.applyEvents(ctx, a, events); err != nil {
		return nil, &eventing.AggregateStoreError{
			Err:           err,
			Op:            eventing.AggregateStoreOpSave,
			AggregateType: agg.AggregateType(),
			AggregateID:   agg.EntityID(),
		}
	}

	return events, r.takeSnapshot(ctx, agg, events[len(events)-1])
}

func (r *AggregateStore) sequenceEvents(ctx context.Context, agg eventing.Aggregate, events []common.Event) error {
	if sequencedAggregate, ok := agg.(SequencesAggregate); ok {

		// Check if the namespace is detached and the stream name is set.
		detachedNamespace := eventing.DetachedNamespaceFromContext(ctx)
		if detachedNamespace && stringsutil.IsBlank(r.streamName) {
			return &eventing.AggregateStoreError{
				Err:           errors.New("stream name is required for sequenced store with detached namespaces"),
				Op:            eventing.AggregateStoreOpSave,
				AggregateType: agg.AggregateType(),
				AggregateID:   agg.EntityID(),
			}
		}

		// If the namespace is not detached, then we can just sequence the events as each subject is unique to the aggregate.
		if !detachedNamespace {
			for i := range events {
				events[i] = eventing.WithSequence(events[i], sequencedAggregate.AggregateSequence()+uint64(i), sequencedAggregate.AggregateSequence()+uint64(i)+1)
			}
			return nil
		} else {
			// If the namespace is detached,then there may be multiple subjects for the same aggregate.

			// Get the subject for each event, then locate the last sequence number for that subject
			// and set the sequence number for the event.
			ctx = eventing.NewContextWithAggregateType(ctx, agg.AggregateType())
			ctx = eventing.NewContextWithAggregateID(ctx, agg.EntityID())
			storeEvents, err := r.store.Load(ctx, agg.EntityID())
			if err != nil {
				if errors.Is(err, eventing.ErrAggregateNotFound) {
					for i := range events {
						events[i] = eventing.WithSequence(events[i], sequencedAggregate.AggregateSequence()+uint64(i), sequencedAggregate.AggregateSequence()+uint64(i)+1)
					}
					return nil
				}

				return &eventing.AggregateStoreError{
					Err:           err,
					Op:            eventing.AggregateStoreOpSave,
					AggregateType: agg.AggregateType(),
					AggregateID:   agg.EntityID(),
				}
			}

			type sequencer interface {
				Sequence() uint64
				Sequenced() bool
			}

			for i := range events {
				event := events[i]
				found := false
				for x := len(storeEvents) - 1; x >= 0; x-- {
					if strings.EqualFold(storeEvents[x].Subject(ctx).Subject(), event.Subject(ctx).Subject()) {
						if sequencer, ok := storeEvents[x].(sequencer); ok && sequencer.Sequenced() {
							events[i] = eventing.WithSequence(event, sequencer.Sequence(), sequencer.Sequence()+1)
						} else {
							return &eventing.AggregateStoreError{
								Err:           errors.New("event store does not support sequenced events"),
								Op:            eventing.AggregateStoreOpSave,
								AggregateType: agg.AggregateType(),
								AggregateID:   agg.EntityID(),
							}
						}
						found = true
						break
					}
				}
				if !found {
					events[i] = eventing.WithSequence(event, 0, 1)
				}
			}
		}
	}

	return nil
}

func (r *AggregateStore) takeSnapshot(ctx context.Context, agg eventing.Aggregate, lastEvent common.Event) error {
	a, ok := agg.(eventing.Snapshotable)
	if !ok || !r.isSnapshotStore {
		return nil
	}

	s, err := r.snapshotStore.LoadSnapshot(ctx, agg.EntityID())
	if err != nil {
		return &eventing.AggregateStoreError{
			Err:           err,
			Op:            eventing.AggregateStoreOpSave,
			AggregateType: agg.AggregateType(),
			AggregateID:   agg.EntityID(),
		}
	}

	version := uint64(0)
	timestamp := time.Now()

	if s != nil {
		version = s.Version
		timestamp = s.Timestamp
	}

	if res := r.snapshotStrategy.ShouldTakeSnapshot(version, timestamp, lastEvent); res {
		if err = r.snapshotStore.SaveSnapshot(ctx, agg.EntityID(), *a.CreateSnapshot()); err != nil {
			return &eventing.AggregateStoreError{
				Err:           err,
				Op:            eventing.AggregateStoreOpSave,
				AggregateType: agg.AggregateType(),
				AggregateID:   agg.EntityID(),
			}
		}
	}

	return nil
}

func (r *AggregateStore) applyEvents(ctx context.Context, a VersionedAggregate, events []common.Event) error {
	for _, event := range events {
		if event.AggregateType() != a.AggregateType() {
			return ErrMismatchedEventType
		}

		if err := a.ApplyEvent(ctx, event); err != nil {
			return errors.Errorf("could not apply event %s: %w", event, err)
		}

		a.SetAggregateVersion(event.Version())

		if sequenced, ok := event.(sequenced); ok && sequenced.Sequenced() {
			if sequencedAggregator, ok := a.(SequencesAggregate); ok {
				sequencedAggregator.SetAggregateSequence(sequenced.Sequence())
			}
		}
	}

	return nil
}

type sequenced interface {
	Sequence() uint64
	Sequenced() bool
}

type contextKey string

var MetadataKey = contextKey("metadata")
