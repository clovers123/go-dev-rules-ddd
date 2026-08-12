 # Git 分支策略

 > **核心原则**: 分支职责明确，合并流程规范。禁止直接在 main 分支上修改。

 ## 分支命名规则

 | 分支类型 | 命名格式 | 基线分支 | 说明 |
 |---------|---------|---------|------|
 | 生产分支 | `main` | — | 稳定版本，禁止直接修改，仅接收 release 和 hotfix 合并 |
 | 开发分支 | `develop` | `main` | 最新开发版本，用于前后端联调 |
 | 功能分支 | `feature/{module}` | `develop` | 新功能开发，命名示例：`feature/user_module`、`feature/cart_module` |
 | 测试分支 | `test` | `develop` | 测试环境，给测试人员使用 |
 | 预发布分支 | `release` | `test` 或 `hotfix` | UAT 测试阶段，不建议直接修改代码 |
 | 紧急修复分支 | `hotfix/{issue}` | `main` | 线上紧急 bug 修复，修复后合并回 main 和 develop |

 ## 合并流程

 ```
 feature/* ──→ develop ──→ test ──→ release ──→ main
                                                     ↑
 hotfix/* ──────────────────────────────────────────┘
 ```

 1. 功能开发：从 `develop` 检出 `feature/{module}`，开发完成后合并回 `develop`
 2. 联调测试：`develop` → `test`
 3. 预发布：`test` → `release`（UAT 测试）
 4. 生产发布：`release` → `main`
 5. 紧急修复：从 `main` 检出 `hotfix/{issue}`，修复后合并到 `main` 和 `develop`

 ## 工作量规则

 - 工作量 < 1 天：可直接在 `develop` 分支开发
 - 工作量 ≥ 1 天：必须从 `develop` 检出 `feature/{module}` 分支开发

 ## 清理与维护

 - `feature/*` 和 `hotfix/*` 分支合并后应及时删除
 - `main` 分支每次更新（合并后）建议打 tag，标注版本号
 - 每次 `git pull` 前，先提交本地修改，避免合并冲突导致代码丢失

 ---
 **执行检查点**: 合并前确认基线分支正确，`main` 分支更新后检查是否已打 tag。
