# app

Unified EXPLODING kittens service launcher

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

## app gameengineserver

Run as gameengineserver service

```
app gameengineserver [flags]
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
      --grpc-port int                             GRPC Port to listen on (default 50054)
  -h, --help                                      help for gameengineserver
      --id string                                 Unique identifier for this services
      --token-signing-key string                  Signing key used for service to service tokens
```

### Environment Variables

- GAMEENGINESERVER_DB_MIGRATIONS_URL :: `gameengineserver.db.migrations.url` Database connection migrations URL
- GAMEENGINESERVER_DB_POSTGRES_CONNECTION_MAX_IDLETIME :: `gameengineserver.db.postgres.max_idletime` Max connection idle time in seconds
- GAMEENGINESERVER_DB_POSTGRES_CONNECTION_MAX_LIFETIME :: `gameengineserver.db.postgres.max_lifetime` Max connection lifetime in seconds
- GAMEENGINESERVER_DB_POSTGRES_MAX_IDLE_CONNECTIONS :: `gameengineserver.db.postgres.max_idle_connections` Maximum number of idle connections
- GAMEENGINESERVER_DB_POSTGRES_MAX_OPEN_CONNECTIONS :: `gameengineserver.db.postgres.max_open_connections` Maximum number of connections
- GAMEENGINESERVER_DB_POSTGRES_TIMEOUT :: `gameengineserver.db.postgres.timeout` Timeout for postgres connection
- GAMEENGINESERVER_DB_READ_URL :: `gameengineserver.db.read.url` Database connection readonly URL
- GAMEENGINESERVER_DB_URL :: `gameengineserver.db.url` Database connection URL
- GAMEENGINESERVER_GRPC_PORT :: `gameengineserver.grpc.port` GRPC Port to listen on
- GAMEENGINESERVER_ID :: `gameengineserver.id` Unique identifier for this services
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

## app userserver

Run as userserver service

```
app userserver [flags]
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
  -h, --help                                      help for userserver
      --id string                                 Unique identifier for this services
      --token-signing-key string                  Signing key used for service to service tokens
```

### Environment Variables

- USERSERVER_DB_MIGRATIONS_URL :: `userserver.db.migrations.url` Database connection migrations URL
- USERSERVER_DB_POSTGRES_CONNECTION_MAX_IDLETIME :: `userserver.db.postgres.max_idletime` Max connection idle time in seconds
- USERSERVER_DB_POSTGRES_CONNECTION_MAX_LIFETIME :: `userserver.db.postgres.max_lifetime` Max connection lifetime in seconds
- USERSERVER_DB_POSTGRES_MAX_IDLE_CONNECTIONS :: `userserver.db.postgres.max_idle_connections` Maximum number of idle connections
- USERSERVER_DB_POSTGRES_MAX_OPEN_CONNECTIONS :: `userserver.db.postgres.max_open_connections` Maximum number of connections
- USERSERVER_DB_POSTGRES_TIMEOUT :: `userserver.db.postgres.timeout` Timeout for postgres connection
- USERSERVER_DB_READ_URL :: `userserver.db.read.url` Database connection readonly URL
- USERSERVER_DB_URL :: `userserver.db.url` Database connection URL
- USERSERVER_GRPC_PORT :: `userserver.grpc.port` GRPC Port to listen on
- USERSERVER_ID :: `userserver.id` Unique identifier for this services
- USERSERVER_SECRETS_TOKEN_SIGNING_KEY :: `userserver.secrets.token_signing_key` Signing key used for service to service tokens
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
