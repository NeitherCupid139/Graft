# User Page Permission Governance Tracking

## Topic

- Topic: `user-page-permission-governance`
- Status: `active`
- Goal: converge the user-management page onto the existing centralized frontend permission visibility and runtime
  guard model without changing backend contracts or business behavior.
- Recovery source:
  - `ai-plan/public/README.md`
  - current RBAC/user management implementation
  - latest `rbac-visibility-governance` archive record
- Worktree: `/home/gewuyou/project/go/Graft-wt/feat/wt-rbac-further-development`
- Branch: `feat/wt-rbac-further-development`

## Scope

- Owned scope:
  - `ai-plan/public/user-page-permission-governance/**`
  - `ai-plan/public/README.md`
  - `web/src/modules/user/**`
  - `web/src/modules/rbac/**`
  - `web/src/store/modules/permission.ts`
  - `web/src/utils/route/**`
  - `web/src/app/bootstrap/**`
  - bounded tests/types only if directly required
- Forbidden scope:
  - `server/**`
  - OpenAPI/generated contract files
  - RBAC backend permission registry
  - unrelated layout/theme/i18n refactors

## Repository Truth

- `AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/design/AI任务追踪与恢复设计.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/public/rbac-visibility-governance/**`

## Governance Guardrails

- No backend API change.
- No backend permission-code change.
- No capability snapshot / denial reason / observability contract.
- No global UI restructuring.
- No user-management business behavior change.
- Missing permission-code gaps are recorded as risks, not solved cross-boundary in this topic.

## Current Recovery Point

- Batch 0 completed topic initialization.
- Topic docs now exist under `ai-plan/public/user-page-permission-governance/**`.
- The current codebase already has the canonical frontend permission building blocks:
  - bootstrap snapshot permission truth in `web/src/store/modules/permission.ts`
  - `v-permission` in `web/src/app/bootstrap/permission-directive.ts`
  - RBAC page reference usage in `web/src/modules/rbac/pages/index.vue`
- The user-management page still carries the main remaining drift from the archived RBAC visibility work:
  - page-local permission computed booleans in `web/src/modules/user/pages/index.vue`
  - visible-but-disabled permission entry points in batch and dropdown actions
  - missing explicit runtime permission guard checks on privileged handlers

## Batch Plan

1. Batch 0: topic initialization and current-state map. Status: completed.
2. Batch 1: user page permission implementation. Status: completed.
3. Batch 2: regression audit and consistency check. Status: pending.
4. Batch 3: archive-ready closeout. Status: pending.

## Batch 0 Findings Summary

- Canonical permission truth:
  - `permissionStore.hasPermission/hasAnyPermission/hasAllPermissions` are the stable helpers
  - `v-permission` is the stable element-visibility primitive
  - route visibility continues to flow from bootstrap menus and mounted runtime routes
- RBAC page reference pattern:
  - create/edit/assign-permission actions already use `v-permission`
  - disabled state remains for business-state constraints, not as the primary permission signal
- User page drift map:
  - header create and empty-state create depend on `canCreateUsers`
  - row assign-roles depends on `canReadUserRoles`
  - row edit depends on `canUpdateUsers`
  - batch assign-roles remains visible but disabled by `!canAssignUserRoles`
  - row dropdown action items remain visible but disabled by `!canDisableUsers` or `!canUpdateUsers`
  - operation-column visibility is aggregated by `canShowOperationColumn`
  - privileged action handlers do not enforce a local runtime guard before API invocation

## Immediate Next Step

- Batch 2 should audit the owned `user`/`rbac` frontend scope for remaining same-pattern permission drift.
- Preferred direction:
  - verify no page-local `canCreate/canUpdate/canDelete/canAssign` wrappers remain in the user page for template
    visibility only
  - verify no dropdown or batch action still uses visible-but-disabled semantics only because of missing permission
  - verify the user page runtime guards are aligned with the existing RBAC visibility strategy without widening scope

## Validation

- Required for Batch 0:
  - `git diff --check`
- Required for Batch 1:
  - `cd web && bun run check`
  - `git diff --check`

## Commit Eligibility Note

- Batch 1 owned scope is clear, but the worktree still includes uncommitted topic-initialization docs from Batch 0 plus
  Batch 1 tracking updates.
- A scoped commit remains possible, but the outer loop should decide whether to:
  - keep Batch 0/1 together in one owned commit, or
  - preserve batch readability and leave commit creation to a later closeout step
