package nats

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/config"
	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	codecJson "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/codec/json"
	consumeroptions "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/consumer_options"
	consumerresetter "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/consumer_resetter"
	consumersync "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/consumer_sync"
	consumerwaiter "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/consumer_waiter"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/ephemeral"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/middleware/sequenced"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"go.uber.org/zap"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

var (
	eventTypeHdr         = "ld-evt-type"
	eventTimeHdr         = "ld-evt-time"
	eventCodecHdr        = "ld-evt-codec"
	EventUnregisteredHdr = "ld-evt-unregistered"
	eventTimeFormat      = time.RFC3339Nano

	eventTimezoneHdr     = "ld-evt-timezone"
	eventServerIDHdr     = "ld-evt-server-id"
	eventJurisdictionHdr = "ld-evt-jurisdiction"
	eventSourceHdr       = "ld-evt-source"
)

// EventBus is a NATS JetStream event bus that delegates handling of published
// events to all matching registered handlers.
type EventBus struct {
	appID         string
	streamName    string
	typeInSubject bool
	busName       string
	pool          *pool.ConnPool
	conn          *nats.Conn
	jetstream     jetstream.JetStream
	stream        jetstream.Stream
	connOpts      []nats.Option
	streamConfig  *jetstream.StreamConfig
	registered    map[common.EventHandlerType]struct{}
	registeredMu  sync.RWMutex
	errCh         chan error
	cctx          context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	codec         eventing.EventCodec
	serviceID     string
	source        string
	unsubscribe   []func()
}

// NewEventBus creates an EventBus, with optional settings.
func NewEventBus(connectionPool *pool.ConnPool, appID string, options ...Option) (*EventBus, error) {
	if connectionPool == nil {
		return nil, errors.WithStack(fmt.Errorf("connection pool is required"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	serviceID := config.Instance().GetString("service_id")
	source := config.Instance().GetString("service")

	b := &EventBus{
		appID:      appID,
		streamName: appID + "_events",
		registered: map[common.EventHandlerType]struct{}{},
		errCh:      make(chan error, 100),
		cctx:       ctx,
		cancel:     cancel,
		codec:      &codecJson.EventCodec{},
		pool:       connectionPool,
		serviceID:  serviceID,
		source:     source,
	}

	// Apply configuration options.
	for _, option := range options {
		if option == nil {
			continue
		}

		if err := option(b); err != nil {
			return nil, errors.WithStack(fmt.Errorf("error while applying option: %w", err))
		}
	}

	// Create the NATS connection.
	var err error
	if b.conn, err = connectionPool.Get(); err != nil {
		return nil, errors.WithStack(fmt.Errorf("could not create NATS connection: %w", err))
	}

	b.busName = fmt.Sprintf("%s::%s::%s", b.appID, b.streamName, b.serviceID)

	b.jetstream, err = jetstream.New(b.conn, jetstream.WithPublishAsyncErrHandler(func(js jetstream.JetStream, msg *nats.Msg, err error) {
		log.Global().ErrorContext(ctx, "eventing: publish error", zap.Error(err), zap.String("bus_name", b.busName))
		b.errCh <- err
	}))
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("could not create Jetstream context: %w", err))
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelCtx()
	if b.stream, err = b.jetstream.Stream(ctx, b.streamName); err == nil {
		return b, nil
	}

	replicas := 1
	if reps := config.Instance().GetInt(config.NatsStreamReplicas); reps > 0 {
		replicas = reps
	}

	storage := jetstream.FileStorage
	if ss := config.Instance().GetString(config.NatsStreamStorage); strings.EqualFold(ss, "memory") {
		storage = jetstream.MemoryStorage
	}

	// Create the stream, which stores messages received on the subject.
	subjects := b.streamName + ".>"
	cfg := jetstream.StreamConfig{
		Name:      b.streamName,
		Subjects:  []string{subjects},
		Storage:   storage,
		Retention: jetstream.LimitsPolicy,
		Replicas:  replicas,
	}

	// Use the custom stream config if provided.
	if b.streamConfig != nil {
		cfg = *b.streamConfig
	}

	if b.stream, err = b.jetstream.CreateStream(ctx, cfg); err != nil {
		return nil, errors.WithStack(fmt.Errorf("could not create NATS stream: %w", err))
	}

	return b, nil
}

// Option is an option setter used to configure creation.
type Option func(*EventBus) error

// WithCodec uses the specified codec for encoding events.
func WithCodec(codec eventing.EventCodec) Option {
	return func(b *EventBus) error {
		b.codec = codec

		return nil
	}
}

func WithStreamName(streamName string) Option {
	return func(b *EventBus) error {
		b.streamName = streamName
		return nil
	}
}

type sequencer interface {
	PreviousSequence() uint64
	Sequenced() bool
}

// HandlerType implements the HandlerType method of the eventing.EventHandler interface.
func (b *EventBus) HandlerType() common.EventHandlerType {
	return "eventbus"
}

// HandleEvent implements the HandleEvent method of the eventing.EventHandler interface.
func (b *EventBus) HandleEvent(ctx context.Context, event common.Event) error {
	if ctx.Err() != nil {
		return errors.WithStack(errors.Wrap(ctx.Err(), fmt.Sprintf("context error, unable to publish; %s", b.busName)))
	}

	//tables.tables.table_created
	//tables.tables.table_deleted
	//subject := fmt.Sprintf("%s.%s.%s.%s", b.streamName, event.AggregateType(), event.AggregateID(), event.EventType())
	//tables.tables.uuid
	subject := fmt.Sprintf("%s.%s", b.streamName, event.Subject(ctx).Subject())

	data, err := b.codec.MarshalEvent(ctx, event)
	if err != nil {
		return errors.WithStack(fmt.Errorf("could not marshal event: %w", err))
	}

	opts := []jetstream.PublishOpt{
		jetstream.WithExpectStream(b.streamName),
	}
	var expectLastSequence uint64
	if sequencer, ok := event.(sequencer); ok && sequencer.Sequenced() {
		expectLastSequence = sequencer.PreviousSequence()
		opts = append(opts, jetstream.WithExpectLastSequencePerSubject(expectLastSequence))
	}

	msg := nats.NewMsg(subject)
	msg.Data = data
	msg.Header.Set(eventTypeHdr, string(event.EventType()))
	msg.Header.Set(eventCodecHdr, nats.JSON_ENCODER)
	msg.Header.Set(eventTimeHdr, event.Timestamp().Format(eventTimeFormat))
	msg.Header.Set(eventServerIDHdr, b.serviceID)
	msg.Header.Set(eventSourceHdr, b.appID)
	if unregistered, ok := event.(common.UnregisteredEvent); ok && unregistered.Unregistered() {
		msg.Header.Set(EventUnregisteredHdr, "true")
	}

	if _, err := b.jetstream.PublishMsg(ctx, msg, opts...); err != nil {
		var natsErr *nats.APIError
		if errors.As(err, &natsErr) {
			if natsErr.ErrorCode == nats.JSErrCodeStreamWrongLastSequence {
				rerr := errors.Wrap(eventing.ErrIncorrectEntityVersion, fmt.Sprintf("could not publish out of sequence event %s %s. expected last sequence: %d. subject: %s", err.Error(), b.busName, expectLastSequence, subject))
				return errors.WithStack(rerr)
			}
		}
		var jsErr *jetstream.APIError
		if errors.As(err, &jsErr) {
			if nats.ErrorCode(jsErr.ErrorCode) == nats.JSErrCodeStreamWrongLastSequence {
				rerr := errors.Wrap(eventing.ErrIncorrectEntityVersion, fmt.Sprintf("could not publish out of sequence event %s %s. expected last sequence: %d. subject: %s", err.Error(), b.busName, expectLastSequence, subject))
				return errors.WithStack(rerr)
			}
		}
		return errors.WithStack(errors.Wrap(err, fmt.Sprintf("could not publish event %s. not a sequence error but expected last sequence: %d. subject: %s", b.busName, expectLastSequence, subject)))
	}

	return nil
}

// AddHandler implements the AddHandler method of the eventing.EventBus interface.
func (b *EventBus) AddHandler(ctx context.Context, m eventing.EventMatcher, h eventing.EventHandler) error {
	if m == nil {
		return errors.WithStack(eventing.ErrMissingMatcher)
	}

	if h == nil {
		return errors.WithStack(eventing.ErrMissingHandler)
	}

	// Check handler existence.
	b.registeredMu.Lock()
	defer b.registeredMu.Unlock()

	if _, ok := b.registered[h.HandlerType()]; ok {
		return errors.WithStack(eventing.ErrHandlerAlreadyAdded)
	}

	// Create a consumer.
	subject := CreateConsumerSubject(b.streamName, m)
	consumerName := fmt.Sprintf("%s_%s", b.appID, h.HandlerType())

	type stoppable interface {
		Stop()
	}

	var sub *nats.Subscription
	var consumer jetstream.Consumer
	var consumerContext stoppable
	var err error
	err = retry.Do(func() error {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		replicas := 1
		if reps := config.Instance().GetInt(config.NatsConsumerReplicas); reps > 0 {
			replicas = reps
		}
		memory := false
		if ss := config.Instance().GetString(config.NatsConsumerStorage); strings.EqualFold(ss, "memory") {
			memory = true
		}
		defaultConfig := jetstream.ConsumerConfig{
			Name:              consumerName,
			DeliverPolicy:     jetstream.DeliverAllPolicy,
			AckPolicy:         jetstream.AckExplicitPolicy,
			AckWait:           10 * time.Second,
			MaxDeliver:        10,
			InactiveThreshold: 60 * time.Minute,
			FilterSubject:     subject,

			Durable:       consumerName,
			Replicas:      replicas,
			MemoryStorage: memory,
		}

		if newoptions, ok := b.handlerIsConsumerOptionsProvider(h, &defaultConfig); ok {
			log.Global().InfoContext(ctx, "eventing: using custom consumer options", zap.String("stream", b.streamName), zap.String("subject", subject), zap.String("consumer", consumerName))
			defaultConfig = *newoptions

			if defaultConfig.Durable != "" {
				defaultConfig.Replicas = replicas
				defaultConfig.MemoryStorage = memory
			}

			if defaultConfig.FilterSubject == "" && len(defaultConfig.FilterSubjects) == 0 {
				defaultConfig.FilterSubject = subject
			}
			if defaultConfig.MaxDeliver == 0 {
				defaultConfig.MaxDeliver = 10
			}
		}

		if newDeliveryPolicy, startSequence, ok := b.handlerIsDeliveryPolicyProvider(h); ok {
			log.Global().InfoContext(ctx, "eventing: using custom delivery policy", zap.String("stream", b.streamName), zap.String("subject", subject), zap.String("consumer", consumerName))
			defaultConfig.DeliverPolicy = newDeliveryPolicy
			defaultConfig.OptStartSeq = startSequence
		}

		if !b.handlerIsEphemeral(h) {
			if b.handlerIsConsumerResetter(h) {
				ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()
				if err := b.jetstream.DeleteConsumer(ctx, b.streamName, consumerName); err != nil {
					if !errors.Is(err, nats.ErrConsumerNotFound) {
						return errors.WithStack(fmt.Errorf("could not delete consumer: %w %s %s", err, b.busName, consumerName))
					}
				} else {
					log.Global().DebugContext(ctx, "eventing: deleted consumer", zap.String("stream", b.streamName), zap.String("consumer", consumerName))
				}
			}

			if b.handlerIsWaiter(h) {
				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()
				// Load all the existing events
				existingEvents, err := LoadRaw(ctx, b.pool, b.streamName, subject, &codecJson.EventCodec{})
				if err != nil {
					return errors.WithStack(errors.Wrap(err, fmt.Sprintf("could not load existing events %s %s", b.busName, consumerName)))
				}
				handler := b.handler(b.cctx, m, h)
				// Deliver the events in order
				for _, event := range existingEvents {
					handler(event)
				}

				if len(existingEvents) != 0 {
					lastEvent := existingEvents[len(existingEvents)-1]
					md, err := lastEvent.Metadata()
					if err != nil {
						return errors.WithStack(fmt.Errorf("unpack: failed to get metadata: %s", err))
					}
					seq := md.Sequence.Stream
					defaultConfig.OptStartSeq = seq + 1
				} else {
					defaultConfig.DeliverPolicy = jetstream.DeliverAllPolicy
				}
			}

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			consumer, err = b.jetstream.CreateOrUpdateConsumer(ctx, b.streamName, defaultConfig)
			if err != nil {
				return errors.WithStack(fmt.Errorf("could not create consumer: %w %s %s", err, b.busName, consumerName))
			}

			ci, err := consumer.Info(ctx)
			if err != nil {
				return errors.WithStack(fmt.Errorf("could not get consumer info: %w %s %s", err, b.busName, consumerName))
			}

			log.Global().InfoContext(ctx, "eventing: created consumer", zap.String("stream", b.streamName), zap.String("consumer", consumerName), zap.String("nats_consumer_name", ci.Name), zap.String("delivery_policy", ci.Config.DeliverPolicy.String()), zap.Uint64("start_sequence", ci.Config.OptStartSeq), zap.String("subject", subject))
		} else {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			deliverPolicy := defaultConfig.DeliverPolicy
			startSequence := defaultConfig.OptStartSeq

			consumer, err = b.jetstream.OrderedConsumer(ctx, b.streamName, jetstream.OrderedConsumerConfig{
				FilterSubjects:    []string{subject},
				DeliverPolicy:     deliverPolicy,
				MaxResetAttempts:  5,
				OptStartSeq:       startSequence,
				InactiveThreshold: 30 * time.Minute,
			})
			if err != nil {
				return errors.WithStack(fmt.Errorf("could not create ordered consumer: %w %s %s", err, b.busName, consumerName))
			}

			ci, err := consumer.Info(ctx)
			if err != nil {
				return errors.WithStack(fmt.Errorf("could not get consumer info: %w %s %s", err, b.busName, consumerName))
			}

			log.Global().InfoContext(ctx, "eventing: created ordered consumer", zap.String("stream", b.streamName), zap.String("consumer", consumerName), zap.String("nats_consumer_name", ci.Name), zap.String("delivery_policy", ci.Config.DeliverPolicy.String()), zap.Uint64("start_sequence", ci.Config.OptStartSeq), zap.String("subject", subject))
		}

		errHandler := jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
			defer b.safeRecover(ctx, "consume error")
			if strings.Contains(err.Error(), "Server Shutdown") {
				if strings.EqualFold(os.Getenv("GO_ENV"), "TEST") {
					log.Global().DebugContext(ctx, "eventing: consume error", zap.Error(err))
				} else {
					log.Global().WarnContext(ctx, "eventing: consume error", zap.Error(err))
				}
				return
			}
			if errors.Is(err, nats.ErrConsumerDeleted) || errors.Is(err, jetstream.ErrConsumerDeleted) {
				log.Global().WarnContext(ctx, "eventing: consumer deleted", zap.Error(err), zap.String("stream", b.streamName), zap.String("consumer", consumerName), zap.String("bus_name", b.busName))
				return
			}
			log.Global().ErrorContext(ctx, "eventing: consume error", zap.Error(err), zap.String("stream", b.streamName), zap.String("consumer", consumerName), zap.String("bus_name", b.busName))
			b.errCh <- err
		})
		_ = errHandler

		if b.handlerIsSynced(h) {
			go func() {
				defer b.safeRecover(ctx, "synced handler")
				msgCtx, err := consumer.Messages()
				if err != nil {
					if strings.Contains(err.Error(), "Server Shutdown") {
						if strings.EqualFold(os.Getenv("GO_ENV"), "TEST") {
							log.Global().DebugContext(ctx, "eventing: consume error", zap.Error(err))
						} else {
							log.Global().WarnContext(ctx, "eventing: consume error", zap.Error(err))
						}
						return
					}
					if errors.Is(err, nats.ErrConsumerDeleted) || errors.Is(err, jetstream.ErrConsumerDeleted) {
						log.Global().WarnContext(ctx, "eventing: consumer deleted", zap.Error(err), zap.String("stream", b.streamName), zap.String("consumer", consumerName), zap.String("bus_name", b.busName))
						return
					}
					log.Global().ErrorContext(ctx, "eventing: consume error", zap.Error(err), zap.String("stream", b.streamName), zap.String("consumer", consumerName), zap.String("bus_name", b.busName))
					b.errCh <- err
				}
				consumerContext = msgCtx
				handler := b.handler(b.cctx, m, h)
				for {
					select {
					case <-b.cctx.Done():
						log.Global().InfoContext(b.cctx, "eventing: closing NATS event bus consumer", zap.String("identifier", b.appID), zap.String("stream_name", b.streamName), zap.String("bus_name", b.busName), zap.String("type", "nats"))
						return
					default:
						msg, err := msgCtx.Next()
						if err != nil {
							if strings.Contains(err.Error(), "Server Shutdown") {
								if strings.EqualFold(os.Getenv("GO_ENV"), "TEST") {
									log.Global().DebugContext(ctx, "eventing: consume error", zap.Error(err))
								} else {
									log.Global().WarnContext(ctx, "eventing: consume error", zap.Error(err))
								}
								return
							}
							log.Global().ErrorContext(ctx, "eventing: consume error", zap.Error(err), zap.String("stream", b.streamName), zap.String("consumer", consumerName), zap.String("bus_name", b.busName))
							b.errCh <- err
							continue
						}
						handler(msg)
					}
				}

			}()
			return nil
		}

		consumerContext, err = consumer.Consume(b.handler(b.cctx, m, h))
		if err != nil {
			return errors.WithStack(fmt.Errorf("could not consume: %w %s %s", err, b.busName, consumerName))
		}

		return nil
	}, retry.Attempts(3), retry.MaxDelay(5*time.Second), retry.LastErrorOnly(true))
	if err != nil {
		return errors.WithStack(fmt.Errorf("could not subscribe to queue: %w (%s %s %s)", err, subject, consumerName, b.streamName))
	}

	// capture the subscription of ephemeral consumers so we can unsubscribe when we exit.
	if b.handlerIsEphemeral(h) {
		b.unsubscribe = append(b.unsubscribe, func() {
			defer b.safeRecover(ctx, "unsubscribe")
			if consumerContext != nil {
				consumerContext.Stop()
			}
		})
	}

	// Register handler.
	b.registered[h.HandlerType()] = struct{}{}

	b.wg.Add(1)

	// Handle until context is cancelled.
	go b.handle(sub)

	return nil
}

func (b *EventBus) safeRecover(ctx context.Context, source string) {
	if r := recover(); r != nil {
		cctx := b.cctx
		if err, ok := r.(error); ok {
			if (strings.Contains(err.Error(), "Server Shutdown") ||
				strings.Contains(err.Error(), "connection closed") ||
				strings.Contains(err.Error(), "context cancelled") ||
				strings.Contains(err.Error(), "send on closed channel")) &&
				(ctx.Err() != nil || (cctx != nil && cctx.Err() != nil)) {
				log.Global().InfoContext(ctx, "eventing: panic recovered", zap.Error(err), zap.String("identifier", b.appID), zap.String("stream_name", b.streamName), zap.String("bus_name", b.busName), zap.String("type", "nats"), zap.String("source", source))
				return
			}
			log.Global().ErrorContext(ctx, "eventing: panic recovered", zap.Error(err), zap.String("type", "nats"), zap.String("bus_name", b.busName), zap.String("source", source))
		} else {
			log.Global().ErrorContext(ctx, "eventing: panic recovered", zap.Any("panic", r), zap.String("type", "nats"), zap.String("bus_name", b.busName), zap.String("source", source))
		}
	}
}

// Errors implements the Errors method of the eventing.EventBus interface.
func (b *EventBus) Errors() <-chan error {
	return b.errCh
}

// handlerIsEphemeral traverses the middleware chain and checks for the
// ephemeral middleware and queries its status.
func (b *EventBus) handlerIsEphemeral(h eventing.EventHandler) bool {
	for {
		if obs, ok := h.(ephemeral.EphemeralHandler); ok {
			return obs.IsEphemeralHandler()
		} else if c, ok := h.(eventing.EventHandlerChain); ok {
			if h = c.InnerHandler(); h != nil {
				continue
			}
		}
		return false
	}
}

// handlerIsWaiter traverses the middleware chain and checks for the
// waiter middleware and queries its status.
func (b *EventBus) handlerIsWaiter(h eventing.EventHandler) bool {
	for {
		if obs, ok := h.(consumerwaiter.ConsumerWaitHandler); ok {
			return obs.WaitForCurrentEvents()
		} else if c, ok := h.(eventing.EventHandlerChain); ok {
			if h = c.InnerHandler(); h != nil {
				continue
			}
		}
		return false
	}
}

// handlerIsSynced traverses the middleware chain and checks for the
// sync middleware and queries its status.
func (b *EventBus) handlerIsSynced(h eventing.EventHandler) bool {
	for {
		if obs, ok := h.(consumersync.ConsumerSyncHandler); ok {
			return obs.SyncEvents()
		} else if c, ok := h.(eventing.EventHandlerChain); ok {
			if h = c.InnerHandler(); h != nil {
				continue
			}
		}
		return false
	}
}

// handlerIsConsumerOptionsProvider traverses the middleware chain and checks for the
// middleware and queries its status.
func (b *EventBus) handlerIsConsumerOptionsProvider(h eventing.EventHandler, defaultOptions *jetstream.ConsumerConfig) (*jetstream.ConsumerConfig, bool) {
	for {
		if obs, ok := h.(consumeroptions.ConsumerOptionsProvider); ok {
			return obs.ConsumerOptions(defaultOptions), true
		} else if c, ok := h.(eventing.EventHandlerChain); ok {
			if h = c.InnerHandler(); h != nil {
				continue
			}
		}
		return nil, false
	}
}

// handlerIsDeliveryPolicyProvider traverses the middleware chain and checks for the
// middleware and queries its status.
func (b *EventBus) handlerIsDeliveryPolicyProvider(h eventing.EventHandler) (jetstream.DeliverPolicy, uint64, bool) {
	for {
		if obs, ok := h.(consumeroptions.DeliveryPolicyProvider); ok {
			return obs.DeliveryPolicy(), obs.StartSequence(), true
		} else if c, ok := h.(eventing.EventHandlerChain); ok {
			if h = c.InnerHandler(); h != nil {
				continue
			}
		}
		return jetstream.DeliverAllPolicy, 0, false
	}
}

// handlerIsConsumerResetter traverses the middleware chain and checks for the
// middleware and queries its status.
func (b *EventBus) handlerIsConsumerResetter(h eventing.EventHandler) bool {
	for {
		if obs, ok := h.(consumerresetter.ConsumerResetterHandler); ok {
			return obs.ResetConsumer()
		} else if c, ok := h.(eventing.EventHandlerChain); ok {
			if h = c.InnerHandler(); h != nil {
				continue
			}
		}
		return false
	}
}

// handlerIsSequenced traverses the middleware chain and checks for the
// middleware and queries its status.
func (b *EventBus) handlerIsSequenced(h eventing.EventHandler) bool {
	for {
		if obs, ok := h.(sequenced.SequencedHandler); ok {
			return obs.Sequenced()
		} else if c, ok := h.(eventing.EventHandlerChain); ok {
			if h = c.InnerHandler(); h != nil {
				continue
			}
		}
		return false
	}
}

// Close implements the Close method of the eventing.EventBus interface.
func (b *EventBus) Close() error {
	defer b.safeRecover(b.cctx, "close")
	b.cancel()
	b.wg.Wait()

	log.Global().InfoContext(b.cctx, "eventing: closing NATS event bus", zap.String("identifier", b.appID), zap.String("stream_name", b.streamName), zap.String("type", "nats"))
	for reg := range b.registered {
		log.Global().InfoContext(b.cctx, "eventing: closing NATS event bus: registered handler", zap.String("identifier", b.appID), zap.String("stream_name", b.streamName), zap.String("type", "nats"), zap.String("handler", string(reg)))
	}

	// unsubscribe any ephemeral subscribers we created.
	for _, unSub := range b.unsubscribe {
		unSub()
	}

	if _, err := b.pool.Put(b.conn); err != nil {
		return fmt.Errorf("could not return connection to pool: %w", err)
	}

	select {
	case _, ok := <-b.errCh:
		if ok {
			close(b.errCh)
		}
	default:
		close(b.errCh)
	}

	return nil
}

// Handles all events coming in on the channel.
func (b *EventBus) handle(sub *nats.Subscription) {
	defer b.wg.Done()

	<-b.cctx.Done()
	if errors.Is(b.cctx.Err(), context.Canceled) {
		log.Global().WarnContext(b.cctx, "eventing: context error in NATS event bus", zap.Error(b.cctx.Err()), zap.String("type", "nats"), zap.String("bus_name", b.busName))
	}
}

func (b *EventBus) handler(ctx context.Context, m eventing.EventMatcher, h eventing.EventHandler) func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		defer b.safeRecover(ctx, "handler")
		md, err := msg.Metadata()
		if err != nil {
			select {
			case b.errCh <- &eventing.EventBusError{Err: err, Ctx: ctx}:
			default:
				log.Global().ErrorContext(ctx, "eventing: missed error in NATS event bus", zap.String("type", "nats"), zap.Error(err), zap.String("bus_name", b.busName))
			}
		}

		opts := []eventing.EventOption{}
		if b.handlerIsSequenced(h) {
			opts = append(opts, eventing.ForSequence(0, md.Sequence.Stream))
		}

		headers := msg.Headers()
		if strings.EqualFold(headers.Get(EventUnregisteredHdr), "true") {
			opts = append(opts, eventing.AsUnregistered())
		}

		ctx := context.WithoutCancel(ctx)
		event, ctx, err := b.codec.UnmarshalEvent(ctx, msg.Data(), opts...)
		if err != nil {
			err = fmt.Errorf("could not unmarshal event: %w", err)
			select {
			case b.errCh <- &eventing.EventBusError{Err: err, Ctx: ctx}:
			default:
				log.Global().ErrorContext(ctx, "eventing: missed error in NATS event bus", zap.String("type", "nats"), zap.Error(err), zap.String("bus_name", b.busName))
			}
			_ = msg.Nak()

			return
		}

		// Ignore non-matching events.
		if !m.Match(event) {
			_ = msg.Ack()

			return
		}

		// Handle the event if it did match.
		if err := h.HandleEvent(ctx, event); err != nil {
			err = fmt.Errorf("could not handle event (%s) %s: %w", h.HandlerType(), b.busName, err)
			select {
			case b.errCh <- &eventing.EventBusError{Err: err, Ctx: ctx, Event: event}:
			default:
				log.Global().ErrorContext(ctx, "eventing: missed error in NATS event bus", zap.String("type", "nats"), zap.Error(err), zap.String("bus_name", b.busName))
			}

			_ = msg.Nak()

			return
		}

		err = msg.Ack()
		if err != nil {
			log.Global().ErrorContext(ctx, "eventing: failed to ack message", zap.Error(err), zap.String("type", "nats"), zap.String("bus_name", b.busName))
		}
	}
}

func CreateConsumerSubject(streamName string, m eventing.EventMatcher) string {
	aggregateMatch := "*"
	//eventMatch := "*"

	switch m := m.(type) {
	case eventing.MatchEvents:
		// Supports only matching one event, otherwise its wildcard.
		//if len(m) == 1 {
		//	eventMatch = m[0].String()
		//}
		return fmt.Sprintf("%s.%s.>", streamName, aggregateMatch)
	case eventing.MatchAggregates:
		// Supports only matching one aggregate, otherwise its wildcard.
		if len(m) == 1 {
			aggregateMatch = m[0].String()
		}
		return fmt.Sprintf("%s.%s.>", streamName, aggregateMatch)
	case eventing.MatchEventSubject:
		tokens := m.GetSubject().Tokens()

		subjectTokens := make([]string, len(tokens))
		for i, token := range tokens {
			if token.Key() == "aggregate_type" {
				subjectTokens[i] = fmt.Sprint(token.Value())
				continue
			}
			subjectTokens[i] = "*"
		}
		return fmt.Sprintf("%s.%s", streamName, strings.Join(subjectTokens, "."))
	case eventing.MatchEventSubjectExact:
		tokens := m.GetSubject().Tokens()
		desiredPos := m.TokenMatcher.TokenPos
		if desiredPos < 0 {
			desiredPos = m.GetSubject().SubjectTokenPosition()
		}

		subjectTokens := make([]string, len(tokens))
		for i, token := range tokens {
			if token.Key() == "aggregate_type" {
				subjectTokens[i] = fmt.Sprint(token.Value())
				continue
			}
			if token.Position() == desiredPos {
				subjectTokens[i] = m.TokenMatcher.ID
				continue
			}
			subjectTokens[i] = "*"
		}
		return fmt.Sprintf("%s.%s", streamName, strings.Join(subjectTokens, "."))
	}

	return fmt.Sprintf("%s.%s.>", streamName, aggregateMatch)
}
