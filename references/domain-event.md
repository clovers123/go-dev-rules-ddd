 # 领域事件规范

 > **核心原则**: 领域事件由领域对象产生，实现跨上下文的松耦合通信。

 ## 设计原则

 - 事件由领域对象（Entity）产生，标记聚合内的状态变更
- 事件在领域层定义，由 `app/ohs/local/appservice` 或 `app/infra` 中的运行时处理器订阅和分发
 - 用于跨上下文异步通信，避免直接的上下文间依赖

 ## 事件定义（领域层）

 ```go
 type LeaveApprovedEvent struct {
     EventID     string
     LeaveID     string
     ApplicantID string
     StartDate   time.Time
     EndDate     time.Time
     OccurredAt  time.Time
 }
 ```

 - 包含事件唯一标识、业务对象 ID、事件发生时间
 - 放在 `domain/event/` 目录下

 ## 事件发布（领域对象内）

 - 领域对象在状态变更的方法中发布事件
 - 事件发布通过接口注入，不直接依赖具体实现

 ```go
 func (a *LeaveApplication) Approve(approval ApprovalRecord) {
     a.Status = LeaveStatusApproved
     a.ApprovalRecords = append(a.ApprovalRecords, approval)

     a.EventPublisher.Publish(LeaveApprovedEvent{
         LeaveID:     a.ID,
         ApplicantID: a.ApplicantID,
         StartDate:   a.LeavePeriod.StartDate,
         EndDate:     a.LeavePeriod.EndDate,
     })
 }
 ```

## 事件订阅（OHS 本地应用服务或 Infra）

- 在 OHS 本地应用服务编排同步用例，在 Infra 处理异步消息订阅与运行时分发
 - 跨上下文事件通过消息中间件或事件总线实现

 ## 适用场景 vs 直接调用

 | 场景 | 推荐方式 |
 |------|---------|
 | 同一上下文内，同步强一致性 | 直接调用领域服务 |
 | 跨团队维护的上下文间通信 | 领域事件异步通信 |
 | 一套逻辑触发多个下游动作 | 领域事件（发布-订阅模式） |
 | 新增订阅者不可影响发布者 | 领域事件 |
 | 非实时处理（如非工作时间审批） | 领域事件异步处理 |

 ---
 **执行检查点**: 确认事件定义放在 `domain/event/`，事件由领域对象产生而非 OHS appservice；跨上下文通信优先使用事件而非直接调用。
