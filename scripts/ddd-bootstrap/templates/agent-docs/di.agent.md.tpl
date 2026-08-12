# cmd/server/di

Role: Uber Fx composition root.

Register providers, repositories, middleware, controllers, route mounts, and lifecycle hooks in the matching file.

Forbidden: business logic or request handling.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
