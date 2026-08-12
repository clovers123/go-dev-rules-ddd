# GO Dev Rules: DDD

一个面向 AI Agent 的 skill，用于强制实施 Go 后端开发规范：**四层 DDD 架构** —— `OHS`、`Domain`、`ACL`、`Infra`，应用服务（Application Service）放在 OHS 内部，而不是独立的 application 层。内置经过实战检验的规则集、项目脚手架 CLI 和架构校验契约。

> 适用于 AI Agent（OpenCode、Claude Code、Cursor 等）。在实现或评审 Go 后端功能、划分限界上下文（Bounded Context）、或搭建新的 DDD 项目时加载本 skill。

## 概览

架构源自两个生产级 Go 后端项目（`water-purifier-backend` 与 `notification-center-backend`）的成熟结构：

| 分层 | 路径 | 职责 |
|------|------|------|
| **OHS** | `app/ohs` | 外部协议与用例编排 |
| **Domain** | `app/domain` | 实体、工厂、领域服务、领域事件与业务不变量 |
| **ACL** | `app/acl` | 端口（Port）、模型映射、持久化与外部系统适配器 |
| **Infra** | `app/infra` | 横切运行时关注点：系统错误、常量、校验、消息、任务、可观测性 |

**核心决策：** 不要创建 `app/application`。应用服务放在 `app/ohs/local/appservice/{bc}`。

## 特性

- **四层依赖约束** —— 针对 import 方向、controller → appservice → ACL port 调用链、Repository 保持薄层的强制护栏。
- **规则参考文档** —— 14 份精选规则文档，覆盖分层映射、DTO/Command/Query/BO/Entity 转换、GORM Gen DAO、PostgreSQL 表结构、Fx DI 注册、Fiber 路由、领域事件、Go 编码规范与 Git 工作流。
- **项目脚手架 CLI** —— `scripts/ddd-bootstrap.sh` 一键生成完整四层骨架、垂直领域与带 SQL + GORM Gen 输入的持久化实体链。
- **架构校验** —— `validate` 子命令对架构违规直接判失败（如存在 `app/application`、Domain 引用了 ACL 实现、Repository 里出现 GORM 链式调用等）。
- **实战约定** —— UUID 主键与 `owned_by`、八个标准系统字段、不建数据库外键、`.gen.go` 生成代码禁止手改、`feature` 不得作为稳定业务字段。

## 安装

将仓库克隆到你的 skills 目录。以 OpenCode（用户级）为例：

```bash
git clone https://github.com/clovers123/go-dev-rules-ddd.git \
  ~/.config/opencode/skills/go-dev-rules-ddd
```

其他平台：将仓库整体（`SKILL.md` 位于根目录）放到对应 Agent 的 skills 路径下，例如 Claude Code 的 `~/.claude/skills/go-dev-rules-ddd`。

## 使用方式

当任务描述匹配时 skill 会自动激活 —— 包括 Go 后端功能开发、DDD 分层、Repository/DAO 设计、PostgreSQL schema 变更、Fx DI 注册或项目脚手架搭建。其强制工作流：

1. 先阅读目标仓库的 `AGENTS.md`、`.agent/rules/*.md` 与最近的 `.agent.md`；仓库本地规则优先于本可复用基线。
2. 将每个变更文件归类为 OHS、Domain、ACL、Infra、DI、配置、SQL 或 GORM Gen。
3. 编辑前先对照四层依赖规则检查 import。
4. 协议转换放在 OHS PL 适配器，持久化转换放在 ACL PL 适配器。
5. 查询逻辑写入具名 GORM Gen 方法；Repository 实现保持薄层。
6. 在 `cmd/server/di/` 对应的分层文件中注册新 provider，在 `routers.go` 中挂载路由。
7. Go 变更执行 `gofmt`，完成前运行 `go test ./...`。
8. 架构正确性需要单独复查 —— 测试通过并不代表分层放置正确。

## 项目脚手架

使用内置包装脚本，不要手写骨架文件：

```bash
# 创建新项目（仅生成 bootstrap 的 health 上下文）
scripts/ddd-bootstrap.sh init \
  --project-name demo \
  --module github.com/acme/demo \
  --output ./demo \
  --db postgres

# 校验骨架
scripts/ddd-bootstrap.sh validate --project-root ./demo

# 添加完整垂直领域（OHS -> Domain -> ACL -> SQL/GORM Gen）
scripts/ddd-bootstrap.sh add-domain \
  --project-root ./demo \
  --domain order \
  --aggregate Order \
  --schema commerce \
  --table orders

# 仅添加持久化实体链（不生成 appservice/controller/router）
scripts/ddd-bootstrap.sh add-entity \
  --project-root ./demo \
  --domain order \
  --entity order_item \
  --schema commerce \
  --table order_items

# 验证通过后移除 bootstrap 的 health 模块
scripts/ddd-bootstrap.sh remove-health --project-root ./demo
```

包装脚本在首次使用时自动编译 Go CLI（位于 `scripts/ddd-bootstrap/`）并缓存二进制。运行 `scripts/ddd-bootstrap.sh --help` 查看完整用法。

### CLI 生成内容

- `init` —— 四个分层、仅 `health` 上下文、`cmd/server/di/` 下的分层 Fx 文件、`configs/config.go` + `configs/config.yml`、SQL、GORM Gen 输入、项目本地规则与目录级 `.agent.md` 文件。
- `add-domain` —— 完整垂直链路：领域实体/service/valueobject/factory/event、appservice、OHS PL + adapter、controller、router、Repository port、ACL PL mapper、PostgreSQL Repository 实现、GORM Gen 接口，以及 DI/router 注册。
- `add-entity` —— 仅持久化实体链；绝不生成子实体的 appservice、controller 或 router。

## 仓库结构

```
SKILL.md                              Skill 入口（frontmatter + 规则）
references/                           供 Agent 消费的规则文档
scripts/ddd-bootstrap.sh              脚手架包装脚本（自动编译 CLI）
scripts/ddd-bootstrap/                Go CLI：init / validate / add-domain / add-entity / remove-health
agents/openai.yaml                    Agent 商店元数据
```

## 规则文档速查

| 文档 | 适用场景 |
|------|----------|
| `architecture-conventions.md`、`ddd-layer-mapping.md` | 新增或移动包之前 |
| `ddd-adapter.md` | DTO/Command/Query/BO/Entity/Model 转换 |
| `repository-dao.md`、`database-system-fields.md`、`sql-file-naming.md` | 持久化或 schema 工作 |
| `di-registration.md` | Fx provider、进程组装、路由注册 |
| `ohs-routing.md` | controller、Fiber、中间件、路由 |
| `feature-field-usage.md` | 读写 `feature` 之前 |
| `domain-event.md` | 跨上下文事件 |
| `go-coding-standards.md` | Go 源码变更 |
| `git-branch-strategy.md`、`git-commit-convention.md` | Git 工作流 |
| `project-bootstrap.md` | 项目脚手架搭建 |
