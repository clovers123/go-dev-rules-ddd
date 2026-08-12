---
trigger: always_on
---

# Repository and GORM Gen Rules

## Development Flow

1. Define named query methods with SQL comments in `cmd/gorm-gen/{bc}/{bc}.{table}.go`.
2. Run `go run ./cmd/gorm-gen`; never hand-edit `.gen.go` files.
3. Expose the generated DAO through `app/acl/adapter/repository/postgres/repository/export.{schema}.{table}.go`.
4. Define the semantic repository interface in `app/acl/port/repository/{bc}`.
5. Implement Entity-to-Model mapping in `app/acl/pl/{bc}`.
6. Implement the repository in `app/acl/adapter/repository/postgres/{bc}`.

## Port Boundary

Repository ports use Domain entities and semantic query values. Do not expose GORM models, generated DAO interfaces, SQL fragments, `*gorm.DB`, or generated `*repository.Query` through a new port.

For a multi-repository transaction, define or reuse a semantic Unit of Work port. Existing repositories may retain transaction pass-through only when repository-local rules explicitly require it; do not spread that implementation type into new Domain APIs.

## Repository Implementation

- Call only named DAO methods obtained through `repository.GetXxxDao(ctx, tx)` or `repository.GetXxxDao(ctx, nil)`.
- Do not use `.Where`, `.Order`, `.Limit`, `.Offset`, `.Like`, `.Or`, `.Join`, `.Raw`, or SQL string assembly.
- Do not access `repository.Q.Xxx` directly.
- Do not assign persistence-model fields in the repository; delegate to ACL PL.
- Do not return DTOs, BOs, Commands, or database models.
- Keep methods near ten lines; move query behavior into GORM Gen and mapping into ACL PL.
- Keep `owned_by` and `is_deleted=false` visible in named query definitions unless a documented global-query exception applies.

```go
func (r *TemplateRepository) FindByCode(ctx context.Context, code string, owner string) (*entity.Template, error) {
    dao := repository.GetTemplateDao(ctx, nil)
    row, err := dao.FindByCode(code, owner)
    if err != nil {
        return nil, err
    }
    return r.adapter.ToEntity(row), nil
}
```

Before changing a repository, confirm the named DAO method exists and the ACL mapper owns all model conversion.
