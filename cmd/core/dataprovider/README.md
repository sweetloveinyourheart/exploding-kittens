# app

Exploding Kittens Data Provider

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

## app dataprovider

Run as dataprovider service

```
app dataprovider [flags]
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
      --grpc-port int                             GRPC Port to listen on (default 50055)
  -h, --help                                      help for dataprovider
      --id string                                 Unique identifier for this services
      --token-signing-key string                  Signing key used for service to service tokens
```

### Environment Variables

- DATAPROVIDER_DB_MIGRATIONS_URL :: `dataprovider.db.migrations.url` Database connection migrations URL
- DATAPROVIDER_DB_POSTGRES_CONNECTION_MAX_IDLETIME :: `dataprovider.db.postgres.max_idletime` Max connection idle time in seconds
- DATAPROVIDER_DB_POSTGRES_CONNECTION_MAX_LIFETIME :: `dataprovider.db.postgres.max_lifetime` Max connection lifetime in seconds
- DATAPROVIDER_DB_POSTGRES_MAX_IDLE_CONNECTIONS :: `dataprovider.db.postgres.max_idle_connections` Maximum number of idle connections
- DATAPROVIDER_DB_POSTGRES_MAX_OPEN_CONNECTIONS :: `dataprovider.db.postgres.max_open_connections` Maximum number of connections
- DATAPROVIDER_DB_POSTGRES_TIMEOUT :: `dataprovider.db.postgres.timeout` Timeout for postgres connection
- DATAPROVIDER_DB_READ_URL :: `dataprovider.db.read.url` Database connection readonly URL
- DATAPROVIDER_DB_URL :: `dataprovider.db.url` Database connection URL
- DATAPROVIDER_GRPC_PORT :: `dataprovider.grpc.port` GRPC Port to listen on
- DATAPROVIDER_ID :: `dataprovider.id` Unique identifier for this services
- DATAPROVIDER_SECRETS_TOKEN_SIGNING_KEY :: `dataprovider.secrets.token_signing_key` Signing key used for service to service tokens
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
