# app/ohs/pl

该目录定义了 **Open Host Service (OHS)** 的契约和适配器，包括 DTO、Command、Query、BO (Business Object) 以及 OHS Adapter。

## OHS Adapter 职责
- **输入转换**：将外部请求的 DTO 转换为 OHS 本地应用服务可接受的 `Command` 或 `Query`。在此过程中可以从 `fiber.Ctx` 提取租户信息、用户信息等。
- **输出转换**：将领域实体 (Entity) 转换为展示给外部的 `BO`。
- **批量转换**：提供 `From{Domain}Entities` 方法，处理列表数据的转换逻辑。

## 规范要求
- **方法命名**：输入转换使用 `FromXXXDTO`，输出转换使用 `FromXXXEntity/Entities`。
- **纯净性**：禁止在此层访问数据库、Redis、Kafka 或任何持久化模型。
- **解耦**：此层是领域层与外部世界的防火墙，确保外部 API 的变动不直接影响领域模型。

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
