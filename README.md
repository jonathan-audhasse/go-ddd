# go-ddd

Domain Driven Design pattern in go

## Architecture

goddd/
├── cmd/
│   ├── api/
│   │   └── main.go                     # wires everything together
│   └── migrate/
│       └── main.go                     # for DB migration
│
├── internal/
│   ├── api/                            # presentation layer
│   │   ├── handler/                    # handlers
│   │   │   ├── health_handler.go
│   │   │   └── user_handler.go
│   │   ├── httperror/                 # controllers
│   │   │   └── httperror.go
│   │   └── middleware/
│   │       └── authentication.go
│   │
│   ├── application/                    # use cases, orchestrates domain
│   │   ├── dto/                        # DTO (Data Transfer Object)
│   │   ├── health/                     # for health check
│   │   ├── transaction/                # the unit of work pattern
│   │   ├── user/                       # user application
│   │   └── service.go                  # the service interfaces
│   │
│   ├── domain/                  # pure business logic, no dependencies
│   │   ├── models/              # entities
│   │   │   ├── user.go
│   │   │   └── post.go
│   │   └── repository/
│   │       └── repository.go
│   │
│   └── infrastructure/                 # ← migrations and config
│       ├── config/
│       │   └── config.go               # cleanenv struct, ReadEnv()
│       ├── logger/
│       │   └── logger.go               # zerolog setup
│       ├── http/
│       │   ├── router/
│       │   │   └── router.go
│       │   └── middleware/
│       │       └── basic_auth.go
│       │
│       └── persistence/
│           ├── postgres/
│           │   ├── db.go               # Open(), connection pool setup
│           │   └── user_repository.go  # implements domain interface
│           └── migrations/
│               ├── migrate.go          # Up(), Down(), Steps()
│               └── files/              # ← migration scripts .sql files
│                   ├── 000001_create_users.up.sql
│                   ├── 000001_create_users.down.sql
│                   ├── 000002_add_posts.up.sql
│                   └── 000002_add_posts.down.sql
│
└── .golangci.yml

## Requirements

The following tools need to be installed in order to set up your local environment:

- [docker](https://docs.docker.com/get-started/get-docker)
- [go](https://go.dev/doc/install) (optional) if you do not develop in the container

## Mount your server

Local usage:

```sh
# to mount your service
make test
```

```sh
# to shut down your service
make test
```

## Tests

To run all tests (unit, integration and e2e tests) run the following command: 

```sh
make test
```

! you need to have [mockery](https://vektra.github.io/mockery/latest) installed locally in order to update the code.

## Help

For more help run:

```sh
make help
```
