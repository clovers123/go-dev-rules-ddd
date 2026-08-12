# app/infra

Role: shared runtime and cross-cutting infrastructure for errors, constants, validation, messaging, jobs, and observability.

Allowed: technical orchestration and reusable primitives. Forbidden: domain entities, product invariants, HTTP DTOs, and persistence model mapping.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`.
