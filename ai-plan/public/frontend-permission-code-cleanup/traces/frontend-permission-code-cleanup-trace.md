# Frontend Permission Code Cleanup Trace

## 2026-05-27 Batch 0 initialized topic docs and recorded permission-code drift

- Reused the inherited startup receipt under root `AGENTS.md` for a `web` round inside `graft-multi-agent-loop`.
- Did not use `graft-multi-agent-batch`.
  - Reason: this round was a single-slice owned-scope audit plus documentation update; internal delegation had no
    practical payoff.
- Used `graft-task-closeout` style acceptance logic for owned-scope verification, validation, and commit eligibility.
- Read:
  - `.ai/environment/tools.ai.yaml`
  - `web/AGENTS.md`
  - `ai-plan/public/README.md`
  - `ai-plan/design/契约治理与魔法值治理规范.md`
  - `ai-plan/design/AI任务追踪与恢复设计.md`
  - `ai-plan/public/rbac-visibility-governance/README.md`
  - `ai-plan/public/user-page-permission-governance/README.md`
  - `web/src/modules/rbac/contract/permissions.ts`
  - `web/src/modules/user/contract/permissions.ts`
  - `web/src/modules/rbac/pages/index.vue`
  - `web/src/modules/user/pages/index.vue`
  - `web/src/store/modules/permission.ts`
- Confirmed the worktree was clean before edits and the current branch is:
  - `feat/wt-rbac-further-development`
- Created the new topic document set:
  - `ai-plan/public/frontend-permission-code-cleanup/README.md`
  - `ai-plan/public/frontend-permission-code-cleanup/todos/frontend-permission-code-cleanup-tracking.md`
  - `ai-plan/public/frontend-permission-code-cleanup/traces/frontend-permission-code-cleanup-trace.md`
- Audited owned-scope permission contracts and confirmed the canonical frontend permission-code map:
  - RBAC canonical value for role-permission assignment is `role.permission.assign`
  - user-role management permissions remain `user.role.read` and `user.role.assign`
  - user-management permissions remain `user.read`, `user.create`, `user.update`, `user.disable`,
    `user.session.read`, and `user.session.revoke`
- Identified the historical alias drift:
  - `RBAC_PERMISSION_CODE.ROLE_PERMISSION_MANAGE` is defined in `web/src/modules/rbac/contract/permissions.ts`
  - it resolves to the same value as `RBAC_PERMISSION_CODE.ROLE_PERMISSION_ASSIGN`
  - it is therefore a duplicate symbolic name rather than a distinct backend permission
- Identified owned-scope alias usage relevant to Batch 1:
  - `web/src/modules/rbac/pages/index.vue`
    - assign-permission button `v-permission` all-of check
    - `canAssignPermissions` computed guard
    - `canShowOperationColumn` computed guard
- Confirmed no runtime behavior change was needed in Batch 0.
- Chose Batch 1 bounds:
  - replace owned-scope alias references with `ROLE_PERMISSION_ASSIGN`
  - remove `ROLE_PERMISSION_MANAGE` after references are eliminated
  - avoid any new alias compatibility layer
- Kept `ai-plan/public/README.md` unchanged in Batch 0.
  - Basis: this round only initializes topic-local recovery material and the provided batch contract did not require a
    recovery-index mutation at initialization time.
