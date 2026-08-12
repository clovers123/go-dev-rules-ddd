# Canonical Four-Layer DDD Architecture

## Canonical Tree

```text
project-root/
├── app/
│   ├── ohs/
│   │   ├── local/appservice/{bc}/       # use cases, transactions, orchestration
│   │   ├── pl/{bc}/                     # DTO, Command, Query, BO
│   │   │   └── adapter/                 # protocol <-> domain mapping
│   │   └── remote/
│   │       ├── controller/{bc}/         # HTTP request orchestration
│   │       ├── routers/                 # one router file per context
│   │       └── middleware/              # protocol middleware
│   ├── domain/{bc}/
│   │   ├── entity/                      # aggregate roots and entities
│   │   ├── valueobject/                 # immutable values
│   │   ├── service/                     # domain services
│   │   ├── factory/                     # complex domain construction
│   │   └── event/                       # domain events
│   ├── acl/
│   │   ├── port/
│   │   │   ├── repository/{bc}/         # repository interfaces
│   │   │   ├── client/                  # external client interfaces
│   │   │   ├── publisher/               # event/message publisher interfaces
│   │   │   ├── subscriber/              # subscriber interfaces
│   │   │   └── security/                # security capability interfaces
│   │   ├── pl/{bc}/                     # domain entity <-> external model mapping
│   │   └── adapter/
│   │       ├── repository/postgres/{bc}/# repository implementations
│   │       ├── repository/postgres/model/
│   │       ├── repository/postgres/repository/
│   │       └── {client,publisher,...}/   # concrete external adapters
│   └── infra/
│       ├── syserrors/                    # centralized system errors
│       ├── constants/                    # technical constants and identifiers
│       ├── enums/                        # shared stable enums
│       ├── validation/                   # shared validation primitives
│       ├── mq/                           # message runtime
│       ├── prosumer/                     # asynchronous consumers/processors
│       └── cron/                         # scheduled jobs
├── cmd/
│   ├── server/di/
│   │   ├── infrastructure.go
│   │   ├── acl.go
│   │   ├── domain.go
│   │   ├── ohs.go
│   │   ├── routers.go
│   │   └── invoke.go
│   └── gorm-gen/{bc}/
├── configs/{config.go,config.yml}
└── sql/{schema}/{schema}.{table}.sql
```

Do not create `app/application`. Application services belong to OHS local because they expose local use cases and orchestrate Domain plus ACL ports.

## Placement Table

| Concern | Location |
|---|---|
| HTTP binding and response | `app/ohs/remote/controller/{bc}` |
| Route method/path/middleware | `app/ohs/remote/routers` |
| Use-case and transaction orchestration | `app/ohs/local/appservice/{bc}` |
| DTO, Command, Query, BO | `app/ohs/pl/{bc}` |
| Protocol-to-domain conversion | `app/ohs/pl/{bc}/adapter` |
| Aggregate root or entity | `app/domain/{bc}/entity` |
| Domain invariant spanning entities | `app/domain/{bc}/service` |
| Repository/external capability interface | `app/acl/port` |
| Domain-to-persistence model conversion | `app/acl/pl/{bc}` |
| PostgreSQL implementation | `app/acl/adapter/repository/postgres/{bc}` |
| Runtime errors/constants/jobs/messaging | `app/infra` |
| Fx composition | `cmd/{process}/di` |

## Import Rules

| Source | Allowed | Forbidden |
|---|---|---|
| OHS controller | OHS PL, OHS appservice, response helpers, Fiber | Repository implementations, DB models, direct domain service calls |
| OHS appservice | OHS PL/adapters, Domain, ACL ports, narrow Infra APIs | Fiber context, ACL implementations, generated DAO/model packages |
| OHS PL adapter | OHS PL types, Domain entities, protocol extraction helpers | Database, Redis, Kafka, repository implementations |
| Domain | Same-context Domain, ACL ports, standard library, narrow shared error/value primitives | OHS, ACL adapters, DB models, generated DAO, Fiber |
| ACL port | Domain semantic types and standard library | Concrete adapters, OHS, GORM models |
| ACL adapter | ACL ports/PL, Domain entities, external SDKs/models | OHS controllers/appservices |
| ACL PL | Domain entities and external models | IO, OHS, business orchestration |
| Infra | Technical libraries; Domain/ACL ports only when implementing runtime jobs | HTTP DTOs, product entity definitions, protocol conversion |

Keep package imports acyclic even though the architectural relationship is port-and-adapter rather than a single linear chain.
