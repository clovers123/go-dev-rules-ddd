# SQL 文件命名与拆分规范

> **核心原则**: 每个数据库表必须拥有独立的 SQL 文件，严禁将多个表的 DDL 写在同一个文件中。

## 1. 命名规则 (🔴 强制)

- **格式**: `{schema}.{table_name}.sql`
- **示例**:
    - `msg.record.sql` (✅ 正确)
    - `msg.tracker.sql` (❌ 错误，tracker 是逻辑模块名而非表名)
    - `app.info.sql` (✅ 正确)

## 2. 存储位置

- 所有 SQL 文件必须存放在项目根目录下的 `sql/{schema}/` 目录中。
- 例如：`sql/msg/msg.record.sql`。

## 3. 内容规范

- 每个文件开头必须包含 `CREATE SCHEMA IF NOT EXISTS {schema};`。
- 必须包含 `CREATE TABLE IF NOT EXISTS "{schema}"."{table_name}"`。
- 必须在该文件中定义该表相关的所有索引（Index）。
- 严禁包含数据插入（INSERT）语句，仅限于 DDL。

## 4. GORM-Gen 联动

- 每新增一个 SQL 文件，必须在 `cmd/gorm-gen/main.go` 的 `executeSQLFile` 调用列表中手动添加该文件路径。
- 必须确保 `go run cmd/gorm-gen/main.go` 能够成功执行并生成模型。
