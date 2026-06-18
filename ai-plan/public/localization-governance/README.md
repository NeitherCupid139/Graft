# Localization Governance

## 当前状态摘要

- 当前主题目标是建立 `Graft` 前后端本地化长期治理，并分阶段把 server 硬编码 i18n 文案迁移到资源文件。
- 当前 loop 状态：`batch-0-authority-reset-and-locale-directory-strategy` 进行中。
- 任务分类为 `server`；本轮只处理 authority reset、目录策略和 `server/internal/i18n` loader 前置收口，不触达 web。
- Canonical design：`ai-plan/design/本地化与i18n治理规范.md`。
- AI 执行 skill：`.agents/skills/graft-localization-governance/SKILL.md`。

## Recovery Receipt

- governance source：root `AGENTS.md`
- task class：`server`
- recovery source：`parent topic`
- authority summary：`ai-plan/design/本地化与i18n治理规范.md` + `server/internal/i18n.Service` facade + embedded locale YAML under `server/internal/i18n/locales/**`

## Owned Scope

允许修改：

- `ai-plan/design/本地化与i18n治理规范.md`
- `ai-plan/public/localization-governance/**`
- `ai-plan/public/README.md`
- `.agents/skills/graft-localization-governance/**`
- 本轮按 batch 允许修改 `server/internal/i18n/**` 与其直接测试，仅用于落实集中 locale 目录策略和 loader 行为。

禁止误触：

- 不得让业务模块直接 import `github.com/nicksnyder/go-i18n`。
- 不得修改 HTTP 错误响应 wire shape。
- 不得删除 JSON Schema 现有 `x-i18n` key。
- 不得把 web 长期 UI 文案所有权迁到 server。
- 不得把 Go 硬编码用户可见文案继续表述为可接受的长期 authority。
- 不得允许业务模块自行 embed、load、validate 或 freeze locale 文件。
- 不得一次性迁移所有 `defaultCatalogEntries`。

## Phase Plan

- Phase 0：现状盘点、规范和 topic 持久化。已完成。
- batch-0：authority reset、主题恢复材料纠偏、集中 locale 目录策略落定，以及 `server/internal/i18n` nested module loader 前置支持。
- slice-1：module registration resource migration。
- slice-2：core default catalog migration。
- slice-3：delete legacy fallbacks and switch to locale resource。
- final：archive readiness and governance sync。

## Current Recovery Point

- 旧的 `ready-for-archive-check` 判定已失效；topic 仍有 server-side authority 与迁移批次未完成。
- 当前 authority 决议：
  - backend 面向用户可见本地化文案的 canonical truth 是 embedded locale YAML。
  - locale 资源的 embed、load、validate、freeze 与 registry construction 只能集中在 `server/internal/i18n`。
  - module 只拥有 namespace/key 语义和调用 `i18n.Service` 的注册边界，不拥有独立 locale 文件加载器。
- 当前 batch 的目标：
  - 更新 design/README/tracking/trace/skill，移除过时 archive-ready 口径和旧 ownership 语言。
  - 将 backend locale 目录策略固定为 `server/internal/i18n/locales/*.yaml` + `server/internal/i18n/locales/modules/*.yaml`。
  - 若 loader 仍只支持顶层 `locales/*.yaml`，则在不改 facade 的前提下补齐 `locales/modules/*.yaml` 支持。

## Validation Targets

```bash
git diff --check
python3 /root/.codex/skills/.system/skill-creator/scripts/quick_validate.py .agents/skills/graft-localization-governance
```

若本轮触达 `server/internal/i18n/**`：

```bash
cd server && go test ./internal/i18n/...
cd server && go build ./cmd/graft
```
