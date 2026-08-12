# OHS, Domain, ACL, Infra Layer Mapping

## 1. OHS

OHS is both the external host boundary and the local use-case boundary.

- `remote/controller`: bind protocol input, call an OHS adapter, invoke an appservice, wrap output.
- `remote/routers`: define paths, methods, middleware, and route names only.
- `remote/middleware`: authentication, authorization, request metadata, protocol guards.
- `pl/{bc}`: define DTO, Command, Query, and BO contracts.
- `pl/{bc}/adapter`: convert protocol contracts to/from domain entities.
- `local/appservice/{bc}`: own use cases, transaction boundaries, workflows, and side-effect orchestration.

Appservices may coordinate multiple domain services and ACL ports. They must not contain core invariants, parse Fiber contexts, or depend on ACL implementations.

## 2. Domain

Domain owns business meaning and behavior.

- `entity`: aggregate roots and entities with identity and behavior.
- `valueobject`: immutable values without identity.
- `service`: business operations that do not naturally belong to one entity.
- `factory`: non-trivial domain construction.
- `event`: domain event definitions.

Domain services may depend on interfaces in `app/acl/port`. They must not import OHS packages, PostgreSQL models, generated DAO code, concrete adapters, Fiber, Redis, Kafka, or external SDK implementations.

## 3. ACL

ACL protects Domain and OHS from infrastructure-specific representations.

- `port`: interfaces expressed with domain semantics.
- `pl`: pure Entity-to-Model or Entity-to-external-record mapping.
- `adapter`: concrete persistence, client, publisher, subscriber, security, storage, or workflow implementations.

Ports must not expose GORM models, generated DAO interfaces, raw SQL fragments, or Fiber types. Repository implementations delegate query execution to named DAO methods and conversion to ACL PL.

## 4. Infra

Infra contains shared runtime and cross-cutting concerns:

- centralized system errors;
- technical constants, identifiers, enums, and validation;
- message runtime, consumers, scheduled jobs, telemetry, registration, and runtime glue.

Infra is not a miscellaneous business layer. Do not define aggregates, DTOs, repositories, or product-specific decision rules there. Infra jobs may depend on Domain types and ACL ports when driving asynchronous or scheduled use cases.

## Dependency Checkpoints

1. Confirm no `app/application` directory exists.
2. Confirm controllers call appservices rather than Domain services or repositories.
3. Confirm appservices depend on ACL ports, not implementations.
4. Confirm Domain imports no OHS or ACL adapter/model package.
5. Confirm repositories use ACL PL mapping and named DAO methods.
6. Confirm Infra contains technical runtime behavior rather than product modeling.
