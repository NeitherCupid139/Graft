# 后端主导的 MVP 闭环收敛计划 Web 跟踪

## Subtopic

- Parent Topic: `mvp-extension-path`
- Subtopic: `web`
- Scope: starter 壳层、真实后端 `menu + route + page + api + permission` 契约挂接、i18n UI surface 与 frontend validation governance

## Goal

- 在后端主导的 MVP 收敛阶段，把 `web` 约束在 starter 壳层收敛与真实契约接线范围内，不再把新增页面或新的前端壳层广度作为近期目标。

## Current Recovery Point

- `web` 现阶段只保留一个可继续收敛的 starter 壳层基线，用于接入真实后端 `auth`、动态菜单、权限门禁和 locale 契约。
- 当前主线不是页面扩张，也不是继续深化独立前端工作台能力；任何 shell 级调整都应服务于真实契约挂接和 mock/demo 清理。
- `signals` 已收敛为文档级候选方案：`Pinia` 继续作为唯一正式共享状态层，当前不进入 `setting/theme` 局部试点，只保留未来最小 POC 的准入与退出规则。
- 前端命令真值保持不变：WSL 场景下继续使用 host Windows Bun，完成态仍以 `bun run check` 零 warning 为门槛。
- 详细前端实现历史保留在 `subtopics/web/traces/web-trace.md`。

## Active Risks

- 如果 `web` 回到页面扩张、长期保留 starter demo/mock 流程，前端会再次偏离“后端主导的 MVP 闭环收敛”主线。
- 如果后端共享契约在收敛期内继续频繁漂移，starter 壳层的真实接线会产生反复返工。
- 混用 WSL Bun 与 host Windows Bun 仍可能破坏当前工作树的前端依赖与 IDE 运行稳定性。

## Latest Validation

- 当前前端恢复基线沿用最近一次 host Windows Bun 完成态校验：
  - `bun run check`
- 该完成态基线要求 `format:check`、`typecheck`、`lint`、`stylelint`、`test:run`、`build` 全部通过且无未处理 warning。
- 本次文档同步没有新增前端运行时校验。
- 本次文档同步通过 `rg`、`sed` 与 `git diff -- ai-plan/design/前端架构设计.md ai-plan/public/mvp-extension-path/subtopics/web` 进行一致性检查。

## Immediate Next Step

- 继续把 starter 壳层挂接到真实后端 `auth + current user + menu + permission + locale` 契约。
- 快速隔离或移除当前阶段不再需要的 mock/demo 入口，避免形成前端自洽假闭环。
- 在真实契约稳定之前，不以新增页面、theme runtime 深化或额外视觉扩张作为当前子主题完成条件。
