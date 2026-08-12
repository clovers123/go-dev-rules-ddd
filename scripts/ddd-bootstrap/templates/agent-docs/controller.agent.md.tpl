# app/ohs/remote/controller

该目录负责处理外部 HTTP 请求，是 **Open Host Service (OHS)** 的最外层。

## 核心流程
1. **Bind**: 使用 `ctx.Bind()` 绑定请求数据到 DTO。
2. **Adapt**: 调用 **OHS Adapter** 将 DTO 转换为领域层的 Command 或 Query。
3. **Execute**: 调用 **Application Service** 执行业务逻辑并获取结果 (BO/Error)。
4. **Respond**: 使用统一响应体包装返回结果。

## 规范要求
- **返回包装**：必须使用 `git.yugeeker.com/SHARED/go-lazy/app/entities` 中的 `resp.ResponseOK(data)` 或 `resp.ResponseFromError(err)`。
- **职责隔离**：禁止直接访问 Repository、Domain Service 或编写业务规则。
- **转换逻辑**：DTO 到 Command/Query 的映射必须委托给 OHS Adapter，Controller 仅做调度。

## 示例
```go
func (c *AppController) Create(ctx fiber.Ctx) error {
    var req pl.CreateAppRequest
    if err := ctx.Bind().Body(&req); err != nil {
        return ctx.JSON(resp.ResponseFromError(err))
    }
    cmd, err := c.adapter.FromCreateAppDTO(ctx, req)
    // ...
    return ctx.JSON(resp.ResponseOK(bo))
}
```

Allowed actions are the role-specific edits described above. Forbidden actions are the blocked edits listed here or in the related rules.

Related rules: `.agent/rules/ddd-layer-mapping.md`, `.agent/rules/ddd-adapter.md`, `.agent/rules/di-registration.md`, `.agent/rules/repository-dao.md`, `.agent/rules/ohs-routing.md`, `.agent/rules/sql-file-naming.md`.
