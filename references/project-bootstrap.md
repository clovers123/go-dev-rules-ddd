# Four-Layer DDD Bootstrap

Use `scripts/ddd-bootstrap.sh`; do not hand-write project skeletons that the CLI can generate.

## Commands

```bash
scripts/ddd-bootstrap.sh init \
  --project-name <name> \
  --module <module-path> \
  --output <directory> \
  --db postgres

scripts/ddd-bootstrap.sh validate --project-root <directory>

scripts/ddd-bootstrap.sh add-domain \
  --project-root <directory> \
  --domain <context> \
  --aggregate <PascalName> \
  --schema <schema> \
  --table <table>

scripts/ddd-bootstrap.sh add-entity \
  --project-root <directory> \
  --domain <context> \
  --entity <entity> \
  --schema <schema> \
  --table <table>

scripts/ddd-bootstrap.sh remove-health --project-root <directory>
```

## Init Contract

`init` generates:

- only the bootstrap `health` context, never an inferred business domain;
- `app/ohs`, `app/domain`, `app/acl`, and `app/infra`;
- no `app/application` directory;
- application services at `app/ohs/local/appservice/health`;
- controllers at `app/ohs/remote/controller/health`;
- router files at `app/ohs/remote/routers`;
- layer-based Fx files under `cmd/server/di`;
- `configs/config.go` and `configs/config.yml`;
- SQL, GORM Gen inputs, local project rules, and directory `.agent.md` files.

Use `--auto-gormgen` for a temporary PostgreSQL container or `--use-existing-db` with explicit DB flags. If the database is unavailable, initialize with `--skip-go-mod-tidy`, report that generated DAO/model files remain pending, and do not claim GORM Gen succeeded.

## Add Domain Contract

`add-domain` generates a full vertical context:

- `app/domain/{bc}/{entity,service,valueobject,factory,event}`;
- `app/ohs/local/appservice/{bc}`;
- `app/ohs/pl/{bc}` and `adapter`;
- `app/ohs/remote/controller/{bc}`;
- `app/ohs/remote/routers/{bc}.go`;
- repository port, ACL PL mapper, PostgreSQL repository implementation;
- GORM Gen interface and one-table SQL file;
- registrations in `acl.go`, `domain.go`, `ohs.go`, and `routers.go`.

## Add Entity Contract

`add-entity` creates only the persistent entity chain:

- Domain entity;
- Repository port and implementation;
- ACL PL mapping;
- SQL and GORM Gen input;
- ACL and generator registrations.

It must not generate an appservice, controller, router, or Domain service for a sub-entity.

## Validation Contract

Validation must fail when:

- `app/application` exists;
- a required four-layer path or DI file is missing;
- a generated vertical flow is incomplete;
- Domain imports OHS or ACL implementations;
- controllers bypass appservices;
- repositories contain GORM chains, raw SQL, or direct model mapping;
- SQL omits tenant, soft-delete, or system fields;
- DI or router registration is missing.

Missing `.gen.go` files may remain warnings until GORM Gen is run. After database generation, rerun validation and `go test ./...`.
