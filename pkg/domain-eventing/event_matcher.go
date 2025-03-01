package eventing

import (
	"strings"

	"github.com/sweetloveinyourheart/exploding-kittens/pkg/domain-eventing/common"
	log "github.com/sweetloveinyourheart/exploding-kittens/pkg/logger"
	"github.com/sweetloveinyourheart/exploding-kittens/pkg/stringsutil"
)

// EventMatcher matches, for example on event types, aggregate types etc.
type EventMatcher interface {
	// Match returns true if the matcher matches an event.
	Match(common.Event) bool
}

// MatchEvents matches any of the event types, nil events never match.
type MatchEvents []common.EventType

// Match implements the Match method of the EventMatcher interface.
func (types MatchEvents) Match(e common.Event) bool {
	for _, t := range types {
		if e != nil && e.EventType() == t {
			return true
		}
	}

	return false
}

// MatchAggregates matches any of the aggregate types, nil events never match.
type MatchAggregates []common.AggregateType

// Match implements the Match method of the EventMatcher interface.
func (types MatchAggregates) Match(e common.Event) bool {
	for _, t := range types {
		if e != nil && e.AggregateType() == t {
			return true
		}
	}

	return false
}

// MatchAny matches any of the matchers.
type MatchAny []EventMatcher

// Match implements the Match method of the EventMatcher interface.
func (matchers MatchAny) Match(e common.Event) bool {
	for _, m := range matchers {
		if m.Match(e) {
			return true
		}
	}

	return false
}

// MatchAll matches all of the matchers.
type MatchAll []EventMatcher

// Match implements the Match method of the EventMatcher interface.
func (matchers MatchAll) Match(e common.Event) bool {
	for _, m := range matchers {
		if !m.Match(e) {
			return false
		}
	}

	return true
}

type MatchEventSubject struct {
	subject       common.EventSubject
	AggregateType common.AggregateType
	EventTypes    []common.EventType
}

func NewMatchEventSubject(subject common.EventSubject, agg common.AggregateType, events ...common.EventType) MatchEventSubject {
	if subject == nil {
		log.Global().Warn("can't match a nil subject")
	}
	return MatchEventSubject{
		subject:       subject,
		AggregateType: agg,
		EventTypes:    events,
	}
}

func (m MatchEventSubject) Match(e common.Event) bool {
	if e != nil && e.AggregateType() == m.AggregateType {
		if len(m.EventTypes) == 0 {
			return true
		}
		for _, t := range m.EventTypes {
			if e.EventType() == t {
				return true
			}
		}
	}

	return false
}

func (m MatchEventSubject) GetSubject() common.EventSubject {
	return m.subject
}

type MatchEventSubjectExact struct {
	subject       common.EventSubject
	AggregateType common.AggregateType
	EventTypes    []common.EventType
	TokenMatcher  TokenMatcher
}

func NewMatchEventSubjectExact(subject common.EventSubject, agg common.AggregateType, aggID string, events ...common.EventType) MatchEventSubjectExact {
	return MatchEventSubjectExact{
		subject:       subject,
		AggregateType: agg,
		EventTypes:    events,
		TokenMatcher: TokenMatcher{
			ID:        aggID,
			TokenPos:  subject.SubjectTokenPosition(),
			TokenName: "aggregate_id",
			MatchFunc: func(token string, id string, e common.Event) bool {
				return e.AggregateID() == id
			},
		},
	}
}

type TokenMatcher struct {
	TokenPos  int
	TokenName string
	ID        string
	MatchFunc func(string, string, common.Event) bool
}

func NewTokenMatcher(tokenName string, id string, matchFunc func(string, string, common.Event) bool) TokenMatcher {
	if stringsutil.IsBlank(tokenName) {
		tokenName = "aggregate_id"
	}
	return TokenMatcher{
		TokenName: tokenName,
		ID:        id,
		MatchFunc: matchFunc,
	}
}

func NewMatchEventSubjectExactForToken(subject common.EventSubject, agg common.AggregateType, matcher TokenMatcher, events ...common.EventType) MatchEventSubjectExact {
	if stringsutil.IsBlank(matcher.TokenName) {
		matcher.TokenName = "aggregate_id"
	}
	matcher.TokenPos = subject.SubjectTokenPosition()
	tokens := subject.Tokens()
	for _, t := range tokens {
		if strings.EqualFold(t.Key(), matcher.TokenName) {
			matcher.TokenPos = t.Position()
		}
	}

	return MatchEventSubjectExact{
		subject:       subject,
		AggregateType: agg,
		EventTypes:    events,
		TokenMatcher:  matcher,
	}
}

func (m MatchEventSubjectExact) Match(e common.Event) bool {
	if e != nil && e.AggregateType() == m.AggregateType {
		if m.TokenMatcher.MatchFunc != nil {
			if !m.TokenMatcher.MatchFunc(m.TokenMatcher.TokenName, m.TokenMatcher.ID, e) {
				return false
			}
		}
		if len(m.EventTypes) == 0 {
			return true
		}
		for _, t := range m.EventTypes {
			if e.EventType() == t {
				return true
			}
		}
	}

	return false
}

func (m MatchEventSubjectExact) GetSubject() common.EventSubject {
	return m.subject
}
