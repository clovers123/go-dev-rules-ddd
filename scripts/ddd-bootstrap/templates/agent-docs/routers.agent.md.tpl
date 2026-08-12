# app/ohs/remote/routers

Role: define HTTP paths, methods, middleware, and route names.

Forbidden: provider construction, business logic, repository access, or DTO mapping.

After adding a router, register it from cmd/server/di/routers.go.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
