# app

Exploding Kittens Client Server

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

## app clientserver

Run as clientserver service

```
app clientserver [flags]
```

### Options

```
      --dataprovider-url string        Dataprovider connection URL (default "http://dataprovider:50055")
      --gameengineserver-url string    Game Engine Server connection URL (default "http://gameengineserver:50054")
      --grpc-port int                  GRPC Port to listen on (default 50051)
  -h, --help                           help for clientserver
      --id string                      Unique identifier for this services
      --nats-consumer-replicas int     Number of times to replicate consumers (default 1)
      --nats-consumer-storage string   Storage type to use for consumers (default "memory")
      --nats-stream-replicas int       Number of times to replicate steams (default 1)
      --nats-stream-storage string     Storage type to use for streams (default "memory")
      --nats-url string                Comma separated list of NATS endpoints (default "nats:4222")
      --token-signing-key string       Signing key used for service to service tokens
      --userserver-url string          Userserver connection URL (default "http://userserver:50052")
```

### Environment Variables

- CLIENTSERVER_DATAPROVIDER_URL :: `clientserver.dataprovider.url` Dataprovider connection URL
- CLIENTSERVER_GAMEENGINESERVER_URL :: `clientserver.gameengineserver.url` Game Engine Server connection URL
- CLIENTSERVER_GRPC_PORT :: `clientserver.grpc.port` GRPC Port to listen on
- CLIENTSERVER_ID :: `clientserver.id` Unique identifier for this services
- CLIENTSERVER_NATS_CONSUMER_REPLICAS :: `clientserver.nats.consumer.replicas` Number of times to replicate consumers
- CLIENTSERVER_NATS_CONSUMER_STORAGE :: `clientserver.nats.consumer.storage` Storage type to use for consumers
- CLIENTSERVER_NATS_STREAM_REPLICAS :: `clientserver.nats.stream.replicas` Number of times to replicate steams
- CLIENTSERVER_NATS_STREAM_STORAGE :: `clientserver.nats.stream.storage` Storage type to use for streams
- CLIENTSERVER_NATS_URL :: `clientserver.nats.url` Comma separated list of NATS endpoints
- CLIENTSERVER_SECRETS_TOKEN_SIGNING_KEY :: `clientserver.secrets.token_signing_key` Signing key used for service to service tokens
- CLIENTSERVER_USERSERVER_URL :: `clientserver.userserver.url` Userserver connection URL
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
