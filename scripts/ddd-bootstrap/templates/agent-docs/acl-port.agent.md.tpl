# app/acl/port

Role: interfaces required by domain and appservice code.

Repository ports use domain entities and semantic parameters.

Forbidden: exposing Fiber, GORM models, or generated DAO types unless the project explicitly keeps tx pass-through in phase one.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
