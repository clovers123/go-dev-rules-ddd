# app/ohs

Role: Open Host Service boundary for {{.ProjectName}}.

Allowed: remote protocol adapters, PL DTO/command/query/BO objects, local appservice use cases.

Forbidden: direct access to ACL adapter implementations, repository details, or database models.

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
