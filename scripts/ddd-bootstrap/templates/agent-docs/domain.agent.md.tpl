# app/domain

Role: entities, domain services, factories, and business invariants.

Allowed: pure domain behavior and semantic validation.

Forbidden: Fiber, controller packages, DTO/BO packages, postgres models, generated DAO types, or infrastructure SDKs.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
