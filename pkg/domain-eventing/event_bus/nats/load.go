package nats

import (
	"context"
	goJson "encoding/json"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	pool "github.com/octu0/nats-pool"

	eventing "github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
)

func lastMsgForSubjectJetStreamCtx(ctx context.Context, js nats.JetStreamContext, stream string, subject string) (*natsStoredMsg, error) {
	rawMsg, err := js.GetLastMsg(stream, subject)
	if err != nil {
		if errors.Is(err, jetstream.ErrNoMessages) || errors.Is(err, jetstream.ErrMsgNotFound) {
			return &natsStoredMsg{}, nil
		} else if errors.Is(err, nats.ErrMsgNotFound) {
			return &natsStoredMsg{}, nil
		}
		return nil, errors.WithStack(err)
	}

	return &natsStoredMsg{
		Sequence: rawMsg.Sequence,
	}, nil
}

func LoadJetStream(ctx context.Context, js jetstream.JetStream, stream string, subject string, codec eventing.EventCodec) ([]common.Event, error) {
	lastMsg, err := lastMsgForSubjectJetStream(ctx, js, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	consumer, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{
		FilterSubjects:    []string{subject},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxResetAttempts:  5,
		InactiveThreshold: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	events := make([]common.Event, 0)
	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return events, nil
	}

	res, err := consumer.FetchNoWait(int(initialPending))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priorSeq := uint64(0)
	for msg := range res.Messages() {
		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		headers := msg.Headers()
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(priorSeq, seq),
		}
		if headers.Get(EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		//event, _, err := codec.UnmarshalEvent(ctx, msg.Data)
		if err != nil {
			return nil, err
		}
		priorSeq = seq

		events = append(events, event)

		if seq == lastMsg.Sequence {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

func LoadJetStreamCtx(ctx context.Context, js nats.JetStreamContext, stream string, subject string, codec eventing.EventCodec) ([]common.Event, error) {
	lastMsg, err := lastMsgForSubjectJetStreamCtx(ctx, js, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	sopts := []nats.SubOpt{
		nats.OrderedConsumer(),
		nats.DeliverAll(),
	}

	sub, err := js.SubscribeSync(subject, sopts...)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe() //nolint

	events := make([]common.Event, 0)
	priorSeq := uint64(0)
	for {
		msg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			return nil, err
		}

		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		headers := msg.Header
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(priorSeq, seq),
		}
		if headers.Get(EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data, eventOpts...)

		if err != nil {
			return nil, err
		}
		priorSeq = seq
		_ = priorSeq

		events = append(events, event)

		if seq == lastMsg.Sequence {
			break
		}
	}

	return events, nil
}

func LoadJetStreamFrom(ctx context.Context, js jetstream.JetStream, stream string, subject string, startSequence uint64, codec eventing.EventCodec) ([]common.Event, error) {
	lastMsg, err := lastMsgForSubjectJetStream(ctx, js, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	if lastMsg.Sequence < startSequence {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	consumer, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{
		FilterSubjects:    []string{subject},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxResetAttempts:  5,
		OptStartSeq:       startSequence,
		InactiveThreshold: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	events := make([]common.Event, 0)
	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return events, nil
	}

	res, err := consumer.FetchNoWait(int(initialPending))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priorSeq := uint64(0)
	for msg := range res.Messages() {
		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		headers := msg.Headers()
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(priorSeq, seq),
		}
		if headers.Get(EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		//event, _, err := codec.UnmarshalEvent(ctx, msg.Data)
		if err != nil {
			return nil, err
		}
		priorSeq = seq

		events = append(events, event)

		if seq == lastMsg.Sequence {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

func Load(ctx context.Context, connectionPool *pool.ConnPool, stream string, subject string, codec eventing.EventCodec) ([]common.Event, error) {
	conn, err := connectionPool.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = connectionPool.Put(conn)
	}()

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	lastMsg, err := lastMsgForSubject(ctx, conn, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	consumer, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{
		FilterSubjects:    []string{subject},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxResetAttempts:  5,
		InactiveThreshold: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	events := make([]common.Event, 0)
	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return events, nil
	}

	res, err := consumer.FetchNoWait(int(initialPending))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priorSeq := uint64(0)
	for msg := range res.Messages() {
		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		headers := msg.Headers()
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(priorSeq, seq),
		}
		if headers.Get(EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		//event, _, err := codec.UnmarshalEvent(ctx, msg.Data)
		if err != nil {
			return nil, err
		}
		priorSeq = seq

		events = append(events, event)

		if seq == lastMsg.Sequence {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

func LoadFrom(ctx context.Context, connectionPool *pool.ConnPool, stream string, subject string, startSequence uint64, codec eventing.EventCodec) ([]common.Event, error) {
	conn, err := connectionPool.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = connectionPool.Put(conn)
	}()

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	lastMsg, err := lastMsgForSubject(ctx, conn, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	if lastMsg.Sequence < startSequence {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	consumer, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{
		FilterSubjects:    []string{subject},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxResetAttempts:  5,
		OptStartSeq:       startSequence,
		InactiveThreshold: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	events := make([]common.Event, 0)
	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return events, nil
	}

	res, err := consumer.FetchNoWait(int(initialPending))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priorSeq := uint64(0)
	for msg := range res.Messages() {

		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		headers := msg.Headers()
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(priorSeq, seq),
		}
		if headers.Get(EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		//event, _, err := codec.UnmarshalEvent(ctx, msg.Data)
		if err != nil {
			return nil, err
		}
		priorSeq = seq

		events = append(events, event)

		if seq == lastMsg.Sequence {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

func LoadBus(ctx context.Context, eventBus eventing.EventBus, stream string, subject string, codec eventing.EventCodec) ([]common.Event, error) {
	if eventBus == nil {
		return nil, errors.WithStack(errors.New("event bus is nil"))
	}

	natsBus, ok := busIsEventBus(eventBus)
	if !ok {
		return nil, errors.WithStack(errors.New("event bus is not nats bus"))
	}

	js := natsBus.jetstream

	lastMsg, err := lastMsgForSubjectJetStream(ctx, js, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	consumer, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{
		FilterSubjects:    []string{subject},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxResetAttempts:  5,
		InactiveThreshold: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	events := make([]common.Event, 0)
	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return events, nil
	}

	res, err := consumer.FetchNoWait(int(initialPending))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priorSeq := uint64(0)
	for msg := range res.Messages() {

		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		headers := msg.Headers()
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(priorSeq, seq),
		}
		if headers.Get(EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		//event, _, err := codec.UnmarshalEvent(ctx, msg.Data)
		if err != nil {
			return nil, err
		}
		priorSeq = seq

		events = append(events, event)

		if seq == lastMsg.Sequence {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

func LoadBusFrom(ctx context.Context, eventBus eventing.EventBus, stream string, subject string, startSequence uint64, codec eventing.EventCodec) ([]common.Event, error) {
	if eventBus == nil {
		return nil, errors.WithStack(errors.New("event bus is nil"))
	}

	natsBus, ok := busIsEventBus(eventBus)
	if !ok {
		return nil, errors.WithStack(errors.New("event bus is not nats bus"))
	}

	js := natsBus.jetstream

	lastMsg, err := lastMsgForSubjectJetStream(ctx, js, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	if lastMsg.Sequence < startSequence {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	consumer, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{
		FilterSubjects:    []string{subject},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxResetAttempts:  5,
		OptStartSeq:       startSequence,
		InactiveThreshold: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	events := make([]common.Event, 0)
	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return events, nil
	}

	res, err := consumer.FetchNoWait(int(initialPending))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	priorSeq := uint64(0)
	for msg := range res.Messages() {

		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		headers := msg.Headers()
		eventOpts := []eventing.EventOption{
			eventing.ForSequence(priorSeq, seq),
		}
		if headers.Get(EventUnregisteredHdr) == "true" {
			eventOpts = append(eventOpts, eventing.AsUnregistered())
		}

		event, _, err := codec.UnmarshalEvent(ctx, msg.Data(), eventOpts...)
		//event, _, err := codec.UnmarshalEvent(ctx, msg.Data)
		if err != nil {
			return nil, err
		}
		priorSeq = seq

		events = append(events, event)

		if seq == lastMsg.Sequence {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

func LoadRaw(ctx context.Context, connectionPool *pool.ConnPool, stream string, subject string, codec eventing.EventCodec) ([]jetstream.Msg, error) {
	conn, err := connectionPool.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = connectionPool.Put(conn)
	}()

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	lastMsg, err := lastMsgForSubject(ctx, conn, stream, subject)
	if err != nil {
		return nil, err
	}

	if lastMsg.Sequence == 0 {
		return nil, nil
	}

	// Ephemeral ordered consumer. read as fast as possible with least overhead.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	consumer, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{
		FilterSubjects:    []string{subject},
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		MaxResetAttempts:  5,
		InactiveThreshold: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	events := make([]jetstream.Msg, 0)
	cachedInfo := consumer.CachedInfo()
	initialPending := cachedInfo.NumPending

	if initialPending == 0 {
		return events, nil
	}

	res, err := consumer.FetchNoWait(int(initialPending))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for msg := range res.Messages() {
		md, err := msg.Metadata()
		if err != nil {
			return nil, fmt.Errorf("unpack: failed to get metadata: %s", err)
		}
		seq := md.Sequence.Stream

		events = append(events, msg)

		if seq == lastMsg.Sequence {
			break
		}
	}

	if res.Error() != nil {
		return nil, errors.WithStack(res.Error())
	}

	return events, nil
}

// LoadLatest fetches the most recent event for a specific subject. The primary use case
// is to use a concrete subject, e.g. "orders.1" corresponding to an
// aggregate/entity identifier. The second use case is to load events for
// a cross-cutting view which can use subject wildcards.
func LoadLatest(ctx context.Context, connectionPool *pool.ConnPool, stream string, subject string, codec eventing.EventCodec) (common.Event, uint64, error) {
	conn, err := connectionPool.Get()
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_, _ = connectionPool.Put(conn)
	}()

	js, err := conn.JetStream()
	if err != nil {
		return nil, 0, err
	}

	rawMsg, err := js.GetLastMsg(stream, subject)
	if err != nil {
		return nil, 0, err
	}
	msg := nats.NewMsg(rawMsg.Subject)
	msg.Data = rawMsg.Data
	msg.Header = rawMsg.Header

	eventOpts := []eventing.EventOption{
		eventing.ForSequence(0, rawMsg.Sequence),
	}
	if rawMsg.Header.Get(EventUnregisteredHdr) == "true" {
		eventOpts = append(eventOpts, eventing.AsUnregistered())
	}

	event, _, err := codec.UnmarshalEvent(ctx, msg.Data, eventOpts...)
	//event, _, err := codec.UnmarshalEvent(ctx, msg.Data, eventing.ForSequence(0, rawMsg.Sequence))
	if err != nil {
		return nil, 0, err
	}

	return event, rawMsg.Sequence, nil
}

// LoadLatestJetStream fetches the most recent event for a specific subject. The primary use case
// is to use a concrete subject, e.g. "orders.1" corresponding to an
// aggregate/entity identifier. The second use case is to load events for
// a cross-cutting view which can use subject wildcards.
func LoadLatestJetStream(ctx context.Context, js jetstream.JetStream, stream string, subject string, codec eventing.EventCodec) (common.Event, uint64, error) {
	jsStream, err := js.Stream(ctx, stream)
	if err != nil {
		return nil, 0, err
	}

	rawMsg, err := jsStream.GetLastMsgForSubject(ctx, subject)
	if err != nil {
		return nil, 0, err
	}
	msg := nats.NewMsg(rawMsg.Subject)
	msg.Data = rawMsg.Data
	msg.Header = rawMsg.Header

	eventOpts := []eventing.EventOption{
		eventing.ForSequence(0, rawMsg.Sequence),
	}
	if rawMsg.Header.Get(EventUnregisteredHdr) == "true" {
		eventOpts = append(eventOpts, eventing.AsUnregistered())
	}

	event, _, err := codec.UnmarshalEvent(ctx, msg.Data, eventOpts...)
	//event, _, err := codec.UnmarshalEvent(ctx, msg.Data, eventing.ForSequence(0, rawMsg.Sequence))
	if err != nil {
		return nil, 0, err
	}

	return event, rawMsg.Sequence, nil
}

// LoadLatestBus fetches the most recent event for a specific subject. The primary use case
// is to use a concrete subject, e.g. "orders.1" corresponding to an
// aggregate/entity identifier. The second use case is to load events for
// a cross-cutting view which can use subject wildcards.
func LoadLatestBus(ctx context.Context, eventBus eventing.EventBus, stream string, subject string, codec eventing.EventCodec) (common.Event, uint64, error) {
	if eventBus == nil {
		return nil, 0, errors.WithStack(errors.New("event bus is nil"))
	}

	natsBus, ok := busIsEventBus(eventBus)
	if !ok {
		return nil, 0, errors.WithStack(errors.New("event bus is not nats bus"))
	}

	js := natsBus.jetstream
	return LoadLatestJetStream(ctx, js, stream, subject, codec)
}

type natsApiError struct {
	Code        int    `json:"code"`
	ErrCode     uint16 `json:"err_code"`
	Description string `json:"description"`
}

type natsGetMsgRequest struct {
	LastBySubject string `json:"last_by_subj"`
}

type natsGetMsgResponse struct {
	Type    string         `json:"type"`
	Error   *natsApiError  `json:"error"`
	Message *natsStoredMsg `json:"message"`
}

type natsStoredMsg struct {
	Sequence uint64 `json:"seq"`
}

func lastMsgForSubjectJetStream(ctx context.Context, js jetstream.JetStream, stream string, subject string) (*natsStoredMsg, error) {
	jsStream, err := js.Stream(ctx, stream)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rawMsg, err := jsStream.GetLastMsgForSubject(ctx, subject)
	if err != nil {
		if errors.Is(err, jetstream.ErrNoMessages) || errors.Is(err, jetstream.ErrMsgNotFound) {
			return &natsStoredMsg{}, nil
		} else if errors.Is(err, nats.ErrMsgNotFound) {
			return &natsStoredMsg{}, nil
		}
		return nil, errors.WithStack(err)
	}

	return &natsStoredMsg{
		Sequence: rawMsg.Sequence,
	}, nil
}

func lastMsgForSubject(ctx context.Context, conn *nats.Conn, stream string, subject string) (*natsStoredMsg, error) {
	rsubject := fmt.Sprintf("$JS.API.STREAM.MSG.GET.%s", stream)

	data, _ := goJson.Marshal(&natsGetMsgRequest{
		LastBySubject: subject,
	})

	msg, err := conn.RequestWithContext(ctx, rsubject, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var rep natsGetMsgResponse
	err = goJson.Unmarshal(msg.Data, &rep)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if rep.Error != nil {
		if rep.Error.Code == 404 {
			return &natsStoredMsg{}, nil
		}
		return nil, fmt.Errorf("%s (%d)", rep.Error.Description, rep.Error.Code)
	}

	return rep.Message, nil
}

func busIsEventBus(bus eventing.EventBus) (*EventBus, bool) {
	for {
		if obs, ok := bus.(*EventBus); ok {
			return obs, true
		} else if c, ok := bus.(eventing.EventBusChain); ok {
			if bus = c.InnerBus(); bus != nil {
				continue
			}
		}
		return nil, false
	}
}
