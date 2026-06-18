# Localization Governance

## 当前状态摘要

- 当前主题目标是建立 `Graft` 前后端本地化长期治理，并分阶段把 server 硬编码 i18n 文案迁移到资源文件。
- 状态：`ready-for-archive-check`。
- 任务分类为 `docs/automation` 起步，Phase 1-5 实际完成路径以 `server` 为主；若后续同步修改 web locale 或跨边界 contract，则升级为 `cross-boundary`。
- Canonical design：`ai-plan/design/本地化与i18n治理规范.md`。
- AI 执行 skill：`.agents/skills/graft-localization-governance/SKILL.md`。

## Recovery Receipt

- governance source：root `AGENTS.md`
- task class：`docs/automation` 起步；Phase 1/2 默认 `server`
- recovery source：`parent topic`
- authority summary：`ai-plan/design/本地化与i18n治理规范.md` + `server/internal/i18n.Service` facade + `web/src/locales/**` aggregation boundary

## Owned Scope

允许修改：

- `ai-plan/design/本地化与i18n治理规范.md`
- `ai-plan/public/localization-governance/**`
- `ai-plan/public/README.md`
- `.agents/skills/graft-localization-governance/**`
- 后续 Phase 按批次允许修改 `server/internal/i18n/**`、`server/internal/dashboard/**`、相关 module i18n 注册点和必要测试。

禁止误触：

- 不得让业务模块直接 import `github.com/nicksnyder/go-i18n`。
- 不得修改 HTTP 错误响应 wire shape。
- 不得删除 JSON Schema 现有 `x-i18n` key。
- 不得把 web 长期 UI 文案所有权迁到 server。
- 不得一次性迁移所有 `defaultCatalogEntries`。

## Phase Plan

- Phase 0：现状盘点、规范和 topic 持久化。已完成。
- Phase 1：server embedded YAML loader、locale 目录、loader 单测，保持 map catalog 和 facade 不变。已完成。
- Phase 2：dashboard quick actions system-config 文案样例迁移。已完成。
- Phase 3：system-config 相关文案批量迁移。已完成。
- Phase 4：菜单、通知、公告、scheduler、container、log explorer 展示文案迁移和治理测试。已完成。
- Phase 5：按真实需求评估是否接入 go-i18n provider。已完成，当前结论为“不引入”。

## Current Recovery Point

- 所有既定 Phase 0-5 batch 已完成，当前没有新的 in-scope pending batch。
- Phase 5 评估结论：
  - 当前 `server/internal/i18n` 已由 facade + map catalog + embedded flat YAML loader 覆盖真实 server 需求。
  - `LookupRequest.TemplateData` 仍是预留字段，没有已落地 plural、复杂模板或翻译平台工作流需求。
  - 现阶段引入 `go-i18n` 只会增加 provider 分叉与测试/治理成本，因此不作为当前 topic 的实现项。
- 外层 loop 下一步应执行 archive-readiness check：
  - 确认 Phase 0-5 commit 与验证证据完整。
  - 确认没有仍需在本主题内完成的 bounded server batch。
  - 若后续出现 plural/template/translation workflow 的真实需求，再以新 topic 或新 batch 重开 provider 评估。

## Validation Targets

```bash
git diff --check
python3 /root/.codex/skills/.system/skill-creator/scripts/quick_validate.py .agents/skills/graft-localization-governance
```

Phase 1 起追加：

```bash
cd server && go test ./internal/i18n/...
cd server && go run ./cmd/graft validate backend --stage lint
cd server && go build ./cmd/graft
```
