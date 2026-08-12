# app/acl

Role: Anti-Corruption Layer.

Ports define capabilities, PL maps between domain and external models, adapters implement external details.

Forbidden: letting upper layers depend on adapter implementation details.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
