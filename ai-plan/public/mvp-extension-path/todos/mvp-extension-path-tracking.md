# 后端主导的 MVP 闭环收敛计划跟踪

## Topic

- Topic: `mvp-extension-path`
- Branch: `feat/mvp-extension-path`
- Scope: 以 `server` 主导完成 MVP 闭环，并推动 `web` 收敛到真实后端契约

## Goal

- 保持 `mvp-extension-path` 作为默认恢复入口，把当前阶段聚焦到后端闭环收敛和前端真实契约对接，而不是继续扩大会话治理或页面广度。

## Repository Truth

- `ai-plan/design/项目设计.md`
- `ai-plan/design/插件与依赖注入设计.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/代码注释与模块文档规范.md`
- `ai-plan/design/AI任务追踪与恢复设计.md`
- `ai-plan/roadmap/MVP实施计划.md`

## Topic Roadmap

- `ai-plan/public/mvp-extension-path/roadmap/backend-mvp-closure-plan.md`

## Subtopics

- `server`
  - Tracking: `ai-plan/public/mvp-extension-path/subtopics/server/todos/server-tracking.md`
  - Trace: `ai-plan/public/mvp-extension-path/subtopics/server/traces/server-trace.md`
  - Use for: backend runtime、plugin lifecycle、registries、Ent/Atlas、event bus、audit、scheduler、auth/RBAC 与跨插件契约稳定化。
- `web`
  - Tracking: `ai-plan/public/mvp-extension-path/subtopics/web/todos/web-tracking.md`
  - Trace: `ai-plan/public/mvp-extension-path/subtopics/web/traces/web-trace.md`
  - Use for: starter 壳层收敛、真实后端 `auth + menu + permission + locale` 契约挂接，以及 mock/demo 清理。

## Parent-Scope Rules

- 父主题只保留跨边界方向、共享风险、共享验证摘要和子主题入口指引。
- 纯 `server` 推进写入 `server` 子主题；纯 `web` 推进写入 `web` 子主题。
- 任何会改变共享契约、校验口径或阶段主线的变更，都要同时更新父主题与相关子主题。

## Current Recovery Point

- 当前阶段正式收敛到“后端主导的 MVP 闭环收敛计划”：先补齐 `server` 的 event bus、audit、scheduler 和稳定插件契约，再让 `web` 挂接这些真实契约。
- `server` 已有 runtime、plugin lifecycle、Ent/Atlas、auth/RBAC 与基础 smoke-validation 路径，但 MVP 还没有在事件流、审计链路、调度能力和跨插件稳定接口上闭环。
- `web` 当前阶段只做 starter 壳层收敛与真实后端契约接线，不再以新页面扩张、主题工作台深化或新的前端治理宽度作为近期主线。
- 较早的拆分前历史保留在 `archive/`，具体实现轨迹保留在各自 `trace` 文件。

## Shared Milestones

- `mvp-extension-path` 已稳定为父主题 + `server` / `web` 子主题的恢复结构。
- 后端已具备最小可运行的 plugin-oriented runtime、迁移链路、认证与 RBAC 基线。
- 前端已具备可继续收敛的 starter 壳层和 host Windows Bun `bun run check` 零 warning 完成门槛。
- 当前主题阶段已改为“后端先闭环，前端跟随真实契约收敛”。

## Shared Risks

- 如果 event bus、audit、scheduler 继续缺位，MVP 会停留在“可运行壳层”而不是“可扩展闭环”。
- 如果 `web` 回到页面扩张或长期保留 mock/demo 依赖，前后端契约会再次漂移。
- `auth`、`menu`、`permission`、`locale` 等共享契约若在后端收敛期内继续无边界扩张，会放大 `web` 接线返工成本。
- disposable PostgreSQL / Redis 校验仍依赖手工准备环境，后续恢复时必须显式说明当前可用的校验入口。

## Shared Validation Summary

- 本次仅同步 topic 级文档方向，没有新增代码或运行时校验。
- 当前跨边界恢复基线沿用 `server` 子主题中最近一次 focused backend validation，以及 `web` 子主题中最近一次 host Windows Bun `bun run check` 完成态校验。
- 本次文档同步通过 `sed`、`rg` 和 `git diff -- ai-plan/public/mvp-extension-path` 进行一致性检查。

## Immediate Next Step

- 按 `server` 主导顺序推进下一阶段：先稳定 event bus 边界，再补最小 audit 路径，再补 scheduler/plugin runtime 闭环，最后冻结当前 `web` 需要消费的 `auth + menu + permission + locale` 契约面。
