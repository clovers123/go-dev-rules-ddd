---
name: go-dev-rules-ddd
description: Enforce Go backend development rules for a four-layer DDD architecture organized as OHS, Domain, ACL, and Infra, with application services inside OHS instead of a standalone application layer. Use when implementing or reviewing Go backend features, defining bounded contexts, placing entities/services/controllers/adapters, designing repository ports and GORM Gen DAOs, changing PostgreSQL schemas, registering Uber Fx dependencies and routes, or scaffolding and validating a new DDD project with init/add-domain/add-entity/remove-health.
---

# GO Dev Rules: DDD

## Overview

Apply a four-layer Go DDD architecture derived from the proven structure of `water-purifier-backend` and `notification-center-backend`:

1. `app/ohs` - external protocols and use-case orchestration.
2. `app/domain` - entities, factories, domain services, events, and invariants.
3. `app/acl` - ports, model mapping, persistence, and external-system adapters.
4. `app/infra` - cross-cutting runtime concerns such as system errors, constants, validation, messaging, jobs, and observability.

Do not create `app/application`. Put application services in `app/ohs/local/appservice/{bc}`.

## Rule Selection

- Read `references/architecture-conventions.md` and `references/ddd-layer-mapping.md` before adding or moving packages.
- Read `references/ddd-adapter.md` for DTO/Command/Query/BO/Entity/Model conversion.
- Read `references/repository-dao.md`, `references/database-system-fields.md`, and `references/sql-file-naming.md` for persistence or schema work.
- Read `references/di-registration.md` for Fx providers, process composition, or route registration.
- Read `references/ohs-routing.md` for controllers, Fiber, middleware, and routes.
- Read `references/feature-field-usage.md` before reading or writing `feature`.
- Read `references/domain-event.md` for cross-context events.
- Read `references/go-coding-standards.md` for Go source changes.
- Read `references/git-branch-strategy.md` and `references/git-commit-convention.md` for Git workflow.
- Read `references/project-bootstrap.md` and use `scripts/ddd-bootstrap.sh` for project scaffolding.

## Mandatory Workflow

1. Read the repository's `AGENTS.md`, `.agent/rules/*.md`, and nearest `.agent.md`; repository-local rules override this reusable baseline.
2. Classify each changed file into OHS, Domain, ACL, Infra, DI, configuration, SQL, or GORM Gen.
3. Check imports against the four-layer dependency rules before editing.
4. Keep protocol conversion in OHS PL adapters and persistence conversion in ACL PL adapters.
5. Put query logic in named GORM Gen methods; keep Repository implementations thin.
6. Register new providers in the matching layer file under `cmd/server/di/` and route mounts in `routers.go`.
7. Run `gofmt` for Go changes and `go test ./...` before completion.
8. Recheck architecture separately from compilation; passing tests does not prove correct layer placement.

## High-Priority Guardrails

- Use only the four top-level application layers: `ohs`, `domain`, `acl`, and `infra`.
- Never create a standalone `app/application` layer.
- Put use-case orchestration and transaction boundaries in `app/ohs/local/appservice/{bc}`.
- Put controllers in `app/ohs/remote/controller/{bc}` and router files in `app/ohs/remote/routers/`.
- Put DTOs, Commands, Queries, BOs, and protocol adapters in `app/ohs/pl/{bc}`.
- Put repository/client/publisher/subscriber/security interfaces in `app/acl/port`; do not expose GORM models or generated DAO types through ports.
- Domain code may depend on ACL port interfaces, but never ACL adapter implementations, PostgreSQL models, OHS packages, Fiber, or generated DAO code.
- Put Entity-to-Model mapping in `app/acl/pl/{bc}` and DTO/Command/Envelope-to-Entity mapping in `app/ohs/pl/{bc}/adapter`.
- Keep `app/infra` for cross-cutting runtime concerns; never use it as a miscellaneous business package.
- Organize Fx DI by layer: `acl.go`, `domain.go`, `infrastructure.go`, `ohs.go`, `routers.go`, and lifecycle/invoke files.
- Use root `configs/config.go` and `configs/config.yml`, matching both reference projects.
- Use UUID for IDs and `owned_by`; include the eight standard system fields; never create database foreign keys.
- Never use `feature` as a stable business field.
- Never hand-edit `.gen.go` files. Regenerate them with `go run ./cmd/gorm-gen`.
- Repository implementations must call named DAO methods through `repository.GetXxxDao(ctx, tx)` or `repository.GetXxxDao(ctx, nil)` and must not contain GORM chains or SQL assembly.
- Keep `owned_by` and `is_deleted=false` explicit in GORM Gen query definitions unless a documented global-query exception applies.
- Define user-facing and SDK errors centrally in `app/infra/syserrors`, using Chinese when required by repository rules.

## Scaffolding

Use the bundled wrapper instead of manually writing skeleton files:

```bash
scripts/ddd-bootstrap.sh init --project-name demo --module github.com/acme/demo --output ./demo --db postgres
scripts/ddd-bootstrap.sh add-domain --project-root ./demo --domain order --aggregate Order --schema commerce --table orders
scripts/ddd-bootstrap.sh add-entity --project-root ./demo --domain order --entity order_item --schema commerce --table order_items
scripts/ddd-bootstrap.sh validate --project-root ./demo
scripts/ddd-bootstrap.sh remove-health --project-root ./demo
```

`init` creates only the bootstrap `health` context. `add-domain` creates the complete OHS -> Domain -> ACL -> SQL/GORM Gen chain. `add-entity` creates only the persistent entity chain and must not create appservices, controllers, or routers.
