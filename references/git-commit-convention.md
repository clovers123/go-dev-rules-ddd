 # Git 提交规范

> **核心原则**: 提交信息必须清晰描述修改内容，类型/范围/描述缺一不可。

## 格式

```
<type>(<scope>): <description>
```

## 要求

1. **类型 (type)**: 必须为以下之一：`feat`, `fix`, `refactor`, `chore`, `docs`, `style`, `test`, `perf`, `build`, `ci`, `revert`。
   - `feat` — 新增功能
   - `fix` — 修复 bug
   - `docs` — 仅文档更改
   - `style` — 不影响代码含义的格式调整（空白、缩进等）
   - `refactor` — 既不修复 bug 也不添加功能的代码重构
   - `perf` — 改进性能的代码更改
   - `test` — 添加或修改测试
   - `chore` — 杂项（构建流程、辅助工具等）
   - `build` — 构建系统或外部依赖变更（webpack、gulp、npm 等）
   - `ci` — CI/CD 流程配置修改（chart、Dockerfile 等）
   - `revert` — 回滚到上一个版本
2. **范围 (scope)**: 必须使用英文，描述修改的代码模块或功能区域（例如：`db`, `api`, `domain`, `ui`, `service` 等）。
3. **描述 (description)**: 必须使用**中文**，清晰、具体地说明修改了什么。禁止使用模糊描述。

## 单次提交注意事项

- 一次提交的问题必须为同一类别
- 单次提交不要超过 3 个问题
- 如果提交后发现不符合规范，使用 `git commit --amend -m "新提交信息"` 修正

## 正确示例

- `feat(db): 在 app_config 表中新增 test_receivers 字段`
- `refactor(service): 优化消息分发逻辑以支持测试白名单绕过配额`
- `fix(ui): 修复模板编辑页面表单回填失效的问题`
- `perf(api): 优化列表查询接口的索引使用，减少全表扫描`
- `build(docker): 升级基础镜像至 golang:1.22`

## 错误示例

- ❌ `feat(数据库): 增加字段` (Scope 必须为英文)
- ❌ `feat(db): add field` (描述必须为中文)
- ❌ `更新代码` (缺少类型和范围)

---
**执行检查点**: 检查提交信息是否符合 `<type>(<scope>): <description>` 格式，type 是否为允许值之一。
