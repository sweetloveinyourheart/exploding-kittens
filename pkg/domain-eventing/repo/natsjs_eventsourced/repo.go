package natsjs_eventsourced

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	codecJson "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/codec/json"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats"
	suppressedloader "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_bus/nats/suppressed_loader"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/event_handler/projector"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/timeutil"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/ttlcache"
)

const DefaultTTL = time.Minute * 60
const DefaultEvictionTimeout = time.Minute * 3

// ErrModelNotSet is when an model factory is not set on the Repo.
var ErrModelNotSet = errors.New("model not set")

// Repo implements an in memory repository of read models.
type Repo[T any, PT eventing.GenericEntity[T]] struct {
	db *ttlcache.Cache[string, []byte]

	streamName   string
	eventSubject common.EventSubject
	projector    projector.Projector[T, PT]
	repoOptions  *repoOptions

	group      *singleflight.Group
	codec      eventing.EventCodec
	entityType string
}

type repoOptions struct {
	eventBus        eventing.EventBus
	jetstream       jetstream.JetStream
	connectionPool  *pool.ConnPool
	subjectIdentity func(ctx context.Context, id string) string
	filterFunc      func(ctx context.Context, id string, event common.Event) bool
}

// NewRepo creates a new Repo.
func NewRepo[T any, PT eventing.GenericEntity[T]](ctx context.Context, streamName string, eventSubject common.EventSubject, projector projector.Projector[T, PT], options ...Option) (*Repo[T, PT], error) {
	cache := ttlcache.New[string, []byte](ttlcache.WithTTL[string, []byte](DefaultTTL))
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Global().Error("panic in repo", zap.Any("recover", r))
			}
		}()
		timer := timeutil.Clock.Timer(DefaultEvictionTimeout)
		stop := func() {
			if !timer.Stop() {
				// drain the timer chan
				select {
				case <-timer.C:
				default:
				}
			}
		}
		defer stop()
		for {
			select {
			case <-timer.C:
				cache.DeleteExpired()
				timer.Reset(DefaultEvictionTimeout)
			case <-ctx.Done():
				timer.Stop()
				return
			}
		}
	}()
	r := &Repo[T, PT]{
		db: cache,

		eventSubject: eventSubject,
		streamName:   streamName,
		projector:    projector,
		repoOptions:  &repoOptions{},

		group: new(singleflight.Group),
		codec: &codecJson.EventCodec{},
	}

	r.entityType = fmt.Sprintf("%T", *new(T))

	for _, option := range options {
		if err := option(r.repoOptions); err != nil {
			return nil, errors.WithStack(fmt.Errorf("error while applying option: %v", err))
		}
	}

	if r.repoOptions.eventBus == nil && r.repoOptions.jetstream == nil && r.repoOptions.connectionPool == nil {
		return nil, errors.WithStack(fmt.Errorf("no event bus or jetstream connection"))
	}

	return r, nil
}

// Option is an option setter used to configure creation.
type Option func(options *repoOptions) error

func WithConnectionPool(pool *pool.ConnPool) Option {
	return func(s *repoOptions) error {
		s.connectionPool = pool

		return nil
	}
}

func WithJetStream(js jetstream.JetStream) Option {
	return func(s *repoOptions) error {
		s.jetstream = js

		return nil
	}
}

func WithEventBus(bus *nats.EventBus) Option {
	return func(s *repoOptions) error {
		s.eventBus = bus

		return nil
	}
}

func WithSubjectIdentity(identityFunc func(ctx context.Context, id string) string) Option {
	return func(s *repoOptions) error {
		s.subjectIdentity = identityFunc

		return nil
	}
}

func WithEventFilter(filterFunc func(ctx context.Context, id string, event common.Event) bool) Option {
	return func(s *repoOptions) error {
		s.filterFunc = filterFunc

		return nil
	}
}

// InnerRepo implements the InnerRepo method of the eventing.ReadRepo interface.
func (r *Repo[T, PT]) InnerRepo(ctx context.Context) eventing.ReadRepo[T, PT] {
	return nil
}

// AdaptFrom tries to convert an eventing.ReadRepo into a Repo by recursively looking at
// inner repos. Returns nil if none was found.
func AdaptFrom[T any, PT eventing.GenericEntity[T]](ctx context.Context, repo eventing.ReadRepo[T, PT]) *Repo[T, PT] {
	if repo == nil {
		return nil
	}

	if r, ok := repo.(*Repo[T, PT]); ok {
		return r
	}

	return AdaptFrom(ctx, repo.InnerRepo(ctx))
}

// Find implements the Find method of the eventing.ReadRepo interface.
func (r *Repo[T, PT]) Find(ctx context.Context, id string) (*T, error) {
	// Fetch entity.
	var subject string
	tokens := r.eventSubject.Tokens()
	tokenAggregateIDPos := r.eventSubject.SubjectTokenPosition()

	subjectID := id
	if r.repoOptions != nil && r.repoOptions.subjectIdentity != nil {
		subjectID = r.repoOptions.subjectIdentity(ctx, id)
	}

	subjectTokens := make([]string, len(tokens))
	for i, token := range tokens {
		if token.Key() == "aggregate_type" {
			subjectTokens[i] = fmt.Sprint(token.Value())
			continue
		}
		if token.Position() == tokenAggregateIDPos {
			subjectTokens[i] = subjectID
			continue
		}
		subjectTokens[i] = "*"
	}
	subject = fmt.Sprintf("%s.%s", r.streamName, strings.Join(subjectTokens, "."))

	loader := suppressedloader.NewSuppressedLoader[string, []byte](ttlcache.LoaderFunc[string, []byte](
		func(cache *ttlcache.Cache[string, []byte], key string) *ttlcache.Item[string, []byte] {
			var ar any
			if r.repoOptions.jetstream != nil {
				events, err := nats.LoadJetStream(ctx, r.repoOptions.jetstream, r.streamName, subject, r.codec)
				if err != nil || len(events) == 0 {
					return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
				}

				ar = new(T)
				for _, event := range events {
					if r.repoOptions.filterFunc != nil && !r.repoOptions.filterFunc(ctx, id, event) {
						continue
					}
					if ar, err = r.projector.Project(ctx, event, ar.(*T)); err != nil {
						return nil
					}
				}
			} else if r.repoOptions.connectionPool != nil {
				events, err := nats.Load(ctx, r.repoOptions.connectionPool, r.streamName, subject, r.codec)
				if err != nil || len(events) == 0 {
					return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
				}

				ar = new(T)
				for _, event := range events {
					if r.repoOptions.filterFunc != nil && !r.repoOptions.filterFunc(ctx, id, event) {
						continue
					}
					if ar, err = r.projector.Project(ctx, event, ar.(*T)); err != nil {
						return nil
					}
				}
			} else if r.repoOptions.eventBus != nil {
				events, err := nats.LoadBus(ctx, r.repoOptions.eventBus, r.streamName, subject, r.codec)
				if err != nil || len(events) == 0 {
					return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
				}

				ar = new(T)
				for _, event := range events {
					if r.repoOptions.filterFunc != nil && !r.repoOptions.filterFunc(ctx, id, event) {
						continue
					}
					if ar, err = r.projector.Project(ctx, event, ar.(*T)); err != nil {
						return nil
					}
				}
			} else {
				return nil
			}

			if ar.(*T) == nil {
				return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
			}

			data, err := json.Marshal(ar)
			if err != nil {
				return cache.CompareAndSet(key, nil, 0, ttlcache.DefaultTTL)
			}

			if versionable, ok := ar.(eventing.Versionable); ok {
				return cache.CompareAndSet(key, data, int64(versionable.AggregateVersion()), ttlcache.DefaultTTL)
			}
			ret := cache.Set(key, data, ttlcache.DefaultTTL)
			return ret
		},
	), r.group)
	entity := r.db.Get(id, ttlcache.WithLoader[string, []byte](loader))
	if entity == nil || len(entity.Value()) == 0 {
		return nil, errors.WithStack(&eventing.RepoError{
			Err:        errors.Wrap(eventing.ErrEntityNotFound, fmt.Sprintf("entity type: %T, subject: %s", *new(T), subject)),
			Op:         eventing.RepoOpFind,
			EntityID:   id,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	var entityValue any = new(T)
	if err := json.Unmarshal(entity.Value(), entityValue); err != nil {
		return nil, errors.WithStack(&eventing.RepoError{
			Err:        fmt.Errorf("could not unmarshal: %w", err),
			Op:         eventing.RepoOpFind,
			EntityID:   id,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	return entityValue.(*T), nil
}

func (r *Repo[T, PT]) loadAll(ctx context.Context) ([]*T, error) {
	var subject string
	tokens := r.eventSubject.Tokens()

	subjectTokens := make([]string, len(tokens))
	for i, token := range tokens {
		if token.Key() == "aggregate_type" {
			subjectTokens[i] = fmt.Sprint(token.Value())
			continue
		}
		subjectTokens[i] = "*"
	}
	subject = fmt.Sprintf("%s.%s", r.streamName, strings.Join(subjectTokens, "."))

	building := make(map[string]*T)
	if r.repoOptions.jetstream != nil {
		events, err := nats.LoadJetStream(ctx, r.repoOptions.jetstream, r.streamName, subject, r.codec)
		if err != nil || len(events) == 0 {
			var jsErr *jetstream.APIError
			if (errors.As(err, &jsErr) && jsErr.Code == 404) || len(events) == 0 {
				return nil, errors.WithStack(eventing.ErrEntityNotFound)
			}
			return nil, errors.Errorf("error loading events: %v", err)
		}

		for _, event := range events {
			ar, ok := building[event.AggregateID()]
			if !ok {
				ar = new(T)
				building[event.AggregateID()] = ar
			}

			if ar, err = r.projector.Project(ctx, event, ar); err != nil {
				return nil, err
			} else {
				building[event.AggregateID()] = ar
			}
		}
	} else if r.repoOptions.connectionPool != nil {
		events, err := nats.Load(ctx, r.repoOptions.connectionPool, r.streamName, subject, r.codec)
		if err != nil || len(events) == 0 {
			var jsErr *jetstream.APIError
			if (errors.As(err, &jsErr) && jsErr.Code == 404) || len(events) == 0 {
				return nil, errors.WithStack(eventing.ErrEntityNotFound)
			}
			return nil, errors.Errorf("error loading events: %v", err)
		}

		for _, event := range events {
			ar, ok := building[event.AggregateID()]
			if !ok {
				ar = new(T)
				building[event.AggregateID()] = ar
			}

			if ar, err = r.projector.Project(ctx, event, ar); err != nil {
				return nil, err
			} else {
				building[event.AggregateID()] = ar
			}
		}
	} else if r.repoOptions.eventBus != nil {
		events, err := nats.LoadBus(ctx, r.repoOptions.eventBus, r.streamName, subject, r.codec)
		if err != nil || len(events) == 0 {
			var jsErr *jetstream.APIError
			if (errors.As(err, &jsErr) && jsErr.Code == 404) || len(events) == 0 {
				return nil, errors.WithStack(eventing.ErrEntityNotFound)
			}
			return nil, errors.Errorf("error loading events: %v", err)
		}

		for _, event := range events {
			ar, ok := building[event.AggregateID()]
			if !ok {
				ar = new(T)
				building[event.AggregateID()] = ar
			}

			if ar, err = r.projector.Project(ctx, event, ar); err != nil {
				return nil, err
			} else {
				building[event.AggregateID()] = ar
			}
		}
	} else {
		return nil, errors.Errorf("no event bus or jetstream connection")
	}

	result := make([]*T, 0)
	for _, v := range building {
		if v == nil {
			continue
		}
		result = append(result, v)
	}

	return result, nil
}

// FindAll implements the FindAll method of the eventing.ReadRepo interface.
func (r *Repo[T, PT]) FindAll(ctx context.Context) ([]*T, error) {
	results, err := r.loadAll(ctx)
	if err != nil {
		if errors.Is(err, eventing.ErrEntityNotFound) {
			var result []*T
			return result, nil
		}
		return nil, errors.WithStack(&eventing.RepoError{
			Err:        fmt.Errorf("could not load all: %w", err),
			Op:         eventing.RepoOpFindAll,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	for _, entity := range results {
		data, err := json.Marshal(entity)
		if err != nil {
			return nil, errors.WithStack(&eventing.RepoError{
				Err:        fmt.Errorf("could not marshal: %w", err),
				Op:         eventing.RepoOpFindAll,
				EntityID:   PT(entity).EntityID(),
				EntityType: fmt.Sprintf("%T", *new(T)),
			})
		}

		var entityValue any = entity
		if versionable, ok := entityValue.(eventing.Versionable); ok {
			r.db.CompareAndSet(PT(entity).EntityID(), data, int64(versionable.AggregateVersion()), ttlcache.DefaultTTL)
		} else {
			r.db.Set(PT(entity).EntityID(), data, ttlcache.DefaultTTL)
		}
	}

	return results, nil
}

// FindAllCached returns all entities from the cache.
func (r *Repo[T, PT]) FindAllCached(ctx context.Context) ([]*T, error) {
	entries := r.db.Items()

	results := make([]*T, 0, len(entries))
	for _, entry := range entries {
		if len(entry.Value()) == 0 {
			continue
		}

		var entityValue any = new(T)
		if err := json.Unmarshal(entry.Value(), entityValue); err != nil {
			return nil, errors.WithStack(&eventing.RepoError{
				Err:        fmt.Errorf("could not unmarshal: %w", err),
				Op:         eventing.RepoOpFindAll,
				EntityID:   entry.Key(),
				EntityType: fmt.Sprintf("%T", *new(T)),
			})
		}
		results = append(results, entityValue.(*T))
	}

	return results, nil
}

// Save implements the Save method of the eventing.WriteRepo interface.
func (r *Repo[T, PT]) Save(ctx context.Context, entity *T) error {
	var gentity any = entity
	id := gentity.(common.Entity).EntityID()
	if id == "" || id == uuid.Nil.String() {
		return errors.WithStack(&eventing.RepoError{
			Err:        fmt.Errorf("missing entity ID"),
			Op:         eventing.RepoOpSave,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	data, err := json.Marshal(entity)
	if err != nil {
		return errors.WithStack(&eventing.RepoError{
			Err:        fmt.Errorf("could not marshal: %w", err),
			Op:         eventing.RepoOpSave,
			EntityID:   id,
			EntityType: fmt.Sprintf("%T", *new(T)),
		})
	}

	var entityValue any = entity
	if versionable, ok := entityValue.(eventing.Versionable); entity != nil && ok {
		r.db.CompareAndSet(id, data, int64(versionable.AggregateVersion()), ttlcache.DefaultTTL)
	} else {
		r.db.Set(id, data, ttlcache.DefaultTTL)
	}

	return nil
}

// Remove implements the Remove method of the eventing.WriteRepo interface.
func (r *Repo[T, PT]) Remove(ctx context.Context, id string) error {
	if ent, ok := r.db.GetAndDelete(id); ok && ent != nil && len(ent.Value()) > 0 {
		return nil
	}

	var e T
	return errors.WithStack(&eventing.RepoError{
		Err:        errors.Wrap(eventing.ErrEntityNotFound, fmt.Sprintf("entity type: %T", e)),
		Op:         eventing.RepoOpRemove,
		EntityID:   id,
		EntityType: fmt.Sprintf("%T", *new(T)),
	})
}

// Close implements the Close method of the eventing.WriteRepo interface.
func (r *Repo[T, PT]) Close() error {
	return nil
}
