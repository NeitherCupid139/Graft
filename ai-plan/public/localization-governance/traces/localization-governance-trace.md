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
- 本主题当前进入 `ready-for-archive-check`，等待外层 loop 根据 commits、验证证据和 pending batch 状态做最终 archive-ready 判断。

## Loop Batch State

```json
{
  "loop_mode": "topic-completion-loop",
  "completed_batches": [
    "phase-0-design-topic-skill-persistence",
    "phase-1-server-loader-foundation",
    "phase-2-dashboard-quick-actions-sample",
    "phase-3-system-config-bulk-migration",
    "phase-4-display-message-migration-governance",
    "phase-5-go-i18n-provider-evaluation"
  ],
  "pending_batches": [],
  "current_batch": "phase-5-go-i18n-provider-evaluation",
  "next_batch": null,
  "closeout_status": "ready-for-archive-check"
}
```
