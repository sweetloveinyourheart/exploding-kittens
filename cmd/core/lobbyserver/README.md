# app

Exploding Kittens Lobby Server

```
app [flags]
```

### Options

```
      --config string              config file (default is $HOME/.EXPLODING-poker/app.yaml)
      --healthcheck-host string    Host to listen on for services that support a health check (default "localhost")
      --healthcheck-port int       Port to listen on for services that support a health check (default 5051)
      --healthcheck-web-port int   Port to listen on for services that support a health check (default 5052)
  -h, --help                       help for app
      --jaeger-url string          URL to send Jaeger data to
      --log-level string           log level to use (default "info")
      --otel-url string            URL to send OpenTelemetry data to (default "localhost:30080")
  -s, --service string             which service to run
```

### Environment Variables

- EXPLODING_KITTENS_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LIVE_DEALER_JAEGER_URL :: `jaeger.url` URL to send Jaeger data to
- LOG_LEVEL :: `log.level` log level to use
- LIVE_DEALER_OTEL_URL :: `otel.url` URL to send OpenTelemetry data to
- EXPLODING_KITTENS_SERVICE :: `service` which service to run
```

## app check

Health check commands

### Synopsis

Commands for running health checks

### Options

```
  -h, --help   help for check
```

### Environment Variables

```

### Options inherited from parent commands

```
      --config string              config file (default is $HOME/.EXPLODING-poker/app.yaml)
      --healthcheck-host string    Host to listen on for services that support a health check (default "localhost")
      --healthcheck-port int       Port to listen on for services that support a health check (default 5051)
      --healthcheck-web-port int   Port to listen on for services that support a health check (default 5052)
      --jaeger-url string          URL to send Jaeger data to
      --log-level string           log level to use (default "info")
      --otel-url string            URL to send OpenTelemetry data to (default "localhost:30080")
  -s, --service string             which service to run
```

### Environment Variables inherited from parent commands

- EXPLODING_KITTENS_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LIVE_DEALER_JAEGER_URL :: `jaeger.url` URL to send Jaeger data to
- LOG_LEVEL :: `log.level` log level to use
- LIVE_DEALER_OTEL_URL :: `otel.url` URL to send OpenTelemetry data to
- EXPLODING_KITTENS_SERVICE :: `service` which service to run
```

## app lobbyserver

Run as lobbyserver service

```
app lobbyserver [flags]
```

### Options

```
      --grpc-port int                  GRPC Port to listen on (default 50053)
  -h, --help                           help for lobbyserver
      --id string                      Unique identifier for this services
      --nats-consumer-replicas int     Number of times to replicate consumers (default 1)
      --nats-consumer-storage string   Storage type to use for consumers (default "memory")
      --nats-stream-replicas int       Number of times to replicate steams (default 1)
      --nats-stream-storage string     Storage type to use for streams (default "memory")
      --nats-url string                Comma separated list of NATS endpoints (default "nats:4222")
      --token-signing-key string       Signing key used for service to service tokens
```

### Environment Variables

- LOBBYSERVER_GRPC_PORT :: `lobbyserver.grpc.port` GRPC Port to listen on
- LOBBYSERVER_ID :: `lobbyserver.id` Unique identifier for this services
- LOBBYSERVER_NATS_CONSUMER_REPLICAS :: `lobbyserver.nats.consumer.replicas` Number of times to replicate consumers
- LOBBYSERVER_NATS_CONSUMER_STORAGE :: `lobbyserver.nats.consumer.storage` Storage type to use for consumers
- LOBBYSERVER_NATS_STREAM_REPLICAS :: `lobbyserver.nats.stream.replicas` Number of times to replicate steams
- LOBBYSERVER_NATS_STREAM_STORAGE :: `lobbyserver.nats.stream.storage` Storage type to use for streams
- LOBBYSERVER_NATS_URL :: `lobbyserver.nats.url` Comma separated list of NATS endpoints
- LOBBYSERVER_SECRETS_TOKEN_SIGNING_KEY :: `lobbyserver.secrets.token_signing_key` Signing key used for service to service tokens
```

### Options inherited from parent commands

```
      --config string              config file (default is $HOME/.EXPLODING-poker/app.yaml)
      --healthcheck-host string    Host to listen on for services that support a health check (default "localhost")
      --healthcheck-port int       Port to listen on for services that support a health check (default 5051)
      --healthcheck-web-port int   Port to listen on for services that support a health check (default 5052)
      --jaeger-url string          URL to send Jaeger data to
      --log-level string           log level to use (default "info")
      --otel-url string            URL to send OpenTelemetry data to (default "localhost:30080")
  -s, --service string             which service to run
```

### Environment Variables inherited from parent commands

- EXPLODING_KITTENS_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LIVE_DEALER_JAEGER_URL :: `jaeger.url` URL to send Jaeger data to
- LOG_LEVEL :: `log.level` log level to use
- LIVE_DEALER_OTEL_URL :: `otel.url` URL to send OpenTelemetry data to
- EXPLODING_KITTENS_SERVICE :: `service` which service to run
```


## Configuration Paths

 - /etc/exploding-kittens/schema.yaml
 - $HOME/.exploding-kittens/schema.yaml
 - ./schema.yaml

### Common

## Testing
```go test ./cmd/exploding-kittens/```
