# app/ohs/local/appservice

Role: use case boundary.

Allowed: transaction orchestration, domain service calls, repository port calls, and side-effect coordination.

Forbidden: parsing Fiber contexts or returning database models.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
