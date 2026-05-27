# RBAC Visibility Governance Tracking

## Topic

- Topic: `rbac-visibility-governance`
- Status: `active`
- Goal: strengthen the existing RBAC visibility closure path without introducing menu CRUD or resource CRUD.
- Worktree: `/home/gewuyou/project/go/Graft-wt/feat/wt-rbac-further-development`
- Branch: `feat/wt-rbac-further-development`

## Scope

- Owned scope:
  - `ai-plan/public/rbac-visibility-governance/**`
  - `ai-plan/public/README.md`
  - `server/plugins/rbac/**`
  - `server/internal/permission/**`
  - `server/internal/menu/**`
  - `server/internal/httpx/**`
  - `server/plugins/user/bootstrap.go`
  - `web/src/store/modules/permission.ts`
  - `web/src/utils/route/**`
  - `web/src/app/bootstrap/**`
  - `web/src/modules/rbac/**`
  - `web/src/modules/access-control/**`
  - bounded OpenAPI/generated contract files only if required

## Repository Truth

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/design/AI任务追踪与恢复设计.md`
- `ai-plan/design/项目设计.md`
- `ai-plan/design/插件与依赖注入设计.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/roadmap/MVP实施计划.md`

## Governance Guardrails

- No menu CRUD.
- No resource CRUD.
- No resource table.
- No migration of menu truth from registry/bootstrap to database CRUD.
- No hand-written API DTO truth that bypasses OpenAPI generated contract.

## Current Recovery Point

- Topic initialized on the dedicated RBAC worktree and branch pair.
- The current implementation direction is Option A only:
  - govern `permission -> bootstrap menus -> dynamic routes -> element visibility -> API guard`
  - avoid menu and resource management expansion
- Batch 1 read-only baseline audit completed.
- The current closure path is present end to end, but the main drift still sits on the `web` side:
  - compatibility-heavy access-control bootstrap normalization
  - frontend-owned menu hierarchy synthesis
  - legacy `title_key` rewriting
  - inconsistent button-level visibility conventions
  - one frontend permission-name alias for the same backend code

## Batch Plan

1. Batch 1: baseline audit and visibility chain map. Status: completed.
2. Batch 2: canonical bootstrap menu and route alignment. Status: next.
3. Batch 3: critical element permission coverage.
4. Batch 4: backend permission-guard consistency audit.
5. Batch 5: capability snapshot observability design.

## Immediate Next Step

- Execute Batch 2 under `graft-multi-agent-loop`.
- Remove compatibility-heavy frontend bootstrap normalization where current backend contracts are already canonical.
- Keep the scope limited to menu path/title_key/hierarchy alignment under Option A only.

## Batch 1 Findings Summary

- Current visibility chain map:
  - permission registry declarations feed backend route guards
  - bootstrap returns granted permissions plus permission-filtered menus
  - `web` permission store persists bootstrap snapshot and builds dynamic routes
  - layouts consume mounted dynamic routes for sidebar/header rendering
  - button-level visibility exists, but is not yet standardized on one mechanism
- Concrete drift points:
  - `web/src/utils/route/bootstrap.ts` still rewrites legacy `/users`, `/roles`, `/permissions` paths into `/access-control/*`
  - `web` still synthesizes access-control root and overview hierarchy instead of consuming one canonical upstream shape
  - `web` still rewrites historical `title_key` values into access-control keys
  - `v-permission` exists, but critical RBAC/user surfaces still depend mostly on per-page computed booleans
  - `web/src/modules/rbac/contract/permissions.ts` still keeps a semantic alias for the same backend permission code
