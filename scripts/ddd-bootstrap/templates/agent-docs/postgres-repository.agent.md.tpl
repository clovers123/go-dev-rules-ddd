# app/acl/adapter/repository/postgres

该目录负责领域仓储 (Repository) 的 PostgreSQL 实现，基于 GORM Gen。

## 核心规则
1. **调用 DAO**：必须通过 `repository.Get{Table}Dao(ctx, tx)` 获取 DAO 实例，以支持事务透传。
2. **委托转换**：Repository 实现类应持有 `I{Domain}AclAdapter`。
    - 在 `Create/Update` 时，调用 Adapter 将 Entity 转为 Model。
    - 在 `Find/List` 时，调用 Adapter 将 Model 转为 Entity。
3. **禁止自留逻辑**：禁止在 Repository 实现中编写字段级别的映射逻辑或手动遍历切片进行转换。

## 示例
```go
func (r *PostgresAppRepository) List(ctx context.Context, tx *repository.Query, ...) ([]*entity.App, error) {
    dao := repository.GetAppDao(ctx, tx)
    records, err := dao.FindPage(...)
    if err != nil {
        return nil, err
    }
    // 正确做法：直接调用 adapter 的批量转换方法
    return r.adapter.ToAppAggregations(records), nil
}
```

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
