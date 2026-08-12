# cmd/gorm-gen

该目录负责项目的 **GORM Gen** 持久层代码生成。它采用模块化设计，每个领域都在其子目录中定义自定义查询接口。

## 目录结构
- `main.go`: 生成器的入口，负责加载配置、连接数据库、配置类型映射并调度各领域的生成逻辑。
- `{domain}/`: 每个领域的 GORM Gen 查询接口目录。
    - `{domain}.{table}.go`: 定义自定义的 SQL 查询方法，仅包含 SQL 注释和方法签名。接口类型名必须是模型名，如 `AppInfo`，不得添加 `Querier` 后缀。

## 操作指南

### 1. 添加新领域的生成逻辑
1. 在 `cmd/gorm-gen/` 下创建领域子目录（如 `cmd/gorm-gen/order/`）。
2. 创建 `{domain}.{table}.go`，定义查询接口：
   ```go
   package order

   import "gorm.io/gen"

   type OrderInfo interface {
       // FindByID
       //
       // SELECT * FROM @@table WHERE id=@id AND owned_by=@ownedBy AND is_deleted=false
       FindByID(id string, ownedBy string) (*gen.T, error)
   }
   ```
3. 在 `main.go` 中注册该领域，包括 Schema 策略、ApplyInterface 和类型映射：
   ```go
   g.WithTableNameStrategy(func(tableName string) string {
       return fmt.Sprintf("order.%s", tableName)
   })
   g.ApplyInterface(func(order.OrderInfo) {},
       g.GenerateModelAs("info", "OrderInfo",
           gen.FieldType("id", "string"),
           gen.FieldType("owned_by", "string"),
           gen.FieldType("is_deleted", "bool")),
   )
   ```

### 2. 规范要求
- **禁止手动修改**：`app/acl/adapter/repository/postgres/repository/` 下生成的 DAO 和 Model 文件严禁手动修改。
- **约束显式化**：在查询接口的 SQL 注释中，必须显式包含 `owned_by` 和 `is_deleted=false` 约束。
- **类型映射集中管理**：所有 `gen.FieldType` 类型覆盖统一在 `main.go` 的 `ApplyInterface` 中配置，不单独建立 Model 文件。
- **模型输出路径**：不要设置 `ModelPkgPath`；仅配置 `OutPath` 为 `./app/acl/adapter/repository/postgres/repository`，让 GORM Gen 自动写入同级 `model/` 目录。

## 执行生成
```bash
go run ./cmd/gorm-gen/main.go
```

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
