# 数据库与实体军规

> **核心原则**: UUID 唯一标识。系统字段必选。

## 1. 必选字段 (🔴 强制)

任何业务表必须包含以下 8 个标准系统底座字段：

- `is_deleted`: 逻辑删除。
- `create_time / update_time / delete_time`: 审计时间。
- `created_by / updated_by / deleted_by`: 审计人。
- `feature`: JSONB 扩展字段。

## 2. 主键与隔离 (🔴 强制)

- **主键**: 必须使用 **UUID** (映射为 string)。不依赖数据库自增。
- **隔离**: 租户隔离字段名为 `owned_by`。所有跨租户查询必须包含 `WHERE owned_by = ?`。

## 3. 外键 (❌ 禁止)

- **绝对禁止** 在数据库层面创建任何 `Foreign Key` 约束。关联一致性由 Domain 层保证。

---
**执行检查点**: 建表 SQL 必须包含上述 8 个系统字段。
