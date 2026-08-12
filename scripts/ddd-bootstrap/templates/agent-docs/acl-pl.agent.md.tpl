# app/acl/pl

该目录是 **Anti-Corruption Layer (ACL) Presentation Layer**，主要负责领域实体 (Entity) 与外部模型（如数据库 Model、第三方 API DTO）之间的转换。

## 核心组件：AclAdapter

每个领域应包含一个 `I{Domain}AclAdapter` 接口及其实现。

### 职责
1. **读路径 (ToEntity)**：将数据库记录 (model) 或外部 DTO 转换为领域实体 (entity)。
2. **写路径 (FromEntity)**：将领域实体转换为可供持久化的数据库记录。
3. **批量转换**：必须提供处理切片的转换方法（如 `ToAggregations`），严禁在 Repository 中手写转换循环。

### 规则
- **禁止**：包含任何业务逻辑、数据库访问或网络调用。
- **强制**：所有的字段映射（包括审计字段、JSONB 序列化、空值处理）必须在此层完成。
- **命名**：统一使用 `{Domain}AclAdapter` 命名，方法名应具名化，如 `ToAppAggregation`。

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
