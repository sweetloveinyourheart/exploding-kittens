package tracing

import "go.opentelemetry.io/otel/attribute"

var (
	EventTypeKey = attribute.Key("event.type")

	EntityTypeKey = attribute.Key("entity.type")

	AggregateTypeKey = attribute.Key("aggregate.type")

	AggregateIDKey = attribute.Key("aggregate.id")

	AggregateVersionKey = attribute.Key("aggregate.version")

	CommandTypeKey = attribute.Key("command.type")
)

func EventType(val string) attribute.KeyValue {
	return EventTypeKey.String(val)
}

func EntityType(val string) attribute.KeyValue {
	return EntityTypeKey.String(val)
}

func AggregateType(val string) attribute.KeyValue {
	return AggregateTypeKey.String(val)
}

func AggregateID(val string) attribute.KeyValue {
	return AggregateIDKey.String(val)
}

func CommandType(val string) attribute.KeyValue {
	return CommandTypeKey.String(val)
}

func AggregateVersion(val int64) attribute.KeyValue {
	return attribute.Int64("aggregate.version", val)
}
