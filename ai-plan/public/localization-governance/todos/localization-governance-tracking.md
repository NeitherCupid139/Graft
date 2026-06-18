# Localization Governance Tracking

## Topic

Localization Governance

## Scope

建立前后端本地化治理规范，新增 AI 执行 skill，并分阶段把 server 侧硬编码 i18n 注册点迁移到资源文件。

## Repository Truth

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/design/本地化与i18n治理规范.md`
- `ai-plan/design/项目设计.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `.agents/skills/graft-localization-governance/SKILL.md`

## Current Recovery Point

- Phase 0-5 已全部完成。
- 当前恢复点：等待外层 loop 执行 archive-readiness check。
- 若未来出现 plural、复杂模板渲染、翻译平台导入导出或新增 locale 的真实需求，再以新 bounded batch 重开 provider 评估。

## Task Checklist

- [x] Phase 0：探索和迁移计划
- [x] Phase 0：设计规范持久化
- [x] Phase 0：public topic 建立
- [x] Phase 0：AI skill 建立
- [x] Phase 1：server embedded YAML loader
- [x] Phase 1：loader 单测
- [x] Phase 2：dashboard quick actions 样例迁移
- [x] Phase 3：system-config 批量迁移
- [x] Phase 4：展示文案迁移和治理测试
- [x] Phase 5：go-i18n provider 评估
