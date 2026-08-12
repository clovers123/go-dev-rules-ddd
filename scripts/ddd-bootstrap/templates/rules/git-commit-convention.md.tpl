# Git 提交规范 (Git Commit Convention)

所有 Git 提交必须遵循以下格式：

`type(scope): 中文描述`

## 规则细节
1. **类型 (type)**: 必须使用英文，遵循 Conventional Commits 标准（例如：`feat`, `fix`, `refactor`, `chore`, `docs`, `style`, `test` 等）。
2. **范围 (scope)**: 必须使用英文，描述修改的代码模块或功能区域（例如：`db`, `api`, `domain`, `ui`, `sdk`, `service` 等）。
3. **描述 (description)**: 必须使用**中文**，清晰、具体地说明修改了什么（例如：“新增了某字段”、“修复了某逻辑问题”）。禁止使用模糊的描述。

## 正确示例
- `feat(db): 在 app_config 表中新增 test_receivers 字段`
- `refactor(service): 优化消息分发逻辑以支持测试白名单绕过配额`
- `fix(ui): 修复模板编辑页面表单回填失效的问题`

## 错误示例
- ❌ `feat(数据库): 增加字段` (Scope 必须为英文)
- ❌ `feat(db): add field` (描述必须为中文)
- ❌ `更新代码` (缺少类型和范围)
