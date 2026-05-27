# Localization Governance Tracking

## Topic

- Topic: `localization-governance`
- Status: `active recovery entry`
- Goal: establish the cross-boundary localization governance baseline before implementation work continues.
- Recovery source: `none`
- Worktree: `/home/gewuyou/project/go/Graft-wt/feat/wt-localization-governance`
- Branch: `feat/wt-localization-governance`

## Scope

- Owned scope:
  - `server/internal/i18n/**`
  - `server/internal/httpx/**`
  - `server/internal/contract/**`
  - `server/plugins/**` when changing localization contract registration or error/message key ownership
  - `web/src/locales/**`
  - `web/src/modules/**` when changing key-first locale consumption
  - `web/src/contracts/**`
  - `web/src/utils/request.ts`
  - `web/src/utils/route/**`
  - `openapi/**` when aligning key-field semantics only
  - `ai-plan/public/**`
  - related design docs when governance truth changes
- Task class: `cross-boundary`

## Repository Truth

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/design/项目设计.md`
- `ai-plan/design/插件与依赖注入设计.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/design/AI任务追踪与恢复设计.md`

## Current Recovery Point

- The prior worktree and branch name were still tied to the archived OpenAPI governance topic and were no longer valid recovery truth for this task.
- The active implementation workspace has now been renamed to:
  - Worktree: `/home/gewuyou/project/go/Graft-wt/feat/wt-localization-governance`
  - Branch: `feat/wt-localization-governance`
- This topic is the new active recovery entry for the project-wide localization governance task.
- Current governance decisions already frozen for the next implementation slice:
  - locale sources must not be assumed to come only from the host app source tree
  - locale keys require owner namespaces
  - `menu title_key`, `permission display_key`, and error `messageKey` are stable key contracts
  - frontend display must consume `key + fallback`
  - backend registry is registration/validation/fallback only, not a UI copy center
  - OpenAPI describes key fields and semantics only, not multilingual copy
  - current static plugin registration is the compile-time equivalent of a future dynamic plugin register API

## Shared Hotspots

- `ai-plan/public/README.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/项目设计.md`
- `server/internal/contract/**`
- `server/internal/i18n/**`
- `web/src/locales/**`
- `web/src/contracts/**`

## Immediate Next Step

- Continue with the first implementation slice under this new topic:
  - align frontend error and title rendering to consistent `key + fallback` consumption
  - identify and freeze missing contract owners for localization keys
  - update governance docs where current repository truth still implies host-app-only locale assumptions
