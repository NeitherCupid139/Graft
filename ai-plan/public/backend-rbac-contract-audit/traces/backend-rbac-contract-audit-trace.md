# Backend RBAC Contract Audit Trace

## 2026-05-27 Batch 0 initialized topic docs and recorded the initial audit inventory

- Reused the inherited startup context under root `AGENTS.md` for a `cross-boundary` retry worker round in
  `graft-multi-agent-loop`.
- Re-read the current-turn startup minimums:
  - `AGENTS.md`
  - `server/AGENTS.md`
  - `web/AGENTS.md`
  - `.ai/environment/tools.ai.yaml`
  - `ai-plan/public/README.md`
- Reused the required recovery sources from the inherited loop context:
  - archived `rbac-visibility-governance`
  - archived `user-page-permission-governance`
  - archived `frontend-permission-code-cleanup`
  - current RBAC backend implementation
  - current RBAC frontend implementation
- Confirmed the worktree was clean before Batch 0 writes.
- Created the new topic document set:
  - `ai-plan/public/backend-rbac-contract-audit/README.md`
  - `ai-plan/public/backend-rbac-contract-audit/todos/backend-rbac-contract-audit-tracking.md`
  - `ai-plan/public/backend-rbac-contract-audit/traces/backend-rbac-contract-audit-trace.md`
- Updated `ai-plan/public/README.md` to register `backend-rbac-contract-audit` as the active recovery topic for the
  current branch/worktree.
- Audited the owned backend RBAC contract surfaces:
  - `server/plugins/rbac/contract/permission.go`
  - `server/plugins/rbac/contract/route.go`
  - `server/plugins/rbac/plugin_registration.go`
  - `server/plugins/rbac/route_registration.go`
  - `server/internal/permission/registry.go`
  - `server/internal/menu/registry.go`
  - `server/internal/httpx/authz.go`
- Audited the owned frontend RBAC/user visibility surfaces:
  - `web/src/modules/rbac/contract/permissions.ts`
  - `web/src/modules/user/contract/permissions.ts`
  - `web/src/modules/rbac/bootstrap-routes.ts`
  - `web/src/modules/user/bootstrap-routes.ts`
  - `web/src/store/modules/permission.ts`
  - `web/src/utils/route/bootstrap.ts`
  - `web/src/modules/rbac/pages/index.vue`
  - `web/src/modules/rbac/pages/permissions/index.vue`
  - `web/src/modules/user/pages/index.vue`
- Recorded the initial inventory required by Batch 0:
  - backend permission registry inventory
  - backend menu declaration inventory
  - backend RBAC API route inventory
  - backend guard inventory
  - frontend permission constant inventory
  - frontend route/menu visibility inventory
  - frontend page/action permission usage inventory
- Drafted the initial RBAC contract audit matrix inside the topic README instead of scattering the first conclusions
  across trace-only notes.
- Recorded two bounded follow-up questions for later batches without widening scope:
  - backend menu `/access-control/overview` currently has no owned frontend page registration in the current read scope
  - frontend route `/access-control/users` is owned in current frontend scope, but its backend menu declaration owner
    is outside the current Batch 0 RBAC backend read scope
- Kept Batch 0 docs-only; no runtime or test code was changed.
