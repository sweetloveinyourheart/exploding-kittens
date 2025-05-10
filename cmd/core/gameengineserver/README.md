# app

Exploding Kittens Game Engine Server

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
      --log-level string           log level to use (default "info")
  -s, --service string             which service to run
```

### Environment Variables

- EXPLODING_KITTENS_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LOG_LEVEL :: `log.level` log level to use
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
      --log-level string           log level to use (default "info")
  -s, --service string             which service to run
```

### Environment Variables inherited from parent commands

- EXPLODING_KITTENS_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LOG_LEVEL :: `log.level` log level to use
- EXPLODING_KITTENS_SERVICE :: `service` which service to run
```

## app gameengineserver

Run as gameengineserver service

```
app gameengineserver [flags]
```

### Options

```
      --dataprovider-url string        Data provider connection URL (default "http://dataprovider:50055")
      --grpc-port int                  GRPC Port to listen on (default 50054)
  -h, --help                           help for gameengineserver
      --id string                      Unique identifier for this services
      --nats-consumer-replicas int     Number of times to replicate consumers (default 1)
      --nats-consumer-storage string   Storage type to use for consumers (default "memory")
      --nats-stream-replicas int       Number of times to replicate steams (default 1)
      --nats-stream-storage string     Storage type to use for streams (default "memory")
      --nats-url string                Comma separated list of NATS endpoints (default "nats:4222")
      --token-signing-key string       Signing key used for service to service tokens
```

### Environment Variables

- GAMEENGINESERVER_DATAPROVIDER_URL :: `gameengineserver.dataprovider.url` Data provider connection URL
- GAMEENGINESERVER_GRPC_PORT :: `gameengineserver.grpc.port` GRPC Port to listen on
- GAMEENGINESERVER_ID :: `gameengineserver.id` Unique identifier for this services
- GAMEENGINESERVER_NATS_CONSUMER_REPLICAS :: `gameengineserver.nats.consumer.replicas` Number of times to replicate consumers
- GAMEENGINESERVER_NATS_CONSUMER_STORAGE :: `gameengineserver.nats.consumer.storage` Storage type to use for consumers
- GAMEENGINESERVER_NATS_STREAM_REPLICAS :: `gameengineserver.nats.stream.replicas` Number of times to replicate steams
- GAMEENGINESERVER_NATS_STREAM_STORAGE :: `gameengineserver.nats.stream.storage` Storage type to use for streams
- GAMEENGINESERVER_NATS_URL :: `gameengineserver.nats.url` Comma separated list of NATS endpoints
- GAMEENGINESERVER_SECRETS_TOKEN_SIGNING_KEY :: `gameengineserver.secrets.token_signing_key` Signing key used for service to service tokens
```

### Options inherited from parent commands

```
      --config string              config file (default is $HOME/.EXPLODING-poker/app.yaml)
      --healthcheck-host string    Host to listen on for services that support a health check (default "localhost")
      --healthcheck-port int       Port to listen on for services that support a health check (default 5051)
      --healthcheck-web-port int   Port to listen on for services that support a health check (default 5052)
      --log-level string           log level to use (default "info")
  -s, --service string             which service to run
```

### Environment Variables inherited from parent commands

- EXPLODING_KITTENS_HEALTHCHECK_HOST :: `healthcheck.host` Host to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_PORT :: `healthcheck.port` Port to listen on for services that support a health check
- EXPLODING_KITTENS_HEALTHCHECK_WEB_PORT :: `healthcheck.web.port` Port to listen on for services that support a health check
- LOG_LEVEL :: `log.level` log level to use
- EXPLODING_KITTENS_SERVICE :: `service` which service to run
```


## Configuration Paths

 - /etc/exploding-kittens/schema.yaml
 - $HOME/.exploding-kittens/schema.yaml
 - ./schema.yaml

### Common

## Testing
```go test ./cmd/exploding-kittens/```
