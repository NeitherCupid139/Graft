# Localization Governance Trace

## 2026-06-18 Phase 0 governance persistence

- 建立本地化长期设计 authority：`ai-plan/design/本地化与i18n治理规范.md`。
- 建立 public recovery topic：`ai-plan/public/localization-governance/README.md`。
- 建立任务追踪：`ai-plan/public/localization-governance/todos/localization-governance-tracking.md`。
- 建立 trace：`ai-plan/public/localization-governance/traces/localization-governance-trace.md`。
- 建立仓库 skill：`.agents/skills/graft-localization-governance/SKILL.md`。
- 确认下一步用 `$graft-multi-agent-loop` 推进 Phase 1 起的 bounded batches。

## 2026-06-18 Phase 1-4 server migration completion

- Phase 1：已落地 `server/internal/i18n` embedded YAML loader、`locales/` 目录与单测，保持 `i18n.Service` facade、
  map catalog、`LookupRequest`、`Freeze`、`RegisteredMessageResources` 等外部 API 不变。
- Phase 2：已迁移 dashboard quick actions system-config 样例文案到 embedded locale YAML。
- Phase 3：已迁移剩余 server-side system-config 展示文案，保持 key 常量、fallback 字段与 JSON Schema `x-i18n`
  元数据不变。
- Phase 4：已迁移菜单、通知、公告、scheduler、container、log explorer 等展示文案到 locale resources，并保持
  provider 不外泄。

## 2026-06-18 Phase 5 go-i18n provider evaluation

- 以当前 `server/internal/i18n` 实现为 authority 完成 provider 评估。
- 证据摘要：
  - 当前 server i18n 结构是 facade + map catalog + embedded flat YAML loader。
  - `LookupRequest.TemplateData` 已预留，但当前调用面没有已落地 plural rules、复杂模板渲染或翻译平台导入导出流程。
  - Phases 1-4 已覆盖当前优先级内的 server locale 资源迁移，未出现必须引入第二 provider 才能继续推进的阻塞。
- 当前决策：不引入 `go-i18n`。
  - 原因是当前收益未被真实需求证明，而 provider 切换会增加实现分叉、测试矩阵和 provider 泄漏风险。
  - 未来只有在 plural、命名模板、翻译平台工作流或新增 locale 成为真实需求时，才重新开启 bounded provider 评估。

## 2026-06-18 Batch 0 authority reset and locale directory strategy

- 识别到 public README、tracking、trace 和 skill 仍保留 `ready-for-archive-check` 与过时 locale ownership 语言，属于 authority drift。
- 本轮将 design/topic/skill 统一重置为以下 backend authority：
  - embedded locale YAML 是 backend 用户可见本地化文案的 canonical truth。
  - `server/internal/i18n` 独占 locale 资源的 embed、load、validate、freeze 与 registry construction。
  - module 不得自持 locale 文件 embed/load 逻辑，也不得在 `server/internal/i18n` 外维护平行 registry。
- 目录策略落定为：
  - `server/internal/i18n/locales/*.yaml` 承载 core-owned namespace。
  - `server/internal/i18n/locales/modules/*.yaml` 承载 module-owned namespace。
  - 两类资源都只能由 `server/internal/i18n` 编译期 embed 并在启动期集中加载。
- 若实现侧仍只支持 `locales/*.yaml`，本轮同步补齐 `locales/modules/*.yaml` loader 支持，但不改 facade、provider exposure 或 wire contract。

## 2026-06-18 Slice 3 delete legacy fallbacks and switch to locale resource

- 删除以下 Go 用户可见 fallback，改由 locale resource authority 提供：
  - `server/internal/httpx/accesslog_explorer.go`：权限 `Name/Description`，菜单 `Title`。
  - `server/internal/logger/explorer.go`：权限 `Name/Description`，菜单 `Title`。
  - `server/modules/audit/module_registration.go`：权限 `Name/Description`，菜单 `Title`。
  - `server/modules/audit/dashboard_widget.go`：Widget/QuickLink/Action/item `Title/Description/Label/empty` fallback。
  - `server/internal/httpx/accesslog_dashboard.go`：Widget item `Title` 与 empty fallback。
  - `server/internal/httpx/accesslog_retention.go`、`server/internal/logger/retention.go`、`server/modules/audit/retention.go`：
    config definition `DomainLabel/GroupLabel/GroupDescription/Title/Description`，job `Title/ShortTitle/Description`，action `Title/Description`。
- 新增 locale authority：
  - `server/internal/i18n/locales/display.{zh-CN,en-US}.yaml`：`dashboard.actions.details`。
  - `server/internal/i18n/locales/modules/rbac.{zh-CN,en-US}.yaml`：`rbac.permissionCatalog.accessLogRead.*`、`appLogRead.*`、`appLogDelete.*`、`auditRead.*`。
- 本轮保留在 Go 的字符串仅作为技术标识：
  - permission code、menu code、module key、route/path、job name、action key、resource name、operation name、query key、状态枚举与内部日志消息。
- 临时例外：无。

## 2026-06-18 Slice 4 residual server visible-copy cleanup

- 删除以下 Go 用户可见 fallback，改为仅保留 locale key authority：
  - `server/modules/rbac/route_registration.go`：菜单 `Title`。
  - `server/internal/moduleruntime/registration.go`：菜单 `Title`。
  - `server/modules/monitor/module.go`：菜单 `Title`、unavailable evidence link `Title`。
  - `server/modules/system-config/module_registration.go`：菜单 `Title`。
  - `server/modules/container/module_registration.go`：菜单 `Title`。
  - `server/modules/announcement/module_registration.go`：菜单 `Title`。
  - `server/modules/user/module_registration.go`：菜单 `Title`。
  - `server/modules/scheduler/module_registration.go`：菜单 `Title`。
  - `server/modules/rbac/dashboard_widget.go`：QuickLink `Title`。
  - `server/modules/monitor/dashboard_widget.go`：Widget / QuickLink / summary / health-item label fallback。
  - `server/modules/system-config/dashboard_widget.go`：QuickLink `Title`。
  - `server/modules/scheduler/dashboard_widget.go`：Widget / QuickLink / action `Title` / `Description` / `Label` fallback。
  - `server/internal/dashboard/quick_actions_config.go`：config definition `DomainLabel` / `GroupLabel` / `GroupDescription` / `Title` / `Description` 与 schema 内嵌 title/description fallback。
  - `server/modules/container/config.go`：config definition `DomainLabel` / `GroupLabel` / `GroupDescription` / `Title` / `Description` 与 schema 内嵌 title/description fallback。
- 新增 locale authority：
  - `server/internal/i18n/locales/display.{zh-CN,en-US}.yaml`：`dashboard.widget.monitorSystemHealth.*`。
- 本轮明确保留在 Go 的技术或非用户可见字符串：
  - permission `Name` / `Description` 当前仍是 RBAC 管理数据与 seeded metadata 的英文稳定文本，不属于本批 residual visible-copy owner。
  - `monitor` health `Detail`、anomaly `Summary`、reason 文本属于运行时动态诊断文本；当前未建立独立 locale resource owner，不在本批直接迁移。
  - route/path、menu code、module key、config key、resource key、event type、dedupe key、enum/status、log message 继续作为技术标识保留。
- 临时例外：
  - `server/modules/scheduler/notification_integration.go`
    - 字段：`PublishNotificationInput.Message`、`PublishNotificationInput.ActionLabel`，以及失败通知的动态拼接 message。
    - 原因：`server/modules/notification/publisher.go` 目前对 `PublishNotificationInput` 强制校验 `Title` 与 `Message` 非空，通知消费链也仍持久化 fallback 文本；本批不允许改 notification wire / facade。
    - 移除条件：后续由 notification authority 批次把发布/持久化/消费链升级为 key-first，可接受空 fallback 或由 source-module locale lookup 生成最终文本。
    - 验证范围：`cd server && go test ./modules/scheduler -run 'TestSchedulerRunSuccessNotifierPublishesLocalizedDisplay|TestSchedulerRunFailureNotifierPublishesFailureNotification'`
  - `server/modules/audit/storeent/repository.go`
    - 字段：`displayTargetLabel()` 返回的 `用户/角色/权限/审计/服务器状态/认证`。
    - 原因：当前需先确认该标签的 canonical consumer authority 是 audit read-model 展示文案还是内部 evidence 分类标签；本批只做分类，不冒然改变 audit evidence 输出。
    - 移除条件：final archive-readiness 批次完成 consumer trace，若判定为用户可见文案则迁移到 locale resource；若判定为内部技术标签则在治理文档中正式归类。
    - 验证范围：依赖后续 audit owner 审计结论。

## Loop Batch State

```json
{
  "loop_mode": "topic-completion-loop",
  "completed_batches": [
    "batch-0-authority-reset-and-locale-directory-strategy"
  ],
  "pending_batches": [
    "slice-1-module-registration-resource-migration",
    "slice-2-core-default-catalog-migration",
    "slice-3-delete-legacy-fallbacks-and-switch-to-locale-resource",
    "final-archive-readiness-and-governance-sync"
  ],
  "current_batch": "batch-0-authority-reset-and-locale-directory-strategy",
  "next_batch": "slice-1-module-registration-resource-migration",
  "closeout_status": "batch-0-in-progress"
}
```
