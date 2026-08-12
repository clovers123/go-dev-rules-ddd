# 数据库与实体军规

> **核心原则**: UUID 唯一标识。系统字段必选。类型严格统一。

## 1. 必选字段 (🔴 强制)

任何业务表必须包含以下 8 个标准系统底座字段：

| 字段 | SQL 类型 | Go GORM tag type | 说明 |
|---|---|---|---|
| `is_deleted` | `BOOLEAN NOT NULL DEFAULT FALSE` | `boolean` | 逻辑删除 |
| `create_time` | `TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()` | `timestamp without time zone` | 创建时间 |
| `update_time` | `TIMESTAMP WITHOUT TIME ZONE` | `timestamp without time zone` | 更新时间 |
| `delete_time` | `TIMESTAMP WITHOUT TIME ZONE` | `timestamp without time zone` | 删除时间 |
| `created_by` | `UUID` | `uuid` | 创建人 |
| `updated_by` | `UUID` | `uuid` | 更新人 |
| `deleted_by` | `UUID` | `uuid` | 删除人 |
| `feature` | `JSONB` | `jsonb` | 扩展字段（不可作为业务字段，仅运维用） |

> **教训**: 早期未约束类型导致 `owned_by`、`created_by` 等字段被误写为 `VARCHAR(36)` / `character varying(36)`，后续不得不编写幂等迁移脚本逐一修复。新增表必须严格按上表类型建列，禁止使用 `VARCHAR` / `character varying` 存储 UUID 值。

## 2. 主键与隔离 (🔴 强制)

- **主键**: 必须使用 **UUID** (SQL 类型 `UUID`，Go 映射为 `string`，GORM tag `type:uuid`)。不依赖数据库自增。
- **隔离**: 租户隔离字段名为 `owned_by`，类型 **必须为 `UUID`**（GORM tag `type:uuid`）。所有跨租户查询必须包含 `WHERE owned_by = ?`。

## 3. 外键 (❌ 禁止)

- **绝对禁止** 在数据库层面创建任何 `Foreign Key` 约束。关联一致性由 Domain 层保证。

## 4. 时间列规范 (🔴 强制)

- 所有时间列统一使用 `TIMESTAMP WITHOUT TIME ZONE`，禁止使用 `TIMESTAMP`（PostgreSQL 中等同于 `TIMESTAMP WITHOUT TIME ZONE`，但显式写出以消除歧义）。
- 默认值使用 `NOW()` 而非 `CURRENT_TIMESTAMP`，保持风格统一。

## 5. gorm-gen model 同步 (🔴 强制)

- 修改 SQL schema 后，必须同步重新生成 gorm-gen model（`app/acl/adapter/repository/postgres/model/*.gen.go`）。
- GORM tag 中的 `type` 必须与 SQL DDL 精确一致：UUID 列写 `type:uuid`，时间列写 `type:timestamp without time zone`。
- 禁止 GORM tag 中出现 `type:character varying(36)` 存储 UUID 值的情况。

---
**执行检查点**: 建表 SQL 必须包含上述 8 个系统字段，且类型严格遵循本规范。生成 model 后核对 GORM tag type 与 DDL 一致。
