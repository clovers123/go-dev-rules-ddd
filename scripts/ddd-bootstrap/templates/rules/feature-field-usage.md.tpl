# Feature 字段使用规范

> **核心原则**: `feature` 是运维扩展字段，不是常规业务建模字段。

## 1. 允许场景

- 临时线上诊断、灰度标记、紧急运维开关。
- 与核心业务协议无关、可随时删除且不影响主流程的数据。

## 2. 禁止场景

- 禁止把前端渲染配置、SDK payload、模板默认展示样式放入 `feature`。
- 禁止把状态追踪、渠道类型、供应商、文件地址、模板 ID、模板编码等稳定业务字段放入 `feature`。
- 禁止让业务逻辑依赖 `feature` 中的 key 判断流程。

## 3. 建模要求

- 新增稳定业务字段时，必须先修改对应 `sql/{schema}/{schema}.{table}.sql`。
- 需要持久化的新字段必须同步更新实体、DTO/BO、ACL Adapter、OHS Adapter，并运行 `go run cmd/gorm-gen/main.go` 生成模型。
- 迁移历史 `feature` 数据时，可以在 SQL 文件中使用 `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` 和一次性 `UPDATE` 回填专用字段。

---
**执行检查点**: 任何新增 `Feature` 赋值或读取前，必须确认它不是常规业务字段；如果是，先建专用列或专用协议字段。
