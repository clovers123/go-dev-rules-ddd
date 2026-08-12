# {{.ProjectName}} Agent Dispatch

Read `.agent/rules/*.md` before planning. Read the nearest `.agent.md` before editing any directory.

Choose the workflow:

- New project: use `ddd-bootstrap init`.
- Add domain: use `ddd-bootstrap add-domain`; it must create the full OHS -> Domain -> ACL -> Infra-aware runtime -> SQL/GORM Gen path and update DI/router.
- Add persistent entity or value object: use `ddd-bootstrap add-entity`; it must not create appservice/controller/router for sub-entities.
- Normal feature: follow `.agent/rules/` and the nearest `.agent.md`; do not bypass generated boundaries.

Hard rules:

- Do not modify `.agent/rules/`.
- Use only `app/ohs`, `app/domain`, `app/acl`, and `app/infra`; never create `app/application`.
- Put application services in `app/ohs/local/appservice`.
- Do not hand-write or patch generated `.gen.go` DAO/model files.
- OHS code must enter through DTO/adapter/appservice and must not bypass adapters.
- Domain code must not import Fiber, controller packages, postgres models, or generated DAO code.
- Repository implementations must call named GORM Gen methods through `repository.GetXDao(ctx, tx)` or `repository.GetXDao(ctx, nil)` and must not write raw GORM chains.
- SQL and GORM Gen queries must keep `owned_by` and `is_deleted=false` visible unless a documented whitelist applies.
- New providers and route mounts must be registered in `cmd/server/di`.
- If this environment has gopls and the gopls MCP configured, use it to inspect
  package APIs, references, and diagnostics when that helps understand the
  project structure.

Health module:

`health` is only system initialization verification, not business code. It may be removed with `ddd-bootstrap remove-health --project-root .` after the project compiles and `validate` passes.

Acceptance:

Run `ddd-bootstrap validate --project-root .`. Resolve every FAIL and review every WARN, including missing generated DAO/model files, raw GORM leakage, adapter bypasses, missing DI/router registration, and missing `owned_by` / `is_deleted` guards.
