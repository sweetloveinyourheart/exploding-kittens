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

## app lobbyserver

Run as lobbyserver service

```
app lobbyserver [flags]
```

### Options

```
      --db-migrations-url string                  Database connection migrations URL
      --db-postgres-connection-max-idletime int   Max connection idle time in seconds (default 180)
      --db-postgres-connection-max-lifetime int   Max connection lifetime in seconds (default 300)
      --db-postgres-max-idle-connections int      Maximum number of idle connections (default 50)
      --db-postgres-max-open-connections int      Maximum number of connections (default 500)
      --db-postgres-timeout int                   Timeout for postgres connection (default 60)
      --db-read-url string                        Database connection readonly URL
      --db-url string                             Database connection URL
      --grpc-port int                             GRPC Port to listen on (default 50052)
  -h, --help                                      help for lobbyserver
      --id string                                 Unique identifier for this services
      --nats-consumer-replicas int                Number of times to replicate consumers (default 1)
      --nats-consumer-storage string              Storage type to use for consumers (default "memory")
      --nats-stream-replicas int                  Number of times to replicate steams (default 1)
      --nats-stream-storage string                Storage type to use for streams (default "memory")
      --nats-url string                           Comma separated list of NATS endpoints (default "nats:4222")
      --token-signing-key string                  Signing key used for service to service tokens
```

### Environment Variables

- LOBBYSERVER_DB_MIGRATIONS_URL :: `lobbyserver.db.migrations.url` Database connection migrations URL
- LOBBYSERVER_DB_POSTGRES_CONNECTION_MAX_IDLETIME :: `lobbyserver.db.postgres.max_idletime` Max connection idle time in seconds
- LOBBYSERVER_DB_POSTGRES_CONNECTION_MAX_LIFETIME :: `lobbyserver.db.postgres.max_lifetime` Max connection lifetime in seconds
- LOBBYSERVER_DB_POSTGRES_MAX_IDLE_CONNECTIONS :: `lobbyserver.db.postgres.max_idle_connections` Maximum number of idle connections
- LOBBYSERVER_DB_POSTGRES_MAX_OPEN_CONNECTIONS :: `lobbyserver.db.postgres.max_open_connections` Maximum number of connections
- LOBBYSERVER_DB_POSTGRES_TIMEOUT :: `lobbyserver.db.postgres.timeout` Timeout for postgres connection
- LOBBYSERVER_DB_READ_URL :: `lobbyserver.db.read.url` Database connection readonly URL
- LOBBYSERVER_DB_URL :: `lobbyserver.db.url` Database connection URL
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
