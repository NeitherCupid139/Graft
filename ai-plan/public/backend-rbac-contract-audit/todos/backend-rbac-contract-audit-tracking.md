# Backend RBAC Contract Audit Tracking

## Topic

- Topic: `backend-rbac-contract-audit`
- Status: `active`
- Goal: audit the current RBAC permission/menu/API/guard contract closure across `server` and `web` without modifying
  runtime code in Batch 0.
- Branch: `feat/wt-rbac-further-development`
- Task class: `cross-boundary`
- Loop mode: `topic-completion-loop`

## Scope

- Owned scope:
  - `ai-plan/public/backend-rbac-contract-audit/**`
  - `ai-plan/public/README.md`
  - read-only audit of:
    - `server/plugins/rbac/**`
    - `server/internal/permission/**`
    - `server/internal/menu/**`
    - `server/internal/httpx/**`
    - `web/src/modules/rbac/**`
    - `web/src/modules/user/**`
    - `web/src/store/modules/permission.ts`
    - `web/src/utils/route/**`
- Forbidden scope:
  - runtime code changes
  - database schema / migrations
  - unrelated plugin code
  - OpenAPI/generated contract mutation without a recorded blocking mismatch
  - capability snapshot, denial reason, row-level permission, org/tenant model expansion

## Repository Truth

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `.ai/environment/tools.ai.yaml`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/design/AI任务追踪与恢复设计.md`
- archived:
  - `rbac-visibility-governance`
  - `user-page-permission-governance`
  - `frontend-permission-code-cleanup`

## Governance Guardrails

- Batch 0 is docs-only.
- Batch 0 must not modify runtime or test code.
- Record confirmed inventory separately from later audit questions.
- If a later batch discovers a real contract mismatch that requires broader change, record it before proposing any fix.

## Current Recovery Point

- Batch 0 completed topic initialization.
- Batch 0 created:
  - `ai-plan/public/backend-rbac-contract-audit/README.md`
  - `ai-plan/public/backend-rbac-contract-audit/todos/backend-rbac-contract-audit-tracking.md`
  - `ai-plan/public/backend-rbac-contract-audit/traces/backend-rbac-contract-audit-trace.md`
- Batch 0 updated `ai-plan/public/README.md` to register this topic as the active recovery entry.
- Batch 0 recorded the initial audit inventory and draft matrix for backend and frontend RBAC contract surfaces.

## Batch Plan

1. Batch 0: topic initialization and audit inventory. Status: completed.
2. Batch 1: backend permission/menu/API/guard audit. Status: pending.
3. Batch 2: frontend permission/route/action audit. Status: pending.
4. Batch 3: cross-boundary consistency audit. Status: pending.
5. Batch 4: MVP-stable decision and archive closeout. Status: pending.

## Batch 0 Findings

- backend RBAC owned permission registry currently exposes nine canonical permission codes:
  - `role.read`
  - `role.create`
  - `role.update`
  - `role.status.update`
  - `role.delete`
  - `role.permission.assign`
  - `permission.read`
  - `user.role.read`
  - `user.role.assign`
- backend RBAC owned menu declarations currently expose four entries:
  - `/access-control`
  - `/access-control/overview`
  - `/access-control/roles`
  - `/access-control/permissions`
- backend RBAC owned route registration currently wires explicit guards for:
  - role read/write routes
  - permission read routes
  - user-role snapshot and mutation routes
- backend guard semantics currently centralize through `httpx.RequirePermission(...)` and return denied
  `permission` detail on `403`.
- frontend owned permission constants currently converge on canonical names in:
  - `web/src/modules/rbac/contract/permissions.ts`
  - `web/src/modules/user/contract/permissions.ts`
- frontend owned bootstrap route registrations currently exist for:
  - `/access-control/users`
  - `/access-control/roles`
  - `/access-control/permissions`
- Batch 0 observed two follow-up consistency questions:
  - `/access-control/overview` backend menu exists but no owned page registration was found in current scope
  - `/access-control/users` frontend route exists but its backend menu owner is outside current Batch 0 RBAC backend
    read scope

## Immediate Next Step

- Execute `batch-1-backend-permission-menu-api-guard-audit`.
- Focus:
  - confirm whether current backend permission registry and route guards fully cover the intended RBAC management
    surfaces in owned scope
  - distinguish confirmed backend gaps from already-accepted cross-plugin ownership split

## Required Validation

- Batch 0:
  - `git diff --check`
- Later cross-boundary implementation batches:
  - required commands depend on whether runtime code changes occur

## Commit Plan

- Batch 0:
  - `docs(rbac-contract-audit): initialize audit topic`
